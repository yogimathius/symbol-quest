package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

func Connect(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

func RunMigrations(db *sql.DB) error {
	migrations := []string{
		`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`,
		
		`CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			email VARCHAR(255) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			subscription_tier VARCHAR(20) DEFAULT 'free',
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW()
		);`,
		
		`CREATE TABLE IF NOT EXISTS card_draws (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			card_id INTEGER NOT NULL,
			card_name VARCHAR(100) NOT NULL,
			draw_date DATE NOT NULL,
			interpretation_basic TEXT,
			interpretation_enhanced TEXT,
			mood VARCHAR(50),
			question TEXT,
			created_at TIMESTAMP DEFAULT NOW()
		);`,
		
		`CREATE TABLE IF NOT EXISTS daily_usage (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			usage_date DATE NOT NULL DEFAULT CURRENT_DATE,
			draws_count INTEGER DEFAULT 0,
			UNIQUE(user_id, usage_date)
		);`,
		
		`CREATE TABLE IF NOT EXISTS subscriptions (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			stripe_subscription_id VARCHAR(255) UNIQUE NOT NULL,
			stripe_customer_id VARCHAR(255) NOT NULL,
			status VARCHAR(50) NOT NULL,
			current_period_start TIMESTAMP,
			current_period_end TIMESTAMP,
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW()
		);`,
		
		`CREATE INDEX IF NOT EXISTS idx_card_draws_user_date ON card_draws(user_id, draw_date);`,
		`CREATE INDEX IF NOT EXISTS idx_daily_usage_user_date ON daily_usage(user_id, usage_date);`,
		`CREATE INDEX IF NOT EXISTS idx_subscriptions_user ON subscriptions(user_id);`,
	}

	for _, migration := range migrations {
		if _, err := db.Exec(migration); err != nil {
			return fmt.Errorf("failed to run migration: %w", err)
		}
	}

	return nil
}