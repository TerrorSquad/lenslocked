FOLDER=/app
# This script should open the /app directory
# Then it should check if the /app/storage directory exists
# If it does not exist, it should create it by making it a symlink to /data/storage
# Also create the /app/images directory if it does not exist
# It should be a link to /data/images

if [ ! -d "$FOLDER" ]; then
  echo "Hey";
fi
