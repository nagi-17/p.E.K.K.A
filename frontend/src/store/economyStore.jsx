import { create } from 'zustand';
import { collectResources } from '../api/economy';
import { useVillageStore } from './villageStore';

export const useEconomyStore=create((set) => ({
    isCollecting: false,
    error: null,

    collect: async (buildingID) => {
        set({ isCollecting: true, error: null });
        try {
            await collectResources(buildingID);
            await useVillageStore.getState().refreshGameData();
            set({ isCollecting: false });
            return true;
        } catch (err) {
            set({ error: err.message, isCollecting: false });
            throw err;
        }
    }
}));
