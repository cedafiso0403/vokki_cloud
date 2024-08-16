echo "Application Start Script..."

cd /opt/vokki_cloud

nohup ./go-vokki-cloud > app.log 2>&1 &

