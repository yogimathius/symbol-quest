import type { TarotCard, UserContext } from '../types/card';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api';

class APIError extends Error {
  constructor(message: string, public status?: number) {
    super(message);
    this.name = 'APIError';
  }
}

class ApiService {
  private getAuthHeaders(): HeadersInit {
    const token = localStorage.getItem('auth_token');
    return {
      'Content-Type': 'application/json',
      ...(token && { Authorization: `Bearer ${token}` }),
    };
  }

  private async handleResponse<T>(response: Response): Promise<T> {
    if (!response.ok) {
      let errorMessage = `HTTP Error ${response.status}`;
      
      try {
        const errorData = await response.json();
        errorMessage = errorData.message || errorMessage;
      } catch {
        // If JSON parsing fails, use status text
        errorMessage = response.statusText || errorMessage;
      }
      
      throw new APIError(errorMessage, response.status);
    }

    try {
      return await response.json();
    } catch {
      // Handle endpoints that don't return JSON
      return {} as T;
    }
  }

  // Authentication
  async register(email: string, password: string): Promise<{ token: string; user: any }> {
    const response = await fetch(`${API_BASE_URL}/auth/register`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email, password }),
    });

    const data = await this.handleResponse<{ token: string; user: any }>(response);
    
    // Store token in localStorage
    if (data.token) {
      localStorage.setItem('auth_token', data.token);
    }
    
    return data;
  }

  async login(email: string, password: string): Promise<{ token: string; user: any }> {
    const response = await fetch(`${API_BASE_URL}/auth/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email, password }),
    });

    const data = await this.handleResponse<{ token: string; user: any }>(response);
    
    // Store token in localStorage
    if (data.token) {
      localStorage.setItem('auth_token', data.token);
    }
    
    return data;
  }

  async logout(): Promise<void> {
    try {
      await fetch(`${API_BASE_URL}/auth/logout`, {
        method: 'POST',
        headers: this.getAuthHeaders(),
      });
    } catch (error) {
      console.warn('Logout request failed:', error);
    } finally {
      localStorage.removeItem('auth_token');
    }
  }

  async getProfile(): Promise<any> {
    const response = await fetch(`${API_BASE_URL}/auth/profile`, {
      method: 'GET',
      headers: this.getAuthHeaders(),
    });

    return this.handleResponse(response);
  }

  // Card draws
  async performDailyDraw(mood: string, question: string): Promise<{ success: boolean; card: any }> {
    const response = await fetch(`${API_BASE_URL}/draws/daily`, {
      method: 'POST',
      headers: this.getAuthHeaders(),
      body: JSON.stringify({ mood, question }),
    });

    return this.handleResponse(response);
  }

  async getDrawHistory(limit: number = 20): Promise<{ draws: any[]; count: number }> {
    const response = await fetch(`${API_BASE_URL}/draws/history?limit=${limit}`, {
      method: 'GET',
      headers: this.getAuthHeaders(),
    });

    return this.handleResponse(response);
  }

  async getTodayStatus(): Promise<{
    has_drawn: boolean;
    can_draw: boolean;
    card: any | null;
    draws_today: number;
    limit: number;
  }> {
    const response = await fetch(`${API_BASE_URL}/draws/today`, {
      method: 'GET',
      headers: this.getAuthHeaders(),
    });

    return this.handleResponse(response);
  }

  // Card meanings
  async getCardMeaning(cardId: number): Promise<{ card: TarotCard }> {
    const response = await fetch(`${API_BASE_URL}/cards/${cardId}/meaning`, {
      method: 'GET',
      headers: { 'Content-Type': 'application/json' },
    });

    return this.handleResponse(response);
  }

  // Enhanced interpretations (premium)
  async getEnhancedInterpretation(
    cardId: number, 
    mood: string, 
    question: string, 
    drawDate?: string
  ): Promise<{ interpretation: string }> {
    const response = await fetch(`${API_BASE_URL}/interpretations/enhanced`, {
      method: 'POST',
      headers: this.getAuthHeaders(),
      body: JSON.stringify({ card_id: cardId, mood, question, draw_date: drawDate }),
    });

    return this.handleResponse(response);
  }

  // Subscriptions
  async createSubscription(): Promise<{ client_secret: string; message: string }> {
    const response = await fetch(`${API_BASE_URL}/subscriptions/create`, {
      method: 'POST',
      headers: this.getAuthHeaders(),
    });

    return this.handleResponse(response);
  }

  async getSubscriptionStatus(): Promise<{
    subscription: any | null;
    status: 'free' | 'premium';
    message?: string;
  }> {
    const response = await fetch(`${API_BASE_URL}/subscriptions/status`, {
      method: 'GET',
      headers: this.getAuthHeaders(),
    });

    return this.handleResponse(response);
  }

  // Health check
  async healthCheck(): Promise<{ status: string }> {
    const response = await fetch(`${API_BASE_URL.replace('/api', '')}/health`, {
      method: 'GET',
      headers: { 'Content-Type': 'application/json' },
    });

    return this.handleResponse(response);
  }

  // Utility methods
  isAuthenticated(): boolean {
    return !!localStorage.getItem('auth_token');
  }

  getToken(): string | null {
    return localStorage.getItem('auth_token');
  }

  clearAuth(): void {
    localStorage.removeItem('auth_token');
  }
}

export const apiService = new ApiService();
export { APIError };