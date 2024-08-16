echo "Starting Before Install Script..."

# Check if Go is installed
if ! command -v go &> /dev/null
then
    echo "Go is not installed. Installing Go..."
    wget https://golang.org/dl/go1.21.0.linux-amd64.tar.gz
    sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
    echo "export PATH=$PATH:/usr/local/go/bin" >> ~/.bashrc
    source ~/.bashrc
    echo "Go installed successfully."
else
    echo "Go is already installed."
fi

# Kill any existing Go processes
PID=$(ps aux | awk '/\/tmp\/go-/' | sed -n 's/  */ /gp' | cut -d ' ' -f 2)
if [ -n "$PID" ]; then
    echo "Killing existing Go process with PID: $PID"
    kill -9 $PID
else
    echo "No existing Go processes found"
fi

# Fetch secrets from AWS Secrets Manager
SUPABASE_API_KEY=$(aws secretsmanager get-secret-value --secret-id SUPABASE_API_KEY --query 'SecretString' --output text)
DB_URL=$(aws secretsmanager get-secret-value --secret-id DB_URL --query 'SecretString' --output text)
FROM_EMAIL=$(aws secretsmanager get-secret-value --secret-id FROM_EMAIL --query 'SecretString' --output text)
FROM_EMAIL_PASSWORD=$(aws secretsmanager get-secret-value --secret-id FROM_EMAIL_PASSWORD --query 'SecretString' --output text)

# Set environment variables
export SUPABASE_API_KEY="${SUPABASE_API_KEY}"
export DB_URL="${DB_URL}"
export FROM_EMAIL="${FROM_EMAIL}"
export FROM_EMAIL_PASSWORD="${FROM_EMAIL_PASSWORD}"

echo "Environment variables set"

# Remove old log file
rm /opt/vokki_cloud/app.log
