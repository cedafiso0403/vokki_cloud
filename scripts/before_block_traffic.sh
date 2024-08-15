PID=$(ps aux | awk '/\/tmp\/go-/' | sed -n 's/  */ /gp' | cut -d ' ' -f 2)

if [ -n "$PID" ]; then
  kill -9 $PID
else
  echo "No existing Go processes found"
fi