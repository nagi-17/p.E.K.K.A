import { apiClient } from "./client";

export async function getArmy() {
    return apiClient('/village/army', {method: 'GET'});
}

export async function trainTroop(troopType, quantity) {
    return apiClient('/village/army/train', {method: 'POST', body: JSON.stringify({troop_type: troopType, quantity: quantity})});
}

export async function upgradeTroop(troopType) {
    return apiClient('/village/lab/troop/upgrade/start', {method: 'POST', body: JSON.stringify({troop_type: troopType})});
}

export async function finishUpgradeTroop(troopType) {
    return apiClient('/village/lab/troop/upgrade/finish', {method: 'POST', body: JSON.stringify({troop_type: troopType})});
}

export async function discardTroops(troopType, quantity) {
    return apiClient('/village/army/discard', {method: 'POST', body: JSON.stringify({troop_type: troopType, quantity: quantity})});
}