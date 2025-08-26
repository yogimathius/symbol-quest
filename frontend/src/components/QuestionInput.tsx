import React from 'react';

interface QuestionInputProps {
  question: string;
  onQuestionChange: (question: string) => void;
  disabled?: boolean;
}

const EXAMPLE_QUESTIONS = [
  "What should I focus on today?",
  "How can I navigate this challenge?",
  "What do I need to know about my relationship?",
  "What's blocking my progress?",
  "How can I find more balance?",
  "What creative inspiration do I need?",
  "How should I approach this decision?",
  "What lesson is this situation teaching me?"
];

export function QuestionInput({ question, onQuestionChange, disabled = false }: QuestionInputProps) {
  const handleQuestionSelect = (exampleQuestion: string) => {
    if (!disabled) {
      onQuestionChange(exampleQuestion);
    }
  };

  const characterCount = question.length;
  const isValid = question.trim().length >= 10;
  const maxLength = 200;

  return (
    <div className="w-full">
      <label className="block text-sm font-medium text-gray-300 mb-3">
        What would you like guidance on?
        <span className="text-xs text-gray-400 ml-2">(minimum 10 characters)</span>
      </label>

      <div className="relative">
        <textarea
          value={question}
          onChange={(e) => onQuestionChange(e.target.value)}
          disabled={disabled}
          placeholder="Ask about anything on your mind - relationships, career, personal growth, decisions..."
          className={`
            w-full px-4 py-3 bg-gray-800/50 border-2 rounded-xl resize-none
            text-white placeholder-gray-400 text-sm leading-relaxed
            transition-all duration-200
            ${isValid 
              ? 'border-green-500/50 focus:border-green-400' 
              : 'border-gray-600 focus:border-purple-400'
            }
            ${disabled 
              ? 'opacity-50 cursor-not-allowed' 
              : 'hover:border-gray-500'
            }
            focus:outline-none focus:ring-2 focus:ring-purple-400/50 focus:ring-offset-2 focus:ring-offset-gray-900
          `}
          rows={3}
          maxLength={maxLength}
          aria-describedby="question-help character-count"
        />
        
        <div 
          id="character-count"
          className={`
            absolute bottom-3 right-3 text-xs 
            ${characterCount > maxLength * 0.9 
              ? 'text-yellow-400' 
              : 'text-gray-500'
            }
          `}
        >
          {characterCount}/{maxLength}
        </div>
      </div>

      {!isValid && question.length > 0 && question.length < 10 && (
        <p className="text-sm text-amber-400 mt-2">
          Please write at least 10 characters for a more meaningful reading.
        </p>
      )}

      {isValid && (
        <p className="text-sm text-green-400 mt-2 flex items-center">
          <svg className="w-4 h-4 mr-1" fill="currentColor" viewBox="0 0 20 20">
            <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
          </svg>
          Your question is ready for a reading
        </p>
      )}

      {!disabled && (
        <>
          <div id="question-help" className="mt-4">
            <p className="text-xs text-gray-400 mb-3">
              Need inspiration? Try one of these questions:
            </p>
            <div className="grid gap-2 sm:grid-cols-2">
              {EXAMPLE_QUESTIONS.map((exampleQ, index) => (
                <button
                  key={index}
                  onClick={() => handleQuestionSelect(exampleQ)}
                  className="text-left text-xs text-purple-300 hover:text-purple-200 
                           bg-gray-800/30 hover:bg-gray-700/50 p-2 rounded-lg 
                           transition-colors duration-150 border border-gray-700 
                           hover:border-gray-600"
                  type="button"
                >
                  "{exampleQ}"
                </button>
              ))}
            </div>
          </div>
        </>
      )}
    </div>
  );
}