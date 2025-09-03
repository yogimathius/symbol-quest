package services

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"symbol-quest/internal/models"

	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/customer"
	"github.com/stripe/stripe-go/v76/subscription"
	"github.com/stripe/stripe-go/v76/webhook"
)

type StripeService struct {
	db               *sql.DB
	webhookSecret    string
}

func NewStripeService(secretKey string) *StripeService {
	stripe.Key = secretKey
	return &StripeService{
		webhookSecret: "", // Will be set later
	}
}

func (s *StripeService) SetDatabase(db *sql.DB) {
	s.db = db
}

func (s *StripeService) SetWebhookSecret(secret string) {
	s.webhookSecret = secret
}

func (s *StripeService) CreateSubscription(userID uuid.UUID, userEmail string) (string, error) {
	// Create or get Stripe customer
	customerParams := &stripe.CustomerParams{
		Email: stripe.String(userEmail),
		Metadata: map[string]string{
			"user_id": userID.String(),
		},
	}
	
	stripeCustomer, err := customer.New(customerParams)
	if err != nil {
		return "", errors.New("failed to create customer: " + err.Error())
	}

	// Get premium price ID - you'll need to create this in Stripe Dashboard
	priceID := "price_premium_monthly" // Replace with actual price ID from Stripe

	// Create subscription
	subscriptionParams := &stripe.SubscriptionParams{
		Customer: stripe.String(stripeCustomer.ID),
		Items: []*stripe.SubscriptionItemsParams{
			{
				Price: stripe.String(priceID),
			},
		},
		PaymentBehavior: stripe.String("default_incomplete"),
		PaymentSettings: &stripe.SubscriptionPaymentSettingsParams{
			SaveDefaultPaymentMethod: stripe.String("on_subscription"),
		},
		Metadata: map[string]string{
			"user_id": userID.String(),
		},
	}

	subscription, err := subscription.New(subscriptionParams)
	if err != nil {
		return "", errors.New("failed to create subscription: " + err.Error())
	}

	// Save subscription to database
	err = s.saveSubscription(userID, subscription)
	if err != nil {
		log.Printf("Failed to save subscription to database: %v", err)
	}

	// Return client secret for frontend to complete payment
	return subscription.LatestInvoice.PaymentIntent.ClientSecret, nil
}

func (s *StripeService) GetSubscriptionStatus(userID uuid.UUID) (*models.Subscription, error) {
	var subscription models.Subscription

	err := s.db.QueryRow(`
		SELECT id, stripe_subscription_id, stripe_customer_id, status,
		       current_period_start, current_period_end, created_at, updated_at
		FROM subscriptions 
		WHERE user_id = $1 AND status IN ('active', 'trialing')
		ORDER BY created_at DESC LIMIT 1
	`, userID).Scan(
		&subscription.ID, &subscription.StripeSubscriptionID,
		&subscription.StripeCustomerID, &subscription.Status,
		&subscription.CurrentPeriodStart, &subscription.CurrentPeriodEnd,
		&subscription.CreatedAt, &subscription.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("no active subscription found")
	}

	if err != nil {
		return nil, err
	}

	subscription.UserID = userID
	return &subscription, nil
}

func (s *StripeService) HandleWebhook(payload []byte, signature string) error {
	if s.webhookSecret == "" {
		return errors.New("webhook secret not configured")
	}

	event, err := webhook.ConstructEvent(payload, signature, s.webhookSecret)
	if err != nil {
		return errors.New("invalid webhook signature: " + err.Error())
	}

	switch event.Type {
	case "customer.subscription.created":
		var subscription stripe.Subscription
		err := json.Unmarshal(event.Data.Raw, &subscription)
		if err != nil {
			return err
		}
		return s.handleSubscriptionCreated(&subscription)

	case "customer.subscription.updated":
		var subscription stripe.Subscription
		err := json.Unmarshal(event.Data.Raw, &subscription)
		if err != nil {
			return err
		}
		return s.handleSubscriptionUpdated(&subscription)

	case "customer.subscription.deleted":
		var subscription stripe.Subscription
		err := json.Unmarshal(event.Data.Raw, &subscription)
		if err != nil {
			return err
		}
		return s.handleSubscriptionDeleted(&subscription)

	case "invoice.payment_succeeded":
		var invoice stripe.Invoice
		err := json.Unmarshal(event.Data.Raw, &invoice)
		if err != nil {
			return err
		}
		return s.handlePaymentSucceeded(&invoice)

	case "invoice.payment_failed":
		var invoice stripe.Invoice
		err := json.Unmarshal(event.Data.Raw, &invoice)
		if err != nil {
			return err
		}
		return s.handlePaymentFailed(&invoice)
	}

	return nil
}

func (s *StripeService) saveSubscription(userID uuid.UUID, subscription *stripe.Subscription) error {
	_, err := s.db.Exec(`
		INSERT INTO subscriptions 
		(id, user_id, stripe_subscription_id, stripe_customer_id, status,
		 current_period_start, current_period_end, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
		ON CONFLICT (stripe_subscription_id) 
		DO UPDATE SET 
			status = $5,
			current_period_start = $6,
			current_period_end = $7,
			updated_at = NOW()
	`, uuid.New(), userID, subscription.ID, subscription.Customer.ID,
		string(subscription.Status),
		subscription.CurrentPeriodStart,
		subscription.CurrentPeriodEnd)

	return err
}

func (s *StripeService) handleSubscriptionCreated(subscription *stripe.Subscription) error {
	userID, err := s.getUserIDFromMetadata(subscription.Metadata)
	if err != nil {
		return err
	}

	return s.saveSubscription(userID, subscription)
}

func (s *StripeService) handleSubscriptionUpdated(subscription *stripe.Subscription) error {
	// Update subscription status
	_, err := s.db.Exec(`
		UPDATE subscriptions SET 
			status = $1,
			current_period_start = $2,
			current_period_end = $3,
			updated_at = NOW()
		WHERE stripe_subscription_id = $4
	`, string(subscription.Status),
		subscription.CurrentPeriodStart,
		subscription.CurrentPeriodEnd,
		subscription.ID)

	if err != nil {
		return err
	}

	// Update user subscription tier
	var userID uuid.UUID
	err = s.db.QueryRow("SELECT user_id FROM subscriptions WHERE stripe_subscription_id = $1", 
		subscription.ID).Scan(&userID)
	if err != nil {
		return err
	}

	tier := "free"
	if subscription.Status == "active" || subscription.Status == "trialing" {
		tier = "premium"
	}

	_, err = s.db.Exec("UPDATE users SET subscription_tier = $1, updated_at = NOW() WHERE id = $2", 
		tier, userID)

	return err
}

func (s *StripeService) handleSubscriptionDeleted(subscription *stripe.Subscription) error {
	// Update subscription status to canceled
	_, err := s.db.Exec(`
		UPDATE subscriptions SET status = 'canceled', updated_at = NOW()
		WHERE stripe_subscription_id = $1
	`, subscription.ID)

	if err != nil {
		return err
	}

	// Update user to free tier
	var userID uuid.UUID
	err = s.db.QueryRow("SELECT user_id FROM subscriptions WHERE stripe_subscription_id = $1", 
		subscription.ID).Scan(&userID)
	if err != nil {
		return err
	}

	_, err = s.db.Exec("UPDATE users SET subscription_tier = 'free', updated_at = NOW() WHERE id = $1", 
		userID)

	return err
}

func (s *StripeService) handlePaymentSucceeded(invoice *stripe.Invoice) error {
	log.Printf("Payment succeeded for subscription: %s", invoice.Subscription.ID)
	return nil
}

func (s *StripeService) handlePaymentFailed(invoice *stripe.Invoice) error {
	log.Printf("Payment failed for subscription: %s", invoice.Subscription.ID)
	return nil
}

func (s *StripeService) getUserIDFromMetadata(metadata map[string]string) (uuid.UUID, error) {
	userIDStr, exists := metadata["user_id"]
	if !exists {
		return uuid.Nil, errors.New("user_id not found in metadata")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.Nil, errors.New("invalid user_id in metadata")
	}

	return userID, nil
}