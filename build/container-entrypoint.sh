#!/bin/sh

# Run the Atlas migration
atlas migrate apply --dir file://build/migrations --url "sqlite://olympics-2024_PARIS.db"

# Start the bot
exec /bin/olympicsBOT
