import { create } from 'zustand';
import { getArmy, trainTroop, upgradeTroop, finishUpgradeTroop, discardTroops } from '../api/army';
import { useVillageStore } from './villageStore';

export const useArmyStore=create((set, get) => ({
    troops: [],
    isLoading: false,
    error: null,

    fetchArmy: async () => {
        set({ isLoading: true, error: null });
        try {
            const data=await getArmy();
            set({ troops: Array.isArray(data) ? data : [], isLoading: false });
        } catch (err) {
            set({ error: err.message, isLoading: false });
        }
    },

    train: async (troopType, quantity) => {
        set({ isLoading: true, error: null });
        try {
            await trainTroop(troopType, quantity);
            await get().fetchArmy();
            await useVillageStore.getState().refreshGameData();
            set({ isLoading: false });
            return true;
        } catch (err) {
            set({ error: err.message, isLoading: false });
            throw err;
        }
    },

    startUpgrade: async (troopType) => {
        set({ isLoading: true, error: null });
        try {
            await upgradeTroop(troopType);
            await get().fetchArmy();
            await useVillageStore.getState().refreshGameData();
            set({ isLoading: false });
            return true;
        } catch (err) {
            set({ error: err.message, isLoading: false });
            throw err;
        }
    },

    finishUpgrade: async (troopType) => {
        set({ isLoading: true, error: null });
        try {
            await finishUpgradeTroop(troopType);
            await get().fetchArmy();
            await useVillageStore.getState().refreshGameData();
            set({ isLoading: false });
            return true;
        } catch (err) {
            set({ error: err.message, isLoading: false });
            throw err;
        }
    },

    discard: async (troopType, quantity) => {
        set({ isLoading: true, error: null });
        try {
            await discardTroops(troopType, quantity);
            await get().fetchArmy();
            await useVillageStore.getState().refreshGameData();
            set({ isLoading: false });
            return true;
        } catch (err) {
            set({ error: err.message, isLoading: false });
            throw err;
        }
    }
}));
