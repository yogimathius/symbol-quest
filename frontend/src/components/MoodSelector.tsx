import React from 'react';
import type { Mood } from '../types/card';
import { MOOD_OPTIONS } from '../hooks/useUserContext';

interface MoodSelectorProps {
  selectedMood: Mood | '';
  onMoodSelect: (mood: Mood) => void;
  disabled?: boolean;
}

const MOOD_EMOJIS: Record<Mood, string> = {
  anxious: '😰',
  excited: '🤩',
  uncertain: '🤔',
  hopeful: '🌟',
  peaceful: '😌',
  frustrated: '😤',
  curious: '🔍',
  contemplative: '🧘'
};

const MOOD_DESCRIPTIONS: Record<Mood, string> = {
  anxious: 'Worried, stressed, or nervous',
  excited: 'Energized, enthusiastic, thrilled',
  uncertain: 'Confused, undecided, questioning',
  hopeful: 'Optimistic, expectant, positive',
  peaceful: 'Calm, serene, balanced',
  frustrated: 'Annoyed, stuck, overwhelmed',
  curious: 'Inquisitive, exploring, wondering',
  contemplative: 'Reflective, thoughtful, introspective'
};

export function MoodSelector({ selectedMood, onMoodSelect, disabled = false }: MoodSelectorProps) {
  return (
    <div className="w-full">
      <label className="block text-sm font-medium text-gray-300 mb-3">
        How are you feeling right now?
      </label>
      
      <div className="grid grid-cols-2 gap-3 sm:grid-cols-4">
        {MOOD_OPTIONS.map((mood) => (
          <button
            key={mood}
            onClick={() => onMoodSelect(mood)}
            disabled={disabled}
            className={`
              p-4 rounded-xl border-2 transition-all duration-200 text-left
              ${selectedMood === mood
                ? 'border-purple-400 bg-purple-900/30 text-white shadow-lg shadow-purple-500/20'
                : 'border-gray-600 bg-gray-800/50 text-gray-300 hover:border-gray-500 hover:bg-gray-700/50'
              }
              ${disabled 
                ? 'opacity-50 cursor-not-allowed' 
                : 'cursor-pointer hover:scale-105 active:scale-95'
              }
              focus:outline-none focus:ring-2 focus:ring-purple-400 focus:ring-offset-2 focus:ring-offset-gray-900
            `}
            aria-pressed={selectedMood === mood}
            type="button"
          >
            <div className="flex items-center space-x-2 mb-1">
              <span className="text-xl" role="img" aria-hidden="true">
                {MOOD_EMOJIS[mood]}
              </span>
              <span className="font-medium capitalize text-sm">
                {mood}
              </span>
            </div>
            <p className="text-xs text-gray-400 leading-tight">
              {MOOD_DESCRIPTIONS[mood]}
            </p>
          </button>
        ))}
      </div>

      {selectedMood && (
        <div className="mt-3 p-3 bg-purple-900/20 border border-purple-400/30 rounded-lg">
          <p className="text-sm text-purple-200">
            <span className="font-medium">Selected:</span> {MOOD_EMOJIS[selectedMood]} {selectedMood} - {MOOD_DESCRIPTIONS[selectedMood]}
          </p>
        </div>
      )}
    </div>
  );
}