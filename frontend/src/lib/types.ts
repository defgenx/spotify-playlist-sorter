export interface Artist {
  id: string;
  name: string;
}

export interface Track {
  id: string;
  name: string;
  artists: Artist[];
  albumName: string;
  albumImage: string;
  primaryGenre: string;
}

export interface TrackMove {
  trackId: string;
  trackName: string;
  artistName: string;
  albumImage: string;
  genre: string;
  fromPlaylist: string;
  toPlaylist: string;
  reason: string;
}

export interface GenreStat {
  genre: string;
  trackCount: number;
  isNew: boolean;
}

export interface SortPlan {
  id: string;
  dryRun: boolean;
  tracksToAdd: TrackMove[];
  tracksToRemove: TrackMove[];
  playlistsToCreate: string[];
  uncategorizedTracks: Track[];
  genreStats: GenreStat[];
}

export interface User {
  id: string;
  displayName: string;
  email: string;
  imageUrl: string;
}

export interface LibraryStats {
  totalTracks: number;
  totalPlaylists: number;
  totalGenres: number;
  analyzedAt?: string;
}

export interface ProgressEvent {
  type: 'progress' | 'complete' | 'error';
  message: string;
  current?: number;
  total?: number;
  percentage?: number;
}

export interface AuthResponse {
  url: string;
}
