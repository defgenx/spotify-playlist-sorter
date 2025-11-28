import type { AuthResponse, User, SortPlan, LibraryStats } from './types';

const API_BASE = '/api';

class ApiClient {
  private async fetch<T>(
    endpoint: string,
    options?: RequestInit
  ): Promise<T> {
    const response = await fetch(`${API_BASE}${endpoint}`, {
      ...options,
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json',
        ...options?.headers,
      },
    });

    if (!response.ok) {
      const error = await response.text();
      throw new Error(error || `HTTP ${response.status}`);
    }

    return response.json();
  }

  // Auth
  async getLoginUrl(): Promise<AuthResponse> {
    return this.fetch<AuthResponse>('/auth/login');
  }

  async handleCallback(code: string): Promise<void> {
    await this.fetch(`/auth/callback?code=${code}`);
  }

  async completeLogin(token: string): Promise<void> {
    await this.fetch(`/auth/complete?token=${token}`);
  }

  async getCurrentUser(): Promise<User> {
    return this.fetch<User>('/auth/me');
  }

  async logout(): Promise<void> {
    await this.fetch('/auth/logout', { method: 'POST' });
  }

  // Library
  async getLibraryStats(): Promise<LibraryStats> {
    return this.fetch<LibraryStats>('/library/stats');
  }

  // Sort
  async createSortPlan(dryRun: boolean, enabledGroups: string[] = [], disabledPlaylists: string[] = []): Promise<SortPlan> {
    return this.fetch<SortPlan>('/sort/plan', {
      method: 'POST',
      body: JSON.stringify({ dryRun, enabledGroups, disabledPlaylists }),
    });
  }

  async executeSortPlan(dryRun: boolean, enabledGroups: string[] = [], disabledPlaylists: string[] = []): Promise<void> {
    await this.fetch(`/sort/execute`, {
      method: 'POST',
      body: JSON.stringify({ dryRun, enabledGroups, disabledPlaylists }),
    });
  }

  // SSE for progress
  createEventSource(endpoint: string): EventSource {
    return new EventSource(`${API_BASE}${endpoint}`, {
      withCredentials: true,
    });
  }
}

export const api = new ApiClient();
