import { useState } from 'react';
import { AlertTriangle, Play, Layers, ChevronDown, ChevronUp, ListMusic } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/Card';
import { Button } from '@/components/ui/Button';
import { Badge } from '@/components/ui/Badge';
import { ChangeSummary } from '@/components/changes/ChangeSummary';
import { ChangeDiff } from '@/components/changes/ChangeDiff';
import { LiveLog } from '@/components/progress/LiveLog';
import { useCreateSortPlan, useExecuteSortPlan } from '@/hooks/useApi';
import { useUIStore } from '@/stores/uiStore';
import { useSSE } from '@/hooks/useSSE';
import type { GroupSuggestion } from '@/lib/types';

export function Changes() {
  const [plan, setPlan] = useState<any>(null);
  const [groupingSuggestions, setGroupingSuggestions] = useState<GroupSuggestion[]>([]);
  const [expandedGroups, setExpandedGroups] = useState<Set<string>>(new Set());
  const [executionEndpoint, setExecutionEndpoint] = useState<string | null>(null);
  const isDryRun = useUIStore((state) => state.isDryRun);
  const enabledGroups = useUIStore((state) => state.enabledGroups);
  const toggleGroup = useUIStore((state) => state.toggleGroup);
  const enableAllGroupsStore = useUIStore((state) => state.enableAllGroups);
  const disableAllGroupsStore = useUIStore((state) => state.disableAllGroups);
  const disabledPlaylists = useUIStore((state) => state.disabledPlaylists);
  const togglePlaylist = useUIStore((state) => state.togglePlaylist);
  const enableAllPlaylists = useUIStore((state) => state.enableAllPlaylists);
  const disableAllPlaylistsStore = useUIStore((state) => state.disableAllPlaylists);
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
      if (result.groupingSuggestions) {
        setGroupingSuggestions(result.groupingSuggestions);
      }
    } catch (error) {
      console.error('Failed to generate plan:', error);
    }
  };

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
    enableAllGroupsStore(groupingSuggestions.map(s => s.parentGenre));
  };

  const disableAllGroups = () => {
    disableAllGroupsStore();
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

      {groupingSuggestions.length > 0 && (
        <Card>
          <CardHeader>
            <div className="flex items-center justify-between">
              <div className="flex items-center space-x-3">
                <div className="p-2 bg-orange-600 rounded-lg">
                  <Layers className="w-5 h-5 text-white" />
                </div>
                <div>
                  <CardTitle>Genre Grouping</CardTitle>
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
                  Click "Refresh Plan" to regenerate with grouping applied.
                </p>
              </div>
            )}
            {groupingSuggestions.map((suggestion) => (
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
                      {suggestion.childGenres.map((child: string) => (
                        <Badge key={child} variant="default">
                          {child}
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
                <div className="flex items-center justify-between">
                  <div className="flex items-center space-x-3">
                    <div className="p-2 bg-blue-600 rounded-lg">
                      <ListMusic className="w-5 h-5 text-white" />
                    </div>
                    <div>
                      <CardTitle>Playlists to Create</CardTitle>
                      <p className="text-sm text-spotify-lightgray mt-1">
                        Click to toggle which playlists will be created
                      </p>
                    </div>
                  </div>
                  <div className="flex gap-2">
                    <Button variant="ghost" size="sm" onClick={enableAllPlaylists}>
                      Enable All
                    </Button>
                    <Button variant="ghost" size="sm" onClick={() => disableAllPlaylistsStore(plan.playlistsToCreate)}>
                      Disable All
                    </Button>
                  </div>
                </div>
              </CardHeader>
              <CardContent>
                {disabledPlaylists.length > 0 && (
                  <div className="bg-yellow-900/20 border border-yellow-600 rounded-lg p-3 mb-4">
                    <p className="text-yellow-400 text-sm">
                      {disabledPlaylists.length} playlist{disabledPlaylists.length !== 1 ? 's' : ''} disabled.
                      Click "Refresh Plan" to update the plan.
                    </p>
                  </div>
                )}
                <div className="flex flex-wrap gap-2">
                  {plan.playlistsToCreate.map((name: string) => {
                    const isEnabled = !disabledPlaylists.includes(name);
                    return (
                      <button
                        key={name}
                        onClick={() => togglePlaylist(name)}
                        className={`inline-flex items-center gap-2 px-3 py-1.5 rounded-full text-sm font-medium transition-all ${
                          isEnabled
                            ? 'bg-blue-600 text-white hover:bg-blue-500'
                            : 'bg-spotify-darkgray text-spotify-lightgray/50 line-through hover:bg-spotify-darkgray/80'
                        }`}
                      >
                        <span className={`w-3 h-3 rounded-full border ${
                          isEnabled ? 'bg-white border-white' : 'border-spotify-lightgray/50'
                        }`} />
                        {name}
                      </button>
                    );
                  })}
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
