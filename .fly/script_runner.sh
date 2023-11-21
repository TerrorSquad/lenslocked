#!/bin/bash

set -e
# This script is used to run scripts in the .fly/scripts directory
# It should execute scripts in order, and if a script fails, it should stop

SCRIPTS=$(ls .fly/scripts/*.sh | sort -V)
for script in $SCRIPTS; do
    echo "Running $script"
    bash "$script"
done
