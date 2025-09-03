package models

import (
	"time"
	"github.com/google/uuid"
)

type User struct {
	ID               uuid.UUID `json:"id" db:"id"`
	Email           string    `json:"email" db:"email"`
	PasswordHash    string    `json:"-" db:"password_hash"`
	SubscriptionTier string    `json:"subscription_tier" db:"subscription_tier"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

type CardDraw struct {
	ID                     uuid.UUID `json:"id" db:"id"`
	UserID                uuid.UUID `json:"user_id" db:"user_id"`
	CardID                int       `json:"card_id" db:"card_id"`
	CardName              string    `json:"card_name" db:"card_name"`
	DrawDate              string    `json:"draw_date" db:"draw_date"`
	InterpretationBasic   string    `json:"interpretation_basic" db:"interpretation_basic"`
	InterpretationEnhanced string   `json:"interpretation_enhanced,omitempty" db:"interpretation_enhanced"`
	Mood                  string    `json:"mood,omitempty" db:"mood"`
	Question              string    `json:"question,omitempty" db:"question"`
	CreatedAt             time.Time `json:"created_at" db:"created_at"`
}

type DailyUsage struct {
	ID         uuid.UUID `json:"id" db:"id"`
	UserID     uuid.UUID `json:"user_id" db:"user_id"`
	UsageDate  string    `json:"usage_date" db:"usage_date"`
	DrawsCount int       `json:"draws_count" db:"draws_count"`
}

type Subscription struct {
	ID                    uuid.UUID  `json:"id" db:"id"`
	UserID               uuid.UUID  `json:"user_id" db:"user_id"`
	StripeSubscriptionID string     `json:"stripe_subscription_id" db:"stripe_subscription_id"`
	StripeCustomerID     string     `json:"stripe_customer_id" db:"stripe_customer_id"`
	Status               string     `json:"status" db:"status"`
	CurrentPeriodStart   *time.Time `json:"current_period_start" db:"current_period_start"`
	CurrentPeriodEnd     *time.Time `json:"current_period_end" db:"current_period_end"`
	CreatedAt            time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at" db:"updated_at"`
}

// Request/Response models
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type DailyDrawRequest struct {
	Mood     string `json:"mood,omitempty"`
	Question string `json:"question,omitempty"`
}

type TarotCard struct {
	ID           int      `json:"id"`
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Keywords     []string `json:"keywords"`
	Element      string   `json:"element,omitempty"`
	Astrology    string   `json:"astrology,omitempty"`
	BasicMeaning string   `json:"basic_meaning"`
	ImageURL     string   `json:"image_url,omitempty"`
}