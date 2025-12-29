#!/bin/bash

# For Execute
# sed -i 's/\r$//' setup.sh && bash setup.sh

echo "--> Preparing directories..."
mkdir -p logs
touch logs/golang.log logs/golang-error.log 2>/dev/null || true

echo "--> Setting default permissions..."
sudo chown -R www:www . 2>/dev/null || true
sudo find . -type d -exec chmod 755 {} \; 2>/dev/null || true
sudo find . -type f -exec chmod 644 {} \; 2>/dev/null || true

echo "--> Binary permission..."
sudo chown www:www runner-app
sudo chmod 755 runner-app

echo "--> Supervisor setup..."
sudo cp ./deploy/supervisor.conf /etc/supervisor/conf.d/nova-cdn_novadev_myid.conf

sudo supervisorctl reread
sudo supervisorctl update
sudo supervisorctl restart nova-cdn_novadev_myid

echo "--> Securing env files..."
sudo chmod 600 .env .env.local .env.production artisan .well-known .git Makefile setup.sh 2>/dev/null

echo "[SUCCESS] This script has been executed successfully."
