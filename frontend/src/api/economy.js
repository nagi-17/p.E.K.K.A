import { apiClient } from "./client";

export async function collectResources(buildingID) {
    return apiClient('/village/collect', {method: 'POST', body: JSON.stringify({building_id: buildingID})});
}