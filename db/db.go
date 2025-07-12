package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DBPool *pgxpool.Pool

func InitDB() {
	var err error
	DBPool, err = pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	defer DBPool.Close()

	// Create the shared updated_at trigger function
	_, err = DBPool.Exec(context.Background(), `
		CREATE OR REPLACE FUNCTION set_updated_at()
		RETURNS TRIGGER AS $$
		BEGIN
			NEW.updated_at = CURRENT_TIMESTAMP;
			RETURN NEW;
		END;
		$$ LANGUAGE plpgsql;
	`)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create set_updated_at function: %v\n", err)
		os.Exit(1)
	}

	// Create tables (including updated_at column)
	tableQueries := map[string]string{
		"users": `
			CREATE TABLE IF NOT EXISTS users (
				id SERIAL PRIMARY KEY,
				auth_provider_title TEXT NOT NULL,
				auth_provider_user_id TEXT NOT NULL,
				auth_provider_email TEXT NOT NULL,
				created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
				UNIQUE(auth_provider_title, auth_provider_email)
			);`,
	}

	for name, query := range tableQueries {
		if _, err := DBPool.Exec(context.Background(), query); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating table %s: %v\n", name, err)
			os.Exit(1)
		}
	}

	// Add updated_at triggers
	for name := range tableQueries {
		triggerQuery := fmt.Sprintf(`
			DO $$
			BEGIN
				IF NOT EXISTS (
					SELECT 1 FROM pg_trigger WHERE tgname = 'set_updated_at_%s'
				) THEN
					CREATE TRIGGER set_updated_at_%s
					BEFORE UPDATE ON %s
					FOR EACH ROW
					EXECUTE FUNCTION set_updated_at();
				END IF;
			END
			$$;
		`, name, name, name)

		if _, err := DBPool.Exec(context.Background(), triggerQuery); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating trigger for table %s: %v\n", name, err)
			os.Exit(1)
		}
	}
}
