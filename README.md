# About

A barebones example of a GoLang web server featuring OAuth 2.0 authentication, session management, and user persistence with PostgreSQL. Features example "Sign in with Google" flow.

# Features

- **Web Server:** Built with standard library Go (`net/http`).
- **OAuth 2.0 Ready:** Simple authentication using the [Goth](https://github.com/markbates/goth) package.
- **Google Sign-In Example:** Example of allowing users to start an auth session with google OAuth.
- **User Persistence:** Stores user information in a PostgreSQL database.
- **Session Management:** Uses [gorilla/sessions](https://github.com/gorilla/sessions) to manage user sessions.
- **Configuration via Environment:** All configuration is managed through a `.env` file.
- **Dockerized Database:** Commands for running PostgreSQL in Docker.

## Prerequisites

Before you begin, ensure you have the following installed:

- [Go](https://go.dev/doc/install)
- [Docker](https://docs.docker.com/get-docker/)
- A Google Cloud Platform account.

## Installation/Running

1. Clone the repo.
2. In the project root, Create a `.env` file and copy the contents of `.env.example` into it.
3. Populate the environment variables with your actual values.
4. Continue to PostgreSQL steps.
5. Finally, in the project root: `go run main.go`

### PostgreSQL

1. Ensure Docker engine is running.
2. Start postgres: `docker run --name your-pg-db -e POSTGRES_PASSWORD=somepassword -p 5431:5432 -d postgres`
3. Stopping: `docker stop your-pg-db`
4. Start over (destroys the db): `docker rm your-pg-db`
5. Using volume: `docker volume create your-pg-db-data`
6. Run with volume: `docker run --name your-pg-db -e POSTGRES_PASSWORD=somepassword -p 5431:5432 -v your-pg-db-data:/var/lib/postgresql/data -d postgres`

In your `.env` file, assign your postgres connection string to the `DATABASE_URL` variable.

Postgres Connection String: `postgres://postgres:somepassword@localhost:5431/postgres`
