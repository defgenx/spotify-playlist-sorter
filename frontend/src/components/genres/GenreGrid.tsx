import type { GenreStat } from '@/lib/types';
import { GenreCard } from './GenreCard';

interface GenreGridProps {
  genres: GenreStat[];
}

export function GenreGrid({ genres }: GenreGridProps) {
  if (genres.length === 0) {
    return (
      <div className="text-center py-12 text-spotify-lightgray">
        No genres found. Run an analysis first.
      </div>
    );
  }

  const sortedGenres = [...genres].sort((a, b) => b.trackCount - a.trackCount);

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
      {sortedGenres.map((genre) => (
        <GenreCard key={genre.genre} genre={genre} />
      ))}
    </div>
  );
}
