#!/bin/bash
# Initialize Let's Encrypt SSL certificate
# Usage: ./init-ssl.sh <domain> <email>

set -e

DOMAIN="${1:?Usage: $0 <domain> <email>}"
EMAIL="${2:?Usage: $0 <domain> <email>}"

echo "=== Obtaining SSL certificate for $DOMAIN ==="

# Create directories
mkdir -p certbot-www

# Start temporary nginx for ACME challenge
docker run -d --name nginx-temp \
  -p 80:80 \
  -v "$(pwd)/certbot-www:/var/www/certbot:ro" \
  nginx:alpine \
  sh -c "echo 'server { listen 80; location /.well-known/acme-challenge/ { root /var/www/certbot; } }' > /etc/nginx/conf.d/default.conf && nginx -g 'daemon off;'"

echo "Waiting for nginx..."
sleep 3

# Request certificate
docker run --rm \
  -v "$(pwd)/certbot-www:/var/www/certbot" \
  -v "lofam_certbot-conf:/etc/letsencrypt" \
  certbot/certbot certonly \
  --webroot \
  --webroot-path=/var/www/certbot \
  --email "$EMAIL" \
  --agree-tos \
  --no-eff-email \
  --cert-name lofam \
  -d "$DOMAIN"

# Stop temporary nginx
docker stop nginx-temp
docker rm nginx-temp

echo ""
echo "=== Certificate obtained ==="
echo "Now run: docker-compose -f docker-compose.letsencrypt.yml up -d"
