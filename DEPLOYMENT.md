# Deployment Guide

This guide covers deploying Lofam to AWS EC2 using GitHub Actions with HTTPS.

## Infrastructure Modes

Choose your SSL termination mode:

| Mode | Description | Best For |
|------|-------------|----------|
| **ALB** (default) | Application Load Balancer with ACM certificate | Production, managed SSL |
| **Let's Encrypt** | Direct EC2 with certbot | Cost-sensitive, self-managed |

### ALB Mode (Recommended)

```
Internet → ALB:443 (SSL) → EC2:80 → nginx → backend/frontend
```

- AWS manages SSL certificates (ACM)
- Auto-renewal, no maintenance
- Easy to add health checks, multiple instances later
- Additional cost (~$16/month for ALB)

### Let's Encrypt Mode

```
Internet → EC2:443 (SSL) → nginx → backend/frontend
```

- Free SSL via Let's Encrypt
- Self-managed with auto-renewal via certbot
- Lower cost
- Requires Elastic IP for stable DNS

## Architecture

### ALB Mode
```
┌─────────────────────────────────────────────────────────┐
│                         AWS                             │
│  ┌─────────────┐      ┌─────────────────────────────┐  │
│  │    ALB      │      │       EC2 Instance          │  │
│  │   :443      │─────▶│  nginx:80 → frontend:3000   │  │
│  │  (ACM SSL)  │      │            → backend:8080   │  │
│  └─────────────┘      └─────────────────────────────┘  │
└─────────────────────────────────────────────────────────┘
```

### Let's Encrypt Mode
```
┌─────────────────────────────────────────────────────────┐
│                    EC2 Instance                         │
│  ┌─────────────────────────────────────────────────┐   │
│  │     nginx:443 (SSL via certbot)                 │   │
│  │    ┌─────────┬─────────┐                        │   │
│  │    │   /*    │  /api/* │                        │   │
│  │    ▼         ▼                                  │   │
│  │ frontend  backend                               │   │
│  │  :3000     :8080                                │   │
│  └─────────────────────────────────────────────────┘   │
│         │               │                              │
│   SSL certs        SQLite volume                       │
│  (Let's Encrypt)                                       │
└────────────────────────────────────────────────────────┘
```

## Prerequisites

1. AWS account with IAM user having EC2/ELBv2/ACM permissions
2. GitHub repository with Actions enabled
3. SSH key pair created in AWS Console
4. **Domain name** (required for SSL)

## Setup

### 1. Create AWS Key Pair

1. Go to AWS Console → EC2 → Key Pairs
2. Create key pair (e.g., `lofam-key`)
3. Download the `.pem` file

### 2. Create IAM User for GitHub Actions

1. Go to AWS Console → IAM → Users → Create user
2. Name: `github-actions`
3. Attach policies:
   - `AmazonEC2FullAccess`
   - `ElasticLoadBalancingFullAccess` (for ALB mode)
   - `AWSCertificateManagerFullAccess` (for ALB mode)
4. Click on the created user → Security credentials tab
5. Click **Create access key**
6. Select "Third-party service"
7. Save both keys (secret is shown only once!)

### 3. Configure GitHub Secrets

Go to your repo → Settings → Secrets and variables → Actions → New repository secret

| Secret | Value |
|--------|-------|
| `AWS_ACCESS_KEY_ID` | Access key from step 2 |
| `AWS_SECRET_ACCESS_KEY` | Secret key from step 2 |
| `AWS_KEY_NAME` | Key pair name (e.g., `lofam-key`) |
| `EC2_SSH_KEY` | Contents of the `.pem` file |

The EC2 host IP is automatically discovered via AWS tags.

### 4. Provision Infrastructure

1. Go to Actions → Infrastructure → Run workflow
2. Configure:
   - **Action**: `provision`
   - **Mode**: `alb` (default) or `letsencrypt`
   - **Domain**: Your domain (required for ALB mode)
3. Click "Run workflow"

This creates (idempotent):
- EC2 t3.micro instance
- Security group(s)
- **ALB mode**: ALB, Target Group, ACM Certificate
- **Let's Encrypt mode**: Elastic IP

All resources are tagged with `Project=lofam`.

### 5. Configure DNS

**ALB Mode:**
1. Get ALB DNS from workflow summary
2. Create CNAME record: `yourdomain.com` → `lofam-alb-xxxxx.us-east-1.elb.amazonaws.com`

**Let's Encrypt Mode:**
1. Get Elastic IP from workflow summary
2. Create A record: `yourdomain.com` → `<elastic-ip>`

### 6. Validate SSL Certificate

**ALB Mode:**
1. Go to AWS Console → ACM → Certificates
2. Find certificate for your domain (status: Pending validation)
3. Click "Create records in Route 53" or manually add the CNAME validation record
4. Wait for status to change to "Issued"
5. Re-run Infrastructure workflow to create HTTPS listener

**Let's Encrypt Mode:**
SSH into the instance and run:
```bash
ssh -i ~/.ssh/lofam-key.pem ec2-user@<elastic-ip>
cd ~/app
git clone https://github.com/your/repo .
./infrastructure/aws/init-ssl.sh yourdomain.com your@email.com
```

### 7. Deploy Application

Push to `main` branch triggers automatic deployment:
1. Runs Go tests
2. Detects infrastructure mode (ALB or Let's Encrypt)
3. Discovers EC2 IP via AWS tag query
4. SSHs to EC2, pulls latest code
5. Runs appropriate docker-compose file

Or manually: Actions → Deploy → Run workflow

## Docker Compose Files

| File | Mode | Description |
|------|------|-------------|
| `docker-compose.prod.yml` | ALB | Simple nginx HTTP proxy (ALB terminates SSL) |
| `docker-compose.letsencrypt.yml` | Let's Encrypt | nginx with SSL + certbot container |

The deploy workflow auto-detects which file to use based on whether an ALB exists.

## SSL Certificate Renewal

**ALB Mode:** Automatic - ACM handles renewal.

**Let's Encrypt Mode:** Automatic via certbot container (checks every 12h).

Manual renewal:
```bash
docker-compose -f docker-compose.letsencrypt.yml exec certbot certbot renew
docker-compose -f docker-compose.letsencrypt.yml exec nginx nginx -s reload
```

## Workflows

### Infrastructure (`infra.yml`)

Manual trigger only. Inputs:
- **action**: `provision` or `destroy`
- **mode**: `alb` (default) or `letsencrypt`
- **domain**: Your domain name (required for ALB)

### Deploy (`deploy.yml`)

Triggers on:
- Push to `main` branch
- Manual dispatch

Auto-detects mode by checking for ALB existence.

## Manual Deployment

```bash
# SSH into instance
ssh -i ~/.ssh/lofam-key.pem ec2-user@<public-ip>
cd ~/app
git pull origin main

# ALB mode
docker-compose -f docker-compose.prod.yml up --build -d

# Let's Encrypt mode
docker-compose -f docker-compose.letsencrypt.yml up --build -d

# View logs
docker-compose -f docker-compose.prod.yml logs -f
```

## Costs

### ALB Mode

| Resource | Free Tier | After Free Tier |
|----------|-----------|-----------------|
| EC2 t3.micro | 750 hrs/month (12 months) | ~$8/month |
| ALB | Not included | ~$16/month |
| ACM Certificate | Free | Free |
| EBS (20GB gp3) | 30GB free (12 months) | ~$2/month |

**Estimated monthly cost**: ~$16 (free tier) or ~$26/month

### Let's Encrypt Mode

| Resource | Free Tier | After Free Tier |
|----------|-----------|-----------------|
| EC2 t3.micro | 750 hrs/month (12 months) | ~$8/month |
| Elastic IP | Free when attached | ~$4/month if unattached |
| EBS (20GB gp3) | 30GB free (12 months) | ~$2/month |
| Let's Encrypt | Free | Free |

**Estimated monthly cost**: $0 (free tier) or ~$10-14/month

## Troubleshooting

### ALB Certificate Pending

```bash
# Check certificate status
aws acm describe-certificate --certificate-arn <arn> --query 'Certificate.Status'

# List validation records needed
aws acm describe-certificate --certificate-arn <arn> \
  --query 'Certificate.DomainValidationOptions[*].ResourceRecord'
```

After adding DNS validation record, re-run Infrastructure workflow.

### ALB Health Check Failing

```bash
# SSH and check nginx responds
curl http://localhost/health

# Check target group health
aws elbv2 describe-target-health --target-group-arn <tg-arn>
```

### Let's Encrypt Certificate Issues

```bash
# Check certificate status
docker-compose -f docker-compose.letsencrypt.yml exec certbot certbot certificates

# View certbot logs
docker-compose -f docker-compose.letsencrypt.yml logs certbot

# Force renewal
docker-compose -f docker-compose.letsencrypt.yml exec certbot certbot renew --force-renewal
```

### SSH Connection Refused

```bash
# Check security group allows port 22
aws ec2 describe-security-groups --filters "Name=tag:Project,Values=lofam"
```

### Docker Not Running

```bash
sudo systemctl status docker
sudo systemctl start docker
```

### Containers Not Starting

```bash
# Check logs
docker-compose -f docker-compose.prod.yml logs

# Rebuild
docker-compose -f docker-compose.prod.yml down
docker-compose -f docker-compose.prod.yml up --build -d
```

## Switching Modes

To switch from one mode to another:

1. Run Infrastructure workflow with `destroy` action
2. Run Infrastructure workflow with `provision` and new mode
3. Update DNS records appropriately
4. Push to main to deploy

**Note**: This will cause downtime. Data in SQLite volume persists on EC2.

## Destroy Infrastructure

1. Go to Actions → Infrastructure → Run workflow
2. Select action: `destroy`
3. Click "Run workflow"

This terminates:
- EC2 instance
- ALB and Target Group (if ALB mode)
- Releases Elastic IP (if Let's Encrypt mode)
- Deletes security groups

**Note**: ACM certificates are preserved (can be reused).
