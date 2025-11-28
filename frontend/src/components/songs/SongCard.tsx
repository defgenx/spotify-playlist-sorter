import { Music } from 'lucide-react';
import type { Track } from '@/lib/types';
import { Card, CardContent } from '@/components/ui/Card';
import { Badge } from '@/components/ui/Badge';

interface SongCardProps {
  track: Track;
}

export function SongCard({ track }: SongCardProps) {
  return (
    <Card className="hover:bg-spotify-gray transition-colors">
      <CardContent className="flex items-center space-x-4 p-4">
        {track.albumImage ? (
          <img
            src={track.albumImage}
            alt={track.albumName}
            className="w-16 h-16 rounded"
          />
        ) : (
          <div className="w-16 h-16 bg-spotify-gray rounded flex items-center justify-center">
            <Music className="w-8 h-8 text-spotify-lightgray" />
          </div>
        )}

        <div className="flex-1 min-w-0">
          <h4 className="text-white font-medium truncate">{track.name}</h4>
          <p className="text-spotify-lightgray text-sm truncate">
            {track.artists.map((a) => a.name).join(', ')}
          </p>
          <p className="text-spotify-lightgray text-xs truncate mt-1">
            {track.albumName}
          </p>
        </div>

        {track.primaryGenre && (
          <Badge variant="default">{track.primaryGenre}</Badge>
        )}
      </CardContent>
    </Card>
  );
}
