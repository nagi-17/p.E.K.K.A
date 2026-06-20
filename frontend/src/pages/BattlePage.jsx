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
            
            return updated;
        });
    };

    const triggerBackendAttack=async (gridX, gridY) => {
        setMatchState('battle');

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

        try {
            const res=await launchAttack(opponent.player_id, gridX, gridY);
            finalResultRef.current=res;
        } catch (err) {
            console.error("Backend attack simulation failed:", err);
        }

        if (battleSceneRef.current) {
            battleSceneRef.current.startBattle();
        }
    };

    const handleSimulationEnd=(visualDamagePercent) => {
        clearInterval(timerIntervalRef.current);
        setMatchState('ended');

        const finalRes=finalResultRef.current || {
            elixir_looted: Math.floor(maxLootElixir * (visualDamagePercent / 100)),
            pancakes_looted: Math.floor(maxLootPancakes * (visualDamagePercent / 100)),
            damage_percent: visualDamagePercent
        };
        setBattleResult(finalRes);
    };

    const handleSkipVisuals=() => {
        if (battleSceneRef.current) {
            battleSceneRef.current.endBattle();
        }
    };

    const maxLootPancakes=opponent ? Math.floor(opponent.pancakes * 0.20) : 0;
    const maxLootElixir=opponent ? Math.floor(opponent.elixir * 0.20) : 0;

    const currentLootPancakes=Math.floor(maxLootPancakes * (damagePercent / 100));
    const currentLootElixir=Math.floor(maxLootElixir * (damagePercent / 100));

    const formatTimer=(seconds) => {
        const m=Math.floor(seconds / 60);
        const s=seconds % 60;
        return `${m}:${s < 10 ? '0' : ''}${s}`;
    };

    const hasAnyTroopsLeft=Object.values(deployRemaining).some(qty => qty > 0);

    return (
        <div style={styles.page}>
            {matchState === 'search' && (
                <div style={styles.searchContainer}>
                    <h1 style={styles.searchTitle}>MULTIPLAYER ATTACK</h1>
                    
                    {!hasTrainedArmy && !loading && (
                        <div style={styles.warningBox}>You have no army! Train troops in your barracks before searching for an opponent.</div>
                    )}

                    {error && <div style={styles.errorBox}>{error}</div>}

                    <div style={styles.searchActions}>
                        <button 
                            style={{...styles.searchBtn,backgroundColor: hasTrainedArmy ? '#2ecc71' : '#7f8c8d',cursor: hasTrainedArmy ? 'pointer' : 'not-allowed'}}
                            onClick={handleFindOpponent} disabled={!hasTrainedArmy || loading}>Find Opponent</button>
                        <button style={styles.cancelBtn} onClick={() => navigate('/')}>Return to Village</button>
                    </div>
                </div>
            )}

            {matchState === 'searching' && (
                <div style={styles.searchingOverlay}>
                    <div style={styles.searchRadar}></div>
                    <h2>Searching for target base...</h2>
                    <p>Scanning skill matchmaking range</p>
                </div>
            )}

            {(matchState === 'found' || matchState === 'battle' || matchState === 'ended') && opponent && (
                <div style={styles.arenaContainer}>
                    <div style={styles.battleHeader}>
                        <div style={styles.opponentStats}>
                            <h3>Opponent (Trophies: {opponent.trophies} 🏆)</h3>
                        </div>
                        
                        <div style={styles.centerClock}>
                            <span style={styles.timerText}>{formatTimer(timer)}</span>
                        </div>

                        <div style={styles.damageTracker}>
                            <span>Damage:</span>
                            <span style={styles.damagePctText}>{Math.round(damagePercent)}%</span>
                        </div>
                    </div>

                    <div style={styles.lootTracker}>
                        <div style={styles.lootRow}>
                            <span>🥞 Pancakes:</span>
                            <span>{currentLootPancakes.toLocaleString()} / {maxLootPancakes.toLocaleString()}</span>
                        </div>
                        <div style={styles.lootRow}>
                            <span>💧 Elixir:</span>
                            <span>{currentLootElixir.toLocaleString()} / {maxLootElixir.toLocaleString()}</span>
                        </div>
                    </div>

                    <div ref={canvasContainerRef} style={styles.canvasContainer}></div>

                    {matchState !== 'ended' && (
                        <div style={styles.hudOverlay}>
                            {matchState === 'found' && (
                                <div style={styles.preBattlePrompt}>Drop a troop on the map to begin the attack!</div>
                            )}

                            <div style={styles.armyBar}>
                                {Object.keys(deployRemaining).map((type) => {
                                    const qty=deployRemaining[type];
                                    if (qty <= 0) return null;
                                    const isSelected=selectedTroopType === type;
                                    const def=TROOP_DEFS[type];

                                    return (
                                        <button key={type} onClick={() => setSelectedTroopType(type)}
                                            style={{...styles.troopButton,border: isSelected ? '4px solid #f1c40f' : '2px solid #34495e',transform: isSelected ? 'scale(1.08)' : 'scale(1)'}}>
                                            <span style={styles.troopEmoji}>{def?.emoji}</span>
                                            <span style={styles.troopQty}>x{qty}</span>
                                            <span style={styles.troopName}>{type}</span>
                                        </button>
                                    );
                                })}
                                {!hasAnyTroopsLeft && (
                                    <div style={styles.noTroopsText}>All Troops Deployed</div>
                                )}
                            </div>

                            <div style={styles.actionRow}>
                                {matchState === 'found' ? (
                                    <button style={styles.surrenderBtn} onClick={() => navigate('/')}>
                                        Surrender</button>
                                ) : (
                                    <button style={styles.skipBtn} onClick={handleSkipVisuals}>
                                        Skip Visuals</button>
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

const styles={
    page: {
        width: '100vw',
        height: '100dvh',
        backgroundColor: '#2c3e50',
        overflow: 'hidden',
        position: 'relative',
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center'
    },
    searchContainer: {
        backgroundColor: '#34495e',
        border: '6px solid #e67e22',
        borderRadius: '16px',
        padding: '2.5rem 2rem',
        width: '90%',
        maxWidth: '500px',
        textAlign: 'center',
        boxShadow: '0 8px 24px rgba(0,0,0,0.5)',
        fontFamily: '"Luckiest Guy", cursive',
        color: 'white',
        display: 'flex',
        flexDirection: 'column',
        gap: '1.5rem'
    },
    searchTitle: {
        fontSize: '2.5rem',
        letterSpacing: '1px',
        textShadow: '2px 2px 0px #000',
        margin: 0
    },
    warningBox: {
        backgroundColor: '#d35400',
        border: '2px solid #e67e22',
        borderRadius: '8px',
        padding: '1rem',
        fontSize: '1rem',
        textShadow: '1px 1px 0px #000'
    },
    errorBox: {
        backgroundColor: '#e74c3c',
        borderRadius: '8px',
        padding: '0.8rem',
        fontSize: '0.95rem'
    },
    searchActions: {
        display: 'flex',
        flexDirection: 'column',
        gap: '0.75rem'
    },
    searchBtn: {
        color: 'white',
        border: 'none',
        borderRadius: '8px',
        padding: '1rem',
        fontSize: '1.5rem',
        fontFamily: 'inherit',
        textShadow: '2px 2px 0px rgba(0,0,0,0.3)',
        boxShadow: '0 4px 0px rgba(0,0,0,0.2)'
    },
    cancelBtn: {
        backgroundColor: '#e74c3c',
        color: 'white',
        border: 'none',
        borderRadius: '8px',
        padding: '1rem',
        fontSize: '1.3rem',
        fontFamily: 'inherit',
        cursor: 'pointer',
        textShadow: '2px 2px 0px #c0392b',
        boxShadow: '0 4px 0px #c0392b'
    },
    searchingOverlay: {
        display: 'flex',
        flexDirection: 'column',
        justifyContent: 'center',
        alignItems: 'center',
        color: 'white',
        gap: '1rem',
        fontFamily: '"Luckiest Guy", cursive'
    },
    searchRadar: {
        width: '80px',
        height: '80px',
        border: '6px solid #f1c40f',
        borderTopColor: 'transparent',
        borderRadius: '50%',
        animation: 'spin 1s linear infinite'
    },
    arenaContainer: {
        width: '100%',
        height: '100%',
        position: 'relative',
        display: 'flex',
        flexDirection: 'column',
        overflow: 'hidden'
    },
    battleHeader: {
        position: 'absolute',
        top: 0,
        left: 0,
        right: 0,
        height: '60px',
        backgroundColor: 'rgba(0, 0, 0, 0.65)',
        display: 'flex',
        justifyContent: 'space-between',
        alignItems: 'center',
        padding: '0 1.5rem',
        zIndex: 10,
        color: 'white',
        fontFamily: '"Luckiest Guy", cursive',
        borderBottom: '2px solid rgba(255,255,255,0.1)'
    },
    opponentStats: {
        flex: 1
    },
    centerClock: {
        flex: 1,
        display: 'flex',
        justifyContent: 'center'
    },
    timerText: {
        backgroundColor: '#e74c3c',
        padding: '0.2rem 1.2rem',
        borderRadius: '20px',
        fontSize: '1.4rem',
        border: '2px solid white',
        textShadow: '1px 1px 0px #000'
    },
    damageTracker: {
        flex: 1,
        display: 'flex',
        justifyContent: 'flex-end',
        gap: '0.5rem',
        fontSize: '1.3rem'
    },
    damagePctText: {
        color: '#f1c40f'
    },
    lootTracker: {
        position: 'absolute',
        top: '70px',
        left: '1rem',
        backgroundColor: 'rgba(0,0,0,0.5)',
        color: 'white',
        padding: '0.5rem 1rem',
        borderRadius: '10px',
        zIndex: 10,
        display: 'flex',
        flexDirection: 'column',
        gap: '0.25rem',
        fontFamily: '"Luckiest Guy", cursive',
        border: '1px solid rgba(255,255,255,0.1)'
    },
    lootRow: {
        display: 'flex',
        justifyContent: 'space-between',
        gap: '1rem',
        fontSize: '1rem'
    },
    canvasContainer: {
        position: 'absolute',
        top: 0,
        left: 0,
        right: 0,
        bottom: 0,
        overflow: 'hidden',
        zIndex: 0
    },
    hudOverlay: {
        position: 'absolute',
        bottom: 0,
        left: 0,
        right: 0,
        padding: '1rem',
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
        gap: '0.75rem',
        zIndex: 10,
        backgroundColor: 'gradient(to top, rgba(0,0,0,0.7), transparent)'
    },
    preBattlePrompt: {
        backgroundColor: 'rgba(230, 126, 34, 0.95)',
        color: 'white',
        padding: '0.4rem 1.5rem',
        borderRadius: '15px',
        fontSize: '1rem',
        fontFamily: '"Luckiest Guy", cursive',
        textShadow: '1px 1px 0px #000',
        animation: 'pulse 1.5s infinite',
        border: '2px solid white'
    },
    armyBar: {
        backgroundColor: 'rgba(0,0,0,0.65)',
        padding: '0.65rem',
        borderRadius: '16px',
        border: '3px solid #34495e',
        display: 'flex',
        gap: '0.8rem',
        width: '90%',
        maxWidth: '650px',
        overflowX: 'auto',
        justifyContent: 'center',
        alignItems: 'center'
    },
    troopButton: {
        backgroundColor: '#2c3e50',
        color: 'white',
        borderRadius: '10px',
        padding: '0.4rem',
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
        minWidth: '70px',
        position: 'relative',
        cursor: 'pointer',
        fontFamily: '"Luckiest Guy", cursive',
        boxShadow: '0 4px 6px rgba(0,0,0,0.3)',
        transition: 'all 0.15s ease'
    },
    troopEmoji: {
        fontSize: '1.8rem'
    },
    troopQty: {
        position: 'absolute',
        top: '-8px',
        right: '-8px',
        backgroundColor: '#e74c3c',
        borderRadius: '50%',
        width: '20px',
        height: '20px',
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center',
        fontSize: '0.75rem',
        border: '1px solid white'
    },
    troopName: {
        fontSize: '0.75rem',
        marginTop: '0.2rem',
        color: '#bdc3c7'
    },
    noTroopsText: {
        color: '#bdc3c7',
        fontSize: '1.2rem',
        fontFamily: '"Luckiest Guy", cursive',
        padding: '0.5rem'
    },
    actionRow: {
        display: 'flex',
        gap: '1rem',
        width: '90%',
        maxWidth: '650px',
        justifyContent: 'center'
    },
    surrenderBtn: {
        flex: 1,
        backgroundColor: '#e74c3c',
        color: 'white',
        border: 'none',
        borderRadius: '8px',
        padding: '0.75rem',
        fontSize: '1.2rem',
        fontFamily: '"Luckiest Guy", cursive',
        cursor: 'pointer',
        textShadow: '1px 1px 0px #c0392b',
        boxShadow: '0 4px 0px #c0392b'
    },
    skipBtn: {
        flex: 1,
        backgroundColor: '#f1c40f',
        color: 'white',
        border: 'none',
        borderRadius: '8px',
        padding: '0.75rem',
        fontSize: '1.2rem',
        fontFamily: '"Luckiest Guy", cursive',
        cursor: 'pointer',
        textShadow: '1px 1px 0px #d35400',
        boxShadow: '0 4px 0px #d35400'
    }
};

if (typeof document !== 'undefined') {
    const styleSheet=document.createElement('style');
    styleSheet.type='text/css';
    styleSheet.innerText=`
        @keyframes spin {
            0% { transform: rotate(0deg); }
            100% { transform: rotate(360deg); }
        }
        @keyframes pulse {
            0% { transform: scale(1); }
            50% { transform: scale(1.05); }
            100% { transform: scale(1); }
        }
    `;
    document.head.appendChild(styleSheet);
}
