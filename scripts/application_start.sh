echo "Application Start Script..."

cd vokki_cloud

nohup go run ./cmd/main.go > app.log 2>&1 &
