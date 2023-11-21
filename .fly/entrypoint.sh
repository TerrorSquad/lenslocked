#!/usr/bin/env sh

set -e
# Run user scripts
# This script is used to run scripts in the .fly/scripts directory
# It should execute scripts in order, and if a script fails, it should stop

SCRIPTS=$(ls .fly/scripts/*.sh | sort -V)
for script in $SCRIPTS; do
    echo "Running $script"
    sh "$script"
done

if [ $# -gt 0 ]; then
    # If we passed a command, run it as root
    exec "$@"
else
    exec supervisord -c /etc/supervisor/supervisord.conf
fi
