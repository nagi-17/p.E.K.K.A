import { apiClient } from "./client";

export async function loginUser(username, password) {
    return apiClient('/login', {method: 'POST', body: JSON.stringify({username, password})});
}

export async function regUser(username, email, password) {
    return apiClient('/register', {method: 'POST', body: JSON.stringify({username, password, email})});
}