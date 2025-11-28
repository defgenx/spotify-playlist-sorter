import { useState } from 'react';
import { Play, Music, Disc, ListMusic } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/Card';
import { Button } from '@/components/ui/Button';
import { ProgressBar } from '@/components/progress/ProgressBar';
import { LiveLog } from '@/components/progress/LiveLog';
import { useSSE } from '@/hooks/useSSE';

interface AnalysisResult {
  totalLikedSongs: number;
  playlists: { id: string; name: string; trackCount: number }[];
  genreBreakdown: Record<string, number>;
}

export function Dashboard() {
  const [isAnalyzing, setIsAnalyzing] = useState(false);
  const [analysisEndpoint, setAnalysisEndpoint] = useState<string | null>(null);
  const [analysisResult, setAnalysisResult] = useState<AnalysisResult | null>(null);
  const [error, setError] = useState<string | null>(null);

  const { events } = useSSE(analysisEndpoint, {
    onComplete: () => {
      setIsAnalyzing(false);
      setAnalysisEndpoint(null);
    },
  });

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
              {analysisResult?.genreBreakdown ? Object.keys(analysisResult.genreBreakdown).length : '-'}
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
    </div>
  );
}
