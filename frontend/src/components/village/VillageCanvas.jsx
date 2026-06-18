import { useEffect, useRef } from 'react';
import { useVillage } from '../../hooks/useVillage';

export default function VillageCanvas() {
    const canvasRef=useRef(null);

    const {loading, error}=useVillage();

    useEffect(()=>{
        if (!loading && canvasRef.current) {
            console.log("Data loaded");
        }
        return ()=>{};
    }, [loading]);

    if (loading) return <div style={styles.temp_text}>Loading Village...</div>;
    if (error) return <div style={styles.temp_text}>Error: {error}</div>;

    return <div ref={canvasRef} style={styles.canvasContainer}></div>;
}

const styles = {
    canvasContainer: { 
        width: '100%', 
        height: '100%', 
        backgroundColor: '#7ec850'
    }, 
    temp_text: { 
        display: 'flex', 
        justifyContent: 'center', 
        alignItems: 'center', 
        height: '100%', 
        color: 'white', 
        fontFamily: '"Luckiest Guy", cursive',
        fontSize: '2rem',
        textShadow: '2px 2px 4px #000'
    }
};