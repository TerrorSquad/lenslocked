#!/usr/bin/env sh
FOLDER=/data/storage/images

if [ ! -d "$FOLDER" ]; then
    mkdir -p /data/storage/images
    ln -s /app/images /data/storage/images
fi
