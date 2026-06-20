import { create } from 'zustand';

export const useUiStore=create((set) => ({
    selectedBuildingId: null,
    isShopOpen: false,
    isArmyOpen: false,
    placementModeId: null, 
    
    moveModeId: null,

    selectBuilding: (id) => set({ selectedBuildingId: id, isShopOpen: false, isArmyOpen: false }),
    clearSelection: () => set({ selectedBuildingId: null }),
    toggleShop: () => set((state) => ({ isShopOpen: !state.isShopOpen, selectedBuildingId: null, isArmyOpen: false })),
    toggleArmy: () => set((state) => ({ isArmyOpen: !state.isArmyOpen, isShopOpen: false, selectedBuildingId: null })),
    
    startPlacement: (buildingDataId) => set({ placementModeId: buildingDataId, isShopOpen: false, selectedBuildingId: null, moveModeId: null, isArmyOpen: false }),
    cancelPlacement: () => set({ placementModeId: null }),

    startMove: (instanceId) => set({ moveModeId: instanceId, selectedBuildingId: null, isShopOpen: false, isArmyOpen: false }),
    cancelMove: () => set({ moveModeId: null })
}));