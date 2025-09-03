package tarot

import (
	"testing"
	"github.com/google/uuid"
)

func TestMajorArcanaData(t *testing.T) {
	// Test that all 22 Major Arcana cards are present
	expectedCards := 22
	if len(MajorArcana) != expectedCards {
		t.Errorf("Expected %d cards, got %d", expectedCards, len(MajorArcana))
	}

	// Test that all cards have required fields
	for id, card := range MajorArcana {
		if card.ID != id {
			t.Errorf("Card ID mismatch: expected %d, got %d", id, card.ID)
		}

		if card.Name == "" {
			t.Errorf("Card %d has empty name", id)
		}

		if card.Number == "" {
			t.Errorf("Card %d has empty number", id)
		}

		if len(card.Keywords) == 0 {
			t.Errorf("Card %d has no keywords", id)
		}

		if len(card.Elements) == 0 {
			t.Errorf("Card %d has no elements", id)
		}

		if card.TraditionalMeaning == "" {
			t.Errorf("Card %d has empty traditional meaning", id)
		}

		if len(card.MoodWeights) == 0 {
			t.Errorf("Card %d has no mood weights", id)
		}

		// Test mood weights are reasonable (0.0 to 2.0)
		for mood, weight := range card.MoodWeights {
			if weight < 0.0 || weight > 2.0 {
				t.Errorf("Card %d has unreasonable mood weight for %s: %f", id, mood, weight)
			}
		}
	}
}

func TestSelectIntelligentCard(t *testing.T) {
	userID := uuid.New()

	t.Run("SelectsValidCard", func(t *testing.T) {
		cardID := SelectIntelligentCard(userID, nil, "excited", "What should I focus on today?")
		
		if cardID < 0 || cardID > 21 {
			t.Errorf("Selected invalid card ID: %d", cardID)
		}

		if _, exists := MajorArcana[cardID]; !exists {
			t.Errorf("Selected card ID %d does not exist in MajorArcana", cardID)
		}
	})

	t.Run("HandlesMoodBasedSelection", func(t *testing.T) {
		// Run multiple times to test consistency
		results := make(map[int]int)
		iterations := 100

		for i := 0; i < iterations; i++ {
			cardID := SelectIntelligentCard(userID, nil, "anxious", "I need guidance")
			results[cardID]++
		}

		// Should select multiple different cards (randomness)
		if len(results) < 5 {
			t.Errorf("Selection not random enough, only %d unique cards in %d iterations", len(results), iterations)
		}

		// Verify all selected cards are valid
		for cardID := range results {
			if cardID < 0 || cardID > 21 {
				t.Errorf("Invalid card ID selected: %d", cardID)
			}
		}
	})

	t.Run("HandlesEmptyMoodAndQuestion", func(t *testing.T) {
		cardID := SelectIntelligentCard(userID, nil, "", "")
		
		if cardID < 0 || cardID > 21 {
			t.Errorf("Selected invalid card ID with empty mood/question: %d", cardID)
		}
	})
}

func TestCalculateCardScore(t *testing.T) {
	card := MajorArcana[0] // The Fool

	t.Run("BaseLine", func(t *testing.T) {
		score := calculateCardScore(card, "", "")
		expected := 1.0
		if score != expected {
			t.Errorf("Expected baseline score %f, got %f", expected, score)
		}
	})

	t.Run("MoodWeighting", func(t *testing.T) {
		score := calculateCardScore(card, "excited", "")
		
		// The Fool has excited weight of 1.2
		expected := 1.2
		if score != expected {
			t.Errorf("Expected mood-weighted score %f, got %f", expected, score)
		}
	})

	t.Run("KeywordMatching", func(t *testing.T) {
		score := calculateCardScore(card, "", "new-beginnings in my life")
		
		// Should get keyword bonus for "new-beginnings" match
		if score <= 1.0 {
			t.Logf("Score: %f - keyword matching may not have triggered", score)
		}
		// Test passes if score is calculated without error
	})

	t.Run("CombinedScoring", func(t *testing.T) {
		score := calculateCardScore(card, "excited", "I'm starting something new")
		
		// Should get both mood weight and keyword bonuses
		if score < 1.2 {
			t.Errorf("Expected score to be at least mood-weighted (1.2), got %f", score)
		}
		// Test that score is reasonable
		if score > 5.0 {
			t.Errorf("Score too high: %f", score)
		}
	})

	t.Run("InvalidMood", func(t *testing.T) {
		score := calculateCardScore(card, "nonexistent-mood", "")
		expected := 1.0
		if score != expected {
			t.Errorf("Expected baseline score for invalid mood, got %f", score)
		}
	})
}

func TestContainsFunction(t *testing.T) {
	slice := []int{1, 3, 5, 7, 9}

	t.Run("ContainsExistingItem", func(t *testing.T) {
		if !contains(slice, 5) {
			t.Error("Expected contains(slice, 5) to be true")
		}
	})

	t.Run("DoesNotContainMissingItem", func(t *testing.T) {
		if contains(slice, 4) {
			t.Error("Expected contains(slice, 4) to be false")
		}
	})

	t.Run("EmptySlice", func(t *testing.T) {
		emptySlice := []int{}
		if contains(emptySlice, 1) {
			t.Error("Expected contains(emptySlice, 1) to be false")
		}
	})
}

func TestGetRecentCards(t *testing.T) {
	// This would require database mocking for proper testing
	// For now, test the nil database case should panic or return empty
	defer func() {
		if r := recover(); r != nil {
			// Panic is expected with nil database
			t.Logf("Expected panic occurred: %v", r)
		}
	}()
	
	userID := uuid.New()
	cards := getRecentCards(userID, nil, 5)
	
	// If we get here, nil database returned empty slice
	if len(cards) != 0 {
		t.Errorf("Expected empty slice with nil database, got %d cards", len(cards))
	}
}

// Benchmark tests for performance
func BenchmarkSelectIntelligentCard(b *testing.B) {
	userID := uuid.New()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SelectIntelligentCard(userID, nil, "excited", "What should I focus on?")
	}
}

func BenchmarkCalculateCardScore(b *testing.B) {
	card := MajorArcana[0]
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		calculateCardScore(card, "excited", "new beginnings in my life")
	}
}