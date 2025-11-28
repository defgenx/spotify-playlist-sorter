import { Music2, Sparkles } from 'lucide-react';
import type { GenreStat } from '@/lib/types';
import { Card, CardContent } from '@/components/ui/Card';
import { Badge } from '@/components/ui/Badge';

interface GenreCardProps {
  genre: GenreStat;
}

export function GenreCard({ genre }: GenreCardProps) {
  return (
    <Card variant={genre.isNew ? 'highlighted' : 'default'}>
      <CardContent className="p-6">
        <div className="flex items-start justify-between mb-4">
          <div className="flex items-center space-x-3">
            <div className="p-3 bg-spotify-green rounded-lg">
              <Music2 className="w-6 h-6 text-white" />
            </div>
            <div>
              <h3 className="text-lg font-bold text-white">{genre.genre}</h3>
              <p className="text-spotify-lightgray text-sm">
                {genre.trackCount} {genre.trackCount === 1 ? 'track' : 'tracks'}
              </p>
            </div>
          </div>
          {genre.isNew && (
            <Badge variant="warning" className="flex items-center gap-1">
              <Sparkles className="w-3 h-3" />
              New
            </Badge>
          )}
        </div>
      </CardContent>
    </Card>
  );
}
