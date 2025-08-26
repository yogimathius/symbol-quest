export interface TarotCard {
  id: number;
  name: string;
  number: string;
  keywords: string[];
  archetypes: string[];
  elements: string[];
  astrology?: string;
  traditionalMeaning: string;
  shadowAspects: string[];
  lightAspects: string[];
  imageryDescription: string;
  colors: string[];
  symbols: string[];
  moodWeights: Record<string, number>;
}

export interface UserContext {
  mood: string;
  question: string;
  timestamp: number;
}

export interface CardDraw {
  card: TarotCard;
  context: UserContext;
  interpretation?: string;
  date: string;
  id: string;
}

export type Mood = 
  | 'anxious' 
  | 'excited' 
  | 'uncertain' 
  | 'hopeful' 
  | 'peaceful' 
  | 'frustrated' 
  | 'curious' 
  | 'contemplative';