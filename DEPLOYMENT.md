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

Docker images are published to Google Artifact Registry.

---

## Option 1: Google Cloud Run (Recommended)

### Prerequisites

1. GCP project with billing enabled
2. `gcloud` CLI installed and authenticated (`gcloud auth login`)

### GCP Setup

1. **Create a GCP project** (or use existing):
   ```bash
   gcloud projects create lofam-12345 --name="Lofam"
   gcloud config set project lofam-12345
   ```
   Note: Project ID must be globally unique. Add a random suffix.

2. **Link billing account** (required for Cloud Run):
   ```bash
   # List available billing accounts
   gcloud billing accounts list

   # Link billing to project
   gcloud billing projects link lofam-12345 \
     --billing-account=YOUR_BILLING_ACCOUNT_ID
   ```

3. **Enable required APIs**:
   ```bash
   gcloud services enable run.googleapis.com
   gcloud services enable artifactregistry.googleapis.com
   ```

4. **Create service account**:
   ```bash
   gcloud iam service-accounts create github-actions \
     --display-name="GitHub Actions"
   ```

5. **Grant permissions**:
   ```bash
   PROJECT_ID=$(gcloud config get-value project)

   # Cloud Run Admin - deploy services
   gcloud projects add-iam-policy-binding $PROJECT_ID \
     --member="serviceAccount:github-actions@$PROJECT_ID.iam.gserviceaccount.com" \
     --role="roles/run.admin"

   # Service Account User - act as service account
   gcloud projects add-iam-policy-binding $PROJECT_ID \
     --member="serviceAccount:github-actions@$PROJECT_ID.iam.gserviceaccount.com" \
     --role="roles/iam.serviceAccountUser"

   # Artifact Registry Admin - push/pull images
   gcloud projects add-iam-policy-binding $PROJECT_ID \
     --member="serviceAccount:github-actions@$PROJECT_ID.iam.gserviceaccount.com" \
     --role="roles/artifactregistry.admin"
   ```

6. **Create service account key**:
   ```bash
   gcloud iam service-accounts keys create key.json \
     --iam-account=github-actions@$PROJECT_ID.iam.gserviceaccount.com

   # View contents to copy to GitHub
   cat key.json

   # Delete after adding to GitHub secrets
   rm key.json
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
3. Push to Artifact Registry
4. Deploy to Cloud Run

### View Deployment

```bash
gcloud run services describe lofam --region us-central1 --format='value(status.url)'
```

### Manual Deployment

```bash
PROJECT_ID=$(gcloud config get-value project)

gcloud run deploy lofam \
  --image us-central1-docker.pkg.dev/$PROJECT_ID/lofam/app:latest \
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
