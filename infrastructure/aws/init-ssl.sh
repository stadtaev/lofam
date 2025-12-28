#!/bin/bash
set -e

DOMAIN="${1:?Usage: $0 <domain> <email>}"
EMAIL="${2:?Usage: $0 <domain> <email>}"

echo "=== Initializing SSL for $DOMAIN ==="

# Create temporary nginx config for initial cert
cat > /tmp/nginx-init.conf << 'EOF'
events { worker_connections 1024; }
http {
    server {
        listen 80;
        server_name _;
        location /.well-known/acme-challenge/ {
            root /var/www/certbot;
        }
        location / {
            return 200 'OK';
            add_header Content-Type text/plain;
        }
    }
}
EOF

# Stop existing containers
docker-compose -f docker-compose.prod.yml down 2>/dev/null || true

# Start temporary nginx for certificate challenge
echo "Starting temporary nginx for ACME challenge..."
docker run -d --name nginx-init \
    -p 80:80 \
    -v /tmp/nginx-init.conf:/etc/nginx/nginx.conf:ro \
    -v lofam_certbot-www:/var/www/certbot \
    nginx:alpine

# Request certificate
echo "Requesting certificate from Let's Encrypt..."
docker run --rm \
    -v lofam_certbot-conf:/etc/letsencrypt \
    -v lofam_certbot-www:/var/www/certbot \
    certbot/certbot certonly \
    --webroot \
    --webroot-path=/var/www/certbot \
    --email "$EMAIL" \
    --agree-tos \
    --no-eff-email \
    -d "$DOMAIN" \
    --cert-name lofam

# Stop temporary nginx
echo "Stopping temporary nginx..."
docker stop nginx-init
docker rm nginx-init

echo ""
echo "=== SSL Certificate obtained ==="
echo "Now start the full stack:"
echo "  docker-compose -f docker-compose.prod.yml up -d"
