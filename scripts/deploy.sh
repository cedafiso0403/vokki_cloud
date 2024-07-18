#!/bin/bash

echo "Deploying Vokki Cloud..."

echo "From email: $FROM_EMAIL"

# Exit script on any error
set -e

# Remove existing project directory
rm -rf ~/vokki_cloud

# Kill any existing Go processes
PID = ps aux | awk '/\/tmp\/go-/' | sed -n 's/  */ /gp' | cut -d ' ' -f 2

if [ -n "$PID" ]; then
  kill -9 $PID
fi

# Clone the repository
cd ~
git clone git@github.com:cedafiso0403/vokki_cloud.git

# Build and run the Go application
cd vokki_cloud
go build ./...
nohup go run ./cmd/vokki_cloud/main.go > app.log 2>&1 &
