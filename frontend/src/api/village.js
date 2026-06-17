import { apiClient } from "./client";

export async function playerInfo() {
    return apiClient('/user/load', {method: 'GET'});
}

export async function getVillage() {
    return apiClient('/village', {method: 'GET'});   
}

export async function placeBuilding(bType, x, y) {
    return apiClient('/village/build', {method: 'POST', body: JSON.stringify({building_type: bType, pos_x: x, pos_y: y})});
}

export async function moveBuilding(buildingID, x, y) {
    return apiClient('/village/move', {method: 'PUT', body: JSON.stringify({building_id: buildingID, new_x: x, new_y: y})});
}

export async function startUpgrade(buildingID) {
    return apiClient('/village/upgrade/start', {method: 'POST', body: JSON.stringify({building_id: buildingID})});
}
