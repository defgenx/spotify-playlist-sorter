import { useMutation, useQuery } from '@tanstack/react-query';
import { api } from '@/lib/api';
import { useUIStore } from '@/stores/uiStore';

export function useLibraryStats() {
  return useQuery({
    queryKey: ['library', 'stats'],
    queryFn: () => api.getLibraryStats(),
  });
}

export function useCreateSortPlan() {
  const isDryRun = useUIStore((state) => state.isDryRun);
  const enabledGroups = useUIStore((state) => state.enabledGroups);
  const disabledPlaylists = useUIStore((state) => state.disabledPlaylists);

  return useMutation({
    mutationFn: () => api.createSortPlan(isDryRun, enabledGroups, disabledPlaylists),
  });
}

export function useExecuteSortPlan() {
  const isDryRun = useUIStore((state) => state.isDryRun);
  const enabledGroups = useUIStore((state) => state.enabledGroups);
  const disabledPlaylists = useUIStore((state) => state.disabledPlaylists);

  return useMutation({
    mutationFn: () => api.executeSortPlan(isDryRun, enabledGroups, disabledPlaylists),
  });
}
