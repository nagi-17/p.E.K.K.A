const SERVER_URL='http://localhost:8080';

export async function apiClient(path, options) {
    if (options===undefined) {
        options={};
    }

    const token=localStorage.getItem('user_jwt');

    const headers={'Content-Type': 'application/json'};
    if (token!==null) {
        headers['Authorization']='Bearer '+token;
    }
    if (options.headers!==undefined) {
        for (const key in options.headers) {
            headers[key]=options.headers[key];
        }
    }
    options.headers=headers

    const response=await fetch(SERVER_URL+path, options);

    if (!response.ok) {
        let errorMessage = `API request failed with status ${response.status}`;
        try {
            const bodyText = await response.text();
            try {
                const errorData = JSON.parse(bodyText);
                errorMessage = errorData.error || errorData.message || errorMessage;
            } catch {
                if (bodyText && bodyText.trim()) {
                    errorMessage = bodyText.trim();
                }
            }
        } catch (err) {}
        throw new Error(errorMessage);
    }

    return response.json();
}