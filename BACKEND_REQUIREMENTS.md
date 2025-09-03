# Symbol Quest - Backend Requirements (Go + Fiber)

## **Current Status: React Frontend Complete (90%), Go Backend Needed**
- Complete React frontend with tarot card drawing interface ✅
- **Missing**: Go + Fiber backend API for user accounts, card draws, OpenAI integration

---

## **Backend Technology: Go + Fiber**

**Why Go + Fiber (Priority Week 1):**
- **Ultra-fast API responses** (<50ms for card draws)
- **Efficient JSON handling** for tarot card data
- **Low resource usage** for cost-effective deployment
- **Simple deployment** (single binary to Fly.io)
- **Excellent concurrency** for multiple user sessions
- **Fast development** (1 week to completion)

---

## **Required API Endpoints**

```go
// Authentication
POST /api/auth/register    // User registration
POST /api/auth/login       // User login
GET /api/auth/profile      // Get user profile
POST /api/auth/logout      // Logout user

// Card draws
POST /api/draws/daily      // Perform daily card draw
GET /api/draws/history     // Get user's draw history
GET /api/draws/today       // Check if daily draw available

// Interpretations (OpenAI integration)
POST /api/interpretations/enhanced  // Get AI interpretation (premium)
GET /api/cards/:id/meaning          // Get basic card meaning

// Subscriptions (Stripe)
POST /api/subscriptions/create      // Create subscription
GET /api/subscriptions/status       // Check subscription status
POST /api/webhooks/stripe           // Stripe webhook handler
```

---

## **Database Schema (PostgreSQL)**

```sql
-- Users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    subscription_tier VARCHAR(20) DEFAULT 'free', -- free, premium
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Card draws table
CREATE TABLE card_draws (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    card_id INTEGER NOT NULL, -- Major Arcana 0-21
    card_name VARCHAR(100) NOT NULL,
    draw_date DATE NOT NULL,
    interpretation_basic TEXT,
    interpretation_enhanced TEXT, -- premium AI interpretation
    created_at TIMESTAMP DEFAULT NOW()
);

-- Usage tracking for freemium limits
CREATE TABLE daily_usage (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    usage_date DATE NOT NULL DEFAULT CURRENT_DATE,
    draws_count INTEGER DEFAULT 0,
    UNIQUE(user_id, usage_date)
);

-- Stripe subscriptions
CREATE TABLE subscriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    stripe_subscription_id VARCHAR(255) UNIQUE NOT NULL,
    stripe_customer_id VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL, -- active, canceled, past_due, etc.
    current_period_start TIMESTAMP,
    current_period_end TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Indexes for performance
CREATE INDEX idx_card_draws_user_date ON card_draws(user_id, draw_date);
CREATE INDEX idx_daily_usage_user_date ON daily_usage(user_id, usage_date);
CREATE INDEX idx_subscriptions_user ON subscriptions(user_id);
```

---

## **Go Backend Architecture**

```go
// main.go
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
    // Load configuration
    cfg := config.Load()
    
    // Initialize database
    db, err := database.Connect(cfg.DatabaseURL)
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    
    // Initialize services
    authService := services.NewAuthService(db, cfg.JWTSecret)
    cardService := services.NewCardService(db)
    openaiService := services.NewOpenAIService(cfg.OpenAIAPIKey)
    stripeService := services.NewStripeService(cfg.StripeSecretKey)
    
    // Initialize handlers
    authHandler := handlers.NewAuthHandler(authService)
    cardHandler := handlers.NewCardHandler(cardService, openaiService)
    subscriptionHandler := handlers.NewSubscriptionHandler(stripeService)
    
    // Initialize Fiber app
    app := fiber.New(fiber.Config{
        ErrorHandler: middleware.ErrorHandler,
    })
    
    // Middleware
    app.Use(logger.New())
    app.Use(helmet.New())
    app.Use(cors.New(cors.Config{
        AllowOrigins: cfg.CORSOrigins,
        AllowHeaders: "Origin, Content-Type, Accept, Authorization",
    }))
    
    // Routes
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
    
    // Start server
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    
    log.Printf("Server starting on port %s", port)
    log.Fatal(app.Listen(":" + port))
}
```

---

## **Card Service Implementation**

```go
// internal/services/card_service.go
package services

import (
    "database/sql"
    "errors"
    "time"
    "symbol-quest/internal/models"
    "symbol-quest/internal/tarot"
    
    "github.com/google/uuid"
)

type CardService struct {
    db *sql.DB
}

func NewCardService(db *sql.DB) *CardService {
    return &CardService{db: db}
}

func (s *CardService) PerformDailyDraw(userID uuid.UUID) (*models.CardDraw, error) {
    // Check if user already drew today
    today := time.Now().Format("2006-01-02")
    var existingDraw models.CardDraw
    
    err := s.db.QueryRow(
        "SELECT id, card_id, card_name FROM card_draws WHERE user_id = $1 AND draw_date = $2",
        userID, today,
    ).Scan(&existingDraw.ID, &existingDraw.CardID, &existingDraw.CardName)
    
    if err == nil {
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
        err = s.db.QueryRow(
            "SELECT draws_count FROM daily_usage WHERE user_id = $1 AND usage_date = $2",
            userID, today,
        ).Scan(&drawsToday)
        
        if err == nil && drawsToday >= 1 {
            return nil, errors.New("daily limit reached - upgrade to premium for unlimited draws")
        }
    }
    
    // Select random card using intelligent algorithm
    cardID := tarot.SelectIntelligentCard(userID, s.db)
    card := tarot.MajorArcana[cardID]
    
    // Create card draw record
    drawID := uuid.New()
    _, err = s.db.Exec(`
        INSERT INTO card_draws (id, user_id, card_id, card_name, draw_date, interpretation_basic, created_at)
        VALUES ($1, $2, $3, $4, $5, $6, NOW())
    `, drawID, userID, cardID, card.Name, today, card.BasicMeaning)
    
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
        ID:                 drawID,
        UserID:            userID,
        CardID:            cardID,
        CardName:          card.Name,
        DrawDate:          today,
        InterpretationBasic: card.BasicMeaning,
    }, nil
}

func (s *CardService) GetDrawHistory(userID uuid.UUID, limit int) ([]models.CardDraw, error) {
    rows, err := s.db.Query(`
        SELECT id, card_id, card_name, draw_date, interpretation_basic, interpretation_enhanced, created_at
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
        var enhancedInterp sql.NullString
        
        err := rows.Scan(
            &draw.ID, &draw.CardID, &draw.CardName, 
            &draw.DrawDate, &draw.InterpretationBasic, 
            &enhancedInterp, &draw.CreatedAt,
        )
        if err != nil {
            return nil, err
        }
        
        if enhancedInterp.Valid {
            draw.InterpretationEnhanced = enhancedInterp.String
        }
        
        draws = append(draws, draw)
    }
    
    return draws, nil
}

func (s *CardService) checkUserLimits(userID uuid.UUID) (bool, error) {
    var tier string
    err := s.db.QueryRow(
        "SELECT subscription_tier FROM users WHERE id = $1",
        userID,
    ).Scan(&tier)
    
    if err != nil {
        return false, err
    }
    
    return tier == "premium", nil
}
```

---

## **Deployment Strategy**

```bash
# Go deployment to Fly.io
fly launch --name symbol-quest-api
fly postgres create symbol-quest-db
fly secrets set DATABASE_URL=postgresql://...
fly secrets set JWT_SECRET=your-jwt-secret
fly secrets set OPENAI_API_KEY=sk-...
fly secrets set STRIPE_SECRET_KEY=sk-...
fly deploy

# Environment variables needed
DATABASE_URL=postgresql://symbol-quest-db.internal:5432/symbol_quest
JWT_SECRET=your-256-bit-secret
OPENAI_API_KEY=sk-proj-...
STRIPE_SECRET_KEY=sk-test-...
STRIPE_WEBHOOK_SECRET=whsec_...
CORS_ORIGINS=https://symbol-quest.vercel.app
```

---

## **Development Timeline (Week 1 Priority)**

**Day 1-2**: Go project setup, database schema, authentication
**Day 3-4**: Card draw logic, OpenAI integration, Stripe payments
**Day 5-6**: Testing, deployment, frontend integration  
**Day 7**: Polish, monitoring, production launch

**Estimated Development**: **1 week to full launch** ✅

---

## **Performance Targets**

- **API Response Time**: <50ms for card draws
- **Authentication**: <30ms JWT validation
- **Database Queries**: <10ms average query time
- **Deployment Size**: <20MB binary
- **Memory Usage**: <100MB RAM
- **Concurrent Users**: 1000+ simultaneous draws