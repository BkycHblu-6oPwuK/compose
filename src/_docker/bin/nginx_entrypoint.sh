#!/bin/bash

CREATE_SIMLINK="$(dirname "$0")/create_simlink.sh"

if [ -f "$CREATE_SIMLINK" ]; then
    bash "$CREATE_SIMLINK"
fi

echo "Starting Nginx..."
nginx -g "daemon off;"