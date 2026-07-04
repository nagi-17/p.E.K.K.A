import styles from './BattleResultModel.module.css';

export default function BattleResultModel({ result, onClose }) {
    if (!result) return null;
    const { elixir_looted, pancakes_looted, damage_percent }=result;

    const damage=Math.round(damage_percent);
    const hasWon=damage >= 50;

    let stars=0;
    if (damage === 100) stars=3;
    else if (damage >= 70) stars=2;
    else if (damage >= 50) stars=1;

    const trophyChange=hasWon ? '+30' : '-20';

    return (
        <div className={styles.overlay}>
            <div className={styles.modal}>
                <h1 className={styles.statusTitle} style={{ color: hasWon ? '#2ecc71' : '#e74c3c' }}>
                    {hasWon ? 'VICTORY!' : 'DEFEAT'}
                </h1>

                <div className={styles.starContainer}>
                    <span className={`${styles.star} ${stars >= 1 ? styles.starWon : styles.starLost}`}>⭐</span>
                    <span className={`${styles.star} ${styles.starBig} ${stars >= 2 ? styles.starWon : styles.starLost}`}>⭐</span>
                    <span className={`${styles.star} ${stars >= 3 ? styles.starWon : styles.starLost}`}>⭐</span>
                </div>

                <div className={styles.damageBox}>
                    <div className={styles.damageText}>Damage Dealt: {damage}%</div>
                    <div className={styles.trophyBadge}>
                        🏆 {trophyChange} Trophies
                    </div>
                </div>

                <div className={styles.lootSection}>
                    <h3 className={styles.lootTitle}>Loot Stolen</h3>
                    <div className={styles.lootGrid}>
                        <div className={styles.lootCard}>
                            <span className={styles.lootEmoji}>🥞</span>
                            <span className={styles.lootAmount}>{pancakes_looted.toLocaleString()}</span>
                        </div>
                        <div className={styles.lootCard}>
                            <span className={styles.lootEmoji}>💧</span>
                            <span className={styles.lootAmount}>{elixir_looted.toLocaleString()}</span>
                        </div>
                    </div>
                </div>

                <button onClick={onClose} className={styles.returnBtn}>
                    Return to Village
                </button>
            </div>
        </div>
    );
}


