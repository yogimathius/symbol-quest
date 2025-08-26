import { describe, it, expect, beforeEach, vi } from 'vitest';
import { 
  selectCard, 
  hasDrawnToday, 
  getTodaysCard, 
  saveCardDraw, 
  getCardHistory 
} from '../card-selection';
import type { UserContext, TarotCard } from '../../types/card';

// Mock localStorage
const localStorageMock = (() => {
  let store: Record<string, string> = {};
  
  return {
    getItem: vi.fn((key: string) => store[key] || null),
    setItem: vi.fn((key: string, value: string) => {
      store[key] = value;
    }),
    clear: vi.fn(() => {
      store = {};
    }),
  };
})();

Object.defineProperty(window, 'localStorage', {
  value: localStorageMock
});

// Mock Math.random for predictable tests
const mockMath = Object.create(global.Math);
mockMath.random = vi.fn(() => 0.5); // Always return 0.5 for predictable results
global.Math = mockMath;

describe('card-selection', () => {
  beforeEach(() => {
    localStorageMock.clear();
    vi.clearAllMocks();
  });

  describe('selectCard', () => {
    it('should return a valid tarot card', () => {
      const context: UserContext = {
        mood: 'hopeful',
        question: 'What should I focus on today?',
        timestamp: Date.now()
      };

      const card = selectCard(context);

      expect(card).toBeDefined();
      expect(card.id).toBeTypeOf('number');
      expect(card.name).toBeTypeOf('string');
      expect(card.keywords).toBeInstanceOf(Array);
      expect(card.archetypes).toBeInstanceOf(Array);
      expect(card.moodWeights).toBeTypeOf('object');
    });

    it('should weight cards based on mood', () => {
      const hopefulContext: UserContext = {
        mood: 'hopeful',
        question: 'general guidance',
        timestamp: Date.now()
      };

      const anxiousContext: UserContext = {
        mood: 'anxious',
        question: 'general guidance',
        timestamp: Date.now()
      };

      // Run multiple selections to check distribution
      const hopefulResults: number[] = [];
      const anxiousResults: number[] = [];

      // Mock Math.random to return different values
      let callCount = 0;
      vi.mocked(Math.random).mockImplementation(() => {
        const values = [0.1, 0.3, 0.5, 0.7, 0.9];
        return values[callCount++ % values.length];
      });

      for (let i = 0; i < 10; i++) {
        hopefulResults.push(selectCard(hopefulContext).id);
        anxiousResults.push(selectCard(anxiousContext).id);
      }

      // Results should be different due to mood weighting
      // (Note: This is a statistical test, so it might occasionally fail)
      expect(hopefulResults).not.toEqual(anxiousResults);
    });

    it('should consider question relevance in card selection', () => {
      const loveContext: UserContext = {
        mood: 'hopeful',
        question: 'How will my relationship progress?',
        timestamp: Date.now()
      };

      const careerContext: UserContext = {
        mood: 'hopeful',
        question: 'Should I change careers?',
        timestamp: Date.now()
      };

      // Run selections multiple times
      const loveResults: number[] = [];
      const careerResults: number[] = [];

      let callCount = 0;
      vi.mocked(Math.random).mockImplementation(() => {
        const values = [0.2, 0.4, 0.6, 0.8];
        return values[callCount++ % values.length];
      });

      for (let i = 0; i < 8; i++) {
        loveResults.push(selectCard(loveContext).id);
        careerResults.push(selectCard(careerContext).id);
      }

      // Different questions should tend to select different cards
      expect(loveResults).not.toEqual(careerResults);
    });

    it('should handle empty or invalid inputs gracefully', () => {
      const invalidContext: UserContext = {
        mood: 'invalid-mood' as any,
        question: '',
        timestamp: Date.now()
      };

      const card = selectCard(invalidContext);
      expect(card).toBeDefined();
      expect(card.id).toBeTypeOf('number');
    });
  });

  describe('hasDrawnToday', () => {
    it('should return false when no card has been drawn', () => {
      expect(hasDrawnToday()).toBe(false);
    });

    it('should return true when a card was drawn today', () => {
      const today = new Date().toDateString();
      const drawData = {
        card: { id: 0, name: 'The Fool' } as TarotCard,
        date: today,
        timestamp: Date.now()
      };

      localStorageMock.setItem('lastCardDraw', JSON.stringify(drawData));
      expect(hasDrawnToday()).toBe(true);
    });

    it('should return false when a card was drawn yesterday', () => {
      const yesterday = new Date(Date.now() - 24 * 60 * 60 * 1000).toDateString();
      const drawData = {
        card: { id: 0, name: 'The Fool' } as TarotCard,
        date: yesterday,
        timestamp: Date.now() - 24 * 60 * 60 * 1000
      };

      localStorageMock.setItem('lastCardDraw', JSON.stringify(drawData));
      expect(hasDrawnToday()).toBe(false);
    });

    it('should handle corrupted localStorage data', () => {
      localStorageMock.setItem('lastCardDraw', 'invalid json');
      expect(hasDrawnToday()).toBe(false);
    });
  });

  describe('getTodaysCard', () => {
    it('should return null when no card drawn today', () => {
      expect(getTodaysCard()).toBeNull();
    });

    it('should return the card drawn today', () => {
      const today = new Date().toDateString();
      const mockCard = { id: 5, name: 'The Hierophant' } as TarotCard;
      const drawData = {
        card: mockCard,
        date: today,
        timestamp: Date.now()
      };

      localStorageMock.setItem('lastCardDraw', JSON.stringify(drawData));
      
      const todaysCard = getTodaysCard();
      expect(todaysCard).toEqual(mockCard);
    });

    it('should return null for yesterday\\'s card', () => {
      const yesterday = new Date(Date.now() - 24 * 60 * 60 * 1000).toDateString();
      const drawData = {
        card: { id: 0, name: 'The Fool' } as TarotCard,
        date: yesterday,
        timestamp: Date.now() - 24 * 60 * 60 * 1000
      };

      localStorageMock.setItem('lastCardDraw', JSON.stringify(drawData));
      expect(getTodaysCard()).toBeNull();
    });
  });

  describe('saveCardDraw', () => {
    it('should save card draw to localStorage', () => {
      const mockCard = { id: 10, name: 'Wheel of Fortune' } as TarotCard;
      const context: UserContext = {
        mood: 'excited',
        question: 'What changes are coming?',
        timestamp: Date.now()
      };
      const interpretation = 'A time of positive change approaches.';

      saveCardDraw(mockCard, context, interpretation);

      expect(localStorageMock.setItem).toHaveBeenCalledWith(
        'lastCardDraw',
        expect.stringContaining('Wheel of Fortune')
      );
      expect(localStorageMock.setItem).toHaveBeenCalledWith(
        'cardHistory',
        expect.stringContaining('Wheel of Fortune')
      );
    });

    it('should save card draw without interpretation', () => {
      const mockCard = { id: 7, name: 'The Chariot' } as TarotCard;
      const context: UserContext = {
        mood: 'determined',
        question: 'How can I achieve my goals?',
        timestamp: Date.now()
      };

      saveCardDraw(mockCard, context);

      const savedData = localStorageMock.setItem.mock.calls
        .find(([key]) => key === 'lastCardDraw')?.[1];
      
      expect(savedData).toBeDefined();
      const parsedData = JSON.parse(savedData!);
      expect(parsedData.card.name).toBe('The Chariot');
      expect(parsedData.context.mood).toBe('determined');
      expect(parsedData.interpretation).toBeUndefined();
    });

    it('should maintain history with proper date formatting', () => {
      const mockCard = { id: 17, name: 'The Star' } as TarotCard;
      const context: UserContext = {
        mood: 'hopeful',
        question: 'What should I hope for?',
        timestamp: Date.now()
      };

      saveCardDraw(mockCard, context);

      const savedHistory = localStorageMock.setItem.mock.calls
        .find(([key]) => key === 'cardHistory')?.[1];
      
      expect(savedHistory).toBeDefined();
      const history = JSON.parse(savedHistory!);
      expect(history).toHaveLength(1);
      expect(history[0].date).toBe(new Date().toDateString());
      expect(history[0].id).toMatch(/^\d+-17$/);
    });
  });

  describe('getCardHistory', () => {
    it('should return empty array when no history exists', () => {
      const history = getCardHistory();
      expect(history).toEqual([]);
    });

    it('should return parsed history from localStorage', () => {
      const mockHistory = [{
        card: { id: 21, name: 'The World' } as TarotCard,
        context: { mood: 'accomplished', question: 'Am I on the right path?', timestamp: Date.now() },
        date: new Date().toDateString(),
        timestamp: Date.now(),
        id: `${Date.now()}-21`
      }];

      localStorageMock.setItem('cardHistory', JSON.stringify(mockHistory));

      const history = getCardHistory();
      expect(history).toEqual(mockHistory);
    });

    it('should handle corrupted history data', () => {
      localStorageMock.setItem('cardHistory', 'invalid json');
      
      const history = getCardHistory();
      expect(history).toEqual([]);
    });

    it('should limit history to 30 entries', () => {
      const mockCard = { id: 1, name: 'The Magician' } as TarotCard;
      const context: UserContext = {
        mood: 'focused',
        question: 'How can I manifest my goals?',
        timestamp: Date.now()
      };

      // Add 35 entries
      for (let i = 0; i < 35; i++) {
        saveCardDraw(mockCard, { ...context, timestamp: Date.now() + i });
      }

      const history = getCardHistory();
      expect(history.length).toBeLessThanOrEqual(30);
    });
  });
});