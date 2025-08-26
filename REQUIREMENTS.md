# Symbol Quest - Detailed MVP Requirements

## Market Reality Check

**Competitor Analysis**: Existing apps like Labyrinthos ($5/month), Golden Thread Tarot ($3/month), and Biddy Tarot (freemium) have 100K-500K users but suffer from:
- Generic AI responses lacking personal depth
- Poor card imagery and UX design
- No symbolic reasoning beyond basic meanings
- Limited progression/gamification elements

**Our Differentiation**: Leverage your ontology server's symbolic reasoning + DevLift's progression mechanics + Freeflow's engagement patterns.

**Week 1 Revenue Target**: $0 (focus on core experience)
**Week 3 Revenue Target**: $500 MRR from 100 beta users at $4.99/month
**6-Month Target**: $15K MRR (3,000 subscribers)

## Requirements

### Core Feature 1: Intelligent Daily Draw (Week 1 Priority)

**User Story**: As someone seeking daily guidance, I want a card that feels personally meaningful, not random.

**Implementation Specifics**:
- **Card Selection Algorithm**: Not random. Use simple user input (mood/question) to weight selection
- **Major Arcana Only**: 22 cards, easier to source quality imagery
- **Card Database Structure**:
  ```json
  {
    "id": 0,
    "name": "The Fool",
    "keywords": ["new-beginnings", "innocence", "spontaneity"],
    "archetypes": ["innocent", "seeker", "beginner"],
    "elements": ["air"],
    "imagery_prompts": "A figure stepping off a cliff with a small bag and white rose",
    "shadow_aspects": ["recklessness", "naivety"],
    "light_aspects": ["faith", "optimism"]
  }
  ```

**MVP Constraints**:
- One card per day per user (enforced by date, not complex logic)
- Mobile-first design (80% of tarot app usage is mobile)
- Loading state max 3 seconds
- Offline capability for previously drawn cards

**Success Metrics**:
- Daily return rate >40% (industry standard: 20-30%)
- Average session time >2 minutes
- Card draw completion rate >90%

### Core Feature 2: Symbolic Reasoning Engine (Week 1-2)

**User Story**: As someone who finds generic tarot interpretations shallow, I want AI that understands symbolic depth and personal context.

**Technical Implementation**:
```javascript
// Prompt Engineering Strategy
const generateInterpretation = (card, userContext) => {
  const prompt = `
  You are interpreting ${card.name} for someone who shared: "${userContext.mood}" and asked: "${userContext.question}"
  
  Card Symbolic Elements:
  - Core Archetype: ${card.archetypes[0]}
  - Element: ${card.elements[0]}
  - Shadow: ${card.shadow_aspects.join(', ')}
  - Light: ${card.light_aspects.join(', ')}
  
  Generate a 3-part interpretation:
  1. SYMBOLIC INSIGHT (2-3 sentences connecting archetype to their situation)
  2. REFLECTION QUESTION (one powerful question for journaling)
  3. TODAY'S GUIDANCE (one specific action they can take)
  
  Tone: Wise but conversational. Avoid fortune-telling. Focus on psychological insight.
  `;
  return openai.complete(prompt);
};
```

**Prompt Engineering Priorities**:
1. **Personal Relevance**: Use mood + question context
2. **Depth over Breadth**: Focus on 2-3 symbolic elements max
3. **Actionable**: Always include a specific "what to do today"
4. **Psychological Frame**: Self-reflection, not prediction

**Fallback Strategy**: Pre-written interpretations for each card (create during week 1)

**Quality Control**:
- Length: 150-250 words (mobile-readable)
- Readability: 8th grade level
- Sentiment: Balanced (not all positive, not all challenging)
- Personalization score: Include user's mood/question in 80%+ of responses

### Revenue Engine: Freemium to Premium Conversion (Week 3)

**Monetization Strategy Based on Competitor Analysis**:

**Free Tier** (Hook them with value):
- 3 AI interpretations per week
- Basic card meanings
- No history/tracking
- Simple card imagery

**Premium Tier** ($4.99/month - competitive with market):
- Unlimited AI interpretations
- Detailed symbolic analysis
- Full interpretation history
- Mood/question context inclusion
- Premium card artwork
- Export interpretations to journal

**Technical Implementation**:
```javascript
// Usage tracking
const checkUsageLimit = (userId) => {
  const weekStart = moment().startOf('week');
  const usageCount = await getUserUsage(userId, weekStart);
  const subscription = await getUserSubscription(userId);
  
  if (subscription?.status === 'active') return true;
  return usageCount < 3;
};

// Conversion optimization
const showUpgradePrompt = (usageCount) => {
  if (usageCount === 2) {
    return "🌟 One more free interpretation this week! Upgrade for unlimited daily guidance.";
  }
  if (usageCount >= 3) {
    return "✨ You've discovered the power of AI tarot! Get unlimited access for $4.99/month.";
  }
};
```

**Stripe Integration**:
- Single price point: $4.99/month
- 7-day free trial (collect payment method upfront)
- Webhooks for subscription changes
- Graceful downgrade (maintain history, limit new interpretations)

**Conversion Funnel Optimization**:
- Week 1-2: Focus on product-market fit
- Week 3: Add payment flow
- Target 10% free-to-paid conversion (industry standard: 5-15%)

### Retention Mechanic: Personal Journey Tracking (Week 2 Priority)

**User Story**: As someone using tarot for personal growth, I want to see patterns in my readings and track my journey over time.

**Implementation Strategy** (Inspired by DevLift's tracking):

```javascript
// Data structure for user journey
const userJourney = {
  daily_draws: [
    {
      date: '2024-01-15',
      card_id: 0,
      user_context: { mood: 'anxious', question: 'career decision' },
      interpretation: '...',
      user_reflection: null, // Added later if they journal
      tags: ['career', 'anxiety'] // Auto-generated from context
    }
  ],
  patterns: {
    most_drawn_card: 'The Hermit',
    common_themes: ['career', 'relationships', 'growth'],
    streak: 12, // Days in a row
    total_draws: 45
  },
  milestones: [
    { type: 'first_draw', date: '2024-01-01' },
    { type: 'week_streak', date: '2024-01-07' },
    { type: 'month_journey', date: '2024-02-01' }
  ]
};
```

**Gamification Elements** (From Freeflow patterns):
- **Streak tracking**: Daily draw streaks with visual progress
- **Card collection**: "Unlock" all 22 major arcana over time
- **Insight badges**: "Deep Reflector" (added personal notes), "Pattern Seeker" (reviewed history)
- **Journey milestones**: 7 days, 30 days, 100 days of daily practice

**History Interface**:
- Calendar view showing cards drawn each day
- Search by theme/keyword
- Pattern insights: "You often draw The Hermit when asking about career"
- Export option: Download your journey as PDF/text

**Retention Metrics**:
- Day 7 retention: >30%
- Day 30 retention: >15%
- Average cards per user: >20 (shows engagement)

### OS-Level Asset: Reusable Symbolic Engine (Week 2-3)

**Developer Story**: Build a symbolic reasoning foundation that powers Symbol Quest and extends to Dream Journal Pro, Philosophy Chat, and future projects.

**Symbolic Engine Architecture**:
```javascript
// Core symbolic reasoning types
interface Symbol {
  id: string;
  name: string;
  archetypes: Archetype[];
  elements: Element[];
  relationships: SymbolRelationship[];
  cultural_contexts: CulturalContext[];
}

interface Archetype {
  name: string; // "The Innocent", "The Sage", "The Hero"
  core_drive: string;
  fears: string[];
  desires: string[];
  shadow_aspects: string[];
  light_aspects: string[];
}

// Reusable across tarot, dreams, mythology, etc.
const symbolicEngine = {
  interpret: (symbols: Symbol[], context: UserContext) => {
    // Combine archetypal patterns
    // Weight by user context
    // Generate coherent narrative
  },
  findPatterns: (userHistory: Symbol[]) => {
    // Identify recurring archetypes
    // Suggest personal themes
  },
  crossReference: (symbol: Symbol, domain: 'tarot' | 'dreams' | 'mythology') => {
    // Enable cross-domain symbolic connections
  }
};
```

**Future Project Applications**:
- **Dream Journal Pro**: Same archetypal analysis for dream symbols
- **Philosophy Chat**: Reference philosophical archetypes
- **Life XP Dashboard**: Symbolic life themes tracking

**Week 1 Minimum Viable Engine**:
- 22 tarot archetypes with relationships
- Basic pattern matching ("you often encounter transformation themes")
- Simple cross-referencing between cards

**Week 3 Enhanced Engine**:
- User-specific archetype affinity scoring
- Context-aware interpretation weighting
- Export symbolic insights for other projects

### User Experience: Mobile-First Design (Week 1 Foundation)

**Design Principle**: 80% of tarot app usage happens on mobile during morning routines or evening reflection.

**Technical Constraints**:
- **Load Time**: <3 seconds on 3G
- **Battery**: Minimal CPU usage, efficient animations
- **Touch**: Minimum 44px touch targets
- **Accessibility**: Screen reader support, high contrast mode

**Card Drawing Experience**:
```javascript
// Animation sequence for card draw
const cardDrawSequence = {
  1: 'Shuffle animation (2 seconds)',
  2: 'Card selection with haptic feedback',
  3: 'Card flip reveal (1.5 seconds)',
  4: 'Slide to interpretation with gentle bounce'
};

// Progressive loading
const loadingStrategy = {
  immediate: 'Card back, shuffle button',
  after_draw: 'Selected card image',
  background: 'AI interpretation generation',
  cached: 'All previously drawn cards'
};
```

**Mobile-Specific Features**:
- **Haptic feedback** on card selection
- **Swipe gestures** to navigate history
- **Share functionality** to social media
- **Background notifications** for daily draw reminders
- **Offline mode** for re-reading past interpretations

**Performance Targets**:
- First Contentful Paint: <1.5s
- Largest Contentful Paint: <2.5s
- Cumulative Layout Shift: <0.1
- First Input Delay: <100ms

**Browser Support**:
- iOS Safari 14+ (primary target)
- Chrome Android 90+ (secondary)
- Progressive Web App capabilities for "Add to Home Screen"

### Security & Privacy Requirements (Week 1-3 Priority)

**Data Protection**:
```typescript
// User data encryption at rest
interface UserDataSecurity {
  encryption: 'AES-256';
  personal_data: 'encrypted_fields';
  api_keys: 'environment_variables_only';
  sessions: 'httpOnly_secure_cookies';
}

// API Security
const apiSecurityHeaders = {
  'Content-Security-Policy': "default-src 'self'",
  'X-Frame-Options': 'DENY',
  'X-Content-Type-Options': 'nosniff',
  'Strict-Transport-Security': 'max-age=31536000'
};
```

**Privacy Compliance**:
- GDPR compliance for EU users
- Data deletion requests within 30 days
- Clear privacy policy and terms of service
- Minimal data collection (email, payment info only)
- User consent for analytics cookies
- Data retention: 2 years inactive account deletion

**Authentication Security**:
- Password requirements: 8+ chars, mixed case, numbers
- Rate limiting: 5 login attempts per 15 minutes
- JWT tokens: 24-hour expiration with refresh
- Two-factor authentication option (future enhancement)
- Secure password reset flow with time-limited tokens

### Testing & Quality Assurance Strategy

**Testing Framework & Coverage**:
```typescript
// Testing stack requirements
const testingStack = {
  unit_testing: 'Jest + React Testing Library',
  integration_testing: 'Cypress for critical user flows',
  api_testing: 'Supertest for backend endpoints',
  performance_testing: 'Lighthouse CI in pipeline',
  coverage_target: '80% code coverage minimum'
};

// Critical test scenarios
const testPriorities = [
  'Card selection algorithm accuracy',
  'AI interpretation generation',
  'Subscription payment flow',
  'Daily usage limit enforcement',
  'Data persistence and retrieval',
  'Mobile responsive behavior'
];
```

**Automated Testing Requirements**:
- Unit tests: All utility functions, hooks, and core logic
- Integration tests: Complete user journeys (signup → card draw → interpretation)
- API tests: All endpoints with authentication and error scenarios
- E2E tests: Payment flow, subscription management, daily limits
- Performance tests: Load testing for 100 concurrent users

**Manual Testing Checklist**:
- Cross-browser testing (iOS Safari, Chrome Android, desktop)
- Accessibility testing with screen readers
- Payment flow testing in Stripe test mode
- Mobile touch interactions and gestures
- Offline functionality validation

### Error Handling & Monitoring Strategy

**Comprehensive Error Handling**:
```typescript
// Error classification and handling
interface ErrorHandlingStrategy {
  client_errors: {
    network_failures: 'Retry with exponential backoff',
    api_rate_limits: 'Queue requests and inform user',
    invalid_input: 'Inline validation with helpful messages',
    session_expiry: 'Auto-refresh tokens or prompt re-login'
  };
  
  server_errors: {
    ai_api_failures: 'Fallback to pre-written interpretations',
    database_errors: 'Graceful degradation with local storage',
    payment_failures: 'Clear error messages and retry options',
    quota_exceeded: 'Upgrade prompts with clear pricing'
  };
}
```

**Monitoring & Observability**:
- Sentry for error tracking and performance monitoring
- Custom metrics: Daily active users, card draw completion rates
- Business metrics: Conversion rates, churn rates, revenue tracking
- Infrastructure monitoring: API response times, database performance
- User experience metrics: Core Web Vitals, user session recordings

**Alerting System**:
- Critical: Payment failures, API downtime (immediate Slack/email)
- Warning: High error rates, slow performance (15-minute delay)
- Info: Daily metrics, user milestones (daily digest)

### Data Management & Compliance

**Database Design & Backup Strategy**:
```sql
-- Data retention and backup requirements
CREATE TABLE user_data_retention (
  user_id UUID PRIMARY KEY,
  created_at TIMESTAMP,
  last_active TIMESTAMP,
  deletion_requested_at TIMESTAMP,
  scheduled_deletion_at TIMESTAMP
);

-- Backup strategy
BACKUP_SCHEDULE = {
  daily_backups: 'Automated at 2 AM UTC',
  weekly_full_backup: 'Sunday complete database dump',
  retention_policy: '30 days daily, 12 weeks weekly',
  disaster_recovery: 'Cross-region backup replication'
};
```

**Data Governance**:
- Personal data minimization: Only collect essential information
- Data anonymization: Remove PII from analytics and logs
- GDPR Article 17: Right to be forgotten implementation
- Data portability: Export user data in JSON format
- Regular security audits and penetration testing

### Compliance & Legal Requirements

**Accessibility Compliance (WCAG 2.1 AA)**:
- Screen reader compatibility for all interactive elements
- Keyboard navigation support throughout the application
- High contrast mode support (4.5:1 minimum ratio)
- Alternative text for all card images and symbols
- Focus indicators visible and prominent
- Zoom support up to 200% without horizontal scrolling

**Legal & Business Compliance**:
- Terms of Service clearly defining usage rights and limitations
- Privacy Policy detailing data collection and usage
- Cookie policy with user consent management
- Age verification (13+ with parental consent for under-18)
- Content moderation for user-generated interpretations
- Trademark compliance for tarot imagery and symbolism

### Operational & DevOps Requirements

**CI/CD Pipeline Requirements**:
```yaml
# GitHub Actions workflow requirements
deployment_pipeline:
  stages:
    - code_quality: 'ESLint, Prettier, TypeScript compilation'
    - testing: 'Unit, integration, E2E test suites'
    - security: 'Dependency scanning, SAST analysis'
    - build: 'Production build with optimization'
    - deploy: 'Staging deployment with smoke tests'
    - production: 'Blue-green deployment with rollback capability'
  
  environment_management:
    development: 'Feature branch deployments'
    staging: 'Pre-production testing environment'
    production: 'Load-balanced, auto-scaling infrastructure'
```

**Infrastructure & Scalability**:
- Auto-scaling based on user traffic patterns
- CDN for static assets and card images
- Database connection pooling and read replicas
- API rate limiting and DDoS protection
- SSL/TLS certificates with automatic renewal
- Container orchestration with health checks

**Maintenance & Support**:
- Weekly dependency updates and security patches
- Monthly performance optimization reviews
- Quarterly user feedback analysis and feature planning
- Annual security audit and compliance review
- 24/7 uptime monitoring with 99.9% availability SLA
- Customer support system with response time commitments
