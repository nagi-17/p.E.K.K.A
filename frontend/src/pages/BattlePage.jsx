import { useEffect, useRef, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import * as PIXI from 'pixi.js';
import { findOpponent, launchAttack } from '../api/battle';
import { getVillage } from '../api/village';
import { getArmy } from '../api/army';
import { createViewport } from '../game/core/Viewport';
import { AssetLoader } from '../game/core/AssetLoader';
import { BattleScene } from '../game/scenes/BattleScene';
import BattleResultModel from '../components/battle/BattleResultModel';
import { TROOP_DEFS } from '../game/gameConfig';
import styles from './BattlePage.module.css';

export default function BattlePage() {
    const navigate=useNavigate();

    const [matchState, setMatchState]=useState('search');
    const [loading, setLoading]=useState(false);
    const [error, setError]=useState(null);
    const [opponent, setOpponent]=useState(null);
    const [opponentBuildings, setOpponentBuildings]=useState([]);
    
    const [deployRemaining, setDeployRemaining]=useState({});
    const [selectedTroopType, setSelectedTroopType]=useState(null);
    const [damagePercent, setDamagePercent]=useState(0);
    const [battleResult, setBattleResult]=useState(null);
    const [timer, setTimer]=useState(180);

    const [hasTrainedArmy, setHasTrainedArmy]=useState(false);
    const [armyCapacity, setArmyCapacity]=useState(0);

    const canvasContainerRef=useRef(null);
    const pixiAppRef=useRef(null);
    const battleSceneRef=useRef(null);
    const finalResultRef=useRef(null);
    const timerIntervalRef=useRef(null);

    const matchStateRef=useRef(matchState);
    matchStateRef.current=matchState;

    useEffect(() => {
        async function checkPlayerArmy() {
            setLoading(true);
            try {
                const trained=await getArmy();
                let totalQuantity=0;
                const remaining={};
                trained.forEach(t => {
                    if (t.quantity > 0) {
                        totalQuantity += t.quantity;
                        remaining[t.troop_type]=t.quantity;
                    }
                });

                setDeployRemaining(remaining);
                setHasTrainedArmy(totalQuantity > 0);

                const keys=Object.keys(remaining);
                if (keys.length > 0) {
                    setSelectedTroopType(keys[0]);
                }
            } catch (err) {
                console.error("Failed to load player army:", err);
            } finally {
                setLoading(false);
            }
        }
        checkPlayerArmy();

        return () => {
            if (pixiAppRef.current) {
                pixiAppRef.current.destroy(true, { children: true });
                pixiAppRef.current=null;
            }
            if (battleSceneRef.current) {
                battleSceneRef.current.destroy();
                battleSceneRef.current=null;
            }
        };
    }, []);

    const handleFindOpponent=async () => {
        setMatchState('searching');
        setError(null);
        try {
            await new Promise(r => setTimeout(r, 1500));
            const opp=await findOpponent();
            const buildings=await getVillage(opp.player_id);
            
            setOpponent(opp);
            setOpponentBuildings(buildings);
            setMatchState('found');

            const battleData = {
                opponent: opp,
                opponentBuildings: buildings,
                deployRemaining: deployRemaining,
                startTimestamp: null,
                deploymentEvents: []
            };
            localStorage.setItem('activeBattle', JSON.stringify(battleData));
        }
        catch (err) {
            setError(err.message || 'Failed to locate opponent.');
            setMatchState('search');
        }
    };

    useEffect(() => {
        let active=true;
        if (matchState !== 'found' && matchState !== 'battle' && matchState !== 'ended') return;
        if (!canvasContainerRef.current) return;

        if (pixiAppRef.current) {
            return () => {
                if (matchStateRef.current !== 'found' && matchStateRef.current !== 'battle' && matchStateRef.current !== 'ended') {
                    active=false;
                    if (pixiAppRef.current) {
                        pixiAppRef.current.destroy(true, { children: true });
                        pixiAppRef.current=null;
                    }
                    if (battleSceneRef.current) {
                        battleSceneRef.current.destroy();
                        battleSceneRef.current=null;
                    }
                }
            };
        }

        let app=null;
        let scene=null;

        async function initPixi() {
            app=new PIXI.Application();
            await app.init({
                resizeTo: canvasContainerRef.current,
                backgroundColor: 0x5a9a3b,
                resolution: window.devicePixelRatio || 1,
                autoDensity: true,
            });

            if (!active) {
                app.destroy(true, { children: true });
                return;
            }

            canvasContainerRef.current.appendChild(app.canvas);
            pixiAppRef.current=app;

            await AssetLoader.loadAssets();

            const viewport=createViewport(app, canvasContainerRef.current);
            app.stage.addChild(viewport);

            scene=new BattleScene(
                (finalDmg) => {
                    handleSimulationEnd(finalDmg);
                },
                (currentDmg) => {
                    setDamagePercent(currentDmg);
                }
            );

            scene.loadDefenderVillage(opponentBuildings);
            viewport.addChild(scene.container);
            battleSceneRef.current=scene;
            scene.deployRemaining=deployRemaining;

            const centerX=(40 * 32) / 2;
            const centerY=(40 * 32) / 2;
            viewport.moveCenter(centerX, centerY);

            if (selectedTroopType) {
                scene.setDeployableTroop(selectedTroopType, handleTroopDeployedOnCanvas);
            }
        }

        initPixi();

        return () => {
            if (matchStateRef.current !== 'found' && matchStateRef.current !== 'battle' && matchStateRef.current !== 'ended') {
                active=false;
                if (app) {
                    app.destroy(true, { children: true });
                    pixiAppRef.current=null;
                }
                if (scene) {
                    scene.destroy();
                    battleSceneRef.current=null;
                }
            }
        };
    }, [matchState, opponentBuildings]);

    useEffect(() => {
        if (battleSceneRef.current) {
            battleSceneRef.current.setDeployableTroop(selectedTroopType, handleTroopDeployedOnCanvas);
        }
    }, [selectedTroopType, matchState]);

    useEffect(() => {
        if (battleSceneRef.current) {
            battleSceneRef.current.deployRemaining=deployRemaining;
        }
    }, [deployRemaining]);

    const handleTroopDeployedOnCanvas=(type, gridX, gridY) => {
        if (matchState === 'found') {
            triggerBackendAttack(gridX, gridY);
        }

        setDeployRemaining(prev => {
            const nextVal=(prev[type] || 0) - 1;
            const updated={ ...prev, [type]: Math.max(0, nextVal) };

            if (updated[type] === 0 && selectedTroopType === type) {
                const available=Object.keys(updated).filter(k => updated[k] > 0);
                setSelectedTroopType(available.length > 0 ? available[0] : null);
            }
            
            const saved = localStorage.getItem('activeBattle');
            if (saved) {
                const data = JSON.parse(saved);
                data.deployRemaining = updated;
                if (battleSceneRef.current) {
                    data.deploymentEvents = battleSceneRef.current.getDeploymentEvents();
                }
                localStorage.setItem('activeBattle', JSON.stringify(data));
            }

            return updated;
        });
    };

    const triggerBackendAttack=(gridX, gridY) => {
        setMatchState('battle');
        const now = Date.now();

        const saved = localStorage.getItem('activeBattle');
        if (saved) {
            const data = JSON.parse(saved);
            data.startTimestamp = now;
            localStorage.setItem('activeBattle', JSON.stringify(data));
        }

        timerIntervalRef.current=setInterval(() => {
            setTimer(prev => {
                if (prev <= 1) {
                    clearInterval(timerIntervalRef.current);
                    if (battleSceneRef.current) {
                        battleSceneRef.current.endBattle();
                    }
                    return 0;
                }
                return prev - 1;
            });
        }, 1000);

        if (battleSceneRef.current) {
            battleSceneRef.current.startBattle();
        }
    };

    const handleSimulationEnd=async (visualDamagePercent) => {
        clearInterval(timerIntervalRef.current);
        
        try {
            const events = battleSceneRef.current ? battleSceneRef.current.getDeploymentEvents() : [];
            const res=await launchAttack(opponent.player_id, events);
            finalResultRef.current=res;
        } catch (err) {
            console.error("Backend attack simulation failed:", err);
            finalResultRef.current={
                elixir_looted: Math.floor(maxLootElixir * (visualDamagePercent / 100)),
                pancakes_looted: Math.floor(maxLootPancakes * (visualDamagePercent / 100)),
                damage_percent: visualDamagePercent
            };
        }

        const finalRes=finalResultRef.current;
        setBattleResult(finalRes);
        setMatchState('ended');
        localStorage.removeItem('activeBattle');
    };

    const handleSkipVisuals=() => {
        if (battleSceneRef.current) {
            battleSceneRef.current.endBattle();
        }
    };

    const maxLootElixir=opponent ? Math.floor(opponent.elixir * 0.20) : 0;
    const maxLootPancakes=opponent ? Math.floor(opponent.pancakes * 0.20) : 0;

    const currentLootElixir = battleResult ? battleResult.elixir_looted : (opponent ? Math.floor(opponent.elixir * (damagePercent / 100) * 0.20) : 0);
    const currentLootPancakes = battleResult ? battleResult.pancakes_looted : (opponent ? Math.floor(opponent.pancakes * (damagePercent / 100) * 0.20) : 0);
    const displayDamagePercent = battleResult ? battleResult.damage_percent : damagePercent;

    const formatTimer=(seconds) => {
        const m=Math.floor(seconds / 60);
        const s=seconds % 60;
        return `${m}:${s < 10 ? '0' : ''}${s}`;
    };

    const hasAnyTroopsLeft=Object.values(deployRemaining).some(qty => qty > 0);

    return (
        <div className={styles.page}>
            {matchState === 'search' && (
                <div className={styles.searchContainer}>
                    <h1 className={styles.searchTitle}>MULTIPLAYER ATTACK</h1>
                    
                    {!hasTrainedArmy && !loading && (
                        <div className={styles.warningBox}>You have no army! Train troops in your barracks before searching for an opponent.</div>
                    )}

                    {error && <div className={styles.errorBox}>{error}</div>}

                    <div className={styles.searchActions}>
                        <button 
                            className={`${styles.searchBtn} ${hasTrainedArmy ? styles.searchBtnActive : styles.searchBtnDisabled}`}
                            onClick={handleFindOpponent} disabled={!hasTrainedArmy || loading}>Find Opponent</button>
                        <button className={styles.cancelBtn} onClick={() => navigate('/')}>Return to Village</button>
                    </div>
                </div>
            )}

            {matchState === 'searching' && (
                <div className={styles.searchingOverlay}>
                    <div className={styles.searchRadar}></div>
                    <h2>Searching for target base...</h2>
                    <p>Scanning skill matchmaking range</p>
                </div>
            )}

            {(matchState === 'found' || matchState === 'battle' || matchState === 'ended') && opponent && (
                <div className={styles.arenaContainer}>
                    <div className={styles.battleHeader}>
                        <div className={styles.opponentStats}>
                            <h3>Opponent (Trophies: {opponent.trophies} 🏆)</h3>
                        </div>
                        
                        <div className={styles.centerClock}>
                            <span className={styles.timerText}>{formatTimer(timer)}</span>
                        </div>

                        <div className={styles.damageTracker}>
                            <span>Damage:</span>
                            <span className={styles.damagePctText}>{Math.round(displayDamagePercent)}%</span>
                        </div>
                    </div>

                    <div className={styles.lootTracker}>
                        <div className={styles.lootRow}>
                            <span>🥞 Pancakes:</span>
                            <span>{currentLootPancakes.toLocaleString()} / {maxLootPancakes.toLocaleString()}</span>
                        </div>
                        <div className={styles.lootRow}>
                            <span>💧 Elixir:</span>
                            <span>{currentLootElixir.toLocaleString()} / {maxLootElixir.toLocaleString()}</span>
                        </div>
                    </div>

                    <div ref={canvasContainerRef} className={styles.canvasContainer}></div>

                    {matchState !== 'ended' && (
                        <div className={styles.hudOverlay}>
                            {matchState === 'found' && (
                                <div className={styles.preBattlePrompt}>Drop a troop on the map to begin the attack!</div>
                            )}

                            <div className={styles.armyBar}>
                                {Object.keys(deployRemaining).map((type) => {
                                    const qty=deployRemaining[type];
                                    if (qty <= 0) return null;
                                    const isSelected=selectedTroopType === type;
                                    const def=TROOP_DEFS[type];

                                    return (
                                        <button key={type} onClick={() => setSelectedTroopType(type)}
                                            className={`${styles.troopButton} ${isSelected ? styles.troopButtonSelected : styles.troopButtonUnselected}`}>
                                            <span className={styles.troopEmoji}>{def?.emoji}</span>
                                            <span className={styles.troopQty}>x{qty}</span>
                                            <span className={styles.troopName}>{type}</span>
                                        </button>
                                    );
                                })}
                                {!hasAnyTroopsLeft && (
                                    <div className={styles.noTroopsText}>All Troops Deployed</div>
                                )}
                            </div>

                            <div className={styles.actionRow}>
                                {matchState === 'found' ? (
                                    <button className={styles.surrenderBtn} onClick={() => {
                                        localStorage.removeItem('activeBattle');
                                        navigate('/');
                                    }}>
                                        Surrender</button>
                                ) : (
                                    <button className={styles.skipBtn} onClick={handleSkipVisuals}>
                                        End Battle</button>
                                )}
                            </div>
                        </div>
                    )}
                </div>
            )}

            {matchState === 'ended' && (
                <BattleResultModel 
                    result={battleResult}
                    onClose={() => navigate('/')}
                />
            )}
        </div>
    );
}

