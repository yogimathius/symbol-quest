package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"symbol-quest/internal/tarot"
)

type OpenAIService struct {
	apiKey string
	client *http.Client
}

type OpenAIRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	MaxTokens int      `json:"max_tokens"`
	Temperature float64 `json:"temperature"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIResponse struct {
	Choices []Choice `json:"choices"`
	Error   *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

type Choice struct {
	Message Message `json:"message"`
}

func NewOpenAIService(apiKey string) *OpenAIService {
	return &OpenAIService{
		apiKey: apiKey,
		client: &http.Client{},
	}
}

func (s *OpenAIService) GenerateEnhancedInterpretation(cardID int, mood, question string) (string, error) {
	if s.apiKey == "" {
		return "", errors.New("OpenAI API key not configured")
	}

	card, exists := tarot.MajorArcana[cardID]
	if !exists {
		return "", errors.New("invalid card ID")
	}

	prompt := s.buildPrompt(card, mood, question)

	req := OpenAIRequest{
		Model: "gpt-3.5-turbo",
		Messages: []Message{
			{
				Role:    "system",
				Content: "You are a wise and compassionate tarot reader who provides personalized, insightful interpretations that blend traditional tarot wisdom with modern psychological insights. Your readings are supportive, empowering, and help people gain clarity and perspective.",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
		MaxTokens:   400,
		Temperature: 0.7,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	httpReq, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+s.apiKey)

	resp, err := s.client.Do(httpReq)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var openaiResp OpenAIResponse
	err = json.Unmarshal(body, &openaiResp)
	if err != nil {
		return "", err
	}

	if openaiResp.Error != nil {
		return "", errors.New("OpenAI API error: " + openaiResp.Error.Message)
	}

	if len(openaiResp.Choices) == 0 {
		return "", errors.New("no response from OpenAI")
	}

	return openaiResp.Choices[0].Message.Content, nil
}

func (s *OpenAIService) buildPrompt(card tarot.Card, mood, question string) string {
	prompt := fmt.Sprintf(`Please provide a personalized tarot interpretation for:

Card: %s (%s)
Traditional Meaning: %s
Keywords: %v
Light Aspects: %v
Shadow Aspects: %v`, 
		card.Name, card.Number, card.TraditionalMeaning, 
		card.Keywords, card.LightAspects, card.ShadowAspects)

	if mood != "" {
		prompt += fmt.Sprintf("\nCurrent Mood: %s", mood)
	}

	if question != "" {
		prompt += fmt.Sprintf("\nQuestion Asked: %s", question)
	}

	prompt += `

Please provide:
1. A personalized interpretation that connects the card's meaning to their mood and question
2. Practical guidance and actionable insights
3. How this card's energy can help them right now
4. A supportive message that empowers them

Keep the tone warm, wise, and encouraging. Focus on personal growth and positive transformation while being honest about any challenges the card might indicate.

Response should be 2-3 paragraphs, around 250-300 words total.`

	return prompt
}