package middleware

import (
	"io"
	"net/http/httptest"
	"symbol-quest/internal/services"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func TestErrorHandler(t *testing.T) {
	app := fiber.New(fiber.Config{
		ErrorHandler: ErrorHandler,
	})

	t.Run("FiberError", func(t *testing.T) {
		app.Get("/test-error", func(c *fiber.Ctx) error {
			return fiber.NewError(fiber.StatusBadRequest, "Test error message")
		})

		req := httptest.NewRequest("GET", "/test-error", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}

		if resp.StatusCode != fiber.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", fiber.StatusBadRequest, resp.StatusCode)
		}

		body, _ := io.ReadAll(resp.Body)
		bodyStr := string(body)
		
		if !contains(bodyStr, "Test error message") {
			t.Errorf("Expected error message in response body, got: %s", bodyStr)
		}
	})

	t.Run("GenericError", func(t *testing.T) {
		app.Get("/test-generic", func(c *fiber.Ctx) error {
			return fiber.ErrInternalServerError
		})

		req := httptest.NewRequest("GET", "/test-generic", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}

		if resp.StatusCode != fiber.StatusInternalServerError {
			t.Errorf("Expected status %d, got %d", fiber.StatusInternalServerError, resp.StatusCode)
		}
	})
}

func TestAuthRequired(t *testing.T) {
	// Create a mock auth service
	mockAuth := &services.AuthService{}
	
	app := fiber.New()
	
	// Set up test route with auth middleware
	app.Get("/protected", AuthRequired(mockAuth), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "success"})
	})

	t.Run("NoAuthHeader", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/protected", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}

		if resp.StatusCode != fiber.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", fiber.StatusUnauthorized, resp.StatusCode)
		}

		body, _ := io.ReadAll(resp.Body)
		bodyStr := string(body)
		
		if !contains(bodyStr, "Authorization header required") {
			t.Errorf("Expected auth header error message, got: %s", bodyStr)
		}
	})

	t.Run("InvalidAuthHeader", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Invalid token format")
		
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}

		if resp.StatusCode != fiber.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", fiber.StatusUnauthorized, resp.StatusCode)
		}

		body, _ := io.ReadAll(resp.Body)
		bodyStr := string(body)
		
		if !contains(bodyStr, "Bearer token required") {
			t.Errorf("Expected bearer token error message, got: %s", bodyStr)
		}
	})

	t.Run("EmptyBearerToken", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer ")
		
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}

		if resp.StatusCode != fiber.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", fiber.StatusUnauthorized, resp.StatusCode)
		}
	})
}

func TestPremiumRequired(t *testing.T) {
	app := fiber.New()

	// Set up test route with premium middleware
	app.Get("/premium", PremiumRequired(), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "premium content"})
	})

	t.Run("NoSubscriptionTier", func(t *testing.T) {
		// Test without setting subscription_tier in locals
		req := httptest.NewRequest("GET", "/premium", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}

		if resp.StatusCode != fiber.StatusForbidden {
			t.Errorf("Expected status %d, got %d", fiber.StatusForbidden, resp.StatusCode)
		}

		body, _ := io.ReadAll(resp.Body)
		bodyStr := string(body)
		
		if !contains(bodyStr, "Premium subscription required") {
			t.Errorf("Expected premium required message, got: %s", bodyStr)
		}
	})

	t.Run("FreeTier", func(t *testing.T) {
		app.Get("/premium-test", func(c *fiber.Ctx) error {
			c.Locals("subscription_tier", "free")
			return PremiumRequired()(c)
		}, func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"message": "should not reach here"})
		})

		req := httptest.NewRequest("GET", "/premium-test", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}

		if resp.StatusCode != fiber.StatusForbidden {
			t.Errorf("Expected status %d for free tier, got %d", fiber.StatusForbidden, resp.StatusCode)
		}
	})
}

func TestMiddlewareIntegration(t *testing.T) {
	app := fiber.New()

	// Mock successful authentication by setting locals
	app.Use("/api", func(c *fiber.Ctx) error {
		// Simulate authenticated user
		c.Locals("user_id", uuid.New().String())
		c.Locals("user_email", "test@example.com")
		c.Locals("subscription_tier", "premium")
		return c.Next()
	})

	app.Get("/api/premium", PremiumRequired(), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "premium success"})
	})

	t.Run("PremiumUserSuccess", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/premium", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}

		if resp.StatusCode != fiber.StatusOK {
			t.Errorf("Expected status %d for premium user, got %d", fiber.StatusOK, resp.StatusCode)
		}

		body, _ := io.ReadAll(resp.Body)
		bodyStr := string(body)
		
		if !contains(bodyStr, "premium success") {
			t.Errorf("Expected premium success message, got: %s", bodyStr)
		}
	})
}

func TestMiddlewareErrorResponses(t *testing.T) {
	// Test that middleware returns proper JSON error responses
	app := fiber.New(fiber.Config{
		ErrorHandler: ErrorHandler,
	})

	app.Get("/test", AuthRequired(nil), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "success"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	// Check Content-Type is JSON
	contentType := resp.Header.Get("Content-Type")
	if !contains(contentType, "application/json") {
		t.Errorf("Expected JSON content type, got: %s", contentType)
	}

	// Check response structure
	body, _ := io.ReadAll(resp.Body)
	bodyStr := string(body)
	
	if !contains(bodyStr, `"error"`) {
		t.Errorf("Expected error field in JSON response, got: %s", bodyStr)
	}
	
	if !contains(bodyStr, `"message"`) {
		t.Errorf("Expected message field in JSON response, got: %s", bodyStr)
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || (len(s) > len(substr) && containsHelper(s, substr)))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}