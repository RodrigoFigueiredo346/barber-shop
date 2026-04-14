package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(databaseURL string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}
	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}
	return pool, nil
}

func Migrate(pool *pgxpool.Pool) error {
	query := `
	CREATE TABLE IF NOT EXISTS clients (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		phone VARCHAR(20) UNIQUE NOT NULL,
		created_at TIMESTAMP DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS appointments (
		id SERIAL PRIMARY KEY,
		client_id INTEGER REFERENCES clients(id),
		date DATE NOT NULL,
		time TIME NOT NULL,
		status VARCHAR(20) DEFAULT 'scheduled',
		created_at TIMESTAMP DEFAULT NOW(),
		UNIQUE(date, time, status)
	);

	CREATE TABLE IF NOT EXISTS schedules (
		id SERIAL PRIMARY KEY,
		day_of_week INTEGER NOT NULL CHECK (day_of_week BETWEEN 0 AND 6),
		start_time TIME NOT NULL,
		end_time TIME NOT NULL,
		active BOOLEAN DEFAULT true,
		UNIQUE(day_of_week)
	);

	CREATE TABLE IF NOT EXISTS blocked_slots (
		id SERIAL PRIMARY KEY,
		date DATE NOT NULL,
		time TIME NOT NULL,
		UNIQUE(date, time)
	);

	CREATE TABLE IF NOT EXISTS settings (
		id SERIAL PRIMARY KEY CHECK (id = 1),
		reminder_minutes INTEGER DEFAULT 60
	);

	INSERT INTO settings (id, reminder_minutes) VALUES (1, 60) ON CONFLICT DO NOTHING;

	CREATE TABLE IF NOT EXISTS services (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		duration INTEGER NOT NULL DEFAULT 30,
		price INTEGER NOT NULL DEFAULT 0,
		active BOOLEAN DEFAULT true
	);

	ALTER TABLE appointments ADD COLUMN IF NOT EXISTS service_id INTEGER REFERENCES services(id);
	`
	_, err := pool.Exec(context.Background(), query)
	return err
}
