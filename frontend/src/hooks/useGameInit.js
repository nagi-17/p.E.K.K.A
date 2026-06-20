import { useState, useEffect } from 'react';
import { getVillage, getPlayerInfo } from '../api/village';
import { useVillageStore } from '../store/villageStore';

export function useGameInit() {
    const [loading, setLoading]=useState(true);
    const [error, setError]=useState(null);
    const setGameData=useVillageStore((state) => state.setGameData);

    useEffect(() => {
        let isMounted=true;
        async function fetchAllData() {
            try {
                const [villageRes, playerRes]=await Promise.all([getVillage(),getPlayerInfo()]);

                if (isMounted) {
                    setGameData(villageRes, playerRes);
                    setLoading(false);
                }
            }
            catch (err) {
                if (isMounted) {
                    setError(err.message);
                    setLoading(false);
                }
            }
        }

        fetchAllData();
        return () => {isMounted=false;};
    }, [setGameData]);

    return { loading, error };
}