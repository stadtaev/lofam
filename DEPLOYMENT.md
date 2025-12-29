# Deployment Guide

Deploy Lofam using GitHub Actions. Two deployment options available:
- **Google Cloud Run** (recommended) - serverless, scales to zero
- **AWS EC2** - traditional VM-based deployment

## Architecture

```
Go backend:80
├── /api/*  → API handlers
└── /*      → Static frontend (SPA)
```

Single container serves both API and frontend.

Docker images are published to GitHub Container Registry (ghcr.io).

---

## Option 1: Google Cloud Run (Recommended)

### Prerequisites

1. GCP project with billing enabled
2. Cloud Run API enabled
3. Service account with Cloud Run Admin role

### GCP Setup

1. **Create a GCP project** (or use existing):
   ```bash
   gcloud projects create lofam-project --name="Lofam"
   gcloud config set project lofam-project
   ```

2. **Enable required APIs**:
   ```bash
   gcloud services enable run.googleapis.com
   ```

3. **Create service account**:
   ```bash
   gcloud iam service-accounts create github-actions \
     --display-name="GitHub Actions"
   ```

4. **Grant permissions**:
   ```bash
   PROJECT_ID=$(gcloud config get-value project)

   gcloud projects add-iam-policy-binding $PROJECT_ID \
     --member="serviceAccount:github-actions@$PROJECT_ID.iam.gserviceaccount.com" \
     --role="roles/run.admin"

   gcloud projects add-iam-policy-binding $PROJECT_ID \
     --member="serviceAccount:github-actions@$PROJECT_ID.iam.gserviceaccount.com" \
     --role="roles/iam.serviceAccountUser"
   ```

5. **Create and download key**:
   ```bash
   gcloud iam service-accounts keys create key.json \
     --iam-account=github-actions@$PROJECT_ID.iam.gserviceaccount.com
   ```

### Configure GitHub Secrets

| Secret | Value |
|--------|-------|
| `GCP_SA_KEY` | Contents of `key.json` (the entire JSON) |

### Workflow Configuration

Edit `.github/workflows/cloud-run.yml` if needed:

```yaml
env:
  GCP_REGION: us-central1        # Change region if needed
  CLOUD_RUN_SERVICE: lofam       # Service name in Cloud Run
```

### Deployment

Push to `main` branch triggers:
1. Go build and tests
2. Build Docker image
3. Push to ghcr.io
4. Deploy to Cloud Run

### View Deployment

```bash
gcloud run services describe lofam --region us-central1 --format='value(status.url)'
```

### Manual Deployment

```bash
gcloud run deploy lofam \
  --image ghcr.io/your-username/lofam:latest \
  --region us-central1 \
  --platform managed \
  --allow-unauthenticated
```

---

## Option 2: AWS EC2

### Prerequisites

1. EC2 instance with Docker installed
2. SSH access to the instance
3. GitHub repository with Actions enabled

### EC2 Setup

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

### Deployment

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

### Troubleshooting

**Check container status:**
```bash
docker-compose -f docker-compose.prod.yml ps
docker-compose -f docker-compose.prod.yml logs
```

**Restart containers:**
```bash
docker-compose -f docker-compose.prod.yml down
docker-compose -f docker-compose.prod.yml up --build -d
```
