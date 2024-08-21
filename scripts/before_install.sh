echo "Starting Before Install Script..."

# Install jq
sudo yum install jq -y

# Check if Go is installed
if ! command -v go &> /dev/null
then
    echo "Go is not installed. Installing Go..."
    wget https://golang.org/dl/go1.21.0.linux-amd64.tar.gz
    sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
    
    # Add Go to PATH for the current session
    export PATH=$PATH:/usr/local/go/bin
    echo "export PATH=$PATH:/usr/local/go/bin" >> ~/.bashrc
    
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
SUPABASE_API_KEY_SCRIPT=$(aws secretsmanager get-secret-value --secret-id prod/vokki_cloud --query "SecretString" --output text | jq -r '.SUPABASE_API_KEY')
DB_URL_SCRIPT=$(aws secretsmanager get-secret-value --secret-id prod/vokki_cloud --query "SecretString" --output text | jq -r '.DB_URL')
FROM_EMAIL_SCRIPT=$(aws secretsmanager get-secret-value --secret-id prod/vokki_cloud --query "SecretString" --output text | jq -r '.FROM_EMAIL')
FROM_EMAIL_PASSWORD_SCRIPT=$(aws secretsmanager get-secret-value --secret-id prod/vokki_cloud --query "SecretString" --output text | jq -r '.FROM_EMAIL_PASSWORD')

# Set environment variables
export SUPABASE_API_KEY="${SUPABASE_API_KEY_SCRIPT}"
export DB_URL="${DB_URL_SCRIPT}"
export FROM_EMAIL="${FROM_EMAIL_SCRIPT}"
export FROM_EMAIL_PASSWORD="${FROM_EMAIL_PASSWORD_SCRIPT}"

echo "Environment variables set"

echo "Before Install Script complete."
