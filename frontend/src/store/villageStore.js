import { create } from 'zustand';

export const useVillageStore=create((set)=>({
    buildings: [], townHallLevel: 1, isLoaded: false,
    setVillage: (data)=>set({buildings: data.buildings || [], townHallLevel: data.townHallLevel || 1, isLoaded: true})
}));