#!/bin/bash

# Exit script on any error
set -e

# Remove existing project directory
rm -rf ~/vokki_cloud

# Kill any existing Go processes
ps aux | awk '/\/tmp\/go-/' | sed -n 's/  */ /gp' | cut -d ' ' -f 2 | xargs kill -9

# Clone the repository
cd ~
git clone git@github.com:cedafiso0403/vokki_cloud.git

# Set environment variables
export SUPABASE_API_KEY="${SUPABASE_API_KEY}"
export DB_URL="${DB_URL}"
export FROM_EMAIL="${FROM_EMAIL}"
export FROM_EMAIL_PASSWORD="${FROM_EMAIL_PASSWORD}"

# Build and run the Go application
cd vokki_cloud
go build ./...
go run ./cmd/vokki_cloud/main.go &
