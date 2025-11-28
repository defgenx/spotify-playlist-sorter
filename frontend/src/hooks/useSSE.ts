import { useEffect, useState, useCallback } from 'react';
import type { ProgressEvent } from '@/lib/types';

interface UseSSEOptions {
  onMessage?: (event: ProgressEvent) => void;
  onError?: (error: Event) => void;
  onComplete?: () => void;
}

export function useSSE(endpoint: string | null, options: UseSSEOptions = {}) {
  const [events, setEvents] = useState<ProgressEvent[]>([]);
  const [isConnected, setIsConnected] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!endpoint) return;

    const eventSource = new EventSource(endpoint, {
      withCredentials: true,
    });

    eventSource.onopen = () => {
      setIsConnected(true);
      setError(null);
    };

    eventSource.onmessage = (event) => {
      try {
        const data: ProgressEvent = JSON.parse(event.data);
        setEvents((prev) => [...prev, data]);

        if (options.onMessage) {
          options.onMessage(data);
        }

        if (data.type === 'complete') {
          eventSource.close();
          setIsConnected(false);
          if (options.onComplete) {
            options.onComplete();
          }
        }

        if (data.type === 'error') {
          eventSource.close();
          setIsConnected(false);
          setError(data.message);
        }
      } catch (err) {
        console.error('Failed to parse SSE event:', err);
      }
    };

    eventSource.onerror = (err) => {
      console.error('SSE error:', err);
      setError('Connection lost');
      setIsConnected(false);
      eventSource.close();

      if (options.onError) {
        options.onError(err);
      }
    };

    return () => {
      eventSource.close();
      setIsConnected(false);
    };
  }, [endpoint]);

  const clear = useCallback(() => {
    setEvents([]);
    setError(null);
  }, []);

  return {
    events,
    isConnected,
    error,
    clear,
  };
}
