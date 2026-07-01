import { apiClient } from "./client";

export async function findOpponent() {
    return apiClient('/battle/matchmake', { method: 'GET' });
}

export async function launchAttack(defenderID, events) {
    return apiClient('/battle/attack', { method: 'POST', body: JSON.stringify({defender_id: defenderID, events: events})});
}