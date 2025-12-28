# Deployment Guide

This guide covers deploying Lofam to AWS EC2 using GitHub Actions with HTTPS.

## Architecture

```
┌─────────────────────────────────────────┐
│              EC2 Instance               │
│  ┌─────────────────────────────────┐    │
│  │     nginx:443 (HTTPS)           │    │
│  │    ┌─────────┬─────────┐        │    │
│  │    │   /*    │  /api/* │        │    │
│  │    ▼         ▼         │        │    │
│  │ frontend  backend      │        │    │
│  │  :3000     :8080       │        │    │
│  └─────────────────────────────────┘    │
│         │               │               │
│   SSL certs        SQLite volume        │
│  (Let's Encrypt)                        │
└─────────────────────────────────────────┘
```

Port 80 redirects to HTTPS. Certificates auto-renew via certbot.

## Prerequisites

1. AWS account with IAM user having EC2 permissions
2. GitHub repository with Actions enabled
3. SSH key pair created in AWS Console
4. **Domain name** pointing to your Elastic IP (required for SSL)

## Setup

### 1. Create AWS Key Pair

1. Go to AWS Console → EC2 → Key Pairs
2. Create key pair (e.g., `lofam-key`)
3. Download the `.pem` file

### 2. Create IAM User for GitHub Actions

1. Go to AWS Console → IAM → Users → Create user
2. Name: `github-actions`
3. Attach policy: `AmazonEC2FullAccess`
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

The EC2 host IP is automatically discovered by querying AWS for the Elastic IP tagged with `Project=lofam`.

### 4. Provision Infrastructure

1. Go to Actions → Infrastructure → Run workflow
2. Select action: `provision`
3. Click "Run workflow"

This creates (idempotent):
- EC2 t3.micro instance
- Security group (ports 22, 80, 443)
- Elastic IP

All resources are tagged with `Project=lofam`.

### 5. Configure DNS

Point your domain to the Elastic IP:
- A record: `yourdomain.com` → `<elastic-ip>`
- Or CNAME if using subdomain

### 6. Initialize SSL Certificate

SSH into the instance and run the SSL init script:

```bash
ssh -i ~/.ssh/lofam-key.pem ec2-user@<elastic-ip>
cd ~/app
./infrastructure/aws/init-ssl.sh yourdomain.com your@email.com
```

This:
1. Starts temporary nginx for ACME challenge
2. Requests certificate from Let's Encrypt
3. Stores certificate in Docker volume

### 7. Deploy Application

Push to `main` branch triggers automatic deployment:
1. Runs Go tests
2. Discovers EC2 IP via AWS tag query
3. SSHs to EC2, pulls latest code
4. Runs `docker-compose -f docker-compose.prod.yml up --build -d`

Or manually: Actions → Deploy → Run workflow

## SSL Certificate Renewal

Certificates auto-renew via the certbot container, which checks every 12 hours.

To manually renew:

```bash
docker-compose -f docker-compose.prod.yml exec certbot certbot renew
docker-compose -f docker-compose.prod.yml exec nginx nginx -s reload
```

## Workflows

### Infrastructure (`infra.yml`)

Manual trigger only. Actions:

- **provision**: Create/verify EC2, security group, Elastic IP
- **destroy**: Terminate all resources tagged `Project=lofam`

### Deploy (`deploy.yml`)

Triggers on:
- Push to `main` branch
- Manual dispatch

Steps:
1. Build and test Go backend
2. SSH to EC2
3. Git pull and docker-compose up

## Manual Deployment

```bash
# SSH into instance
ssh -i ~/.ssh/lofam-key.pem ec2-user@<public-ip>

# Navigate to app
cd ~/app

# Pull and deploy
git pull origin main
docker-compose -f docker-compose.prod.yml up --build -d

# View logs
docker-compose -f docker-compose.prod.yml logs -f
```

## Costs

| Resource | Free Tier | After Free Tier |
|----------|-----------|-----------------|
| EC2 t3.micro | 750 hrs/month for 12 months | ~$8/month |
| Elastic IP | Free when attached | ~$4/month if unattached |
| EBS (20GB gp3) | 30GB free for 12 months | ~$2/month |

**Estimated monthly cost**: $0 (free tier) or ~$10-14/month

SSL certificates from Let's Encrypt are **free**.

## Troubleshooting

### SSL Certificate Issues

```bash
# Check certificate status
docker-compose -f docker-compose.prod.yml exec certbot certbot certificates

# View certbot logs
docker-compose -f docker-compose.prod.yml logs certbot

# Force renewal
docker-compose -f docker-compose.prod.yml exec certbot certbot renew --force-renewal
```

### SSH Connection Refused

```bash
# Check security group allows port 22
aws ec2 describe-security-groups --filters "Name=tag:Project,Values=lofam"
```

### Docker Not Running

```bash
# SSH into instance and check
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

### HTTPS Not Working

```bash
# Check nginx config
docker-compose -f docker-compose.prod.yml exec nginx nginx -t

# Check certificates exist
docker-compose -f docker-compose.prod.yml exec nginx ls -la /etc/letsencrypt/live/lofam/
```

## Destroy Infrastructure

To tear down all resources:

1. Go to Actions → Infrastructure → Run workflow
2. Select action: `destroy`
3. Click "Run workflow"

This terminates:
- EC2 instance
- Releases Elastic IP
- Deletes security group

**Note**: SSL certificates in Docker volumes are lost when instance is terminated.
