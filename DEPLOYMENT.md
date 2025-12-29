# Deployment Guide

Deploy Lofam to an EC2 instance using GitHub Actions.

## Architecture

```
Go backend:80
├── /api/*  → API handlers
└── /*      → Static frontend (SPA)
```

Single container serves both API and frontend.

## Prerequisites

1. EC2 instance with Docker installed
2. SSH access to the instance
3. GitHub repository with Actions enabled

## EC2 Setup

### Option A: GitHub Actions (Recommended)

1. Go to Actions → Infrastructure → Run workflow
2. Select action: `create`
3. Wait for completion, note the Public IP from summary
4. Add `EC2_HOST` secret with the Public IP

### Option B: Manual Launch

1. Go to AWS Console → EC2 → Launch Instance
2. Use AMI: `ami-0f99d7be7d4273bba`
3. Choose instance type: t2.micro
4. Create or select a key pair
5. Configure security group:
   - SSH (22) from your IP
   - HTTP (80) from anywhere
6. Launch and install Docker:
   ```bash
   sudo yum update -y
   sudo yum install -y docker git
   sudo systemctl enable docker
   sudo systemctl start docker
   sudo usermod -aG docker ec2-user
   ```

### Configure GitHub Secrets

| Secret | Value |
|--------|-------|
| `EC2_HOST` | Your EC2 public IP |
| `EC2_SSH_KEY` | Contents of your `.pem` private key file |

## Deployment

### Automatic (CI/CD)

Push to `main` branch triggers:
1. Go build and tests
2. SSH to EC2
3. Git pull and docker-compose up

### Manual

```bash
ssh -i ~/.ssh/your-key.pem ec2-user@<ec2-host>
cd ~/app
git clone https://github.com/your/repo .  # First time only
git pull origin main
docker-compose -f docker-compose.prod.yml up --build -d
```

## Troubleshooting

### Check container status
```bash
docker-compose -f docker-compose.prod.yml ps
docker-compose -f docker-compose.prod.yml logs
```

### Restart containers
```bash
docker-compose -f docker-compose.prod.yml down
docker-compose -f docker-compose.prod.yml up --build -d
```
