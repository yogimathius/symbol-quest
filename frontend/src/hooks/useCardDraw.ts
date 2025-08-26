import { useState, useCallback, useEffect } from 'react';
import { selectCard, saveCardDraw, hasDrawnToday, getTodaysCard } from '../utils/card-selection';
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
    const checkTodaysDraw = () => {
      const hasDrawn = hasDrawnToday();
      setDrawnToday(hasDrawn);
      
      if (hasDrawn) {
        const card = getTodaysCard();
        setTodaysCard(card);
        setSelectedCard(card);
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
      if (hasDrawnToday()) {
        const existing = getTodaysCard();
        if (existing) {
          setSelectedCard(existing);
          setDrawnToday(true);
          return;
        }
      }

      // Add delay for card drawing animation
      await new Promise(resolve => setTimeout(resolve, 2000));

      // Select card using intelligent algorithm
      const card = selectCard(context);
      
      // Save the draw
      saveCardDraw(card, context);
      
      // Update state
      setSelectedCard(card);
      setTodaysCard(card);
      setDrawnToday(true);
      
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