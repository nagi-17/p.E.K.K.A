import { apiClient } from "./client";
import { mapBuildingData, mapPlayerStats } from "../utils/mapper";

export async function getPlayerInfo() {
    const rawData=await apiClient('/user/load', {method: 'GET'});
    return mapPlayerStats(rawData);
}

export async function getVillage(playerID) {
    const url = playerID ? `/village?player_id=${playerID}` : '/village';
    const rawData=await apiClient(url, {method: 'GET'});
    if (Array.isArray(rawData)) {
        return rawData.map(mapBuildingData);
    }
    return [];
}

export async function placeBuilding(data) {
    return apiClient('/village/build', {method: 'POST', body: JSON.stringify(data)});
}

export async function moveBuilding(buildingID, x, y) {
    return apiClient('/village/move', {method: 'PUT', body: JSON.stringify({building_id: buildingID, new_x: x, new_y: y})});
}

export async function startUpgrade(buildingID) {
    return apiClient('/village/upgrade/start', {method: 'POST', body: JSON.stringify({building_id: buildingID})});
}

export async function finishUpgrade(buildingID) {
    return apiClient('/village/upgrade/finish', {method: 'POST', body: JSON.stringify({ building_id: buildingID }),});
}

export async function collectResource(buildingID) {
    return apiClient('/village/collect', {method: 'POST', body: JSON.stringify({ building_id: buildingID })});
}

export async function cancelUpgrade(buildingID) {
    return apiClient('/village/upgrade/cancel', {method: 'POST', body: JSON.stringify({ building_id: buildingID })});
}
