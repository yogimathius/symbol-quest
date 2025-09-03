package services

import (
	"database/sql"
	"errors"
	"symbol-quest/internal/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	db        *sql.DB
	jwtSecret []byte
}

func NewAuthService(db *sql.DB, jwtSecret string) *AuthService {
	return &AuthService{
		db:        db,
		jwtSecret: []byte(jwtSecret),
	}
}

func (s *AuthService) Register(email, password string) (*models.User, error) {
	// Check if user already exists
	var existingUser models.User
	err := s.db.QueryRow("SELECT id FROM users WHERE email = $1", email).Scan(&existingUser.ID)
	if err == nil {
		return nil, errors.New("user already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Create user
	userID := uuid.New()
	_, err = s.db.Exec(`
		INSERT INTO users (id, email, password_hash, subscription_tier, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
	`, userID, email, string(hashedPassword), "free")

	if err != nil {
		return nil, errors.New("failed to create user")
	}

	user := &models.User{
		ID:               userID,
		Email:           email,
		SubscriptionTier: "free",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	return user, nil
}

func (s *AuthService) Login(email, password string) (*models.User, string, error) {
	var user models.User
	var passwordHash string

	err := s.db.QueryRow(`
		SELECT id, email, password_hash, subscription_tier, created_at, updated_at
		FROM users WHERE email = $1
	`, email).Scan(
		&user.ID, &user.Email, &passwordHash, &user.SubscriptionTier,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	if err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	// Generate JWT token
	token, err := s.GenerateToken(user.ID, user.Email, user.SubscriptionTier)
	if err != nil {
		return nil, "", errors.New("failed to generate token")
	}

	return &user, token, nil
}

func (s *AuthService) GetUserByID(userID uuid.UUID) (*models.User, error) {
	var user models.User

	err := s.db.QueryRow(`
		SELECT id, email, subscription_tier, created_at, updated_at
		FROM users WHERE id = $1
	`, userID).Scan(
		&user.ID, &user.Email, &user.SubscriptionTier,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		return nil, errors.New("user not found")
	}

	return &user, nil
}

func (s *AuthService) GenerateToken(userID uuid.UUID, email, subscriptionTier string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":           userID.String(),
		"email":            email,
		"subscription_tier": subscriptionTier,
		"exp":              time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 days
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func (s *AuthService) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func (s *AuthService) UpdateSubscriptionTier(userID uuid.UUID, tier string) error {
	_, err := s.db.Exec(`
		UPDATE users SET subscription_tier = $1, updated_at = NOW()
		WHERE id = $2
	`, tier, userID)

	return err
}