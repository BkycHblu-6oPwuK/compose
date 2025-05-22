#!/bin/bash

CREATE_SIMLINK="$(dirname "$0")/create_simlink.sh"

if [ -f "$CREATE_SIMLINK" ]; then
    bash "$CREATE_SIMLINK"
fi

exec supervisord -c /etc/supervisor/conf.d/supervisord.conf
