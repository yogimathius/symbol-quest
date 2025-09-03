import { useEffect, useState } from 'react';
import { CardDrawInterface } from './components/CardDrawInterface';
import { ErrorBoundary } from './components/ErrorBoundary';
import { Auth } from './components/Auth';
import { AuthProvider, useAuth } from './contexts/AuthContext';
import { reportUserAction } from './utils/error-monitoring';
import './App.css';

function AppContent() {
  const [showAuth, setShowAuth] = useState(false);
  const { user, isAuthenticated, logout } = useAuth();

  return (
    <div className="App">
      {/* Header with authentication */}
      <header className="fixed top-0 right-0 p-4 z-40">
        {isAuthenticated ? (
          <div className="flex items-center gap-4 text-white">
            <div className="text-sm">
              <span className="text-purple-300">Welcome, </span>
              <span>{user?.email}</span>
              {user?.subscription_tier === 'premium' && (
                <span className="ml-2 px-2 py-1 bg-yellow-500/20 text-yellow-300 rounded text-xs">
                  Premium
                </span>
              )}
            </div>
            <button
              onClick={logout}
              className="text-purple-300 hover:text-white transition-colors text-sm"
            >
              Sign Out
            </button>
          </div>
        ) : (
          <button
            onClick={() => setShowAuth(true)}
            className="bg-gradient-to-r from-purple-600 to-indigo-600 hover:from-purple-700 hover:to-indigo-700 text-white px-4 py-2 rounded-lg transition-all duration-200"
          >
            Sign In
          </button>
        )}
      </header>

      <CardDrawInterface />
      
      {showAuth && <Auth onClose={() => setShowAuth(false)} />}
    </div>
  );
}

function App() {
  useEffect(() => {
    // Initialize error monitoring
    reportUserAction('app_load', 'App', {
      timestamp: new Date().toISOString(),
      url: window.location.href
    });

    // Set up performance monitoring
    const handleLoad = () => {
      reportUserAction('app_ready', 'App');
    };

    if (document.readyState === 'loading') {
      document.addEventListener('DOMContentLoaded', handleLoad);
    } else {
      handleLoad();
    }

    return () => {
      document.removeEventListener('DOMContentLoaded', handleLoad);
    };
  }, []);

  return (
    <ErrorBoundary>
      <AuthProvider>
        <AppContent />
      </AuthProvider>
    </ErrorBoundary>
  );
}

export default App;