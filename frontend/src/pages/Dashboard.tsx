import { useState } from 'react';
import { Play, Music, Disc, ListMusic, Layers, ChevronDown, ChevronUp } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/Card';
import { Button } from '@/components/ui/Button';
import { Badge } from '@/components/ui/Badge';
import { ProgressBar } from '@/components/progress/ProgressBar';
import { LiveLog } from '@/components/progress/LiveLog';
import { useSSE } from '@/hooks/useSSE';
import { useUIStore } from '@/stores/uiStore';
import { useNavigate } from 'react-router-dom';

interface GroupSuggestion {
  parentGenre: string;
  childGenres: string[];
  totalTracks: number;
  playlistsToMerge: number;
}

interface GenreGroup {
  parent: string;
  children: string[];
  count: number;
}

interface AnalysisResult {
  totalLikedSongs: number;
  playlists: { id: string; name: string; trackCount: number }[];
  genreDistribution: Record<string, number>;
  groupingSuggestions: GroupSuggestion[];
  genreGroups: Record<string, GenreGroup>;
}

export function Dashboard() {
  const navigate = useNavigate();
  const [isAnalyzing, setIsAnalyzing] = useState(false);
  const [analysisEndpoint, setAnalysisEndpoint] = useState<string | null>(null);
  const [analysisResult, setAnalysisResult] = useState<AnalysisResult | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [expandedGroups, setExpandedGroups] = useState<Set<string>>(new Set());

  const enabledGroups = useUIStore((state) => state.enabledGroups);
  const toggleGroup = useUIStore((state) => state.toggleGroup);
  const enableAllGroupsStore = useUIStore((state) => state.enableAllGroups);
  const disableAllGroupsStore = useUIStore((state) => state.disableAllGroups);

  const { events } = useSSE(analysisEndpoint, {
    onComplete: () => {
      setIsAnalyzing(false);
      setAnalysisEndpoint(null);
    },
  });

  const toggleExpanded = (parent: string) => {
    setExpandedGroups(prev => {
      const newSet = new Set(prev);
      if (newSet.has(parent)) {
        newSet.delete(parent);
      } else {
        newSet.add(parent);
      }
      return newSet;
    });
  };

  const enableAllGroups = () => {
    if (analysisResult?.groupingSuggestions) {
      enableAllGroupsStore(analysisResult.groupingSuggestions.map(s => s.parentGenre));
    }
  };

  const disableAllGroups = () => {
    disableAllGroupsStore();
  };

  const handleAnalyze = async () => {
    setIsAnalyzing(true);
    setError(null);
    setAnalysisEndpoint('/api/events?type=analysis');

    try {
      const response = await fetch('/api/library/analysis', {
        credentials: 'include',
      });

      if (!response.ok) {
        throw new Error(`Analysis failed: ${response.statusText}`);
      }

      const result = await response.json();
      setAnalysisResult(result);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Analysis failed');
    } finally {
      setIsAnalyzing(false);
      setAnalysisEndpoint(null);
    }
  };

  const currentEvent = events[events.length - 1];
  const showProgress = isAnalyzing && currentEvent?.current && currentEvent?.total;

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold text-white mb-2">Dashboard</h1>
        <p className="text-spotify-lightgray">
          Analyze your Spotify library and organize tracks by genre
        </p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <Card variant="highlighted">
          <CardHeader>
            <div className="flex items-center space-x-3">
              <div className="p-2 bg-spotify-green rounded-lg">
                <Music className="w-5 h-5 text-white" />
              </div>
              <CardTitle>Total Tracks</CardTitle>
            </div>
          </CardHeader>
          <CardContent>
            <p className="text-3xl font-bold text-white">
              {analysisResult?.totalLikedSongs ?? '-'}
            </p>
            <p className="text-sm text-spotify-lightgray mt-1">
              {analysisResult ? 'Liked songs in library' : 'Run analysis to see stats'}
            </p>
          </CardContent>
        </Card>

        <Card variant="highlighted">
          <CardHeader>
            <div className="flex items-center space-x-3">
              <div className="p-2 bg-blue-600 rounded-lg">
                <ListMusic className="w-5 h-5 text-white" />
              </div>
              <CardTitle>Playlists</CardTitle>
            </div>
          </CardHeader>
          <CardContent>
            <p className="text-3xl font-bold text-white">
              {analysisResult?.playlists?.length ?? '-'}
            </p>
            <p className="text-sm text-spotify-lightgray mt-1">
              {analysisResult ? 'User-created playlists' : 'Run analysis to see stats'}
            </p>
          </CardContent>
        </Card>

        <Card variant="highlighted">
          <CardHeader>
            <div className="flex items-center space-x-3">
              <div className="p-2 bg-purple-600 rounded-lg">
                <Disc className="w-5 h-5 text-white" />
              </div>
              <CardTitle>Genres</CardTitle>
            </div>
          </CardHeader>
          <CardContent>
            <p className="text-3xl font-bold text-white">
              {analysisResult?.genreDistribution ? Object.keys(analysisResult.genreDistribution).length : '-'}
            </p>
            <p className="text-sm text-spotify-lightgray mt-1">
              {analysisResult ? 'Unique genres detected' : 'Run analysis to see stats'}
            </p>
          </CardContent>
        </Card>
      </div>

      {error && (
        <Card variant="highlighted" className="border-red-600">
          <CardContent className="p-4">
            <p className="text-red-500">{error}</p>
          </CardContent>
        </Card>
      )}

      {analysisResult?.groupingSuggestions && analysisResult.groupingSuggestions.length > 0 && (
        <Card>
          <CardHeader>
            <div className="flex items-center justify-between">
              <div className="flex items-center space-x-3">
                <div className="p-2 bg-orange-600 rounded-lg">
                  <Layers className="w-5 h-5 text-white" />
                </div>
                <div>
                  <CardTitle>Genre Grouping Suggestions</CardTitle>
                  <p className="text-sm text-spotify-lightgray mt-1">
                    Reduce playlists by grouping similar genres together
                  </p>
                </div>
              </div>
              <div className="flex gap-2">
                <Button variant="ghost" size="sm" onClick={enableAllGroups}>
                  Enable All
                </Button>
                <Button variant="ghost" size="sm" onClick={disableAllGroups}>
                  Disable All
                </Button>
              </div>
            </div>
          </CardHeader>
          <CardContent className="space-y-3">
            {enabledGroups.length > 0 && (
              <div className="bg-green-900/20 border border-green-600 rounded-lg p-3 mb-4">
                <p className="text-green-400 text-sm">
                  {enabledGroups.length} group{enabledGroups.length !== 1 ? 's' : ''} enabled.
                  Sub-genres will be merged into parent playlists when sorting.
                </p>
              </div>
            )}
            {analysisResult.groupingSuggestions.map((suggestion) => (
              <div
                key={suggestion.parentGenre}
                className={`border rounded-lg p-4 transition-all ${
                  enabledGroups.includes(suggestion.parentGenre)
                    ? 'border-spotify-green bg-spotify-green/10'
                    : 'border-spotify-lightgray/30 bg-spotify-darkgray/50'
                }`}
              >
                <div className="flex items-center justify-between">
                  <div className="flex items-center space-x-3">
                    <button
                      onClick={() => toggleGroup(suggestion.parentGenre)}
                      className={`w-5 h-5 rounded border-2 flex items-center justify-center transition-all ${
                        enabledGroups.includes(suggestion.parentGenre)
                          ? 'bg-spotify-green border-spotify-green'
                          : 'border-spotify-lightgray/50'
                      }`}
                    >
                      {enabledGroups.includes(suggestion.parentGenre) && (
                        <svg className="w-3 h-3 text-white" viewBox="0 0 12 12">
                          <path
                            fill="currentColor"
                            d="M10 3L4.5 8.5 2 6l-.75.75 3.25 3.25 6.25-6.25z"
                          />
                        </svg>
                      )}
                    </button>
                    <div>
                      <h4 className="font-semibold text-white">{suggestion.parentGenre}</h4>
                      <p className="text-sm text-spotify-lightgray">
                        {suggestion.totalTracks} tracks · {suggestion.playlistsToMerge} sub-genres
                      </p>
                    </div>
                  </div>
                  <button
                    onClick={() => toggleExpanded(suggestion.parentGenre)}
                    className="text-spotify-lightgray hover:text-white transition-colors"
                  >
                    {expandedGroups.has(suggestion.parentGenre) ? (
                      <ChevronUp className="w-5 h-5" />
                    ) : (
                      <ChevronDown className="w-5 h-5" />
                    )}
                  </button>
                </div>
                {expandedGroups.has(suggestion.parentGenre) && (
                  <div className="mt-3 pt-3 border-t border-spotify-lightgray/20">
                    <p className="text-xs text-spotify-lightgray mb-2">Sub-genres that will be merged:</p>
                    <div className="flex flex-wrap gap-2">
                      {suggestion.childGenres.map((child) => (
                        <Badge key={child} variant="default">
                          {child} ({analysisResult.genreDistribution[child] || 0})
                        </Badge>
                      ))}
                    </div>
                  </div>
                )}
              </div>
            ))}
          </CardContent>
        </Card>
      )}

      <Card>
        <CardHeader>
          <CardTitle>Library Analysis</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <p className="text-spotify-lightgray">
            Analyze your entire Spotify library to identify genres and prepare for sorting.
            This process may take a few minutes depending on your library size.
          </p>

          {showProgress && (
            <ProgressBar
              current={currentEvent.current || 0}
              total={currentEvent.total || 100}
              message={currentEvent.message}
            />
          )}

          <Button
            onClick={handleAnalyze}
            disabled={isAnalyzing}
            isLoading={isAnalyzing}
            className="gap-2"
          >
            <Play className="w-4 h-4" />
            {isAnalyzing ? 'Analyzing...' : 'Start Analysis'}
          </Button>
        </CardContent>
      </Card>

      {events.length > 0 && <LiveLog events={events} />}

      {analysisResult && (
        <div className="flex justify-end">
          <Button onClick={() => navigate('/changes')} className="gap-2">
            Continue to Changes
            <ChevronDown className="w-4 h-4 rotate-[-90deg]" />
          </Button>
        </div>
      )}
    </div>
  );
}
