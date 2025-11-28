import { useState } from 'react';
import { AlertTriangle, Play } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/Card';
import { Button } from '@/components/ui/Button';
import { Badge } from '@/components/ui/Badge';
import { ChangeSummary } from '@/components/changes/ChangeSummary';
import { ChangeDiff } from '@/components/changes/ChangeDiff';
import { LiveLog } from '@/components/progress/LiveLog';
import { useCreateSortPlan, useExecuteSortPlan } from '@/hooks/useApi';
import { useUIStore } from '@/stores/uiStore';
import { useSSE } from '@/hooks/useSSE';

export function Changes() {
  const [plan, setPlan] = useState<any>(null);
  const [executionEndpoint, setExecutionEndpoint] = useState<string | null>(null);
  const isDryRun = useUIStore((state) => state.isDryRun);
  const createPlan = useCreateSortPlan();
  const executePlan = useExecuteSortPlan();

  const { events, isConnected } = useSSE(executionEndpoint, {
    onComplete: () => {
      setExecutionEndpoint(null);
      // Refresh the plan after execution
      handleGeneratePlan();
    },
  });

  const handleGeneratePlan = async () => {
    try {
      const result = await createPlan.mutateAsync();
      setPlan(result);
    } catch (error) {
      console.error('Failed to generate plan:', error);
    }
  };

  const handleExecute = async () => {
    if (!plan || isDryRun) return;

    try {
      await executePlan.mutateAsync(plan.id);
      setExecutionEndpoint(`/api/events?type=execution&planId=${plan.id}`);
    } catch (error) {
      console.error('Failed to execute plan:', error);
    }
  };

  const isExecuting = executePlan.isPending || isConnected;

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-white mb-2">Changes</h1>
          <p className="text-spotify-lightgray">
            Review and apply changes to your playlists
          </p>
        </div>

        <Button
          onClick={handleGeneratePlan}
          isLoading={createPlan.isPending}
        >
          {plan ? 'Refresh Plan' : 'Generate Plan'}
        </Button>
      </div>

      {!isDryRun && (
        <Card variant="highlighted" className="border-red-600">
          <CardContent className="flex items-start space-x-3 p-4">
            <AlertTriangle className="w-5 h-5 text-red-500 mt-0.5 flex-shrink-0" />
            <div>
              <h4 className="text-white font-bold mb-1">Live Mode Active</h4>
              <p className="text-sm text-spotify-lightgray">
                Changes will be applied directly to your Spotify account.
                Toggle "Dry Run" in the header to preview changes first.
              </p>
            </div>
          </CardContent>
        </Card>
      )}

      {!plan && !createPlan.isPending && (
        <Card>
          <CardHeader>
            <CardTitle>No Plan Generated</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-spotify-lightgray mb-4">
              Click "Generate Plan" to analyze your library and preview changes.
            </p>
          </CardContent>
        </Card>
      )}

      {plan && (
        <>
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-3">
              <Badge variant={isDryRun ? 'success' : 'danger'}>
                {isDryRun ? 'Preview Mode' : 'Live Mode'}
              </Badge>
              <span className="text-sm text-spotify-lightgray">
                Plan ID: {plan.id}
              </span>
            </div>

            <Button
              onClick={handleExecute}
              disabled={isDryRun || isExecuting}
              isLoading={isExecuting}
              variant="primary"
              className="gap-2"
            >
              <Play className="w-4 h-4" />
              {isExecuting ? 'Executing...' : 'Execute Changes'}
            </Button>
          </div>

          <ChangeSummary
            tracksToAdd={plan.tracksToAdd?.length || 0}
            tracksToRemove={plan.tracksToRemove?.length || 0}
            playlistsToCreate={plan.playlistsToCreate?.length || 0}
          />

          {plan.playlistsToCreate && plan.playlistsToCreate.length > 0 && (
            <Card>
              <CardHeader>
                <CardTitle>New Playlists to Create</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="flex flex-wrap gap-2">
                  {plan.playlistsToCreate.map((name: string) => (
                    <Badge key={name} variant="info">
                      {name}
                    </Badge>
                  ))}
                </div>
              </CardContent>
            </Card>
          )}

          <ChangeDiff
            tracksToAdd={plan.tracksToAdd || []}
            tracksToRemove={plan.tracksToRemove || []}
          />

          {events.length > 0 && <LiveLog events={events} />}
        </>
      )}
    </div>
  );
}
