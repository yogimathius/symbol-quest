import React, { useState } from 'react';
import { MoodSelector } from './MoodSelector';
import { QuestionInput } from './QuestionInput';
import { CardDisplay } from './CardDisplay';
import { useCardDraw } from '../hooks/useCardDraw';
import { useUserContext } from '../hooks/useUserContext';
import type { Mood } from '../types/card';

export function CardDrawInterface() {
  const [step, setStep] = useState<'input' | 'drawing' | 'result'>('input');
  const { mood, question, setMood, setQuestion, isValid, reset } = useUserContext();
  const { selectedCard, isDrawing, hasDrawnToday, drawCard, todaysCard, error } = useCardDraw();

  const handleDrawCard = async () => {
    if (!isValid) return;

    setStep('drawing');
    
    await drawCard({
      mood: mood as Mood,
      question,
      timestamp: Date.now()
    });

    setStep('result');
  };

  const handleNewReading = () => {
    reset();
    setStep('input');
  };

  const handleBackToInput = () => {
    setStep('input');
  };

  // Show today's card if already drawn
  if (hasDrawnToday && todaysCard && step === 'input') {
    return (
      <div className="min-h-screen bg-gradient-to-br from-gray-900 via-purple-900 to-blue-900 p-4">
        <div className="max-w-lg mx-auto py-8">
          {/* Header */}
          <div className="text-center mb-8">
            <h1 className="text-3xl font-bold text-white mb-2">Symbol Quest</h1>
            <p className="text-purple-200">Your daily card has been drawn</p>
          </div>

          {/* Today's Card */}
          <div className="mb-8">
            <CardDisplay card={todaysCard} showDetails={true} />
          </div>

          {/* Message */}
          <div className="bg-purple-900/30 border border-purple-400/30 rounded-xl p-6 mb-6">
            <div className="flex items-center space-x-2 mb-3">
              <svg className="w-5 h-5 text-purple-300" fill="currentColor" viewBox="0 0 20 20">
                <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
              </svg>
              <h3 className="text-lg font-medium text-white">Today's Guidance</h3>
            </div>
            <p className="text-purple-200 text-sm leading-relaxed">
              You've already drawn your card for today. Take time to reflect on its message and how it applies to your current situation. 
              Come back tomorrow for your next daily reading.
            </p>
          </div>

          {/* Action Button */}
          <button
            onClick={handleNewReading}
            className="w-full px-6 py-3 bg-gradient-to-r from-purple-600 to-blue-600 text-white font-medium 
                     rounded-xl transition-all duration-200 hover:from-purple-500 hover:to-blue-500
                     focus:outline-none focus:ring-2 focus:ring-purple-400 focus:ring-offset-2 focus:ring-offset-gray-900
                     disabled:opacity-50 disabled:cursor-not-allowed"
          >
            Reflect on Today's Message
          </button>
        </div>
      </div>
    );
  }

  // Input step
  if (step === 'input') {
    return (
      <div className="min-h-screen bg-gradient-to-br from-gray-900 via-purple-900 to-blue-900 p-4">
        <div className="max-w-lg mx-auto py-8">
          {/* Header */}
          <div className="text-center mb-8">
            <h1 className="text-3xl font-bold text-white mb-2">Symbol Quest</h1>
            <p className="text-purple-200">Daily guidance through symbolic wisdom</p>
          </div>

          {/* Error Display */}
          {error && (
            <div className="bg-red-900/30 border border-red-400/30 rounded-xl p-4 mb-6">
              <div className="flex items-center space-x-2">
                <svg className="w-5 h-5 text-red-300 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
                  <path fillRule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clipRule="evenodd" />
                </svg>
                <p className="text-red-200 text-sm">{error}</p>
              </div>
            </div>
          )}

          {/* Form */}
          <div className="space-y-8">
            {/* Mood Selection */}
            <MoodSelector
              selectedMood={mood}
              onMoodSelect={setMood}
              disabled={isDrawing}
            />

            {/* Question Input */}
            <QuestionInput
              question={question}
              onQuestionChange={setQuestion}
              disabled={isDrawing}
            />

            {/* Draw Button */}
            <button
              onClick={handleDrawCard}
              disabled={!isValid || isDrawing}
              className={`
                w-full px-6 py-4 font-medium rounded-xl text-white
                transition-all duration-300 transform
                ${isValid && !isDrawing
                  ? 'bg-gradient-to-r from-purple-600 to-blue-600 hover:from-purple-500 hover:to-blue-500 hover:scale-105 shadow-lg shadow-purple-500/25'
                  : 'bg-gray-700 cursor-not-allowed opacity-50'
                }
                focus:outline-none focus:ring-2 focus:ring-purple-400 focus:ring-offset-2 focus:ring-offset-gray-900
                disabled:transform-none disabled:shadow-none
              `}
            >
              {isDrawing ? (
                <div className="flex items-center justify-center space-x-2">
                  <svg className="animate-spin w-5 h-5 text-white" fill="none" viewBox="0 0 24 24">
                    <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                    <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                  </svg>
                  <span>Drawing your card...</span>
                </div>
              ) : isValid ? (
                'Draw Your Daily Card'
              ) : (
                'Please complete both fields above'
              )}
            </button>

            {/* Help Text */}
            <div className="text-center text-xs text-gray-400 leading-relaxed">
              <p>
                Your daily card is selected through an intelligent algorithm that considers your mood and question to provide meaningful guidance. 
                You can draw one card per day.
              </p>
            </div>
          </div>
        </div>
      </div>
    );
  }

  // Drawing step
  if (step === 'drawing') {
    return (
      <div className="min-h-screen bg-gradient-to-br from-gray-900 via-purple-900 to-blue-900 p-4">
        <div className="max-w-lg mx-auto py-8">
          {/* Header */}
          <div className="text-center mb-8">
            <h1 className="text-3xl font-bold text-white mb-2">Symbol Quest</h1>
            <p className="text-purple-200">Drawing your card...</p>
          </div>

          {/* Drawing Animation */}
          <div className="flex flex-col items-center justify-center min-h-[400px] space-y-8">
            <div className="relative">
              {/* Animated Card Stack */}
              <div className="space-y-2">
                {[0, 1, 2].map((index) => (
                  <div
                    key={index}
                    className={`
                      w-32 h-48 bg-gradient-to-br from-purple-900 via-blue-900 to-indigo-900 
                      rounded-lg border border-purple-400/50 shadow-lg
                      animate-shuffle-${index} transform origin-center
                    `}
                    style={{
                      transform: `rotate(${(index - 1) * 5}deg) translateY(${index * -4}px)`,
                      animationDelay: `${index * 0.2}s`
                    }}
                  />
                ))}
              </div>

              {/* Sparkle Effects */}
              <div className="absolute inset-0 pointer-events-none">
                {[...Array(6)].map((_, i) => (
                  <div
                    key={i}
                    className="absolute w-2 h-2 bg-yellow-300 rounded-full animate-sparkle"
                    style={{
                      left: `${20 + Math.random() * 60}%`,
                      top: `${20 + Math.random() * 60}%`,
                      animationDelay: `${i * 0.3}s`
                    }}
                  />
                ))}
              </div>
            </div>

            {/* Loading Text */}
            <div className="text-center space-y-2">
              <div className="flex items-center justify-center space-x-2">
                <svg className="animate-spin w-6 h-6 text-purple-300" fill="none" viewBox="0 0 24 24">
                  <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                  <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                <span className="text-lg text-white font-medium">Selecting your card...</span>
              </div>
              <p className="text-purple-200 text-sm">
                Considering your {mood} mood and question about guidance
              </p>
            </div>
          </div>
        </div>
      </div>
    );
  }

  // Result step
  if (step === 'result' && selectedCard) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-gray-900 via-purple-900 to-blue-900 p-4">
        <div className="max-w-lg mx-auto py-8">
          {/* Header */}
          <div className="text-center mb-8">
            <h1 className="text-3xl font-bold text-white mb-2">Your Daily Card</h1>
            <p className="text-purple-200">Drawn with intention and wisdom</p>
          </div>

          {/* Card Display */}
          <div className="mb-8">
            <CardDisplay card={selectedCard} isRevealing={true} showDetails={true} />
          </div>

          {/* Card Meaning */}
          <div className="bg-gray-800/50 border border-gray-600 rounded-xl p-6 mb-6">
            <h3 className="text-lg font-medium text-white mb-3">Traditional Meaning</h3>
            <p className="text-gray-300 text-sm leading-relaxed mb-4">
              {selectedCard.traditionalMeaning}
            </p>

            {/* Light and Shadow Aspects */}
            <div className="grid gap-4 sm:grid-cols-2">
              <div>
                <h4 className="text-sm font-medium text-green-400 mb-2 flex items-center">
                  <svg className="w-4 h-4 mr-1" fill="currentColor" viewBox="0 0 20 20">
                    <path fillRule="evenodd" d="M10 2L13.09 8.26L20 9L15 13.74L16.18 20.02L10 16.77L3.82 20.02L5 13.74L0 9L6.91 8.26L10 2Z" />
                  </svg>
                  Light Aspects
                </h4>
                <div className="flex flex-wrap gap-1">
                  {selectedCard.lightAspects.map((aspect, index) => (
                    <span key={index} className="text-xs text-green-200 bg-green-900/30 px-2 py-1 rounded">
                      {aspect.replace('-', ' ')}
                    </span>
                  ))}
                </div>
              </div>

              <div>
                <h4 className="text-sm font-medium text-orange-400 mb-2 flex items-center">
                  <svg className="w-4 h-4 mr-1" fill="currentColor" viewBox="0 0 20 20">
                    <path fillRule="evenodd" d="M3 4a1 1 0 011-1h12a1 1 0 011 1v2a1 1 0 01-1 1H4a1 1 0 01-1-1V4zM3 10a1 1 0 011-1h6a1 1 0 011 1v6a1 1 0 01-1 1H4a1 1 0 01-1-1v-6zM14 9a1 1 0 00-1 1v6a1 1 0 001 1h2a1 1 0 001-1v-6a1 1 0 00-1-1h-2z" />
                  </svg>
                  Shadow Aspects
                </h4>
                <div className="flex flex-wrap gap-1">
                  {selectedCard.shadowAspects.map((aspect, index) => (
                    <span key={index} className="text-xs text-orange-200 bg-orange-900/30 px-2 py-1 rounded">
                      {aspect.replace('-', ' ')}
                    </span>
                  ))}
                </div>
              </div>
            </div>
          </div>

          {/* Context Display */}
          <div className="bg-purple-900/20 border border-purple-400/30 rounded-xl p-4 mb-6">
            <h3 className="text-sm font-medium text-purple-300 mb-2">Your Reading Context</h3>
            <div className="text-xs text-purple-200 space-y-1">
              <p><span className="font-medium">Mood:</span> {mood}</p>
              <p><span className="font-medium">Question:</span> "{question}"</p>
              <p><span className="font-medium">Date:</span> {new Date().toLocaleDateString()}</p>
            </div>
          </div>

          {/* Actions */}
          <div className="space-y-3">
            <button
              onClick={handleBackToInput}
              className="w-full px-6 py-3 bg-gray-700 hover:bg-gray-600 text-white font-medium 
                       rounded-xl transition-colors duration-200
                       focus:outline-none focus:ring-2 focus:ring-gray-400 focus:ring-offset-2 focus:ring-offset-gray-900"
            >
              Reflect on This Reading
            </button>
            
            <p className="text-center text-xs text-gray-400">
              Come back tomorrow for your next daily card
            </p>
          </div>
        </div>
      </div>
    );
  }

  return null;
}