import { Component } from 'react';
import type { ErrorInfo, ReactNode } from 'react';

interface Props {
  children: ReactNode;
}

interface State {
  hasError: boolean;
  error: Error | null;
  errorInfo: ErrorInfo | null;
}

export class ErrorBoundary extends Component<Props, State> {
  public state: State = {
    hasError: false,
    error: null,
    errorInfo: null
  };

  public static getDerivedStateFromError(error: Error): State {
    // Update state so the next render will show the fallback UI
    return { hasError: true, error, errorInfo: null };
  }

  public componentDidCatch(error: Error, errorInfo: ErrorInfo) {
    console.error('ErrorBoundary caught an error:', error, errorInfo);
    
    this.setState({
      error,
      errorInfo
    });

    // Log to monitoring service in production
    if (process.env.NODE_ENV === 'production') {
      // TODO: Send to monitoring service (Sentry, LogRocket, etc.)
      this.logErrorToService(error, errorInfo);
    }
  }

  private logErrorToService(error: Error, errorInfo: ErrorInfo) {
    // This would integrate with your monitoring service
    const errorData = {
      message: error.message,
      stack: error.stack,
      componentStack: errorInfo.componentStack,
      timestamp: new Date().toISOString(),
      userAgent: navigator.userAgent,
      url: window.location.href
    };

    // Send to monitoring service
    console.error('Error logged to monitoring:', errorData);
  }

  private handleRetry = () => {
    this.setState({ hasError: false, error: null, errorInfo: null });
  };

  private handleReload = () => {
    window.location.reload();
  };

  public render() {
    if (this.state.hasError) {
      return (
        <div className="min-h-screen bg-gradient-to-br from-gray-900 via-red-900 to-gray-900 flex items-center justify-center p-4">
          <div className="max-w-md w-full bg-gray-800/90 border border-red-500/50 rounded-xl p-6 backdrop-blur-sm">
            <div className="text-center space-y-4">
              {/* Error Icon */}
              <div className="w-16 h-16 mx-auto bg-red-500/20 rounded-full flex items-center justify-center">
                <svg className="w-8 h-8 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 19c-.77.833.192 2.5 1.732 2.5z" />
                </svg>
              </div>

              {/* Error Title */}
              <div>
                <h1 className="text-xl font-bold text-white mb-2">Something went wrong</h1>
                <p className="text-gray-300 text-sm leading-relaxed">
                  We encountered an unexpected error while loading your tarot reading. This has been automatically reported to our team.
                </p>
              </div>

              {/* Error Details (Development only) */}
              {process.env.NODE_ENV === 'development' && this.state.error && (
                <div className="bg-gray-900/50 border border-gray-600 rounded-lg p-4 text-left">
                  <h3 className="text-sm font-medium text-red-400 mb-2">Error Details (Development)</h3>
                  <div className="space-y-2">
                    <div>
                      <p className="text-xs font-medium text-gray-400">Message:</p>
                      <p className="text-xs text-gray-300 font-mono break-all">{this.state.error.message}</p>
                    </div>
                    {this.state.error.stack && (
                      <div>
                        <p className="text-xs font-medium text-gray-400">Stack:</p>
                        <pre className="text-xs text-gray-300 font-mono overflow-auto max-h-24 whitespace-pre-wrap">
                          {this.state.error.stack}
                        </pre>
                      </div>
                    )}
                  </div>
                </div>
              )}

              {/* Action Buttons */}
              <div className="space-y-3">
                <button
                  onClick={this.handleRetry}
                  className="w-full px-4 py-3 bg-purple-600 hover:bg-purple-500 text-white font-medium 
                           rounded-lg transition-colors duration-200
                           focus:outline-none focus:ring-2 focus:ring-purple-400 focus:ring-offset-2 focus:ring-offset-gray-800"
                >
                  Try Again
                </button>
                
                <button
                  onClick={this.handleReload}
                  className="w-full px-4 py-2 bg-gray-700 hover:bg-gray-600 text-gray-300 text-sm
                           rounded-lg transition-colors duration-200
                           focus:outline-none focus:ring-2 focus:ring-gray-400 focus:ring-offset-2 focus:ring-offset-gray-800"
                >
                  Reload Page
                </button>
              </div>

              {/* Contact Support */}
              <div className="pt-4 border-t border-gray-700">
                <p className="text-xs text-gray-400">
                  If this problem persists, please{' '}
                  <a 
                    href="mailto:support@symbolquest.app" 
                    className="text-purple-400 hover:text-purple-300 underline"
                  >
                    contact support
                  </a>
                </p>
              </div>
            </div>
          </div>
        </div>
      );
    }

    return this.props.children;
  }
}