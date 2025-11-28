import { useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/Card';
import { Button } from '@/components/ui/Button';
import { GenreGrid } from '@/components/genres/GenreGrid';
import { useCreateSortPlan } from '@/hooks/useApi';
import { RefreshCw } from 'lucide-react';

export function Genres() {
  const [plan, setPlan] = useState<any>(null);
  const createPlan = useCreateSortPlan();

  const handleLoadGenres = async () => {
    try {
      const result = await createPlan.mutateAsync();
      setPlan(result);
    } catch (error) {
      console.error('Failed to load genres:', error);
    }
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-white mb-2">Genres</h1>
          <p className="text-spotify-lightgray">
            View all genres detected in your library
          </p>
        </div>

        <Button
          onClick={handleLoadGenres}
          isLoading={createPlan.isPending}
          className="gap-2"
        >
          <RefreshCw className="w-4 h-4" />
          {plan ? 'Refresh' : 'Load Genres'}
        </Button>
      </div>

      {!plan && !createPlan.isPending && (
        <Card>
          <CardHeader>
            <CardTitle>No Data</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-spotify-lightgray mb-4">
              Click "Load Genres" to analyze your library and view genre statistics.
            </p>
          </CardContent>
        </Card>
      )}

      {plan && plan.genreStats && (
        <>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <Card variant="highlighted">
              <CardHeader>
                <CardTitle>Total Genres</CardTitle>
              </CardHeader>
              <CardContent>
                <p className="text-3xl font-bold text-white">
                  {plan.genreStats.length}
                </p>
              </CardContent>
            </Card>

            <Card variant="highlighted">
              <CardHeader>
                <CardTitle>New Genres</CardTitle>
              </CardHeader>
              <CardContent>
                <p className="text-3xl font-bold text-white">
                  {plan.genreStats.filter((g: any) => g.isNew).length}
                </p>
              </CardContent>
            </Card>
          </div>

          <GenreGrid genres={plan.genreStats} />
        </>
      )}

      {plan && plan.uncategorizedTracks && plan.uncategorizedTracks.length > 0 && (
        <Card>
          <CardHeader>
            <CardTitle className="text-yellow-500">
              Uncategorized Tracks ({plan.uncategorizedTracks.length})
            </CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-spotify-lightgray">
              Some tracks could not be categorized into a genre. These will be skipped during sorting.
            </p>
          </CardContent>
        </Card>
      )}
    </div>
  );
}
