import { apiClient } from "./client";

export async function trainTroop(troopType, quantity) {
    return apiClient('/village/army/train', {method: 'POST', body: JSON.stringify({troop_type: troopType, quantity: quantity})});
}

export async function upgradeTroop(troopType) {
    return apiClient('/village/lab/troop/upgrade/start', {method: 'POST', body: JSON.stringify({troop_type: troopType})});
}