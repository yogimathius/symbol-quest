package services

import (
	"database/sql"
	"testing"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	// Import test database driver
	_ "github.com/lib/pq"
)

func TestAuthService_GenerateToken(t *testing.T) {
	service := &AuthService{
		jwtSecret: []byte("test-secret-key"),
	}

	userID := uuid.New()
	email := "test@example.com"
	tier := "free"

	token, err := service.GenerateToken(userID, email, tier)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	if token == "" {
		t.Error("Generated token is empty")
	}

	// Validate token can be parsed
	claims, err := service.ValidateToken(token)
	if err != nil {
		t.Fatalf("Failed to validate generated token: %v", err)
	}

	if claims["user_id"] != userID.String() {
		t.Errorf("Expected user_id %s, got %s", userID.String(), claims["user_id"])
	}

	if claims["email"] != email {
		t.Errorf("Expected email %s, got %s", email, claims["email"])
	}

	if claims["subscription_tier"] != tier {
		t.Errorf("Expected subscription_tier %s, got %s", tier, claims["subscription_tier"])
	}
}

func TestAuthService_ValidateToken(t *testing.T) {
	service := &AuthService{
		jwtSecret: []byte("test-secret-key"),
	}

	t.Run("ValidToken", func(t *testing.T) {
		userID := uuid.New()
		token, _ := service.GenerateToken(userID, "test@example.com", "free")

		claims, err := service.ValidateToken(token)
		if err != nil {
			t.Fatalf("Expected valid token to pass validation: %v", err)
		}

		if claims["user_id"] != userID.String() {
			t.Errorf("Claims user_id mismatch")
		}
	})

	t.Run("InvalidToken", func(t *testing.T) {
		_, err := service.ValidateToken("invalid.token.here")
		if err == nil {
			t.Error("Expected invalid token to fail validation")
		}
	})

	t.Run("EmptyToken", func(t *testing.T) {
		_, err := service.ValidateToken("")
		if err == nil {
			t.Error("Expected empty token to fail validation")
		}
	})

	t.Run("WrongSecret", func(t *testing.T) {
		// Generate token with different secret
		wrongService := &AuthService{jwtSecret: []byte("wrong-secret")}
		token, _ := wrongService.GenerateToken(uuid.New(), "test@example.com", "free")

		// Try to validate with correct service
		_, err := service.ValidateToken(token)
		if err == nil {
			t.Error("Expected token with wrong secret to fail validation")
		}
	})
}

func TestPasswordHashing(t *testing.T) {
	password := "test-password-123"

	// Test bcrypt hashing (similar to what AuthService.Register does)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	// Test password verification (similar to what AuthService.Login does)
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		t.Errorf("Password verification failed: %v", err)
	}

	// Test wrong password fails
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte("wrong-password"))
	if err == nil {
		t.Error("Expected wrong password to fail verification")
	}
}

func TestTokenExpiration(t *testing.T) {
	service := &AuthService{
		jwtSecret: []byte("test-secret-key"),
	}

	userID := uuid.New()
	token, err := service.GenerateToken(userID, "test@example.com", "free")
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	claims, err := service.ValidateToken(token)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}

	// Check expiration is set (should be 7 days from now)
	exp, ok := claims["exp"]
	if !ok {
		t.Error("Token missing expiration claim")
		return
	}

	// Convert to time and check it's in the future
	expTime := time.Unix(int64(exp.(float64)), 0)
	if expTime.Before(time.Now()) {
		t.Error("Token expiration is in the past")
	}

	// Check it's approximately 7 days from now (within 1 minute tolerance)
	expectedExp := time.Now().Add(time.Hour * 24 * 7)
	diff := expTime.Sub(expectedExp)
	if diff > time.Minute || diff < -time.Minute {
		t.Errorf("Token expiration not approximately 7 days: %v", diff)
	}
}

// Mock database tests (would require actual DB setup for integration tests)
type mockDB struct {
	users map[string]mockUser
}

type mockUser struct {
	id           uuid.UUID
	email        string
	passwordHash string
	tier         string
}

func (m *mockDB) QueryRow(query string, args ...interface{}) *sql.Row {
	// This would require more sophisticated mocking
	// For now, return nil to test error handling
	return nil
}

func TestAuthServiceIntegration(t *testing.T) {
	// Skip if not running integration tests
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}

	// These tests would require a real database connection
	// For now, just test the service creation
	service := NewAuthService(nil, "test-secret")
	if service == nil {
		t.Error("NewAuthService returned nil")
	}

	if string(service.jwtSecret) != "test-secret" {
		t.Error("JWT secret not set correctly")
	}
}