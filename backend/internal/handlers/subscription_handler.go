package handlers

import (
	"io"
	"symbol-quest/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type SubscriptionHandler struct {
	stripeService *services.StripeService
}

func NewSubscriptionHandler(stripeService *services.StripeService) *SubscriptionHandler {
	return &SubscriptionHandler{stripeService: stripeService}
}

func (h *SubscriptionHandler) Create(c *fiber.Ctx) error {
	userIDStr := c.Locals("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid user ID",
		})
	}

	userEmail := c.Locals("user_email").(string)
	if userEmail == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "User email not found",
		})
	}

	clientSecret, err := h.stripeService.CreateSubscription(userID, userEmail)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to create subscription: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"client_secret": clientSecret,
		"message":       "Subscription created successfully",
	})
}

func (h *SubscriptionHandler) Status(c *fiber.Ctx) error {
	userIDStr := c.Locals("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid user ID",
		})
	}

	subscription, err := h.stripeService.GetSubscriptionStatus(userID)
	if err != nil {
		// No active subscription found
		return c.JSON(fiber.Map{
			"subscription": nil,
			"status":      "free",
			"message":     "No active subscription",
		})
	}

	return c.JSON(fiber.Map{
		"subscription": subscription,
		"status":      "premium",
	})
}

func (h *SubscriptionHandler) StripeWebhook(c *fiber.Ctx) error {
	signature := c.Get("Stripe-Signature")
	if signature == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Missing Stripe signature",
		})
	}

	payload, err := io.ReadAll(c.Context().RequestBodyStream())
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to read request body",
		})
	}

	err = h.stripeService.HandleWebhook(payload, signature)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Webhook processing failed: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"received": true,
	})
}