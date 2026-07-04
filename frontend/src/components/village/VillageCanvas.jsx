import { useEffect, useRef } from 'react';
import { useGameInit } from '../../hooks/useGameInit';
import { GameApp } from '../../game/core/GameApp';
import { useVillageStore } from '../../store/villageStore';
import styles from './VillageCanvas.module.css';

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
        return <div className={styles.tempText}>Loading Village</div>;
    }
    if (error) {
        return <div className={styles.tempText}>Error: {error}</div>;
    }

    return <div ref={canvasRef} className={styles.canvasContainer}></div>;
}
