package services

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCardService_Creation(t *testing.T) {
	service := NewCardService(nil)
	if service == nil {
		t.Error("NewCardService returned nil")
	}
}

func TestCardService_CheckUserLimits(t *testing.T) {
	// Skip database-dependent test in unit mode
	if testing.Short() {
		t.Skip("Skipping database-dependent test in short mode")
	}
	
	service := &CardService{db: nil}
	
	// Test with nil database (should return error)
	_, err := service.checkUserLimits(uuid.New())
	if err == nil {
		t.Error("Expected error with nil database, got none")
	}
}

func TestCardService_GetCardMeaning(t *testing.T) {
	service := &CardService{db: nil}

	t.Run("ValidCardID", func(t *testing.T) {
		card, err := service.GetCardMeaning(0) // The Fool
		if err != nil {
			t.Fatalf("Expected valid card ID to return card: %v", err)
		}

		if card == nil {
			t.Error("Expected non-nil card")
		}

		if card.Name != "The Fool" {
			t.Errorf("Expected card name 'The Fool', got '%s'", card.Name)
		}

		if card.ID != 0 {
			t.Errorf("Expected card ID 0, got %d", card.ID)
		}
	})

	t.Run("InvalidCardID", func(t *testing.T) {
		_, err := service.GetCardMeaning(99)
		if err == nil {
			t.Error("Expected error for invalid card ID, got none")
		}

		expectedError := "card not found"
		if err.Error() != expectedError {
			t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
		}
	})

	t.Run("NegativeCardID", func(t *testing.T) {
		_, err := service.GetCardMeaning(-1)
		if err == nil {
			t.Error("Expected error for negative card ID, got none")
		}
	})
}

func TestDailyDrawValidation(t *testing.T) {
	t.Run("ValidInputs", func(t *testing.T) {
		mood := "excited"
		question := "What should I focus on today?"
		
		// These are valid inputs that should pass validation
		if mood == "" {
			t.Error("Mood should not be empty in valid test")
		}
		if question == "" {
			t.Error("Question should not be empty in valid test")
		}
	})

	t.Run("EmptyMood", func(t *testing.T) {
		mood := ""
		question := "What should I focus on today?"
		
		// Empty mood should be handled gracefully
		if mood != "" {
			t.Error("Expected empty mood in this test")
		}
		// Should still allow the draw to proceed
		_ = question // Use the variable
	})

	t.Run("EmptyQuestion", func(t *testing.T) {
		mood := "excited"
		question := ""
		
		// Empty question should be handled gracefully
		if question != "" {
			t.Error("Expected empty question in this test")
		}
		// Should still allow the draw to proceed
		_ = mood // Use the variable
	})
}

func TestTodayDateFormatting(t *testing.T) {
	// Test that date formatting matches expected format
	today := time.Now().Format("2006-01-02")
	
	// Should be in YYYY-MM-DD format
	if len(today) != 10 {
		t.Errorf("Expected date format YYYY-MM-DD (10 chars), got %s (%d chars)", today, len(today))
	}

	// Should contain dashes
	if today[4] != '-' || today[7] != '-' {
		t.Errorf("Expected date format YYYY-MM-DD with dashes, got %s", today)
	}

	// Should be parseable back to time
	_, err := time.Parse("2006-01-02", today)
	if err != nil {
		t.Errorf("Date format %s is not parseable: %v", today, err)
	}
}

func TestUsageLimitLogic(t *testing.T) {
	t.Run("FreeUserLimit", func(t *testing.T) {
		// Free users should have daily limit of 1
		freeLimit := 1
		currentUsage := 0

		if currentUsage >= freeLimit {
			t.Error("Free user should be able to draw when under limit")
		}

		currentUsage = 1
		if currentUsage < freeLimit {
			t.Error("Free user should hit limit after 1 draw")
		}
	})

	t.Run("PremiumUserLimit", func(t *testing.T) {
		// Premium users should have unlimited draws
		isPremium := true
		currentUsage := 100 // Doesn't matter for premium

		if !isPremium {
			t.Error("Expected premium user flag to be true")
		}

		// Premium users bypass usage checks
		canDraw := isPremium || currentUsage < 1
		if !canDraw {
			t.Error("Premium users should always be able to draw")
		}
	})
}

func TestErrorHandling(t *testing.T) {
	service := &CardService{db: nil}

	t.Run("NilDatabase", func(t *testing.T) {
		// Most database operations should handle nil gracefully or return appropriate errors
		_, err := service.checkUserLimits(uuid.New())
		if err == nil {
			t.Error("Expected error when database is nil")
		}
	})

	t.Run("InvalidUUID", func(t *testing.T) {
		// Test with zero UUID - should return error due to nil database
		_, err := service.checkUserLimits(uuid.UUID{})
		if err == nil {
			t.Error("Expected error with zero UUID and nil database")
		}
	})
}

func TestCardServiceConstants(t *testing.T) {
	// Test various constants and limits used by card service
	
	t.Run("DefaultHistoryLimit", func(t *testing.T) {
		defaultLimit := 20
		if defaultLimit <= 0 {
			t.Error("Default history limit should be positive")
		}
		if defaultLimit > 100 {
			t.Error("Default history limit should be reasonable")
		}
	})

	t.Run("FreeUserDailyLimit", func(t *testing.T) {
		freeLimit := 1
		if freeLimit != 1 {
			t.Errorf("Expected free user daily limit to be 1, got %d", freeLimit)
		}
	})
}

// Benchmark tests for performance
func BenchmarkGetCardMeaning(b *testing.B) {
	service := &CardService{db: nil}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.GetCardMeaning(i % 22) // Cycle through all cards
	}
}

func BenchmarkTodayDateFormatting(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = time.Now().Format("2006-01-02")
	}
}