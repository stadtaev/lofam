#!/bin/bash
set -e

# GCP Setup Script for Lofam
# Run this script to configure GCP for Cloud Run deployment

echo "=== Lofam GCP Setup ==="
echo ""

# Check if gcloud is installed
if ! command -v gcloud &> /dev/null; then
    echo "Error: gcloud CLI is not installed"
    echo "Install from: https://cloud.google.com/sdk/docs/install"
    exit 1
fi

# Check if logged in
if ! gcloud auth list --filter=status:ACTIVE --format="value(account)" | grep -q .; then
    echo "Not logged in. Running gcloud auth login..."
    gcloud auth login
fi

# Prompt for project ID
read -p "Enter GCP Project ID (must be globally unique, e.g., lofam-12345): " PROJECT_ID

if [ -z "$PROJECT_ID" ]; then
    echo "Error: Project ID is required"
    exit 1
fi

echo ""
echo "=== Step 1: Create Project ==="
if gcloud projects describe "$PROJECT_ID" &> /dev/null; then
    echo "Project $PROJECT_ID already exists"
else
    echo "Creating project $PROJECT_ID..."
    gcloud projects create "$PROJECT_ID" --name="Lofam"
fi

gcloud config set project "$PROJECT_ID"
echo "Active project: $PROJECT_ID"

echo ""
echo "=== Step 2: Link Billing Account ==="
BILLING_ACCOUNTS=$(gcloud billing accounts list --format="value(name)")

if [ -z "$BILLING_ACCOUNTS" ]; then
    echo "Error: No billing accounts found. Create one at https://console.cloud.google.com/billing"
    exit 1
fi

echo "Available billing accounts:"
gcloud billing accounts list

read -p "Enter Billing Account ID (e.g., 012345-6789AB-CDEF01): " BILLING_ACCOUNT

if [ -z "$BILLING_ACCOUNT" ]; then
    echo "Error: Billing account is required for Cloud Run"
    exit 1
fi

gcloud billing projects link "$PROJECT_ID" --billing-account="$BILLING_ACCOUNT"
echo "Billing linked successfully"

echo ""
echo "=== Step 3: Enable APIs ==="
echo "Enabling Cloud Run API..."
gcloud services enable run.googleapis.com

echo "Enabling Artifact Registry API..."
gcloud services enable artifactregistry.googleapis.com

echo "APIs enabled successfully"

echo ""
echo "=== Step 4: Create Service Account ==="
SA_NAME="github-actions"
SA_EMAIL="$SA_NAME@$PROJECT_ID.iam.gserviceaccount.com"

if gcloud iam service-accounts describe "$SA_EMAIL" &> /dev/null; then
    echo "Service account $SA_NAME already exists"
else
    echo "Creating service account $SA_NAME..."
    gcloud iam service-accounts create "$SA_NAME" \
        --display-name="GitHub Actions"
fi

echo ""
echo "=== Step 5: Grant Permissions ==="
echo "Granting Cloud Run Admin..."
gcloud projects add-iam-policy-binding "$PROJECT_ID" \
    --member="serviceAccount:$SA_EMAIL" \
    --role="roles/run.admin" \
    --quiet

echo "Granting Service Account User..."
gcloud projects add-iam-policy-binding "$PROJECT_ID" \
    --member="serviceAccount:$SA_EMAIL" \
    --role="roles/iam.serviceAccountUser" \
    --quiet

echo "Granting Artifact Registry Admin..."
gcloud projects add-iam-policy-binding "$PROJECT_ID" \
    --member="serviceAccount:$SA_EMAIL" \
    --role="roles/artifactregistry.admin" \
    --quiet

echo "Granting Storage Admin..."
gcloud projects add-iam-policy-binding "$PROJECT_ID" \
    --member="serviceAccount:$SA_EMAIL" \
    --role="roles/storage.admin" \
    --quiet

echo "Permissions granted successfully"

echo ""
echo "=== Step 6: Create Service Account Key ==="
KEY_FILE="gcp-sa-key.json"

if [ -f "$KEY_FILE" ]; then
    read -p "$KEY_FILE already exists. Overwrite? (y/N): " OVERWRITE
    if [ "$OVERWRITE" != "y" ] && [ "$OVERWRITE" != "Y" ]; then
        echo "Skipping key creation"
        KEY_FILE=""
    fi
fi

if [ -n "$KEY_FILE" ]; then
    gcloud iam service-accounts keys create "$KEY_FILE" \
        --iam-account="$SA_EMAIL"
    echo ""
    echo "Key saved to $KEY_FILE"
fi

echo ""
echo "=== Setup Complete ==="
echo ""
echo "Next steps:"
echo "1. Go to GitHub repo → Settings → Secrets and variables → Actions"
echo "2. Create secret: GCP_SA_KEY"
echo "3. Paste the contents of $KEY_FILE as the value"
echo "4. Delete $KEY_FILE after adding to GitHub"
echo ""
echo "To view the key contents:"
echo "  cat $KEY_FILE"
echo ""
echo "To delete after adding to GitHub:"
echo "  rm $KEY_FILE"
