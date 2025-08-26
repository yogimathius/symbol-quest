import type { TarotCard, UserContext, Mood } from '../types/card';
import majorArcanaData from '../data/major-arcana.json';

/**
 * Weighted random selection from an array of items with weights
 */
function weightedRandomSelect<T extends { weight: number }>(items: T[]): T {
  const totalWeight = items.reduce((sum, item) => sum + item.weight, 0);
  let random = Math.random() * totalWeight;
  
  for (const item of items) {
    random -= item.weight;
    if (random <= 0) {
      return item;
    }
  }
  
  // Fallback to last item (shouldn't happen with proper weights)
  return items[items.length - 1];
}

/**
 * Calculate semantic similarity between user question and card keywords
 * Simple implementation using word overlap and synonyms
 */
function calculateQuestionRelevance(question: string, keywords: string[]): number {
  const questionWords = question.toLowerCase()
    .replace(/[^a-z\s]/g, '') // Remove punctuation
    .split(/\s+/)
    .filter(word => word.length > 2); // Filter short words
  
  if (questionWords.length === 0) return 1.0;
  
  // Check for direct keyword matches
  const directMatches = keywords.filter(keyword => 
    questionWords.some(word => 
      word.includes(keyword.toLowerCase()) || keyword.toLowerCase().includes(word)
    )
  ).length;
  
  // Semantic keyword mapping for common tarot themes
  const semanticKeywords: Record<string, string[]> = {
    'love': ['relationship', 'romance', 'partner', 'dating', 'heart', 'marriage'],
    'career': ['work', 'job', 'profession', 'business', 'employment', 'money'],
    'change': ['transition', 'transformation', 'new', 'different', 'shift'],
    'growth': ['development', 'progress', 'improvement', 'learning', 'evolve'],
    'decision': ['choice', 'choose', 'decide', 'option', 'path', 'direction'],
    'spirituality': ['spiritual', 'soul', 'purpose', 'meaning', 'faith', 'divine'],
    'creativity': ['creative', 'art', 'imagination', 'inspiration', 'express'],
    'health': ['wellness', 'healing', 'body', 'mind', 'energy', 'balance']
  };
  
  // Check for semantic matches
  let semanticMatches = 0;
  for (const [theme, synonyms] of Object.entries(semanticKeywords)) {
    const hasThemeKeyword = keywords.some(k => k.toLowerCase().includes(theme));
    const hasQuestionSynonym = questionWords.some(word => 
      synonyms.some(synonym => word.includes(synonym) || synonym.includes(word))
    );
    
    if (hasThemeKeyword && hasQuestionSynonym) {
      semanticMatches += 0.5; // Weight semantic matches less than direct matches
    }
  }
  
  const totalMatches = directMatches + semanticMatches;
  const maxPossibleMatches = Math.min(keywords.length, questionWords.length);
  
  // Return relevance score between 1.0 (no match) and 2.0 (perfect match)
  return 1.0 + Math.min(totalMatches / Math.max(maxPossibleMatches, 1), 1.0);
}

/**
 * Intelligent card selection based on user mood and question context
 */
export function selectCard(userContext: UserContext): TarotCard {
  const { mood, question } = userContext;
  const cards = majorArcanaData.cards as TarotCard[];
  
  // Apply mood-based weighting
  const moodWeightedCards = cards.map(card => ({
    ...card,
    weight: card.moodWeights[mood as Mood] || 1.0
  }));
  
  // Apply question-based weighting
  const questionWeightedCards = moodWeightedCards.map(card => {
    const questionBonus = calculateQuestionRelevance(question, card.keywords);
    
    return {
      ...card,
      weight: card.weight * questionBonus
    };
  });
  
  // Add subtle randomization to prevent always selecting the highest weighted card
  const finalWeightedCards = questionWeightedCards.map(card => ({
    ...card,
    weight: card.weight * (0.7 + Math.random() * 0.6) // Randomize weight by ±30%
  }));
  
  return weightedRandomSelect(finalWeightedCards);
}

/**
 * Check if user has already drawn a card today
 */
export function hasDrawnToday(): boolean {
  const today = new Date().toDateString();
  const lastDraw = localStorage.getItem('lastCardDraw');
  
  if (!lastDraw) return false;
  
  try {
    const drawData = JSON.parse(lastDraw);
    return drawData.date === today;
  } catch {
    return false;
  }
}

/**
 * Get today's card if already drawn
 */
export function getTodaysCard(): TarotCard | null {
  if (!hasDrawnToday()) return null;
  
  const lastDraw = localStorage.getItem('lastCardDraw');
  if (!lastDraw) return null;
  
  try {
    const drawData = JSON.parse(lastDraw);
    return drawData.card;
  } catch {
    return null;
  }
}

/**
 * Save card draw to local storage
 */
export function saveCardDraw(card: TarotCard, context: UserContext, interpretation?: string): void {
  const drawData = {
    card,
    context,
    interpretation,
    date: new Date().toDateString(),
    timestamp: Date.now(),
    id: `${Date.now()}-${card.id}`
  };
  
  localStorage.setItem('lastCardDraw', JSON.stringify(drawData));
  
  // Also save to history
  const history = getCardHistory();
  history.unshift(drawData);
  
  // Keep only last 30 draws
  const limitedHistory = history.slice(0, 30);
  localStorage.setItem('cardHistory', JSON.stringify(limitedHistory));
}

/**
 * Get user's card draw history
 */
export function getCardHistory(): Array<{
  card: TarotCard;
  context: UserContext;
  interpretation?: string;
  date: string;
  timestamp: number;
  id: string;
}> {
  const history = localStorage.getItem('cardHistory');
  if (!history) return [];
  
  try {
    return JSON.parse(history);
  } catch {
    return [];
  }
}