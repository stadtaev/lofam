# Deployment Guide

Deploy Lofam to an EC2 instance using GitHub Actions.

## Prerequisites

1. EC2 instance running Amazon Linux 2023
2. Docker and docker-compose installed on EC2
3. SSH access to the instance
4. GitHub repository with Actions enabled

## EC2 Setup

### 1. Launch EC2 Instance

1. Go to AWS Console → EC2 → Launch Instance
2. Select Amazon Linux 2023 AMI
3. Choose instance type (t3.micro for free tier)
4. Create or select a key pair
5. Configure security group:
   - SSH (22) from your IP
   - HTTP (80) from anywhere
   - HTTPS (443) from anywhere (if using SSL)
6. Launch instance

### 2. Install Docker

SSH into your instance and run:

```bash
sudo dnf update -y
sudo dnf install -y docker git
sudo systemctl enable docker
sudo systemctl start docker
sudo usermod -aG docker ec2-user

# Install docker-compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# Create app directory
mkdir -p ~/app

# Log out and back in for docker group to take effect
exit
```

### 3. Configure GitHub Secrets

Go to your repo → Settings → Secrets and variables → Actions → New repository secret

| Secret | Value |
|--------|-------|
| `EC2_HOST` | Your EC2 public IP or domain |
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

## Production Stack

```
nginx:80 → frontend:3000 (Next.js)
         → backend:8080  (Go API via /api/*)
```

## Adding SSL (Let's Encrypt)

1. Point your domain A record to EC2 public IP
2. SSH into the instance and run:
   ```bash
   cd ~/app
   ./init-ssl.sh yourdomain.com your@email.com
   ```
3. Switch to SSL compose file:
   ```bash
   docker-compose -f docker-compose.letsencrypt.yml up -d
   ```

Certificates auto-renew via certbot container.

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

### Check Docker service
```bash
sudo systemctl status docker
```
