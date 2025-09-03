package main

import (
	"log"
	"os"
	"symbol-quest/internal/config"
	"symbol-quest/internal/database"
	"symbol-quest/internal/handlers"
	"symbol-quest/internal/middleware"
	"symbol-quest/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Run database migrations
	if err := database.RunMigrations(db); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	authService := services.NewAuthService(db, cfg.JWTSecret)
	cardService := services.NewCardService(db)
	openaiService := services.NewOpenAIService(cfg.OpenAIAPIKey)
	stripeService := services.NewStripeService(cfg.StripeSecretKey)
	stripeService.SetDatabase(db)
	stripeService.SetWebhookSecret(cfg.StripeWebhookSecret)

	authHandler := handlers.NewAuthHandler(authService)
	cardHandler := handlers.NewCardHandler(cardService, openaiService)
	subscriptionHandler := handlers.NewSubscriptionHandler(stripeService)

	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
	})

	app.Use(logger.New())
	app.Use(helmet.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: cfg.CORSOrigins,
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	api := app.Group("/api")

	// Auth routes
	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/logout", authHandler.Logout)
	auth.Get("/profile", middleware.AuthRequired(authService), authHandler.Profile)

	// Card draw routes
	draws := api.Group("/draws", middleware.AuthRequired(authService))
	draws.Post("/daily", cardHandler.DailyDraw)
	draws.Get("/history", cardHandler.History)
	draws.Get("/today", cardHandler.TodayStatus)

	// Interpretation routes
	interpretations := api.Group("/interpretations", middleware.AuthRequired(authService))
	interpretations.Post("/enhanced", middleware.PremiumRequired(), cardHandler.EnhancedInterpretation)

	// Card info routes
	cards := api.Group("/cards")
	cards.Get("/:id/meaning", cardHandler.BasicMeaning)

	// Subscription routes
	subscriptions := api.Group("/subscriptions", middleware.AuthRequired(authService))
	subscriptions.Post("/create", subscriptionHandler.Create)
	subscriptions.Get("/status", subscriptionHandler.Status)

	// Webhook routes
	webhooks := api.Group("/webhooks")
	webhooks.Post("/stripe", subscriptionHandler.StripeWebhook)

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}