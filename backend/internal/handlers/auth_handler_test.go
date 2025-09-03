package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"symbol-quest/internal/services"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestAuthHandler_Creation(t *testing.T) {
	mockAuthService := &services.AuthService{}
	handler := NewAuthHandler(mockAuthService)
	
	if handler == nil {
		t.Error("NewAuthHandler returned nil")
	}
	
	if handler.authService != mockAuthService {
		t.Error("AuthHandler not initialized with correct service")
	}
}

func TestAuthHandler_Register_Validation(t *testing.T) {
	mockAuthService := &services.AuthService{}
	handler := NewAuthHandler(mockAuthService)
	
	app := fiber.New()
	app.Post("/register", handler.Register)

	t.Run("InvalidJSON", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/register", bytes.NewBufferString("invalid json"))
		req.Header.Set("Content-Type", "application/json")
		
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}

		if resp.StatusCode != fiber.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", fiber.StatusBadRequest, resp.StatusCode)
		}
	})

	t.Run("MissingEmail", func(t *testing.T) {
		reqBody := map[string]string{
			"password": "validpassword123",
		}
		jsonBody, _ := json.Marshal(reqBody)
		
		req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}

		if resp.StatusCode != fiber.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", fiber.StatusBadRequest, resp.StatusCode)
		}

		body, _ := io.ReadAll(resp.Body)
		bodyStr := string(body)
		
		if !contains(bodyStr, "Email and password are required") {
			t.Errorf("Expected validation error message, got: %s", bodyStr)
		}
	})

	t.Run("MissingPassword", func(t *testing.T) {
		reqBody := map[string]string{
			"email": "test@example.com",
		}
		jsonBody, _ := json.Marshal(reqBody)
		
		req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}

		if resp.StatusCode != fiber.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", fiber.StatusBadRequest, resp.StatusCode)
		}
	})

	t.Run("ShortPassword", func(t *testing.T) {
		reqBody := map[string]string{
			"email":    "test@example.com",
			"password": "short",
		}
		jsonBody, _ := json.Marshal(reqBody)
		
		req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}

		if resp.StatusCode != fiber.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", fiber.StatusBadRequest, resp.StatusCode)
		}

		body, _ := io.ReadAll(resp.Body)
		bodyStr := string(body)
		
		if !contains(bodyStr, "Password must be at least 8 characters") {
			t.Errorf("Expected password length error, got: %s", bodyStr)
		}
	})

	t.Run("EmptyRequest", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/register", bytes.NewBufferString("{}"))
		req.Header.Set("Content-Type", "application/json")
		
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}

		if resp.StatusCode != fiber.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", fiber.StatusBadRequest, resp.StatusCode)
		}
	})
}

func TestAuthHandler_Login_Validation(t *testing.T) {
	mockAuthService := &services.AuthService{}
	handler := NewAuthHandler(mockAuthService)
	
	app := fiber.New()
	app.Post("/login", handler.Login)

	t.Run("InvalidJSON", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/login", bytes.NewBufferString("invalid json"))
		req.Header.Set("Content-Type", "application/json")
		
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}

		if resp.StatusCode != fiber.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", fiber.StatusBadRequest, resp.StatusCode)
		}
	})

	t.Run("MissingCredentials", func(t *testing.T) {
		reqBody := map[string]string{}
		jsonBody, _ := json.Marshal(reqBody)
		
		req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}

		if resp.StatusCode != fiber.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", fiber.StatusBadRequest, resp.StatusCode)
		}

		body, _ := io.ReadAll(resp.Body)
		bodyStr := string(body)
		
		if !contains(bodyStr, "Email and password are required") {
			t.Errorf("Expected validation error message, got: %s", bodyStr)
		}
	})
}

func TestAuthHandler_Logout(t *testing.T) {
	mockAuthService := &services.AuthService{}
	handler := NewAuthHandler(mockAuthService)
	
	app := fiber.New()
	app.Post("/logout", handler.Logout)

	req := httptest.NewRequest("POST", "/logout", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Expected status %d, got %d", fiber.StatusOK, resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	bodyStr := string(body)
	
	if !contains(bodyStr, "Logged out successfully") {
		t.Errorf("Expected logout success message, got: %s", bodyStr)
	}
}

func TestAuthHandler_Profile_Validation(t *testing.T) {
	mockAuthService := &services.AuthService{}
	handler := NewAuthHandler(mockAuthService)
	
	app := fiber.New()
	app.Get("/profile", func(c *fiber.Ctx) error {
		// Simulate middleware setting invalid user_id
		c.Locals("user_id", "invalid-uuid")
		return handler.Profile(c)
	})

	req := httptest.NewRequest("GET", "/profile", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if resp.StatusCode != fiber.StatusBadRequest {
		t.Errorf("Expected status %d for invalid UUID, got %d", fiber.StatusBadRequest, resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	bodyStr := string(body)
	
	if !contains(bodyStr, "Invalid user ID") {
		t.Errorf("Expected invalid user ID error, got: %s", bodyStr)
	}
}

func TestRequestValidation(t *testing.T) {
	// Test various request validation scenarios
	
	t.Run("EmailValidation", func(t *testing.T) {
		validEmails := []string{
			"test@example.com",
			"user.name@domain.co.uk",
			"test+filter@gmail.com",
		}
		
		for _, email := range validEmails {
			if email == "" {
				t.Errorf("Valid email should not be empty: %s", email)
			}
			if !contains(email, "@") {
				t.Errorf("Valid email should contain @: %s", email)
			}
		}
	})

	t.Run("PasswordValidation", func(t *testing.T) {
		validPasswords := []string{
			"password123",
			"mySecurePass!",
			"12345678",
		}
		
		for _, password := range validPasswords {
			if len(password) < 8 {
				t.Errorf("Valid password should be at least 8 chars: %s", password)
			}
		}

		invalidPasswords := []string{
			"short",
			"1234567",
			"",
		}
		
		for _, password := range invalidPasswords {
			if len(password) >= 8 {
				t.Errorf("Invalid password should be less than 8 chars: %s", password)
			}
		}
	})
}

func TestJSONResponseStructure(t *testing.T) {
	mockAuthService := &services.AuthService{}
	handler := NewAuthHandler(mockAuthService)
	
	app := fiber.New()
	app.Post("/register", handler.Register)

	// Test that error responses have consistent structure
	reqBody := map[string]string{
		"email":    "",
		"password": "",
	}
	jsonBody, _ := json.Marshal(reqBody)
	
	req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	body, _ := io.ReadAll(resp.Body)
	
	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		t.Fatalf("Response is not valid JSON: %v", err)
	}

	// Check required fields in error response
	if _, exists := response["error"]; !exists {
		t.Error("Error response missing 'error' field")
	}

	if _, exists := response["message"]; !exists {
		t.Error("Error response missing 'message' field")
	}

	if response["error"] != true {
		t.Error("Error response 'error' field should be true")
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	if len(substr) == 0 {
		return true
	}
	if len(s) < len(substr) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}