import { create } from 'zustand';
import { persist } from 'zustand/middleware';

interface UIState {
  isDryRun: boolean;
  toggleDryRun: () => void;
  setDryRun: (value: boolean) => void;
  enabledGroups: string[];
  setEnabledGroups: (groups: string[]) => void;
  toggleGroup: (group: string) => void;
  enableAllGroups: (groups: string[]) => void;
  disableAllGroups: () => void;
  disabledPlaylists: string[];
  togglePlaylist: (playlist: string) => void;
  enableAllPlaylists: () => void;
  disableAllPlaylists: (playlists: string[]) => void;
}

export const useUIStore = create<UIState>()(
  persist(
    (set) => ({
      isDryRun: true,
      toggleDryRun: () => set((state) => ({ isDryRun: !state.isDryRun })),
      setDryRun: (value) => set({ isDryRun: value }),
      enabledGroups: [],
      setEnabledGroups: (groups) => set({ enabledGroups: groups }),
      toggleGroup: (group) => set((state) => {
        const newGroups = state.enabledGroups.includes(group)
          ? state.enabledGroups.filter(g => g !== group)
          : [...state.enabledGroups, group];
        return { enabledGroups: newGroups };
      }),
      enableAllGroups: (groups) => set({ enabledGroups: groups }),
      disableAllGroups: () => set({ enabledGroups: [] }),
      disabledPlaylists: [],
      togglePlaylist: (playlist) => set((state) => {
        const newDisabled = state.disabledPlaylists.includes(playlist)
          ? state.disabledPlaylists.filter(p => p !== playlist)
          : [...state.disabledPlaylists, playlist];
        return { disabledPlaylists: newDisabled };
      }),
      enableAllPlaylists: () => set({ disabledPlaylists: [] }),
      disableAllPlaylists: (playlists) => set({ disabledPlaylists: playlists }),
    }),
    {
      name: 'ui-storage',
    }
  )
);
