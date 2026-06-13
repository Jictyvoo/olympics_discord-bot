#!/bin/bash
# remove_backticks.sh
# Remove backticks from SQL migration files

MIGRATIONS_DIRECTORY=$1

if [ -d "$MIGRATIONS_DIRECTORY" ]; then
  for file in "$MIGRATIONS_DIRECTORY"/*.sql; do
    if [ -f "$file" ]; then
      sed -i.bak 's/`//g' "$file" && rm "$file.bak"
      echo "Processed $file"
    fi
  done
else
  echo "Directory $MIGRATIONS_DIRECTORY does not exist."
fi
