#!/bin/bash
set -e

# Usage: ./provision.sh [--mode alb|letsencrypt] --domain <domain>
# ALB mode (default): Uses ALB with ACM certificate
# Let's Encrypt mode: Direct EC2 with certbot

PROJECT="lofam"
REGION="${AWS_REGION:-us-east-1}"
INSTANCE_TYPE="t3.micro"
KEY_NAME="${KEY_NAME:?KEY_NAME required}"
MODE="alb"
DOMAIN=""

# Parse arguments
while [[ $# -gt 0 ]]; do
  case $1 in
    --mode) MODE="$2"; shift 2 ;;
    --domain) DOMAIN="$2"; shift 2 ;;
    *) echo "Unknown option: $1"; exit 1 ;;
  esac
done

if [ "$MODE" = "alb" ] && [ -z "$DOMAIN" ]; then
  echo "Error: --domain required for ALB mode"
  echo "Usage: $0 --mode alb --domain example.com"
  exit 1
fi

echo "=== Provisioning $PROJECT infrastructure (mode: $MODE) ==="

# Get default VPC
VPC_ID=$(aws ec2 describe-vpcs \
  --filters "Name=is-default,Values=true" \
  --query 'Vpcs[0].VpcId' \
  --output text)
echo "Using VPC: $VPC_ID"

# Get subnets (need at least 2 for ALB)
SUBNET_IDS=$(aws ec2 describe-subnets \
  --filters "Name=vpc-id,Values=$VPC_ID" \
  --query 'Subnets[*].SubnetId' \
  --output text)
SUBNET_ARRAY=($SUBNET_IDS)
echo "Using subnets: ${SUBNET_ARRAY[0]}, ${SUBNET_ARRAY[1]}"

# Get latest Amazon Linux 2023 AMI
AMI_ID=$(aws ec2 describe-images \
  --owners amazon \
  --filters "Name=name,Values=al2023-ami-*-x86_64" "Name=state,Values=available" \
  --query 'Images | sort_by(@, &CreationDate) | [-1].ImageId' \
  --output text)
echo "Using AMI: $AMI_ID"

# EC2 Security Group
EC2_SG_ID=$(aws ec2 describe-security-groups \
  --filters "Name=tag:Project,Values=$PROJECT" "Name=tag:Component,Values=ec2" \
  --query 'SecurityGroups[0].GroupId' \
  --output text 2>/dev/null || echo "None")

if [ "$EC2_SG_ID" = "None" ] || [ -z "$EC2_SG_ID" ]; then
  echo "Creating EC2 security group..."
  EC2_SG_ID=$(aws ec2 create-security-group \
    --group-name "$PROJECT-ec2-sg" \
    --description "$PROJECT EC2 security group" \
    --vpc-id "$VPC_ID" \
    --query 'GroupId' \
    --output text)
  aws ec2 create-tags --resources "$EC2_SG_ID" \
    --tags "Key=Project,Value=$PROJECT" "Key=Component,Value=ec2"

  # SSH from anywhere
  aws ec2 authorize-security-group-ingress --group-id "$EC2_SG_ID" \
    --protocol tcp --port 22 --cidr 0.0.0.0/0

  if [ "$MODE" = "letsencrypt" ]; then
    # Direct HTTP/HTTPS access
    aws ec2 authorize-security-group-ingress --group-id "$EC2_SG_ID" \
      --protocol tcp --port 80 --cidr 0.0.0.0/0
    aws ec2 authorize-security-group-ingress --group-id "$EC2_SG_ID" \
      --protocol tcp --port 443 --cidr 0.0.0.0/0
  fi

  echo "Created EC2 security group: $EC2_SG_ID"
else
  echo "EC2 security group exists: $EC2_SG_ID"
fi

# ALB-specific resources
if [ "$MODE" = "alb" ]; then
  # ALB Security Group
  ALB_SG_ID=$(aws ec2 describe-security-groups \
    --filters "Name=tag:Project,Values=$PROJECT" "Name=tag:Component,Values=alb" \
    --query 'SecurityGroups[0].GroupId' \
    --output text 2>/dev/null || echo "None")

  if [ "$ALB_SG_ID" = "None" ] || [ -z "$ALB_SG_ID" ]; then
    echo "Creating ALB security group..."
    ALB_SG_ID=$(aws ec2 create-security-group \
      --group-name "$PROJECT-alb-sg" \
      --description "$PROJECT ALB security group" \
      --vpc-id "$VPC_ID" \
      --query 'GroupId' \
      --output text)
    aws ec2 create-tags --resources "$ALB_SG_ID" \
      --tags "Key=Project,Value=$PROJECT" "Key=Component,Value=alb"

    # HTTP/HTTPS from anywhere
    aws ec2 authorize-security-group-ingress --group-id "$ALB_SG_ID" \
      --protocol tcp --port 80 --cidr 0.0.0.0/0
    aws ec2 authorize-security-group-ingress --group-id "$ALB_SG_ID" \
      --protocol tcp --port 443 --cidr 0.0.0.0/0

    # Allow EC2 to receive traffic from ALB
    aws ec2 authorize-security-group-ingress --group-id "$EC2_SG_ID" \
      --protocol tcp --port 80 --source-group "$ALB_SG_ID"

    echo "Created ALB security group: $ALB_SG_ID"
  else
    echo "ALB security group exists: $ALB_SG_ID"
  fi

  # ACM Certificate
  CERT_ARN=$(aws acm list-certificates \
    --query "CertificateSummaryList[?DomainName=='$DOMAIN'].CertificateArn | [0]" \
    --output text 2>/dev/null || echo "None")

  if [ "$CERT_ARN" = "None" ] || [ -z "$CERT_ARN" ]; then
    echo "Requesting ACM certificate for $DOMAIN..."
    CERT_ARN=$(aws acm request-certificate \
      --domain-name "$DOMAIN" \
      --validation-method DNS \
      --query 'CertificateArn' \
      --output text)
    echo "Certificate ARN: $CERT_ARN"
    echo ""
    echo "⚠️  ACTION REQUIRED: Validate certificate via DNS"
    echo "   Go to AWS Console → ACM → $DOMAIN → Create DNS record"
    echo "   Or run: aws acm describe-certificate --certificate-arn $CERT_ARN"
    echo ""
  else
    echo "ACM certificate exists: $CERT_ARN"
  fi
fi

# EC2 Instance
INSTANCE_ID=$(aws ec2 describe-instances \
  --filters "Name=tag:Project,Values=$PROJECT" "Name=instance-state-name,Values=running,pending,stopped" \
  --query 'Reservations[0].Instances[0].InstanceId' \
  --output text 2>/dev/null || echo "None")

if [ "$INSTANCE_ID" = "None" ] || [ -z "$INSTANCE_ID" ]; then
  echo "Creating EC2 instance..."

  INSTANCE_ID=$(aws ec2 run-instances \
    --image-id "$AMI_ID" \
    --instance-type "$INSTANCE_TYPE" \
    --key-name "$KEY_NAME" \
    --security-group-ids "$EC2_SG_ID" \
    --subnet-id "${SUBNET_ARRAY[0]}" \
    --associate-public-ip-address \
    --block-device-mappings '[{"DeviceName":"/dev/xvda","Ebs":{"VolumeSize":20,"VolumeType":"gp3"}}]' \
    --user-data '#!/bin/bash
dnf update -y
dnf install -y docker git
systemctl enable docker
systemctl start docker
usermod -aG docker ec2-user
curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
chmod +x /usr/local/bin/docker-compose
mkdir -p /home/ec2-user/app
chown ec2-user:ec2-user /home/ec2-user/app' \
    --query 'Instances[0].InstanceId' \
    --output text)

  aws ec2 create-tags --resources "$INSTANCE_ID" \
    --tags "Key=Project,Value=$PROJECT" "Key=Name,Value=$PROJECT"

  echo "Created instance: $INSTANCE_ID"
  echo "Waiting for instance to be running..."
  aws ec2 wait instance-running --instance-ids "$INSTANCE_ID"
else
  echo "Instance exists: $INSTANCE_ID"
fi

# Get instance private IP (for ALB target group)
PRIVATE_IP=$(aws ec2 describe-instances \
  --instance-ids "$INSTANCE_ID" \
  --query 'Reservations[0].Instances[0].PrivateIpAddress' \
  --output text)

# ALB resources
if [ "$MODE" = "alb" ]; then
  # Target Group
  TG_ARN=$(aws elbv2 describe-target-groups \
    --names "$PROJECT-tg" \
    --query 'TargetGroups[0].TargetGroupArn' \
    --output text 2>/dev/null || echo "None")

  if [ "$TG_ARN" = "None" ] || [ -z "$TG_ARN" ]; then
    echo "Creating target group..."
    TG_ARN=$(aws elbv2 create-target-group \
      --name "$PROJECT-tg" \
      --protocol HTTP \
      --port 80 \
      --vpc-id "$VPC_ID" \
      --target-type instance \
      --health-check-path "/" \
      --health-check-interval-seconds 30 \
      --query 'TargetGroups[0].TargetGroupArn' \
      --output text)
    echo "Created target group: $TG_ARN"
  else
    echo "Target group exists: $TG_ARN"
  fi

  # Register instance with target group
  aws elbv2 register-targets \
    --target-group-arn "$TG_ARN" \
    --targets "Id=$INSTANCE_ID" 2>/dev/null || true

  # Application Load Balancer
  ALB_ARN=$(aws elbv2 describe-load-balancers \
    --names "$PROJECT-alb" \
    --query 'LoadBalancers[0].LoadBalancerArn' \
    --output text 2>/dev/null || echo "None")

  if [ "$ALB_ARN" = "None" ] || [ -z "$ALB_ARN" ]; then
    echo "Creating Application Load Balancer..."
    ALB_ARN=$(aws elbv2 create-load-balancer \
      --name "$PROJECT-alb" \
      --subnets "${SUBNET_ARRAY[0]}" "${SUBNET_ARRAY[1]}" \
      --security-groups "$ALB_SG_ID" \
      --scheme internet-facing \
      --type application \
      --query 'LoadBalancers[0].LoadBalancerArn' \
      --output text)

    # Wait for ALB to be active
    echo "Waiting for ALB to be active..."
    aws elbv2 wait load-balancer-available --load-balancer-arns "$ALB_ARN"

    echo "Created ALB: $ALB_ARN"
  else
    echo "ALB exists: $ALB_ARN"
  fi

  # HTTP Listener (redirect to HTTPS)
  HTTP_LISTENER=$(aws elbv2 describe-listeners \
    --load-balancer-arn "$ALB_ARN" \
    --query "Listeners[?Port==\`80\`].ListenerArn | [0]" \
    --output text 2>/dev/null || echo "None")

  if [ "$HTTP_LISTENER" = "None" ] || [ -z "$HTTP_LISTENER" ]; then
    echo "Creating HTTP listener (redirect to HTTPS)..."
    aws elbv2 create-listener \
      --load-balancer-arn "$ALB_ARN" \
      --protocol HTTP \
      --port 80 \
      --default-actions 'Type=redirect,RedirectConfig={Protocol=HTTPS,Port=443,StatusCode=HTTP_301}'
  fi

  # HTTPS Listener
  HTTPS_LISTENER=$(aws elbv2 describe-listeners \
    --load-balancer-arn "$ALB_ARN" \
    --query "Listeners[?Port==\`443\`].ListenerArn | [0]" \
    --output text 2>/dev/null || echo "None")

  if [ "$HTTPS_LISTENER" = "None" ] || [ -z "$HTTPS_LISTENER" ]; then
    # Check if certificate is validated
    CERT_STATUS=$(aws acm describe-certificate \
      --certificate-arn "$CERT_ARN" \
      --query 'Certificate.Status' \
      --output text)

    if [ "$CERT_STATUS" = "ISSUED" ]; then
      echo "Creating HTTPS listener..."
      aws elbv2 create-listener \
        --load-balancer-arn "$ALB_ARN" \
        --protocol HTTPS \
        --port 443 \
        --certificates "CertificateArn=$CERT_ARN" \
        --default-actions "Type=forward,TargetGroupArn=$TG_ARN"
    else
      echo "⚠️  Certificate not yet validated (status: $CERT_STATUS)"
      echo "   HTTPS listener will be created after validation"
      echo "   Re-run this script after validating the certificate"
    fi
  fi

  # Get ALB DNS name
  ALB_DNS=$(aws elbv2 describe-load-balancers \
    --load-balancer-arns "$ALB_ARN" \
    --query 'LoadBalancers[0].DNSName' \
    --output text)
fi

# Elastic IP (for Let's Encrypt mode or SSH access)
if [ "$MODE" = "letsencrypt" ]; then
  EIP_ALLOC=$(aws ec2 describe-addresses \
    --filters "Name=tag:Project,Values=$PROJECT" \
    --query 'Addresses[0].AllocationId' \
    --output text 2>/dev/null || echo "None")

  if [ "$EIP_ALLOC" = "None" ] || [ -z "$EIP_ALLOC" ]; then
    echo "Allocating Elastic IP..."
    EIP_ALLOC=$(aws ec2 allocate-address --domain vpc --query 'AllocationId' --output text)
    aws ec2 create-tags --resources "$EIP_ALLOC" --tags "Key=Project,Value=$PROJECT"
  fi

  ASSOC_ID=$(aws ec2 describe-addresses \
    --allocation-ids "$EIP_ALLOC" \
    --query 'Addresses[0].AssociationId' \
    --output text 2>/dev/null || echo "None")

  if [ "$ASSOC_ID" = "None" ] || [ -z "$ASSOC_ID" ]; then
    aws ec2 associate-address --instance-id "$INSTANCE_ID" --allocation-id "$EIP_ALLOC"
  fi

  PUBLIC_IP=$(aws ec2 describe-addresses \
    --allocation-ids "$EIP_ALLOC" \
    --query 'Addresses[0].PublicIp' \
    --output text)
else
  # Get public IP from instance directly
  PUBLIC_IP=$(aws ec2 describe-instances \
    --instance-ids "$INSTANCE_ID" \
    --query 'Reservations[0].Instances[0].PublicIpAddress' \
    --output text)
fi

echo ""
echo "=== Done ==="
echo "Mode: $MODE"
echo "Instance ID: $INSTANCE_ID"
echo "Public IP: $PUBLIC_IP"
echo "SSH: ssh -i ~/.ssh/$KEY_NAME.pem ec2-user@$PUBLIC_IP"

if [ "$MODE" = "alb" ]; then
  echo ""
  echo "ALB DNS: $ALB_DNS"
  echo ""
  echo "Next steps:"
  echo "1. Create DNS CNAME: $DOMAIN → $ALB_DNS"
  echo "2. Validate ACM certificate (if not done)"
  echo "3. Deploy: docker-compose -f docker-compose.prod.yml up -d"
else
  echo ""
  echo "Next steps:"
  echo "1. Create DNS A record: $DOMAIN → $PUBLIC_IP"
  echo "2. Run SSL init: ./init-ssl.sh $DOMAIN your@email.com"
  echo "3. Deploy: docker-compose -f docker-compose.letsencrypt.yml up -d"
fi

# Output for GitHub Actions
if [ -n "$GITHUB_OUTPUT" ]; then
  echo "instance_id=$INSTANCE_ID" >> "$GITHUB_OUTPUT"
  echo "public_ip=$PUBLIC_IP" >> "$GITHUB_OUTPUT"
  echo "mode=$MODE" >> "$GITHUB_OUTPUT"
  [ "$MODE" = "alb" ] && echo "alb_dns=$ALB_DNS" >> "$GITHUB_OUTPUT"
fi
