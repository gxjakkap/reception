#!/bin/bash


MIGRATION_NAME=""

echo "Generating migration: $MIGRATION_NAME..."

atlas migrate diff "$MIGRATION_NAME" --env gorm

if [ $? -eq 0 ]; then
  echo "Successfully generated migration files in ./migrations"
else
  echo "Error: Failed to generate migration files."
  exit 1
fi
