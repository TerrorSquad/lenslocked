#!/usr/bin/env sh
set -e
STORAGE_PATH=/data/storage
STORAGE_IMAGES_PATH=/data/storage/images

APP_IMAGES_PATH=/app/images

# Create storage directory
if [ ! -d "$STORAGE_PATH" ]; then
    mkdir -p $STORAGE_PATH
fi

if [ ! -d "$APP_IMAGES_PATH" ]; then
    ln -s $STORAGE_IMAGES_PATH $APP_IMAGES_PATH
    echo "Created $APP_IMAGES_PATH and linked to $STORAGE_IMAGES_PATH"
fi
