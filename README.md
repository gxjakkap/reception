# Reception (WIP)

A modular, general-purpose Discord bot built with Go. 

## Features
- Reaction Roles
- Welcome Message
- Ticket System

## Prerequisites

- Docker (with docker compose)

## Installation & Setup

### Local Setup
1. Clone the repository:
    ```bash
    git clone https://github.com/gxjakkap/reception.git
    cd reception
    ```

2. Install dependencies:
    ```bash
    go mod download
    ```

3. Configure `.env` (use `.env.example` as a template).
4. Run the bot:
    ```bash
    go run main.go
    ```

### Docker Compose (Recommended)
This method spins up the bot, a PostgreSQL database, and handles migrations automatically.

1. Ensure `.env` is configured.
2. Start the stack:
    ```bash
    docker compose up -d
    ```

**What happens on startup?**
-   **db**: Starts a PostgreSQL 17 instance.
-   **atlas**: Waits for the DB to be healthy, then automatically applies all SQL migrations from the `/migrations` folder.
-   **bot**: Waits for migrations to finish, then starts the bot.

## License
[Apache License 2.0](LICENSE)
