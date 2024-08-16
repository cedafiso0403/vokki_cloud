echo "Starting Before Install Script..."


PID=$(ps aux | awk '/\/tmp\/go-/' | sed -n 's/  */ /gp' | cut -d ' ' -f 2)

if [ -n "$PID" ]; then
  kill -9 $PID
else
  echo "No existing Go processes found"
fi


SUPABASE_API_KEY=$(aws secretsmanager get-secret-value --secret-id SUPABASE_API_KEY --query 'SecretString' --output text)
DB_URL=$(aws secretsmanager get-secret-value --secret-id  --query DB_URL --output text)
FROM_EMAIL=$(aws secretsmanager get-secret-value --secret-id FROM_EMAIL --query 'SecretString' --output text)
FROM_EMAIL_PASSWORD=$(aws secretsmanager get-secret-value --secret-id FROM_EMAIL_PASSWORD --query 'SecretString' --output text)

export SUPABASE_API_KEY="${SUPABASE_API_KEY}"
export DB_URL="${DB_URL}"
export FROM_EMAIL="${FROM_EMAIL}"
export FROM_EMAIL_PASSWORD="${FROM_EMAIL_PASSWORD}"