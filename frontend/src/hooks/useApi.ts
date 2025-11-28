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

  return useMutation({
    mutationFn: () => api.createSortPlan(isDryRun),
  });
}

export function useExecuteSortPlan() {
  return useMutation({
    mutationFn: (planId: string) => api.executeSortPlan(planId),
  });
}
