import { useEffect, useState, useRef } from 'react';
import { TROOP_DEFS } from '../../game/gameConfig';

export default function TrainTroopCard({ troop, onTrain, onDiscard, onUpgradeStart, onUpgradeFinish, labLevel, elixirBalance, isCampFull }) {
    const { troop_type, level, quantity, upgrade_complete_at }=troop;
    const def=TROOP_DEFS[troop_type];
    if (!def) return null;

    const currentStats=def.levels[level] || def.levels[1];
    const nextLevel=level < 4 ? level + 1 : null;
    const nextStats=nextLevel ? def.levels[nextLevel] : null;

    const [timeLeft, setTimeLeft]=useState(0);
    const [totalUpgradeTime, setTotalUpgradeTime]=useState(0);
    const attemptedFinalizeRef=useRef(false);

    useEffect(() => {
        if (!upgrade_complete_at) {
            setTimeLeft(0);
            attemptedFinalizeRef.current=false;
            return;
        }

        const endTime=new Date(upgrade_complete_at).getTime();
        const startUpgradeTime=currentStats.upgradeTime || 60;
        setTotalUpgradeTime(startUpgradeTime);

        let interval;

        const updateTimer=() => {
            const now=Date.now();
            const diff=endTime - now;
            if (diff <= 0) {
                setTimeLeft(0);
                if (interval) clearInterval(interval);
                if (!attemptedFinalizeRef.current) {
                    attemptedFinalizeRef.current=true;
                    onUpgradeFinish(troop_type);
                }
            } else {
                attemptedFinalizeRef.current=false;
                setTimeLeft(Math.ceil(diff / 1000));
            }
        };

        updateTimer();
        interval=setInterval(updateTimer, 1000);
        return () => clearInterval(interval);
    }, [upgrade_complete_at, level]);

    const isUpgrading=timeLeft > 0;
    const canUpgrade=nextLevel && labLevel >= nextStats.labReq && elixirBalance >= nextStats.upgradeCost && !isUpgrading;

    const formatTime=(seconds) => {
        if (seconds <= 0) return '0s';
        const m=Math.floor(seconds / 60);
        const s=seconds % 60;
        if (m > 0) return `${m}m ${s}s`;
        return `${s}s`;
    };

    return (
        <div style={styles.card}>
            <div style={styles.badge}>{quantity}</div>
            
            <div style={styles.header}>
                <div style={styles.emojiContainer}><span style={styles.emoji}>{def.emoji}</span></div>
                <div style={styles.titleInfo}>
                    <h3 style={styles.troopName}>{troop_type}</h3>
                    <span style={styles.levelBadge}>Lvl {level}</span>
                </div>
            </div>

            <div style={styles.statsContainer}>
                <div style={styles.statRow}>
                    <span>❤️ HP: {currentStats.hp}</span>
                    <span>⚔️ DPS: {currentStats.dps}</span>
                </div>
                <div style={styles.statRow}>
                    <span>⛺ Space: {def.space}</span>
                    {nextLevel && <span style={styles.nextLvl}>Next: Lvl {nextLevel}</span>}
                </div>
            </div>

            {isUpgrading ? (
                <div style={styles.upgradeProgressContainer}>
                    <div style={styles.progressLabel}>
                        <span>Lab Upgrading...</span>
                        <span>{formatTime(timeLeft)}</span>
                    </div>
                    <div style={styles.progressBarBg}>
                        <div 
                            style={{ 
                                ...styles.progressBarFill, 
                                width: `${Math.max(0, Math.min(100, (1 - timeLeft / totalUpgradeTime) * 100))}%` 
                            }} 
                        />
                    </div>
                </div>
            ) : (
                <div style={styles.actionContainer}>
                    <div style={{ display: 'flex', flexDirection: 'column', flex: 1.2, gap: '4px', width: '100%' }}>
                        <button 
                            style={{
                                ...styles.trainBtn,
                                backgroundColor: isCampFull ? '#7f8c8d' : '#f1c40f',
                                cursor: isCampFull ? 'not-allowed' : 'pointer',
                                width: '100%'
                            }}
                            disabled={isCampFull}
                            onClick={() => onTrain(troop_type)}
                        >
                            {isCampFull ? 'Full' : 'Train (+1)'}
                        </button>
                        {quantity > 0 && (
                            <button
                                style={{
                                    ...styles.trainBtn,
                                    backgroundColor: '#e74c3c',
                                    cursor: 'pointer',
                                    width: '100%',
                                    fontSize: '0.8rem',
                                    padding: '0.4rem 0.2rem'
                                }}
                                onClick={() => onDiscard(troop_type)}
                            >
                                Discard (-1)
                            </button>
                        )}
                    </div>

                    {nextLevel ? (
                        <button
                            style={{
                                ...styles.upgradeBtn,
                                backgroundColor: canUpgrade ? '#9b59b6' : '#7f8c8d',
                                cursor: canUpgrade ? 'pointer' : 'not-allowed'
                            }}
                            disabled={!canUpgrade}
                            onClick={() => onUpgradeStart(troop_type)}
                            title={
                                !labLevel 
                                    ? 'Requires Laboratory' 
                                    : labLevel < nextStats.labReq 
                                    ? `Requires Lab Lvl ${nextStats.labReq}` 
                                    : elixirBalance < nextStats.upgradeCost 
                                    ? 'Insufficient Elixir' 
                                    : `Upgrade to Lvl ${nextLevel}`
                            }
                        >
                            <span style={styles.upgradeBtnText}>Lvl 🧪</span>
                            <span style={styles.upgradeCostText}>{nextStats.upgradeCost.toLocaleString()} 💧</span>
                        </button>
                    ) : (
                        <div style={styles.maxBadge}>MAX</div>
                    )}
                </div>
            )}
        </div>
    );
}

const styles={
    card: {
        position: 'relative',
        backgroundColor: '#2c3e50',
        color: 'white',
        border: '3px solid #34495e',
        borderRadius: '12px',
        padding: '0.85rem',
        display: 'flex',
        flexDirection: 'column',
        gap: '0.6rem',
        boxShadow: '0 4px 6px rgba(0,0,0,0.3)',
        transition: 'transform 0.15s ease',
        fontFamily: '"Luckiest Guy", cursive',
    },
    badge: {
        position: 'absolute',
        top: '-10px',
        right: '-10px',
        backgroundColor: '#e74c3c',
        color: 'white',
        borderRadius: '50%',
        width: '28px',
        height: '28px',
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center',
        fontSize: '1.1rem',
        border: '2px solid white',
        boxShadow: '0 2px 4px rgba(0,0,0,0.4)',
        textShadow: '1px 1px 0px #000'
    },
    header: {
        display: 'flex',
        alignItems: 'center',
        gap: '0.75rem'
    },
    emojiContainer: {
        backgroundColor: '#34495e',
        borderRadius: '8px',
        width: '45px',
        height: '45px',
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center',
        border: '2px solid #1a252f'
    },
    emoji: {
        fontSize: '2rem'
    },
    titleInfo: {
        display: 'flex',
        flexDirection: 'column',
        gap: '0.15rem'
    },
    troopName: {
        margin: 0,
        fontSize: '1.2rem',
        letterSpacing: '0.5px',
        color: '#f39c12',
        textShadow: '1px 1px 0px #000'
    },
    levelBadge: {
        fontSize: '0.85rem',
        color: '#bdc3c7'
    },
    statsContainer: {
        backgroundColor: 'rgba(0, 0, 0, 0.2)',
        borderRadius: '8px',
        padding: '0.4rem 0.6rem',
        fontSize: '0.85rem',
        display: 'flex',
        flexDirection: 'column',
        gap: '0.2rem',
        border: '1px solid rgba(255,255,255,0.05)'
    },
    statRow: {
        display: 'flex',
        justifyContent: 'space-between'
    },
    nextLvl: {
        color: '#9b59b6'
    },
    actionContainer: {
        display: 'flex',
        gap: '0.4rem',
        marginTop: 'auto'
    },
    trainBtn: {
        flex: 1.2,
        color: 'white',
        border: 'none',
        borderRadius: '6px',
        padding: '0.5rem',
        fontSize: '0.95rem',
        fontFamily: 'inherit',
        textShadow: '1px 1px 0px #000',
        boxShadow: '0 2px 0px rgba(0,0,0,0.2)',
    },
    upgradeBtn: {
        flex: 1.8,
        color: 'white',
        border: 'none',
        borderRadius: '6px',
        padding: '0.5rem',
        fontSize: '0.8rem',
        fontFamily: 'inherit',
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
        justifyContent: 'center',
        textShadow: '1px 1px 0px #000',
        boxShadow: '0 2px 0px rgba(0,0,0,0.2)',
    },
    upgradeBtnText: {
        fontSize: '0.85rem'
    },
    upgradeCostText: {
        fontSize: '0.7rem',
        marginTop: '0.05rem'
    },
    maxBadge: {
        flex: 1,
        backgroundColor: '#16a085',
        color: 'white',
        borderRadius: '6px',
        padding: '0.5rem',
        fontSize: '1rem',
        textAlign: 'center',
        textShadow: '1px 1px 0px #000',
        border: '2px solid #117a65'
    },
    upgradeProgressContainer: {
        display: 'flex',
        flexDirection: 'column',
        gap: '0.25rem',
        marginTop: 'auto'
    },
    progressLabel: {
        display: 'flex',
        justifyContent: 'space-between',
        fontSize: '0.75rem',
        color: '#bdc3c7'
    },
    progressBarBg: {
        backgroundColor: '#1a252f',
        borderRadius: '4px',
        height: '8px',
        overflow: 'hidden',
        border: '1px solid #34495e'
    },
    progressBarFill: {
        backgroundColor: '#9b59b6',
        height: '100%',
        borderRadius: '4px',
        transition: 'width 0.3s ease'
    }
};
