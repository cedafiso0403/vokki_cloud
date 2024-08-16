echo "Application Start Script..."

cd /opt/vokki_cloud


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


nohup go run ./cmd/main.go


echo "Application Start Script Completed"

