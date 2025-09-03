package tarot

import (
	"database/sql"
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Card struct {
	ID              int                    `json:"id"`
	Name            string                 `json:"name"`
	Number          string                 `json:"number"`
	Keywords        []string               `json:"keywords"`
	Archetypes      []string               `json:"archetypes"`
	Elements        []string               `json:"elements"`
	Astrology       string                 `json:"astrology"`
	TraditionalMeaning string              `json:"traditional_meaning"`
	ShadowAspects   []string               `json:"shadow_aspects"`
	LightAspects    []string               `json:"light_aspects"`
	MoodWeights     map[string]float64     `json:"mood_weights"`
}

var MajorArcana = map[int]Card{
	0: {
		ID: 0, Name: "The Fool", Number: "0",
		Keywords: []string{"new-beginnings", "innocence", "spontaneity", "faith", "potential"},
		Archetypes: []string{"innocent", "seeker", "beginner"},
		Elements: []string{"air"}, Astrology: "Uranus",
		TraditionalMeaning: "New beginnings, innocence, spontaneity, leap of faith",
		ShadowAspects: []string{"recklessness", "naivety", "foolishness", "poor judgment"},
		LightAspects: []string{"faith", "optimism", "adventure", "trust", "openness"},
		MoodWeights: map[string]float64{
			"anxious": 0.3, "excited": 1.2, "uncertain": 1.1, "hopeful": 1.3,
			"peaceful": 0.8, "frustrated": 0.7, "curious": 1.2, "contemplative": 0.9,
		},
	},
	1: {
		ID: 1, Name: "The Magician", Number: "I",
		Keywords: []string{"manifestation", "power", "skill", "concentration", "action"},
		Archetypes: []string{"creator", "magician", "alchemist"},
		Elements: []string{"fire", "air"}, Astrology: "Mercury",
		TraditionalMeaning: "Manifestation, resourcefulness, power, inspired action",
		ShadowAspects: []string{"manipulation", "poor planning", "unused talents"},
		LightAspects: []string{"willpower", "desire", "creation", "manifestation"},
		MoodWeights: map[string]float64{
			"anxious": 0.8, "excited": 1.3, "uncertain": 0.9, "hopeful": 1.2,
			"peaceful": 0.7, "frustrated": 1.1, "curious": 1.0, "contemplative": 0.8,
		},
	},
	2: {
		ID: 2, Name: "The High Priestess", Number: "II",
		Keywords: []string{"intuition", "sacred-knowledge", "divine-feminine", "subconscious"},
		Archetypes: []string{"wise-woman", "oracle", "mystic"},
		Elements: []string{"water"}, Astrology: "Moon",
		TraditionalMeaning: "Intuition, sacred knowledge, divine feminine, the subconscious mind",
		ShadowAspects: []string{"secrets", "withdrawn", "silence", "repressed-feelings"},
		LightAspects: []string{"intuitive", "wise", "serene", "understanding"},
		MoodWeights: map[string]float64{
			"anxious": 1.1, "excited": 0.6, "uncertain": 1.2, "hopeful": 0.9,
			"peaceful": 1.3, "frustrated": 0.8, "curious": 1.1, "contemplative": 1.4,
		},
	},
	3: {
		ID: 3, Name: "The Empress", Number: "III",
		Keywords: []string{"fertility", "femininity", "beauty", "nature", "abundance"},
		Archetypes: []string{"mother", "creator", "nurturer"},
		Elements: []string{"earth"}, Astrology: "Venus",
		TraditionalMeaning: "Fertility, femininity, beauty, nature, abundance",
		ShadowAspects: []string{"creative-block", "dependence", "smothering", "lack"},
		LightAspects: []string{"motherhood", "fertility", "sensuality", "creativity"},
		MoodWeights: map[string]float64{
			"anxious": 0.7, "excited": 1.1, "uncertain": 0.8, "hopeful": 1.2,
			"peaceful": 1.3, "frustrated": 0.6, "curious": 0.9, "contemplative": 1.0,
		},
	},
	4: {
		ID: 4, Name: "The Emperor", Number: "IV",
		Keywords: []string{"authority", "father-figure", "structure", "control", "leadership"},
		Archetypes: []string{"ruler", "father", "leader"},
		Elements: []string{"fire"}, Astrology: "Aries",
		TraditionalMeaning: "Authority, father-figure, structure, control",
		ShadowAspects: []string{"domination", "excessive-control", "rigidity", "lack-of-compassion"},
		LightAspects: []string{"leadership", "logic", "stability", "security"},
		MoodWeights: map[string]float64{
			"anxious": 1.0, "excited": 0.8, "uncertain": 1.1, "hopeful": 1.0,
			"peaceful": 0.7, "frustrated": 1.2, "curious": 0.8, "contemplative": 0.9,
		},
	},
	5: {
		ID: 5, Name: "The Hierophant", Number: "V",
		Keywords: []string{"spiritual-wisdom", "religious-beliefs", "conformity", "tradition"},
		Archetypes: []string{"teacher", "guide", "traditionalist"},
		Elements: []string{"earth"}, Astrology: "Taurus",
		TraditionalMeaning: "Spiritual wisdom, religious beliefs, conformity, tradition, institutions",
		ShadowAspects: []string{"restriction", "challenging-the-status-quo", "personal-beliefs"},
		LightAspects: []string{"education", "knowledge", "beliefs", "conformity"},
		MoodWeights: map[string]float64{
			"anxious": 1.0, "excited": 0.7, "uncertain": 1.1, "hopeful": 0.9,
			"peaceful": 1.2, "frustrated": 0.8, "curious": 1.0, "contemplative": 1.3,
		},
	},
	6: {
		ID: 6, Name: "The Lovers", Number: "VI",
		Keywords: []string{"love", "harmony", "relationships", "values-alignment", "choices"},
		Archetypes: []string{"lover", "partner", "chooser"},
		Elements: []string{"air"}, Astrology: "Gemini",
		TraditionalMeaning: "Love, harmony, relationships, values alignment",
		ShadowAspects: []string{"disharmony", "imbalance", "misalignment-of-values", "indecision"},
		LightAspects: []string{"love", "unity", "relationships", "partnerships"},
		MoodWeights: map[string]float64{
			"anxious": 0.9, "excited": 1.2, "uncertain": 1.3, "hopeful": 1.2,
			"peaceful": 1.1, "frustrated": 0.7, "curious": 1.0, "contemplative": 1.0,
		},
	},
	7: {
		ID: 7, Name: "The Chariot", Number: "VII",
		Keywords: []string{"control", "willpower", "success", "determination", "direction"},
		Archetypes: []string{"warrior", "victor", "driver"},
		Elements: []string{"water"}, Astrology: "Cancer",
		TraditionalMeaning: "Control, willpower, success, determination, direction",
		ShadowAspects: []string{"lack-of-control", "lack-of-direction", "aggression"},
		LightAspects: []string{"control", "willpower", "victory", "assertion"},
		MoodWeights: map[string]float64{
			"anxious": 0.8, "excited": 1.1, "uncertain": 0.9, "hopeful": 1.2,
			"peaceful": 0.6, "frustrated": 1.3, "curious": 0.9, "contemplative": 0.7,
		},
	},
	8: {
		ID: 8, Name: "Strength", Number: "VIII",
		Keywords: []string{"strength", "courage", "patience", "control", "compassion"},
		Archetypes: []string{"healer", "saint", "tamer"},
		Elements: []string{"fire"}, Astrology: "Leo",
		TraditionalMeaning: "Strength, courage, patience, control, compassion",
		ShadowAspects: []string{"self-doubt", "lack-of-confidence", "inadequacy"},
		LightAspects: []string{"strength", "courage", "patience", "control"},
		MoodWeights: map[string]float64{
			"anxious": 1.2, "excited": 1.0, "uncertain": 1.1, "hopeful": 1.1,
			"peaceful": 1.2, "frustrated": 1.3, "curious": 0.9, "contemplative": 1.0,
		},
	},
	9: {
		ID: 9, Name: "The Hermit", Number: "IX",
		Keywords: []string{"soul-searching", "seeking-inner-guidance", "looking-inward"},
		Archetypes: []string{"sage", "seeker", "guide"},
		Elements: []string{"earth"}, Astrology: "Virgo",
		TraditionalMeaning: "Soul searching, seeking inner guidance, looking inward",
		ShadowAspects: []string{"isolation", "loneliness", "withdrawal", "paranoia"},
		LightAspects: []string{"self-reflection", "introspection", "guidance", "solitude"},
		MoodWeights: map[string]float64{
			"anxious": 1.1, "excited": 0.5, "uncertain": 1.3, "hopeful": 0.8,
			"peaceful": 1.2, "frustrated": 1.0, "curious": 1.2, "contemplative": 1.4,
		},
	},
	10: {
		ID: 10, Name: "Wheel of Fortune", Number: "X",
		Keywords: []string{"change", "cycles", "fate", "turning-point", "luck"},
		Archetypes: []string{"gambler", "opportunist", "fatalist"},
		Elements: []string{"fire"}, Astrology: "Jupiter",
		TraditionalMeaning: "Change, cycles, fate, turning point, good luck",
		ShadowAspects: []string{"lack-of-control", "clinging-to-the-past", "bad-luck"},
		LightAspects: []string{"good-luck", "karma", "life-cycles", "destiny"},
		MoodWeights: map[string]float64{
			"anxious": 1.0, "excited": 1.2, "uncertain": 1.3, "hopeful": 1.2,
			"peaceful": 0.8, "frustrated": 1.1, "curious": 1.1, "contemplative": 1.0,
		},
	},
	11: {
		ID: 11, Name: "Justice", Number: "XI",
		Keywords: []string{"justice", "fairness", "truth", "cause-and-effect", "law"},
		Archetypes: []string{"judge", "arbiter", "seeker-of-truth"},
		Elements: []string{"air"}, Astrology: "Libra",
		TraditionalMeaning: "Justice, fairness, truth, cause and effect, law",
		ShadowAspects: []string{"unfairness", "lack-of-accountability", "dishonesty"},
		LightAspects: []string{"justice", "truth", "fairness", "integrity"},
		MoodWeights: map[string]float64{
			"anxious": 1.0, "excited": 0.8, "uncertain": 1.1, "hopeful": 1.0,
			"peaceful": 1.1, "frustrated": 1.2, "curious": 1.0, "contemplative": 1.2,
		},
	},
	12: {
		ID: 12, Name: "The Hanged Man", Number: "XII",
		Keywords: []string{"suspension", "restriction", "letting-go", "sacrifice"},
		Archetypes: []string{"martyr", "sacrificer", "suspended-one"},
		Elements: []string{"water"}, Astrology: "Neptune",
		TraditionalMeaning: "Suspension, restriction, letting go, sacrifice",
		ShadowAspects: []string{"delays", "resistance", "stalling", "needless-sacrifice"},
		LightAspects: []string{"letting-go", "surrendering", "new-perspective", "sacrifice"},
		MoodWeights: map[string]float64{
			"anxious": 1.2, "excited": 0.4, "uncertain": 1.3, "hopeful": 0.7,
			"peaceful": 1.1, "frustrated": 1.3, "curious": 1.1, "contemplative": 1.4,
		},
	},
	13: {
		ID: 13, Name: "Death", Number: "XIII",
		Keywords: []string{"endings", "beginnings", "change", "transformation", "transition"},
		Archetypes: []string{"transformer", "ender", "renewer"},
		Elements: []string{"water"}, Astrology: "Scorpio",
		TraditionalMeaning: "Endings, beginnings, change, transformation, transition",
		ShadowAspects: []string{"resistance-to-change", "repeating-negative-patterns"},
		LightAspects: []string{"transformation", "renewal", "metamorphosis", "release"},
		MoodWeights: map[string]float64{
			"anxious": 1.3, "excited": 0.6, "uncertain": 1.2, "hopeful": 0.8,
			"peaceful": 0.7, "frustrated": 1.1, "curious": 1.0, "contemplative": 1.3,
		},
	},
	14: {
		ID: 14, Name: "Temperance", Number: "XIV",
		Keywords: []string{"balance", "moderation", "patience", "purpose", "meaning"},
		Archetypes: []string{"alchemist", "angel", "mixer"},
		Elements: []string{"fire"}, Astrology: "Sagittarius",
		TraditionalMeaning: "Balance, moderation, patience, purpose",
		ShadowAspects: []string{"imbalance", "excess", "self-indulgence", "clashing"},
		LightAspects: []string{"balance", "moderation", "patience", "purpose"},
		MoodWeights: map[string]float64{
			"anxious": 1.1, "excited": 0.8, "uncertain": 1.0, "hopeful": 1.1,
			"peaceful": 1.3, "frustrated": 1.2, "curious": 1.0, "contemplative": 1.2,
		},
	},
	15: {
		ID: 15, Name: "The Devil", Number: "XV",
		Keywords: []string{"bondage", "addiction", "sexuality", "materialism", "playfulness"},
		Archetypes: []string{"shadow", "tempter", "bound-one"},
		Elements: []string{"earth"}, Astrology: "Capricorn",
		TraditionalMeaning: "Bondage, addiction, sexuality, materialism, playfulness",
		ShadowAspects: []string{"addiction", "materialism", "playfulness", "powerlessness"},
		LightAspects: []string{"humor", "sexuality", "passion", "commitment"},
		MoodWeights: map[string]float64{
			"anxious": 1.2, "excited": 1.1, "uncertain": 1.0, "hopeful": 0.6,
			"peaceful": 0.5, "frustrated": 1.3, "curious": 1.2, "contemplative": 1.0,
		},
	},
	16: {
		ID: 16, Name: "The Tower", Number: "XVI",
		Keywords: []string{"sudden-change", "upheaval", "chaos", "revelation", "awakening"},
		Archetypes: []string{"destroyer", "awakener", "revolutionary"},
		Elements: []string{"fire"}, Astrology: "Mars",
		TraditionalMeaning: "Sudden change, upheaval, chaos, revelation, awakening",
		ShadowAspects: []string{"disaster", "upheaval", "trauma", "sudden-change"},
		LightAspects: []string{"revelation", "awakening", "breakthrough", "disaster"},
		MoodWeights: map[string]float64{
			"anxious": 1.4, "excited": 0.8, "uncertain": 1.3, "hopeful": 0.5,
			"peaceful": 0.3, "frustrated": 1.2, "curious": 1.1, "contemplative": 1.0,
		},
	},
	17: {
		ID: 17, Name: "The Star", Number: "XVII",
		Keywords: []string{"hope", "faith", "purpose", "renewal", "spirituality"},
		Archetypes: []string{"star", "wisher", "hope-bringer"},
		Elements: []string{"air"}, Astrology: "Aquarius",
		TraditionalMeaning: "Hope, faith, purpose, renewal, spirituality",
		ShadowAspects: []string{"lack-of-faith", "despair", "self-trust", "disconnection"},
		LightAspects: []string{"hope", "faith", "purpose", "renewal"},
		MoodWeights: map[string]float64{
			"anxious": 0.8, "excited": 1.1, "uncertain": 0.9, "hopeful": 1.4,
			"peaceful": 1.3, "frustrated": 0.7, "curious": 1.0, "contemplative": 1.2,
		},
	},
	18: {
		ID: 18, Name: "The Moon", Number: "XVIII",
		Keywords: []string{"illusion", "fear", "anxiety", "subconscious", "intuition"},
		Archetypes: []string{"dreamer", "intuitive", "shadow-walker"},
		Elements: []string{"water"}, Astrology: "Pisces",
		TraditionalMeaning: "Illusion, fear, anxiety, subconscious, intuition",
		ShadowAspects: []string{"fear", "anxiety", "confusion", "illusion"},
		LightAspects: []string{"intuition", "dreams", "subconscious", "mystery"},
		MoodWeights: map[string]float64{
			"anxious": 1.4, "excited": 0.6, "uncertain": 1.3, "hopeful": 0.7,
			"peaceful": 0.8, "frustrated": 1.1, "curious": 1.2, "contemplative": 1.3,
		},
	},
	19: {
		ID: 19, Name: "The Sun", Number: "XIX",
		Keywords: []string{"joy", "success", "celebration", "positivity", "vitality"},
		Archetypes: []string{"child", "celebrant", "optimist"},
		Elements: []string{"fire"}, Astrology: "Sun",
		TraditionalMeaning: "Joy, success, celebration, positivity, vitality",
		ShadowAspects: []string{"inner-child", "feeling-down", "lack-of-enthusiasm"},
		LightAspects: []string{"joy", "success", "vitality", "enlightenment"},
		MoodWeights: map[string]float64{
			"anxious": 0.6, "excited": 1.4, "uncertain": 0.7, "hopeful": 1.3,
			"peaceful": 1.2, "frustrated": 0.5, "curious": 1.1, "contemplative": 0.8,
		},
	},
	20: {
		ID: 20, Name: "Judgement", Number: "XX",
		Keywords: []string{"judgement", "rebirth", "inner-calling", "forgiveness"},
		Archetypes: []string{"judge", "awakener", "caller"},
		Elements: []string{"fire"}, Astrology: "Pluto",
		TraditionalMeaning: "Judgement, rebirth, inner calling, forgiveness",
		ShadowAspects: []string{"harsh-judgement", "self-doubt", "lack-of-self-awareness"},
		LightAspects: []string{"judgement", "rebirth", "inner-calling", "forgiveness"},
		MoodWeights: map[string]float64{
			"anxious": 1.0, "excited": 1.1, "uncertain": 1.2, "hopeful": 1.1,
			"peaceful": 1.0, "frustrated": 1.0, "curious": 1.1, "contemplative": 1.3,
		},
	},
	21: {
		ID: 21, Name: "The World", Number: "XXI",
		Keywords: []string{"completion", "accomplishment", "travel", "success", "fulfillment"},
		Archetypes: []string{"achiever", "completion", "wholeness"},
		Elements: []string{"earth"}, Astrology: "Saturn",
		TraditionalMeaning: "Completion, accomplishment, travel, success, fulfillment",
		ShadowAspects: []string{"incomplete", "no-closure", "stagnation", "failed-goals"},
		LightAspects: []string{"completion", "accomplishment", "success", "fulfillment"},
		MoodWeights: map[string]float64{
			"anxious": 0.7, "excited": 1.2, "uncertain": 0.8, "hopeful": 1.2,
			"peaceful": 1.2, "frustrated": 0.6, "curious": 1.0, "contemplative": 1.1,
		},
	},
}

func SelectIntelligentCard(userID uuid.UUID, db *sql.DB, mood string, question string) int {
	rand.Seed(time.Now().UnixNano())
	
	// Get user's recent cards to avoid repeats
	var recentCards []int
	if db != nil {
		recentCards = getRecentCards(userID, db, 5)
	}
	
	var bestCardID int
	var bestScore float64 = 0
	
	for cardID, card := range MajorArcana {
		// Skip recently drawn cards
		if contains(recentCards, cardID) {
			continue
		}
		
		score := calculateCardScore(card, mood, question)
		
		// Add some randomness
		randomFactor := 0.8 + rand.Float64()*0.4 // 0.8 to 1.2
		score *= randomFactor
		
		if score > bestScore {
			bestScore = score
			bestCardID = cardID
		}
	}
	
	// Fallback to random if no good match
	if bestScore == 0 {
		return rand.Intn(22)
	}
	
	return bestCardID
}

func calculateCardScore(card Card, mood string, question string) float64 {
	score := 1.0
	
	// Apply mood weights
	if mood != "" {
		if weight, exists := card.MoodWeights[strings.ToLower(mood)]; exists {
			score *= weight
		}
	}
	
	// Apply question keyword matching
	if question != "" {
		questionLower := strings.ToLower(question)
		
		// Check keywords
		for _, keyword := range card.Keywords {
			if strings.Contains(questionLower, strings.ToLower(keyword)) {
				score *= 1.2
			}
		}
		
		// Check traditional meaning
		meaningWords := strings.Fields(strings.ToLower(card.TraditionalMeaning))
		for _, word := range meaningWords {
			if len(word) > 3 && strings.Contains(questionLower, word) {
				score *= 1.1
			}
		}
	}
	
	return score
}

func getRecentCards(userID uuid.UUID, db *sql.DB, limit int) []int {
	if db == nil {
		return []int{}
	}
	
	rows, err := db.Query(`
		SELECT card_id FROM card_draws 
		WHERE user_id = $1 
		ORDER BY created_at DESC 
		LIMIT $2
	`, userID, limit)
	
	if err != nil {
		return []int{}
	}
	defer rows.Close()
	
	var recentCards []int
	for rows.Next() {
		var cardID int
		if err := rows.Scan(&cardID); err == nil {
			recentCards = append(recentCards, cardID)
		}
	}
	
	return recentCards
}

func contains(slice []int, item int) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}