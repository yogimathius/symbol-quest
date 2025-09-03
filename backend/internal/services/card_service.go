package services

import (
	"database/sql"
	"errors"
	"symbol-quest/internal/models"
	"symbol-quest/internal/tarot"
	"time"

	"github.com/google/uuid"
)

type CardService struct {
	db *sql.DB
}

func NewCardService(db *sql.DB) *CardService {
	return &CardService{db: db}
}

func (s *CardService) PerformDailyDraw(userID uuid.UUID, mood, question string) (*models.CardDraw, error) {
	// Check if user already drew today
	today := time.Now().Format("2006-01-02")
	var existingDraw models.CardDraw

	err := s.db.QueryRow(`
		SELECT id, card_id, card_name, interpretation_basic, mood, question, created_at
		FROM card_draws 
		WHERE user_id = $1 AND draw_date = $2
	`, userID, today).Scan(
		&existingDraw.ID, &existingDraw.CardID, &existingDraw.CardName,
		&existingDraw.InterpretationBasic, &existingDraw.Mood,
		&existingDraw.Question, &existingDraw.CreatedAt,
	)

	if err == nil {
		existingDraw.UserID = userID
		existingDraw.DrawDate = today
		return &existingDraw, errors.New("daily draw already completed")
	}

	if err != sql.ErrNoRows {
		return nil, err
	}

	// Check subscription for usage limits
	isUnlimited, err := s.checkUserLimits(userID)
	if err != nil {
		return nil, err
	}

	if !isUnlimited {
		// Check daily usage for free users
		var drawsToday int
		err = s.db.QueryRow(`
			SELECT COALESCE(draws_count, 0) FROM daily_usage 
			WHERE user_id = $1 AND usage_date = $2
		`, userID, today).Scan(&drawsToday)

		if err == nil && drawsToday >= 1 {
			return nil, errors.New("daily limit reached - upgrade to premium for unlimited draws")
		}
	}

	// Select intelligent card
	cardID := tarot.SelectIntelligentCard(userID, s.db, mood, question)
	card, exists := tarot.MajorArcana[cardID]
	if !exists {
		return nil, errors.New("invalid card selected")
	}

	// Create card draw record
	drawID := uuid.New()
	_, err = s.db.Exec(`
		INSERT INTO card_draws (id, user_id, card_id, card_name, draw_date, 
		                       interpretation_basic, mood, question, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW())
	`, drawID, userID, cardID, card.Name, today, card.TraditionalMeaning, mood, question)

	if err != nil {
		return nil, err
	}

	// Update daily usage
	_, err = s.db.Exec(`
		INSERT INTO daily_usage (user_id, usage_date, draws_count)
		VALUES ($1, $2, 1)
		ON CONFLICT (user_id, usage_date)
		DO UPDATE SET draws_count = daily_usage.draws_count + 1
	`, userID, today)

	if err != nil {
		return nil, err
	}

	return &models.CardDraw{
		ID:                  drawID,
		UserID:             userID,
		CardID:             cardID,
		CardName:           card.Name,
		DrawDate:           today,
		InterpretationBasic: card.TraditionalMeaning,
		Mood:               mood,
		Question:           question,
		CreatedAt:          time.Now(),
	}, nil
}

func (s *CardService) GetDrawHistory(userID uuid.UUID, limit int) ([]models.CardDraw, error) {
	if limit <= 0 {
		limit = 20
	}

	rows, err := s.db.Query(`
		SELECT id, card_id, card_name, draw_date, interpretation_basic, 
		       COALESCE(interpretation_enhanced, ''), COALESCE(mood, ''), 
		       COALESCE(question, ''), created_at
		FROM card_draws 
		WHERE user_id = $1 
		ORDER BY created_at DESC 
		LIMIT $2
	`, userID, limit)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var draws []models.CardDraw
	for rows.Next() {
		var draw models.CardDraw
		err := rows.Scan(
			&draw.ID, &draw.CardID, &draw.CardName,
			&draw.DrawDate, &draw.InterpretationBasic,
			&draw.InterpretationEnhanced, &draw.Mood,
			&draw.Question, &draw.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		draw.UserID = userID
		draws = append(draws, draw)
	}

	return draws, nil
}

func (s *CardService) GetTodayStatus(userID uuid.UUID) (map[string]interface{}, error) {
	today := time.Now().Format("2006-01-02")

	var drawID uuid.UUID
	var cardID int
	var cardName string

	err := s.db.QueryRow(`
		SELECT id, card_id, card_name
		FROM card_draws 
		WHERE user_id = $1 AND draw_date = $2
	`, userID, today).Scan(&drawID, &cardID, &cardName)

	if err == sql.ErrNoRows {
		return map[string]interface{}{
			"has_drawn":    false,
			"can_draw":     true,
			"card":         nil,
			"draws_today":  0,
			"limit":        1,
		}, nil
	}

	if err != nil {
		return nil, err
	}

	// Get usage count
	var drawsToday int
	s.db.QueryRow(`
		SELECT COALESCE(draws_count, 0) FROM daily_usage 
		WHERE user_id = $1 AND usage_date = $2
	`, userID, today).Scan(&drawsToday)

	card := tarot.MajorArcana[cardID]

	return map[string]interface{}{
		"has_drawn":   true,
		"can_draw":    false,
		"card": map[string]interface{}{
			"id":   cardID,
			"name": cardName,
			"traditional_meaning": card.TraditionalMeaning,
		},
		"draws_today": drawsToday,
		"limit":       1,
	}, nil
}

func (s *CardService) GetCardMeaning(cardID int) (*tarot.Card, error) {
	if card, exists := tarot.MajorArcana[cardID]; exists {
		return &card, nil
	}
	return nil, errors.New("card not found")
}

func (s *CardService) SaveEnhancedInterpretation(userID uuid.UUID, drawDate string, interpretation string) error {
	_, err := s.db.Exec(`
		UPDATE card_draws 
		SET interpretation_enhanced = $1
		WHERE user_id = $2 AND draw_date = $3
	`, interpretation, userID, drawDate)

	return err
}

func (s *CardService) checkUserLimits(userID uuid.UUID) (bool, error) {
	var tier string
	err := s.db.QueryRow("SELECT subscription_tier FROM users WHERE id = $1", userID).Scan(&tier)

	if err != nil {
		return false, err
	}

	return tier == "premium", nil
}