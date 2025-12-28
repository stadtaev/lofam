#!/bin/bash
set -e

PROJECT="lofam"
REGION="${AWS_REGION:-us-east-1}"
INSTANCE_TYPE="t3.micro"
KEY_NAME="${KEY_NAME:?KEY_NAME required}"

echo "=== Provisioning $PROJECT infrastructure ==="

# Get latest Amazon Linux 2023 AMI
AMI_ID=$(aws ec2 describe-images \
  --owners amazon \
  --filters "Name=name,Values=al2023-ami-*-x86_64" "Name=state,Values=available" \
  --query 'Images | sort_by(@, &CreationDate) | [-1].ImageId' \
  --output text)
echo "Using AMI: $AMI_ID"

# Security Group
SG_ID=$(aws ec2 describe-security-groups \
  --filters "Name=tag:Project,Values=$PROJECT" \
  --query 'SecurityGroups[0].GroupId' \
  --output text 2>/dev/null || echo "None")

if [ "$SG_ID" = "None" ] || [ -z "$SG_ID" ]; then
  echo "Creating security group..."
  SG_ID=$(aws ec2 create-security-group \
    --group-name "$PROJECT-sg" \
    --description "$PROJECT security group" \
    --query 'GroupId' \
    --output text)

  aws ec2 create-tags --resources "$SG_ID" --tags "Key=Project,Value=$PROJECT"

  # SSH
  aws ec2 authorize-security-group-ingress --group-id "$SG_ID" \
    --protocol tcp --port 22 --cidr 0.0.0.0/0
  # HTTP
  aws ec2 authorize-security-group-ingress --group-id "$SG_ID" \
    --protocol tcp --port 80 --cidr 0.0.0.0/0
  # HTTPS
  aws ec2 authorize-security-group-ingress --group-id "$SG_ID" \
    --protocol tcp --port 443 --cidr 0.0.0.0/0

  echo "Created security group: $SG_ID"
else
  echo "Security group exists: $SG_ID"
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
    --security-group-ids "$SG_ID" \
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

  aws ec2 create-tags --resources "$INSTANCE_ID" --tags "Key=Project,Value=$PROJECT" "Key=Name,Value=$PROJECT"

  echo "Created instance: $INSTANCE_ID"
  echo "Waiting for instance to be running..."
  aws ec2 wait instance-running --instance-ids "$INSTANCE_ID"
else
  echo "Instance exists: $INSTANCE_ID"
fi

# Elastic IP
EIP_ALLOC=$(aws ec2 describe-addresses \
  --filters "Name=tag:Project,Values=$PROJECT" \
  --query 'Addresses[0].AllocationId' \
  --output text 2>/dev/null || echo "None")

if [ "$EIP_ALLOC" = "None" ] || [ -z "$EIP_ALLOC" ]; then
  echo "Allocating Elastic IP..."
  EIP_ALLOC=$(aws ec2 allocate-address --domain vpc --query 'AllocationId' --output text)
  aws ec2 create-tags --resources "$EIP_ALLOC" --tags "Key=Project,Value=$PROJECT"
  echo "Created EIP: $EIP_ALLOC"
fi

# Associate EIP if not already
ASSOC_ID=$(aws ec2 describe-addresses \
  --allocation-ids "$EIP_ALLOC" \
  --query 'Addresses[0].AssociationId' \
  --output text 2>/dev/null || echo "None")

if [ "$ASSOC_ID" = "None" ] || [ -z "$ASSOC_ID" ]; then
  echo "Associating Elastic IP..."
  aws ec2 associate-address --instance-id "$INSTANCE_ID" --allocation-id "$EIP_ALLOC"
fi

PUBLIC_IP=$(aws ec2 describe-addresses \
  --allocation-ids "$EIP_ALLOC" \
  --query 'Addresses[0].PublicIp' \
  --output text)

echo ""
echo "=== Done ==="
echo "Instance ID: $INSTANCE_ID"
echo "Public IP: $PUBLIC_IP"
echo "SSH: ssh -i ~/.ssh/$KEY_NAME.pem ec2-user@$PUBLIC_IP"

# Output for GitHub Actions
if [ -n "$GITHUB_OUTPUT" ]; then
  echo "instance_id=$INSTANCE_ID" >> "$GITHUB_OUTPUT"
  echo "public_ip=$PUBLIC_IP" >> "$GITHUB_OUTPUT"
fi
