import type { Track } from '@/lib/types';
import { SongCard } from './SongCard';

interface SongListProps {
  tracks: Track[];
  emptyMessage?: string;
}

export function SongList({ tracks, emptyMessage = 'No tracks found' }: SongListProps) {
  if (tracks.length === 0) {
    return (
      <div className="text-center py-12 text-spotify-lightgray">
        {emptyMessage}
      </div>
    );
  }

  return (
    <div className="space-y-3">
      {tracks.map((track) => (
        <SongCard key={track.id} track={track} />
      ))}
    </div>
  );
}
