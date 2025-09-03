package handlers

import (
	"strconv"
	"symbol-quest/internal/models"
	"symbol-quest/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type CardHandler struct {
	cardService   *services.CardService
	openaiService *services.OpenAIService
}

func NewCardHandler(cardService *services.CardService, openaiService *services.OpenAIService) *CardHandler {
	return &CardHandler{
		cardService:   cardService,
		openaiService: openaiService,
	}
}

func (h *CardHandler) DailyDraw(c *fiber.Ctx) error {
	userIDStr := c.Locals("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid user ID",
		})
	}

	var req models.DailyDrawRequest
	if err := c.BodyParser(&req); err != nil {
		// If body parsing fails, continue with empty mood/question
		req.Mood = ""
		req.Question = ""
	}

	draw, err := h.cardService.PerformDailyDraw(userID, req.Mood, req.Question)
	if err != nil {
		if err.Error() == "daily draw already completed" {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error":   true,
				"message": "You have already drawn your card for today",
				"card":    draw,
			})
		}
		if err.Error() == "daily limit reached - upgrade to premium for unlimited draws" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error":   true,
				"message": "Daily limit reached. Upgrade to premium for unlimited draws.",
				"upgrade_required": true,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"card":    draw,
	})
}

func (h *CardHandler) History(c *fiber.Ctx) error {
	userIDStr := c.Locals("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid user ID",
		})
	}

	limitStr := c.Query("limit", "20")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 20
	}

	draws, err := h.cardService.GetDrawHistory(userID, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to fetch draw history",
		})
	}

	return c.JSON(fiber.Map{
		"draws": draws,
		"count": len(draws),
	})
}

func (h *CardHandler) TodayStatus(c *fiber.Ctx) error {
	userIDStr := c.Locals("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid user ID",
		})
	}

	status, err := h.cardService.GetTodayStatus(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to get today's status",
		})
	}

	return c.JSON(status)
}

func (h *CardHandler) BasicMeaning(c *fiber.Ctx) error {
	cardIDStr := c.Params("id")
	cardID, err := strconv.Atoi(cardIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid card ID",
		})
	}

	card, err := h.cardService.GetCardMeaning(cardID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "Card not found",
		})
	}

	return c.JSON(fiber.Map{
		"card": models.TarotCard{
			ID:           card.ID,
			Name:         card.Name,
			Keywords:     card.Keywords,
			Element:      card.Elements[0], // Take first element
			Astrology:    card.Astrology,
			BasicMeaning: card.TraditionalMeaning,
		},
	})
}

func (h *CardHandler) EnhancedInterpretation(c *fiber.Ctx) error {
	userIDStr := c.Locals("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid user ID",
		})
	}

	var req struct {
		CardID   int    `json:"card_id"`
		Mood     string `json:"mood"`
		Question string `json:"question"`
		DrawDate string `json:"draw_date"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
		})
	}

	// Generate enhanced interpretation using OpenAI
	interpretation, err := h.openaiService.GenerateEnhancedInterpretation(
		req.CardID, req.Mood, req.Question,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to generate enhanced interpretation: " + err.Error(),
		})
	}

	// Save the enhanced interpretation to the database
	if req.DrawDate != "" {
		err = h.cardService.SaveEnhancedInterpretation(userID, req.DrawDate, interpretation)
		if err != nil {
			// Log the error but don't fail the request
			// The user still gets their interpretation
		}
	}

	return c.JSON(fiber.Map{
		"interpretation": interpretation,
	})
}