#!/bin/bash

echo "Deploying Vokki Cloud..."

# Exit script on any error
set -e

# Check if the script received 4 arguments
if [ "$#" -ne 4 ]; then
  echo "Error: This script requires 4 arguments."
  echo "Usage: $0 <SUPABASE_API_KEY> <DB_URL> <FROM_EMAIL> <FROM_EMAIL_PASSWORD>"
  exit 1
fi

# Assign arguments to variables
SUPABASE_API_KEY=$1
DB_URL=$2
FROM_EMAIL=$3
FROM_EMAIL_PASSWORD=$4

# Remove existing project directory
rm -rf ~/vokki_cloud

# Kill any existing Go processes
PID=$(ps aux | awk '/\/tmp\/go-/' | sed -n 's/  */ /gp' | cut -d ' ' -f 2)

if [ -n "$PID" ]; then
  kill -9 $PID
else
  echo "No existing Go processes found"
fi

# Clone the repository
cd ~
git clone git@github.com:cedafiso0403/vokki_cloud.git

# Set environment variables
export SUPABASE_API_KEY="${SUPABASE_API_KEY}"
export DB_URL="${DB_URL}"
export FROM_EMAIL="${FROM_EMAIL}"
export FROM_EMAIL_PASSWORD="${FROM_EMAIL_PASSWORD}"

echo "Repository cloned"

# Build and run the Go application
cd vokki_cloud
go build ./...
nohup go run ./cmd/vokki_cloud/main.go > app.log 2>&1 &
