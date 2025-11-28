import { ArrowRight, Music } from 'lucide-react';
import type { TrackMove as TrackMoveType } from '@/lib/types';
import { Card, CardContent } from '@/components/ui/Card';
import { Badge } from '@/components/ui/Badge';

interface TrackMoveProps {
  move: TrackMoveType;
  variant?: 'add' | 'remove';
}

export function TrackMove({ move, variant = 'add' }: TrackMoveProps) {
  return (
    <Card className="hover:bg-spotify-gray transition-colors">
      <CardContent className="flex items-center space-x-4 p-4">
        {move.albumImage ? (
          <img
            src={move.albumImage}
            alt={move.trackName}
            className="w-12 h-12 rounded"
          />
        ) : (
          <div className="w-12 h-12 bg-spotify-gray rounded flex items-center justify-center">
            <Music className="w-6 h-6 text-spotify-lightgray" />
          </div>
        )}

        <div className="flex-1 min-w-0">
          <h4 className="text-white font-medium truncate">{move.trackName}</h4>
          <p className="text-spotify-lightgray text-sm truncate">{move.artistName}</p>

          <div className="flex items-center space-x-2 mt-2 text-xs">
            <span className="text-spotify-lightgray">{move.fromPlaylist}</span>
            <ArrowRight className="w-3 h-3 text-spotify-lightgray" />
            <span className="text-white font-medium">{move.toPlaylist}</span>
          </div>

          {move.reason && (
            <p className="text-spotify-lightgray text-xs mt-1 italic">{move.reason}</p>
          )}
        </div>

        <div className="flex flex-col items-end space-y-2">
          <Badge variant={variant === 'add' ? 'success' : 'danger'}>
            {variant === 'add' ? 'Add' : 'Remove'}
          </Badge>
          <Badge variant="default">{move.genre}</Badge>
        </div>
      </CardContent>
    </Card>
  );
}
