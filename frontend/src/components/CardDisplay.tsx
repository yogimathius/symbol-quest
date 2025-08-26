import type { TarotCard } from '../types/card';

interface CardDisplayProps {
  card: TarotCard;
  isRevealing?: boolean;
  showDetails?: boolean;
}

export function CardDisplay({ card, isRevealing = false, showDetails = true }: CardDisplayProps) {
  return (
    <div className="w-full max-w-md mx-auto">
      {/* Card Container */}
      <div 
        className={`
          relative perspective-1000 transition-transform duration-700
          ${isRevealing ? 'animate-card-flip' : ''}
        `}
      >
        {/* Card Back (visible during reveal) */}
        {isRevealing && (
          <div className="absolute inset-0 backface-hidden">
            <div className="w-full aspect-[2/3] bg-gradient-to-br from-purple-900 via-blue-900 to-indigo-900 
                          rounded-xl border-2 border-purple-400/50 shadow-2xl shadow-purple-500/20
                          flex items-center justify-center">
              <div className="text-center space-y-4">
                <div className="w-16 h-16 mx-auto rounded-full bg-gradient-to-br from-yellow-300 to-yellow-500 
                              flex items-center justify-center">
                  <svg className="w-8 h-8 text-purple-900" fill="currentColor" viewBox="0 0 20 20">
                    <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z" />
                  </svg>
                </div>
                <div className="text-sm text-purple-200 font-medium">Symbol Quest</div>
              </div>
            </div>
          </div>
        )}

        {/* Card Front */}
        <div className={`w-full ${isRevealing ? 'rotateY-180' : ''}`}>
          <div className="w-full aspect-[2/3] bg-gradient-to-br from-gray-800 via-gray-900 to-black 
                        rounded-xl border-2 border-gray-600 shadow-2xl shadow-black/40
                        flex flex-col overflow-hidden">
            
            {/* Card Header */}
            <div className="bg-gradient-to-r from-purple-800 to-blue-800 px-4 py-3 border-b border-gray-600">
              <div className="flex items-center justify-between">
                <span className="text-xs text-purple-200 font-medium">{card.number}</span>
                <div className="flex items-center space-x-1">
                  {card.elements.map((element, index) => (
                    <span key={index} className="text-xs text-blue-200 bg-blue-900/50 px-2 py-1 rounded">
                      {element}
                    </span>
                  ))}
                </div>
              </div>
              <h2 className="text-lg font-bold text-white mt-1">{card.name}</h2>
            </div>

            {/* Card Image Placeholder */}
            <div className="flex-1 flex items-center justify-center bg-gradient-to-b from-gray-700 to-gray-800 p-6">
              <div className="text-center space-y-4">
                {/* Symbolic Representation */}
                <div className="w-20 h-20 mx-auto bg-gradient-to-br from-yellow-400 via-orange-400 to-red-400 
                              rounded-full flex items-center justify-center shadow-lg">
                  <span className="text-2xl font-bold text-white">
                    {card.number === '0' ? '∞' : card.number}
                  </span>
                </div>
                
                {/* Imagery Description */}
                <p className="text-xs text-gray-300 leading-relaxed px-2">
                  {card.imageryDescription}
                </p>

                {/* Key Symbols */}
                <div className="flex flex-wrap gap-1 justify-center">
                  {card.symbols.slice(0, 4).map((symbol, index) => (
                    <span key={index} className="text-xs text-gray-400 bg-gray-700 px-2 py-1 rounded">
                      {symbol}
                    </span>
                  ))}
                </div>
              </div>
            </div>

            {/* Card Footer */}
            {showDetails && (
              <div className="bg-gray-800 px-4 py-3 border-t border-gray-600 space-y-2">
                {/* Keywords */}
                <div>
                  <p className="text-xs text-gray-400 font-medium mb-1">Keywords:</p>
                  <div className="flex flex-wrap gap-1">
                    {card.keywords.slice(0, 4).map((keyword, index) => (
                      <span key={index} className="text-xs text-purple-200 bg-purple-900/30 px-2 py-1 rounded">
                        {keyword.replace('-', ' ')}
                      </span>
                    ))}
                  </div>
                </div>

                {/* Archetypes */}
                <div>
                  <p className="text-xs text-gray-400 font-medium mb-1">Archetypes:</p>
                  <div className="flex flex-wrap gap-1">
                    {card.archetypes.map((archetype, index) => (
                      <span key={index} className="text-xs text-blue-200 bg-blue-900/30 px-2 py-1 rounded">
                        {archetype}
                      </span>
                    ))}
                  </div>
                </div>
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Card Title for Screen Readers */}
      <div className="sr-only">
        <h2>{card.name} - {card.traditionalMeaning}</h2>
        <p>Keywords: {card.keywords.join(', ')}</p>
        <p>Archetypes: {card.archetypes.join(', ')}</p>
      </div>
    </div>
  );
}