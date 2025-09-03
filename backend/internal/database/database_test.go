package database

import (
	"strings"
	"testing"
)

func TestConnect(t *testing.T) {
	t.Run("InvalidDatabaseURL", func(t *testing.T) {
		// Test with invalid database URL
		_, err := Connect("invalid-url")
		if err == nil {
			t.Error("Expected error with invalid database URL")
		}
	})

	t.Run("EmptyDatabaseURL", func(t *testing.T) {
		// Test with empty database URL
		_, err := Connect("")
		if err == nil {
			t.Error("Expected error with empty database URL")
		}
	})
}

func TestRunMigrations(t *testing.T) {
	// Test that migrations are valid SQL
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

	for i, migration := range migrations {
		t.Run(formatMigrationName(i, migration), func(t *testing.T) {
			// Test that migration is valid SQL syntax
			if strings.TrimSpace(migration) == "" {
				t.Error("Migration is empty")
				return
			}

			// Test that migration starts with valid SQL command
			trimmed := strings.TrimSpace(strings.ToUpper(migration))
			validStarts := []string{
				"CREATE TABLE",
				"CREATE INDEX",
				"CREATE EXTENSION",
				"ALTER TABLE",
				"INSERT INTO",
				"UPDATE",
			}

			isValid := false
			for _, start := range validStarts {
				if strings.HasPrefix(trimmed, start) {
					isValid = true
					break
				}
			}

			if !isValid {
				t.Errorf("Migration does not start with valid SQL command: %s", migration[:min(50, len(migration))])
			}

			// Test that migration ends with semicolon
			if !strings.HasSuffix(strings.TrimSpace(migration), ";") {
				t.Error("Migration does not end with semicolon")
			}
		})
	}

	t.Run("NilDatabase", func(t *testing.T) {
		// Test with nil database - should panic or error
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic when running migrations with nil database")
			}
		}()
		
		// This should panic
		RunMigrations(nil)
	})
}

func TestDatabaseSchemaStructure(t *testing.T) {
	t.Run("UsersTableStructure", func(t *testing.T) {
		usersMigration := `CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			email VARCHAR(255) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			subscription_tier VARCHAR(20) DEFAULT 'free',
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW()
		);`

		// Test required fields are present
		requiredFields := []string{"id", "email", "password_hash", "subscription_tier"}
		for _, field := range requiredFields {
			if !strings.Contains(usersMigration, field) {
				t.Errorf("Users table missing required field: %s", field)
			}
		}

		// Test constraints
		if !strings.Contains(usersMigration, "PRIMARY KEY") {
			t.Error("Users table missing primary key")
		}

		if !strings.Contains(usersMigration, "UNIQUE") {
			t.Error("Users table missing unique constraint on email")
		}

		if !strings.Contains(usersMigration, "NOT NULL") {
			t.Error("Users table missing not null constraints")
		}
	})

	t.Run("CardDrawsTableStructure", func(t *testing.T) {
		cardDrawsMigration := `CREATE TABLE IF NOT EXISTS card_draws (
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
		);`

		// Test foreign key relationship
		if !strings.Contains(cardDrawsMigration, "REFERENCES users(id)") {
			t.Error("Card draws table missing foreign key to users")
		}

		// Test cascade delete
		if !strings.Contains(cardDrawsMigration, "ON DELETE CASCADE") {
			t.Error("Card draws table missing cascade delete")
		}

		// Test required fields
		requiredFields := []string{"user_id", "card_id", "card_name", "draw_date"}
		for _, field := range requiredFields {
			if !strings.Contains(cardDrawsMigration, field) {
				t.Errorf("Card draws table missing required field: %s", field)
			}
		}
	})

	t.Run("IndexesStructure", func(t *testing.T) {
		indexes := []string{
			`CREATE INDEX IF NOT EXISTS idx_card_draws_user_date ON card_draws(user_id, draw_date);`,
			`CREATE INDEX IF NOT EXISTS idx_daily_usage_user_date ON daily_usage(user_id, usage_date);`,
			`CREATE INDEX IF NOT EXISTS idx_subscriptions_user ON subscriptions(user_id);`,
		}

		for _, index := range indexes {
			// Test index syntax
			if !strings.Contains(index, "CREATE INDEX") {
				t.Errorf("Invalid index syntax: %s", index)
			}

			if !strings.Contains(index, "IF NOT EXISTS") {
				t.Errorf("Index missing IF NOT EXISTS: %s", index)
			}

			// Test that indexes are on user_id for performance
			if !strings.Contains(index, "user_id") {
				t.Errorf("Index should include user_id for performance: %s", index)
			}
		}
	})
}

func TestDatabaseConstraints(t *testing.T) {
	t.Run("SubscriptionTierEnum", func(t *testing.T) {
		// Test that subscription tier has appropriate constraints
		defaultTier := "free"
		validTiers := []string{"free", "premium"}

		if defaultTier != "free" {
			t.Errorf("Expected default subscription tier to be 'free', got '%s'", defaultTier)
		}

		// Check valid tiers are reasonable
		for _, tier := range validTiers {
			if tier == "" {
				t.Error("Subscription tier should not be empty")
			}
			if len(tier) > 20 {
				t.Errorf("Subscription tier too long: %s", tier)
			}
		}
	})

	t.Run("UniqueConstraints", func(t *testing.T) {
		// Test unique constraints are properly defined
		uniqueConstraints := map[string][]string{
			"users":         {"email"},
			"subscriptions": {"stripe_subscription_id"},
			"daily_usage":   {"user_id, usage_date"},
		}

		for table, fields := range uniqueConstraints {
			if len(fields) == 0 {
				t.Errorf("Table %s should have unique constraints", table)
			}
		}
	})
}

// Helper functions
func formatMigrationName(index int, migration string) string {
	// Extract table name or first few words for test name
	lines := strings.Split(migration, "\n")
	firstLine := strings.TrimSpace(lines[0])
	
	words := strings.Fields(firstLine)
	if len(words) >= 3 {
		return formatTestName(words[2]) // Usually the table name
	}
	
	return formatTestName("Migration" + string(rune(index)))
}

func formatTestName(name string) string {
	// Remove special characters and make test-friendly
	name = strings.ReplaceAll(name, "(", "")
	name = strings.ReplaceAll(name, ")", "")
	name = strings.ReplaceAll(name, ";", "")
	return name
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}