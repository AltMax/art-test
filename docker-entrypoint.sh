#!/bin/sh
set -e

echo "migrate_common"
./migrate_common

echo "unit_service"
exec ./unit_service