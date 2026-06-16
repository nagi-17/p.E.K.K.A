const DATABASE_URL='http://localhost:8080';

export async function apiClient(path, query) {
    if (query===undefined) {
        query={};
    }

    const token=localStorage.getItem('user_jwt');

    const headers={'Content-Type': 'application/json'};
    if (token!==null) {
        headers['Authorization']='Bearer '+token;
    }
    if (query.headers!==undefined) {
        for (const key in query.headers) {
            headers[key]=query.headers[key];
        }
    }
    query.headers=headers

    const response=await fetch(DATABASE_URL+path, query);

    if (!response.ok) {
        let errorData={};
        try {
            errorData=await response.json();
        }
        catch(err){}
        throw new Error(errorData.message||'API request failed');
    }

    return response.json();
}