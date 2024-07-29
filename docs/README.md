# Olympics PARIS2024 - Event Bot

![Olympics Data Fetcher Logo](olympics_data_fecher-logo.png)

## Overview

The **Olympics PARIS2024 - Event Bot** is designed to provide up-to-date information on the Olympic events happening in Paris 2024. This bot fetches data every hour, caches it, and stores it in a SQLite database. It then uses this data to manage notifications for subscribed Discord channels, keeping users informed about current and upcoming events.

## Features

- **Data Fetching**: Downloads event information every hour.
- **Caching**: Efficiently caches fetched data to minimize server load.
- **Database Management**: Stores event data in a SQLite database.
- **Discord Integration**: Manages notifications for subscribed Discord channels.
- **Written in Go**: Utilizes the Go programming language for development.

## Installation

### Prerequisites

- Go 1.22+
- Docker

## Usage

1. **Start the bot:**
   The bot will automatically start fetching data every hour and managing notifications.

2. **Subscribe to notifications:**
   Use the bot commands in your Discord server to subscribe to event notifications.

## Configuration

### Environment Variables

- `DISCORD_TOKEN`: Your Discord bot token.
- `DISCORD_CLIENT_ID`: The client ID for discord bot, used to generate an invitation link
- `DATABASE_URL`: SQLite database URL (default: `sqlite3://data/events.db`).

## Dependencies

- Go
  - `libcurl`: For HTTP requests.
  - `sqlc`: For SQL query generation.
  - `atlasgo`: For migration diff + apply
- SQLite
- Docker
- Discord API

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

## License

This project is licensed under the MIT License.
