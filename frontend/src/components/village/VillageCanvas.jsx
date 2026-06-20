import { useEffect, useRef } from 'react';
import { useGameInit } from '../../hooks/useGameInit';
import { GameApp } from '../../game/core/GameApp';
import { useVillageStore } from '../../store/villageStore';

export default function VillageCanvas() {
    const canvasRef=useRef(null);
    const gameAppRef=useRef(null);
    const {loading, error}=useGameInit();
    const buildings=useVillageStore(function(state){
        return state.buildings;
    });

    useEffect(()=>{
        let isMounted=true;

        async function initGame() {
            if (!loading && canvasRef.current && !gameAppRef.current) {
                const game=new GameApp();
                await game.init(canvasRef.current);
                
                if (isMounted) {
                    game.mount(canvasRef.current);
                    gameAppRef.current=game;
                    game.scene.loadVillage(buildings);
                }
                else {
                    game.destroy();
                }
            }
        }

        initGame();
        return ()=>{
            isMounted=false;
            if (gameAppRef.current) {
                gameAppRef.current.destroy();
                gameAppRef.current=null;
            }
        };
    }, [loading]);

    useEffect(() => {
        if (gameAppRef.current && gameAppRef.current.scene && !loading) {
            gameAppRef.current.scene.loadVillage(buildings);
        }
    }, [buildings, loading]);

    if (loading) {
        return <div style={styles.temp_text}>Loading Village</div>;
    }
    if (error) {
        return <div style={styles.temp_text}>Error: {error}</div>;
    }

    return <div ref={canvasRef} style={styles.canvasContainer}></div>;
}

const styles={
    canvasContainer: { 
        width: '100%', 
        height: '100%', 
        overflow: 'hidden'
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