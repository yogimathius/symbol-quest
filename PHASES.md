# Symbol Quest - Sprint Planning & Implementation

## Sprint 1: Core Card Experience (Days 1-7)

**Daily Success Metrics**: Card draw completion rate >90%, mobile load time <3s

### Day 1-2: Project Foundation
```bash
# Exact setup commands
npx create-react-app symbol-quest --template typescript
cd symbol-quest
npm install @tailwindcss/forms axios react-router-dom
npm install -D @types/node
```

**File Structure**:
```
src/
├── components/
│   ├── CardDraw.tsx           # Main card drawing interface
│   ├── CardDisplay.tsx        # Individual card presentation
│   └── MobileLayout.tsx       # Responsive wrapper
├── data/
│   ├── major-arcana.json      # 22 cards with symbolic data
│   └── symbolic-engine.ts     # Core reasoning functions
├── hooks/
│   ├── useCardDraw.ts         # Card selection logic
│   └── useUserContext.ts      # Mood/question state
└── utils/
    ├── card-selection.ts      # Weighted selection algorithm
    └── storage.ts             # LocalStorage helpers
```

### Day 3-4: Card Data & Selection
**Priority**: Create rich card database and intelligent selection

```typescript
// major-arcana.json structure
{
  "cards": [
    {
      "id": 0,
      "name": "The Fool",
      "number": "0",
      "keywords": ["new-beginnings", "innocence", "spontaneity", "faith"],
      "archetypes": ["innocent", "seeker", "beginner"],
      "elements": ["air"],
      "astrology": "Uranus",
      "traditional_meaning": "New beginnings, innocence, spontaneity",
      "shadow_aspects": ["recklessness", "naivety", "foolishness"],
      "light_aspects": ["faith", "optimism", "adventure"],
      "imagery_description": "A young person stepping off a cliff with a small bag and white rose",
      "colors": ["yellow", "light-blue", "white"],
      "symbols": ["cliff", "rose", "bag", "sun", "mountains"],
      "mood_weights": {
        "anxious": 0.3,
        "excited": 1.2,
        "uncertain": 1.1,
        "hopeful": 1.3
      }
    }
    // ... 21 more cards
  ]
}
```

**Selection Algorithm**:
```typescript
const selectCard = (userMood: string, userQuestion: string) => {
  const cards = majorArcana.cards;
  
  // Weight cards based on mood
  const weightedCards = cards.map(card => ({
    ...card,
    weight: card.mood_weights[userMood] || 1.0
  }));
  
  // Add question context weighting
  const questionWeighted = weightedCards.map(card => {
    const questionBonus = card.keywords.some(keyword => 
      userQuestion.toLowerCase().includes(keyword)
    ) ? 1.2 : 1.0;
    
    return {
      ...card,
      weight: card.weight * questionBonus
    };
  });
  
  // Weighted random selection
  return weightedRandomSelect(questionWeighted);
};
```

### Day 5-7: UI/UX Implementation
**Priority**: Mobile-first card drawing experience

**Component: CardDraw.tsx**
```typescript
const CardDraw = () => {
  const [userMood, setUserMood] = useState('');
  const [userQuestion, setUserQuestion] = useState('');
  const [selectedCard, setSelectedCard] = useState(null);
  const [isDrawing, setIsDrawing] = useState(false);
  
  const handleDraw = async () => {
    setIsDrawing(true);
    
    // 2-second shuffle animation
    await new Promise(resolve => setTimeout(resolve, 2000));
    
    const card = selectCard(userMood, userQuestion);
    setSelectedCard(card);
    setIsDrawing(false);
    
    // Save to localStorage for daily limit
    const today = new Date().toDateString();
    localStorage.setItem(`daily-card-${today}`, JSON.stringify({
      card,
      mood: userMood,
      question: userQuestion,
      timestamp: Date.now()
    }));
  };
  
  return (
    <div className="min-h-screen bg-gradient-to-b from-purple-900 to-blue-900 p-4">
      {/* Mood selection */}
      {/* Question input */}
      {/* Draw button with animation */}
      {/* Card display */}
    </div>
  );
};
```

**Success Criteria for Week 1**:
- ✅ All 22 cards have complete data
- ✅ Mood/question input affects card selection observably
- ✅ Mobile experience loads <3 seconds
- ✅ Daily card persistence works
- ✅ Smooth animations on card reveal

## Phase 2: AI Integration (Week 2)

**Goal**: AI-powered interpretations and user accounts

### Tasks:

- Integrate OpenAI API for interpretation generation
- Create prompt engineering system for tarot interpretations
- Implement user authentication (email/password)
- Build user dashboard and card history
- Create interpretation display and storage
- Add loading states and error handling

**Deliverable**: Full AI interpretation system with user accounts

## Phase 3: Monetization & Polish (Week 3)

**Goal**: Subscription system and production readiness

### Tasks:

- Integrate Stripe for subscription payments
- Implement free trial and subscription logic
- Create subscription management interface
- Add daily card limitation and reset logic
- Implement symbolic reasoning engine structure
- Polish UI/UX and add animations
- Set up production deployment
- Add analytics and error tracking

**Deliverable**: Production-ready MVP with subscription system

## Phase 4: Post-Launch Enhancements (Future)

**Goal**: User engagement and feature expansion

### Potential Features:

- Multi-card spreads (3-card, Celtic Cross)
- Personalized interpretation based on user history
- Social sharing of daily cards
- Meditation and reflection prompts
- Integration with calendar and reminder systems
- Advanced symbolic pattern recognition

## Technical Architecture

### Frontend:

- React with TypeScript
- Tailwind CSS for styling
- React Router for navigation
- Axios for API calls

### Backend:

- Node.js/Express API
- PostgreSQL database
- JWT authentication
- OpenAI API integration
- Stripe webhook handling

### Infrastructure:

- Vercel/Netlify for frontend
- Railway/Heroku for backend
- Stripe for payments
- Sentry for error tracking
