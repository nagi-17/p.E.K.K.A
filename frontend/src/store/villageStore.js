import { create } from 'zustand';
import {getVillage, getPlayerInfo} from '../api/village';

export const useVillageStore=create((set) => ({
    buildings: [], 
    playerStats: null,
    isLoaded: false,
    setGameData: (buildingsData, statsData) => set({ 
        buildings: Array.isArray(buildingsData) ? buildingsData : [], 
        playerStats: statsData || null,
        isLoaded: true 
    }),
    setVillage: (data) => set({ buildings: Array.isArray(data) ? data : [] }),
    setPlayerStats: (stats) => set({ playerStats: stats }),
    refreshGameData: async () => {
        try {
            const [buildings, playerStats]=await Promise.all([getVillage(), getPlayerInfo()]);
            set({ buildings, playerStats });
        }
        catch (err) {
            console.error("Failed to refresh game data:", err);
        }
    }
}));