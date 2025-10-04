#!/bin/bash

# A script to build, deploy, and restart the goHome application on a Raspberry Pi.

# --- Configuration ---
# Set the username and hostname for your Raspberry Pi.
REMOTE_USER="admin"
REMOTE_HOST="raspberrypi.local" # Or use the IP address

# Set the path to the project directory on the remote machine.
REMOTE_PATH="/home/admin/go-home/"

# Set the name of your service on the remote machine.
SERVICE_NAME="gohome.service"

# Set the name of the binary to be built.
BINARY_NAME="goHome-rpi"

# --- Script Logic ---

# Exit immediately if any command fails
set -e

echo "--- 1. Building the Go application for Raspberry Pi (linux/arm64)..."
GOOS=linux GOARCH=arm64 go build -o ${BINARY_NAME} -ldflags="-s -w" .
echo "Build complete."

echo ""
echo "--- 2. Stopping the remote service: ${SERVICE_NAME}..."
ssh ${REMOTE_USER}@${REMOTE_HOST} "sudo systemctl stop ${SERVICE_NAME}"
echo "Service stopped."

echo ""
echo "--- 3. Copying the new binary to the Raspberry Pi..."
scp ./${BINARY_NAME} ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_PATH}
echo "Binary copied successfully."

echo ""
echo "--- 4. Starting the remote service and checking its status..."
ssh ${REMOTE_USER}@${REMOTE_HOST} "sudo systemctl start ${SERVICE_NAME} && echo '--- Service Status: ---' && sudo systemctl status ${SERVICE_NAME}"

echo ""
echo "✅ Deployment finished successfully!"