import { create } from 'zustand';
import { persist } from 'zustand/middleware';

interface UIState {
  isDryRun: boolean;
  toggleDryRun: () => void;
  setDryRun: (value: boolean) => void;
}

export const useUIStore = create<UIState>()(
  persist(
    (set) => ({
      isDryRun: true,
      toggleDryRun: () => set((state) => ({ isDryRun: !state.isDryRun })),
      setDryRun: (value) => set({ isDryRun: value }),
    }),
    {
      name: 'ui-storage',
    }
  )
);
