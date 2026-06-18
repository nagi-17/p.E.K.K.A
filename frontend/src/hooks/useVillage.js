import { useState, useEffect } from 'react';
import { getVillage } from '../api/village';
import { useVillageStore } from '../store/villageStore';

export function useVillage() {
    const [loading, setLoading]=useState(true);
    const [error, setError]=useState(null);
    const setVillage=useVillageStore(function(state) { 
        return state.setVillage;
    });

    useEffect(function() {
        async function fetchVillageData() {
            try {
                const data=await getVillage();
                setVillage(data);
                setLoading(false);
            }
            catch (err) {
                setError(err.message);
                setLoading(false);
            }
        }

        fetchVillageData();
    }, [setVillage]);

    return {loading, error};
}