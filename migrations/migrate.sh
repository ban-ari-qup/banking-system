#!/bin/bash

# Скрипт для применения миграций
set -e

DB_NAME="mybank"
DB_USER="postgres"
DB_HOST="localhost"
DB_PORT="5432"

echo "Applying database migrations..."

for migration_file in migrations/*.sql; do
    echo "Applying: $migration_file"
    psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f "$migration_file"
done

echo "All migrations applied successfully!"