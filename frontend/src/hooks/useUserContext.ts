import { useState, useCallback } from 'react';
import type { Mood } from '../types/card';

export interface UseUserContextReturn {
  mood: Mood | '';
  question: string;
  setMood: (mood: Mood | '') => void;
  setQuestion: (question: string) => void;
  isValid: boolean;
  reset: () => void;
}

const MOOD_OPTIONS: Mood[] = [
  'anxious',
  'excited', 
  'uncertain',
  'hopeful',
  'peaceful',
  'frustrated',
  'curious',
  'contemplative'
];

export function useUserContext(): UseUserContextReturn {
  const [mood, setMoodState] = useState<Mood | ''>('');
  const [question, setQuestionState] = useState('');

  const setMood = useCallback((newMood: Mood | '') => {
    setMoodState(newMood);
  }, []);

  const setQuestion = useCallback((newQuestion: string) => {
    // Trim and limit question length
    const trimmed = newQuestion.trim().slice(0, 200);
    setQuestionState(trimmed);
  }, []);

  const reset = useCallback(() => {
    setMoodState('');
    setQuestionState('');
  }, []);

  const isValid = mood !== '' && question.trim().length >= 10;

  return {
    mood,
    question,
    setMood,
    setQuestion,
    isValid,
    reset
  };
}

export { MOOD_OPTIONS };