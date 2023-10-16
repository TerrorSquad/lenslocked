#!/usr/bin/env bash

ARCH=$(uname -m)
OS=$(uname -s)

URL="https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-"

if [[ "$OS" == "Linux" ]]; then
    URL+="linux-"
elif [[ "$OS" == "Darwin" ]]; then
    URL+="macos-"
else
    echo "OS not supported"
    exit 1
fi

if [[ "$ARCH" == "x86_64" ]]; then
    URL+="x64"
elif [[ "$ARCH" == "arm64" ]]; then
    URL+="arm64"
else
    echo "OS not supported"
    exit 1
fi

if [[ ! -f "tailwindcss" ]]; then
    curl -L $URL -o tailwindcss
    chmod +x tailwindcss
    exit 0
fi