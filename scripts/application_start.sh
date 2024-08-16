echo "Application Start Script..."

cd /opt/vokki_cloud


if ! command -v go &> /dev/null
then
    echo "No Go installed."
else
    echo "Go is already installed."
fi

nohup go run ./cmd/main.go


echo "Application Start Script Completed"

