import { useEffect, useRef } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/Card';
import type { ProgressEvent } from '@/lib/types';
import { CheckCircle, XCircle, Loader2 } from 'lucide-react';

interface LiveLogProps {
  events: ProgressEvent[];
}

export function LiveLog({ events }: LiveLogProps) {
  const endRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    endRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [events]);

  if (events.length === 0) {
    return null;
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle>Activity Log</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="bg-spotify-black rounded-lg p-4 max-h-96 overflow-y-auto font-mono text-sm space-y-2">
          {events.map((event, idx) => (
            <div key={idx} className="flex items-start space-x-2">
              {event.type === 'complete' && (
                <CheckCircle className="w-4 h-4 text-green-500 mt-0.5 flex-shrink-0" />
              )}
              {event.type === 'error' && (
                <XCircle className="w-4 h-4 text-red-500 mt-0.5 flex-shrink-0" />
              )}
              {event.type === 'progress' && (
                <Loader2 className="w-4 h-4 text-spotify-green animate-spin mt-0.5 flex-shrink-0" />
              )}
              <span className={
                event.type === 'error'
                  ? 'text-red-400'
                  : event.type === 'complete'
                  ? 'text-green-400'
                  : 'text-spotify-lightgray'
              }>
                {event.message}
              </span>
            </div>
          ))}
          <div ref={endRef} />
        </div>
      </CardContent>
    </Card>
  );
}
