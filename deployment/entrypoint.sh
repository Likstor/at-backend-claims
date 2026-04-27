#!/bin/sh

./at-migrator
if [ $? -ne 0 ]; then
    echo "migrator error"
    exit 1
fi

./at-backend