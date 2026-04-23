#!/bin/sh
set -e

TARGET=$1

if [ "$TARGET" != "blue" ] && [ "$TARGET" != "green" ]; then
    echo "Usage: $0 blue|green"
    exit 1
fi

perl -i -pe "s/server app_(blue|green):8080;/server app_$TARGET:8080;/" deploy/nginx/nginx.conf
docker compose exec nginx nginx -s reload

echo "Switched to: $TARGET"
