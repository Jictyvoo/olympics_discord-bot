#!/bin/sh

# Default value for DATABASE_PATH if not set
DATABASE_PATH=${DATABASE_PATH:-"/app/data/olympics-2024_PARIS.db"}

# Run the Atlas migration
atlas migrate apply --dir file://build/migrations --url "sqlite://${DATABASE_PATH}"

# Start the bot
exec /bin/olympicsBOT
