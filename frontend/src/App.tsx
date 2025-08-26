import { useEffect } from 'react';
import { CardDrawInterface } from './components/CardDrawInterface';
import { ErrorBoundary } from './components/ErrorBoundary';
import { reportUserAction } from './utils/error-monitoring';
import './App.css';

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
      <div className="App">
        <CardDrawInterface />
      </div>
    </ErrorBoundary>
  );
}

export default App;