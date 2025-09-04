import { useState, useCallback, useEffect } from 'react';
import { selectCard, saveCardDraw, hasDrawnToday, getTodaysCard } from '../utils/card-selection';
import { apiService, APIError } from '../services/api';
import type { TarotCard, UserContext } from '../types/card';

export interface UseCardDrawReturn {
  selectedCard: TarotCard | null;
  isDrawing: boolean;
  hasDrawnToday: boolean;
  drawCard: (context: UserContext) => Promise<void>;
  resetDraw: () => void;
  todaysCard: TarotCard | null;
  error: string | null;
}

export function useCardDraw(): UseCardDrawReturn {
  const [selectedCard, setSelectedCard] = useState<TarotCard | null>(null);
  const [isDrawing, setIsDrawing] = useState(false);
  const [drawnToday, setDrawnToday] = useState(false);
  const [todaysCard, setTodaysCard] = useState<TarotCard | null>(null);
  const [error, setError] = useState<string | null>(null);

  // Check if user has already drawn today on mount
  useEffect(() => {
    const checkTodaysDraw = async () => {
      if (apiService.isAuthenticated()) {
        try {
          const status = await apiService.getTodayStatus();
          setDrawnToday(status.has_drawn);
          
          if (status.has_drawn && status.card) {
            const card: TarotCard = {
              id: status.card.id,
              name: status.card.name,
              number: status.card.number || '0',
              // Map the API response to TarotCard format
              keywords: [],
              archetypes: [],
              elements: [],
              astrology: '',
              traditionalMeaning: status.card.traditional_meaning || '',
              shadowAspects: [],
              lightAspects: [],
              imageryDescription: '',
              colors: [],
              symbols: [],
              moodWeights: {},
            };
            setTodaysCard(card);
            setSelectedCard(card);
          }
        } catch (error) {
          console.error('Failed to check today status:', error);
          // Fallback to local storage
          const hasDrawn = hasDrawnToday();
          setDrawnToday(hasDrawn);
          
          if (hasDrawn) {
            const card = getTodaysCard();
            setTodaysCard(card);
            setSelectedCard(card);
          }
        }
      } else {
        // Use local storage when not authenticated
        const hasDrawn = hasDrawnToday();
        setDrawnToday(hasDrawn);
        
        if (hasDrawn) {
          const card = getTodaysCard();
          setTodaysCard(card);
          setSelectedCard(card);
        }
      }
    };

    checkTodaysDraw();
  }, []);

  const drawCard = useCallback(async (context: UserContext) => {
    try {
      setError(null);
      setIsDrawing(true);

      // Validate context
      if (!context.mood || !context.question.trim()) {
        throw new Error('Please provide both your current mood and a question.');
      }

      // Check if already drawn today
      if (apiService.isAuthenticated()) {
        try {
          const status = await apiService.getTodayStatus();
          if (status.has_drawn) {
            setError('You have already drawn your card for today');
            return;
          }
        } catch (error) {
          console.error('Failed to check today status:', error);
        }
      } else if (hasDrawnToday()) {
        const existing = getTodaysCard();
        if (existing) {
          setSelectedCard(existing);
          setDrawnToday(true);
          return;
        }
      }

      // Add delay for card drawing animation
      await new Promise(resolve => setTimeout(resolve, 2000));

      if (apiService.isAuthenticated()) {
        // Use API for authenticated users
        try {
          const response = await apiService.performDailyDraw(context.mood, context.question);
          
          if (response.success && response.card) {
            const card: TarotCard = {
              id: response.card.card_id,
              name: response.card.card_name,
              number: response.card.number || '0',
              keywords: [],
              archetypes: [],
              elements: [],
              astrology: '',
              traditionalMeaning: response.card.interpretation_basic || '',
              shadowAspects: [],
              lightAspects: [],
              imageryDescription: '',
              colors: [],
              symbols: [],
              moodWeights: {},
            };
            
            setSelectedCard(card);
            setTodaysCard(card);
            setDrawnToday(true);
          }
        } catch (error) {
          if (error instanceof APIError && error.status === 409) {
            setError('You have already drawn your card for today');
            return;
          }
          if (error instanceof APIError && error.status === 403) {
            setError('Daily limit reached. Upgrade to premium for unlimited draws.');
            return;
          }
          throw error;
        }
      } else {
        // Fallback to local algorithm for unauthenticated users
        const card = selectCard(context);
        saveCardDraw(card, context);
        
        setSelectedCard(card);
        setTodaysCard(card);
        setDrawnToday(true);
      }
      
    } catch (err) {
      setError(err instanceof Error ? err.message : 'An error occurred while drawing your card.');
      console.error('Card draw error:', err);
    } finally {
      setIsDrawing(false);
    }
  }, []);

  const resetDraw = useCallback(() => {
    // Only allow reset if no card drawn today (for development/testing)
    if (!hasDrawnToday()) {
      setSelectedCard(null);
      setTodaysCard(null);
      setDrawnToday(false);
      setError(null);
    }
  }, []);

  return {
    selectedCard,
    isDrawing,
    hasDrawnToday: drawnToday,
    drawCard,
    resetDraw,
    todaysCard,
    error
  };
}