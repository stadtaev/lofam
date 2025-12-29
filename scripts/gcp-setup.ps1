# GCP Setup Script for Lofam
# Run this script to configure GCP for Cloud Run deployment

$ErrorActionPreference = "Stop"

Write-Host "=== Lofam GCP Setup ===" -ForegroundColor Cyan
Write-Host ""

# Check if gcloud is installed
try {
    $null = Get-Command gcloud -ErrorAction Stop
} catch {
    Write-Host "Error: gcloud CLI is not installed" -ForegroundColor Red
    Write-Host "Install from: https://cloud.google.com/sdk/docs/install"
    exit 1
}

# Check if logged in
$activeAccount = gcloud auth list --filter="status:ACTIVE" --format="value(account)" 2>$null
if (-not $activeAccount) {
    Write-Host "Not logged in. Running gcloud auth login..."
    gcloud auth login
}

# Prompt for project ID
$PROJECT_ID = Read-Host "Enter GCP Project ID (must be globally unique, e.g., lofam-12345)"

if (-not $PROJECT_ID) {
    Write-Host "Error: Project ID is required" -ForegroundColor Red
    exit 1
}

Write-Host ""
Write-Host "=== Step 1: Create Project ===" -ForegroundColor Yellow

$projectExists = gcloud projects describe $PROJECT_ID 2>$null
if ($LASTEXITCODE -eq 0) {
    Write-Host "Project $PROJECT_ID already exists"
} else {
    Write-Host "Creating project $PROJECT_ID..."
    gcloud projects create $PROJECT_ID --name="Lofam"
}

gcloud config set project $PROJECT_ID
Write-Host "Active project: $PROJECT_ID" -ForegroundColor Green

Write-Host ""
Write-Host "=== Step 2: Link Billing Account ===" -ForegroundColor Yellow

Write-Host "Available billing accounts:"
gcloud billing accounts list

$BILLING_ACCOUNT = Read-Host "Enter Billing Account ID (e.g., 012345-6789AB-CDEF01)"

if (-not $BILLING_ACCOUNT) {
    Write-Host "Error: Billing account is required for Cloud Run" -ForegroundColor Red
    exit 1
}

gcloud billing projects link $PROJECT_ID --billing-account=$BILLING_ACCOUNT
Write-Host "Billing linked successfully" -ForegroundColor Green

Write-Host ""
Write-Host "=== Step 3: Enable APIs ===" -ForegroundColor Yellow

Write-Host "Enabling Cloud Run API..."
gcloud services enable run.googleapis.com

Write-Host "Enabling Artifact Registry API..."
gcloud services enable artifactregistry.googleapis.com

Write-Host "APIs enabled successfully" -ForegroundColor Green

Write-Host ""
Write-Host "=== Step 4: Create Service Account ===" -ForegroundColor Yellow

$SA_NAME = "github-actions"
$SA_EMAIL = "$SA_NAME@$PROJECT_ID.iam.gserviceaccount.com"

$saExists = gcloud iam service-accounts describe $SA_EMAIL 2>$null
if ($LASTEXITCODE -eq 0) {
    Write-Host "Service account $SA_NAME already exists"
} else {
    Write-Host "Creating service account $SA_NAME..."
    gcloud iam service-accounts create $SA_NAME --display-name="GitHub Actions"
}

Write-Host ""
Write-Host "=== Step 5: Grant Permissions ===" -ForegroundColor Yellow

Write-Host "Granting Cloud Run Admin..."
gcloud projects add-iam-policy-binding $PROJECT_ID `
    --member="serviceAccount:$SA_EMAIL" `
    --role="roles/run.admin" `
    --quiet

Write-Host "Granting Service Account User..."
gcloud projects add-iam-policy-binding $PROJECT_ID `
    --member="serviceAccount:$SA_EMAIL" `
    --role="roles/iam.serviceAccountUser" `
    --quiet

Write-Host "Granting Artifact Registry Admin..."
gcloud projects add-iam-policy-binding $PROJECT_ID `
    --member="serviceAccount:$SA_EMAIL" `
    --role="roles/artifactregistry.admin" `
    --quiet

Write-Host "Granting Storage Admin..."
gcloud projects add-iam-policy-binding $PROJECT_ID `
    --member="serviceAccount:$SA_EMAIL" `
    --role="roles/storage.admin" `
    --quiet

Write-Host "Permissions granted successfully" -ForegroundColor Green

Write-Host ""
Write-Host "=== Step 6: Create Service Account Key ===" -ForegroundColor Yellow

$KEY_FILE = "gcp-sa-key.json"

if (Test-Path $KEY_FILE) {
    $overwrite = Read-Host "$KEY_FILE already exists. Overwrite? (y/N)"
    if ($overwrite -ne "y" -and $overwrite -ne "Y") {
        Write-Host "Skipping key creation"
        $KEY_FILE = $null
    }
}

if ($KEY_FILE) {
    gcloud iam service-accounts keys create $KEY_FILE --iam-account=$SA_EMAIL
    Write-Host ""
    Write-Host "Key saved to $KEY_FILE" -ForegroundColor Green
}

Write-Host ""
Write-Host "=== Setup Complete ===" -ForegroundColor Cyan
Write-Host ""
Write-Host "Next steps:" -ForegroundColor Yellow
Write-Host "1. Go to GitHub repo -> Settings -> Secrets and variables -> Actions"
Write-Host "2. Create secret: GCP_SA_KEY"
Write-Host "3. Paste the contents of $KEY_FILE as the value"
Write-Host "4. Delete $KEY_FILE after adding to GitHub"
Write-Host ""
Write-Host "To view the key contents:"
Write-Host "  Get-Content $KEY_FILE"
Write-Host ""
Write-Host "To delete after adding to GitHub:"
Write-Host "  Remove-Item $KEY_FILE"
