import { create } from 'zustand';

export const useVillageStore=create((set)=>({
    buildings: [], townHallLevel: 1, isLoaded: false,
    setVillage: (data)=>set({buildings: Array.isArray(data)?data:[], isLoaded: true})
}));