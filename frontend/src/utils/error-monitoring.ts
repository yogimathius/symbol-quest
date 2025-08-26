// Error monitoring and reporting utilities

export interface ErrorReport {
  message: string;
  stack?: string;
  url: string;
  userAgent: string;
  timestamp: string;
  userId?: string;
  sessionId?: string;
  additionalContext?: Record<string, any>;
}

export interface PerformanceMetric {
  name: string;
  value: number;
  timestamp: string;
  url: string;
}

class ErrorMonitor {
  private sessionId: string;
  private userId?: string;

  constructor() {
    this.sessionId = this.generateSessionId();
    this.setupGlobalErrorHandlers();
    this.setupPerformanceMonitoring();
  }

  private generateSessionId(): string {
    return `session_${Date.now()}_${Math.random().toString(36).substring(2, 15)}`;
  }

  public setUserId(userId: string) {
    this.userId = userId;
  }

  private setupGlobalErrorHandlers() {
    // Handle uncaught JavaScript errors
    window.addEventListener('error', (event) => {
      this.reportError({
        message: event.message,
        stack: event.error?.stack,
        url: window.location.href,
        userAgent: navigator.userAgent,
        timestamp: new Date().toISOString(),
        userId: this.userId,
        sessionId: this.sessionId,
        additionalContext: {
          type: 'javascript_error',
          filename: event.filename,
          lineno: event.lineno,
          colno: event.colno
        }
      });
    });

    // Handle unhandled promise rejections
    window.addEventListener('unhandledrejection', (event) => {
      this.reportError({
        message: `Unhandled promise rejection: ${event.reason}`,
        stack: event.reason?.stack,
        url: window.location.href,
        userAgent: navigator.userAgent,
        timestamp: new Date().toISOString(),
        userId: this.userId,
        sessionId: this.sessionId,
        additionalContext: {
          type: 'unhandled_promise_rejection',
          reason: event.reason
        }
      });
    });
  }

  private setupPerformanceMonitoring() {
    // Monitor Core Web Vitals
    if ('web-vital' in window) {
      // This would integrate with web-vitals library in production
      this.monitorCoreWebVitals();
    }

    // Monitor page load performance
    window.addEventListener('load', () => {
      setTimeout(() => {
        const perfData = performance.getEntriesByType('navigation')[0] as PerformanceNavigationTiming;
        if (perfData) {
          this.reportPerformance({
            name: 'page_load_time',
            value: perfData.loadEventEnd - perfData.fetchStart,
            timestamp: new Date().toISOString(),
            url: window.location.href
          });
        }
      }, 0);
    });
  }

  private monitorCoreWebVitals() {
    // Monitor First Contentful Paint (FCP)
    const observer = new PerformanceObserver((list) => {
      for (const entry of list.getEntries()) {
        if (entry.name === 'first-contentful-paint') {
          this.reportPerformance({
            name: 'first_contentful_paint',
            value: entry.startTime,
            timestamp: new Date().toISOString(),
            url: window.location.href
          });
        }
      }
    });

    try {
      observer.observe({ entryTypes: ['paint'] });
    } catch (error) {
      console.warn('Performance Observer not supported:', error);
    }
  }

  public reportError(errorReport: ErrorReport) {
    console.error('Error reported:', errorReport);

    // In production, send to monitoring service
    if (process.env.NODE_ENV === 'production') {
      this.sendToMonitoringService('error', errorReport);
    }

    // Store locally for debugging
    this.storeErrorLocally(errorReport);
  }

  public reportPerformance(metric: PerformanceMetric) {
    console.info('Performance metric:', metric);

    // In production, send to monitoring service
    if (process.env.NODE_ENV === 'production') {
      this.sendToMonitoringService('performance', metric);
    }
  }

  public captureUserInteraction(action: string, element?: string, additionalData?: Record<string, any>) {
    const interactionData = {
      action,
      element,
      timestamp: new Date().toISOString(),
      url: window.location.href,
      userId: this.userId,
      sessionId: this.sessionId,
      ...additionalData
    };

    console.info('User interaction:', interactionData);

    // Send to analytics service in production
    if (process.env.NODE_ENV === 'production') {
      this.sendToMonitoringService('interaction', interactionData);
    }
  }

  private async sendToMonitoringService(type: 'error' | 'performance' | 'interaction', data: any) {
    try {
      // This would be replaced with actual monitoring service endpoint
      const endpoint = this.getMonitoringEndpoint(type);
      
      await fetch(endpoint, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(data),
      });
    } catch (error) {
      console.warn('Failed to send data to monitoring service:', error);
      // Fallback to local storage
      this.storeErrorLocally({ 
        ...data, 
        type, 
        sendError: error instanceof Error ? error.message : String(error) 
      });
    }
  }

  private getMonitoringEndpoint(type: string): string {
    // In production, these would be real endpoints
    const endpoints: Record<string, string> = {
      error: '/api/monitoring/errors',
      performance: '/api/monitoring/performance',
      interaction: '/api/monitoring/interactions'
    };
    return endpoints[type] || '/api/monitoring/generic';
  }

  private storeErrorLocally(data: any) {
    try {
      const storageKey = 'symbolquest_error_logs';
      const existingLogs = JSON.parse(localStorage.getItem(storageKey) || '[]');
      
      // Keep only last 50 errors
      const updatedLogs = [data, ...existingLogs].slice(0, 50);
      
      localStorage.setItem(storageKey, JSON.stringify(updatedLogs));
    } catch (error) {
      console.warn('Failed to store error locally:', error);
    }
  }

  public getStoredErrors(): any[] {
    try {
      return JSON.parse(localStorage.getItem('symbolquest_error_logs') || '[]');
    } catch {
      return [];
    }
  }

  public clearStoredErrors() {
    localStorage.removeItem('symbolquest_error_logs');
  }
}

// Singleton instance
export const errorMonitor = new ErrorMonitor();

// Convenience functions for common use cases
export const reportError = (error: Error, context?: Record<string, any>) => {
  errorMonitor.reportError({
    message: error.message,
    stack: error.stack,
    url: window.location.href,
    userAgent: navigator.userAgent,
    timestamp: new Date().toISOString(),
    additionalContext: context
  });
};

export const reportUserAction = (action: string, element?: string, data?: Record<string, any>) => {
  errorMonitor.captureUserInteraction(action, element, data);
};

export const setUserId = (userId: string) => {
  errorMonitor.setUserId(userId);
};

// Network error handler for API calls
export const withErrorHandling = async <T>(
  apiCall: () => Promise<T>,
  context?: Record<string, any>
): Promise<T> => {
  try {
    return await apiCall();
  } catch (error) {
    reportError(error instanceof Error ? error : new Error(String(error)), {
      type: 'api_call',
      ...context
    });
    throw error;
  }
};