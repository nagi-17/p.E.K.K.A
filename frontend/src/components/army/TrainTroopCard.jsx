import { useEffect, useState, useRef } from 'react';
import { TROOP_DEFS } from '../../game/gameConfig';
import styles from './TrainTroopCard.module.css';

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
        <div className={styles.card}>
            <div className={styles.badge}>{quantity}</div>
            
            <div className={styles.header}>
                <div className={styles.emojiContainer}><span className={styles.emoji}>{def.emoji}</span></div>
                <div className={styles.titleInfo}>
                    <h3 className={styles.troopName}>{troop_type}</h3>
                    <span className={styles.levelBadge}>Lvl {level}</span>
                </div>
            </div>

            <div className={styles.statsContainer}>
                <div className={styles.statRow}>
                    <span>❤️ HP: {currentStats.hp}</span>
                    <span>⚔️ DPS: {currentStats.dps}</span>
                </div>
                <div className={styles.statRow}>
                    <span>⛺ Space: {def.space}</span>
                    {nextLevel && <span className={styles.nextLvl}>Next: Lvl {nextLevel}</span>}
                </div>
            </div>

            {isUpgrading ? (
                <div className={styles.upgradeProgressContainer}>
                    <div className={styles.progressLabel}>
                        <span>Lab Upgrading...</span>
                        <span>{formatTime(timeLeft)}</span>
                    </div>
                    <div className={styles.progressBarBg}>
                        <div 
                            className={styles.progressBarFill}
                            style={{ 
                                width: `${Math.max(0, Math.min(100, (1 - timeLeft / totalUpgradeTime) * 100))}%` 
                            }} 
                        />
                    </div>
                </div>
            ) : (
                <div className={styles.actionContainer}>
                    <div style={{ display: 'flex', flexDirection: 'column', flex: 1.2, gap: '4px', width: '100%' }}>
                        <button 
                            className={`${styles.trainBtn} ${isCampFull ? styles.trainBtnExpensive : styles.trainBtnAffordable}`}
                            disabled={isCampFull}
                            onClick={() => onTrain(troop_type)}
                        >
                            {isCampFull ? 'Full' : 'Train (+1)'}
                        </button>
                        {quantity > 0 && (
                            <button
                                className={styles.discardBtn}
                                onClick={() => onDiscard(troop_type)}
                            >
                                Discard (-1)
                            </button>
                        )}
                    </div>

                    {nextLevel ? (
                        <button
                            className={`${styles.upgradeBtn} ${canUpgrade ? styles.upgradeBtnAffordable : styles.upgradeBtnExpensive}`}
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
                            <span className={styles.upgradeBtnText}>Lvl 🧪</span>
                            <span className={styles.upgradeCostText}>{nextStats.upgradeCost.toLocaleString()} 💧</span>
                        </button>
                    ) : (
                        <div className={styles.maxBadge}>MAX</div>
                    )}
                </div>
            )}
        </div>
    );
}


