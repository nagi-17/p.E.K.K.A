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
        <div style={styles.overlay}>
            <div style={styles.modal}>
                <h1 style={{ ...styles.statusTitle, color: hasWon ? '#2ecc71' : '#e74c3c' }}>
                    {hasWon ? 'VICTORY!' : 'DEFEAT'}
                </h1>

                <div style={styles.starContainer}>
                    <span style={{ ...styles.star, color: stars >= 1 ? '#f1c40f' : '#7f8c8d' }}>⭐</span>
                    <span style={{ ...styles.star, color: stars >= 2 ? '#f1c40f' : '#7f8c8d', fontSize: '4.5rem' }}>⭐</span>
                    <span style={{ ...styles.star, color: stars >= 3 ? '#f1c40f' : '#7f8c8d' }}>⭐</span>
                </div>

                <div style={styles.damageBox}>
                    <div style={styles.damageText}>Damage Dealt: {damage}%</div>
                    <div style={styles.trophyBadge}>
                        🏆 {trophyChange} Trophies
                    </div>
                </div>

                <div style={styles.lootSection}>
                    <h3 style={styles.lootTitle}>Loot Stolen</h3>
                    <div style={styles.lootGrid}>
                        <div style={styles.lootCard}>
                            <span style={styles.lootEmoji}>🥞</span>
                            <span style={styles.lootAmount}>{pancakes_looted.toLocaleString()}</span>
                        </div>
                        <div style={styles.lootCard}>
                            <span style={styles.lootEmoji}>💧</span>
                            <span style={styles.lootAmount}>{elixir_looted.toLocaleString()}</span>
                        </div>
                    </div>
                </div>

                <button onClick={onClose} style={styles.returnBtn}>
                    Return to Village
                </button>
            </div>
        </div>
    );
}

const styles={
    overlay: {
        position: 'absolute',
        top: 0,
        left: 0,
        right: 0,
        bottom: 0,
        backgroundColor: 'rgba(0, 0, 0, 0.8)',
        zIndex: 200,
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center',
        fontFamily: '"Luckiest Guy", cursive'
    },
    modal: {
        backgroundColor: '#ecf0f1',
        width: '85%',
        maxWidth: '450px',
        borderRadius: '20px',
        border: '6px solid #e67e22',
        padding: '2rem 1.5rem',
        boxShadow: '0 12px 30px rgba(0,0,0,0.6)',
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
        textAlign: 'center',
        gap: '1.2rem'
    },
    statusTitle: {
        fontSize: '3rem',
        margin: 0,
        textShadow: '2px 2px 0px #000',
        letterSpacing: '2px'
    },
    starContainer: {
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center',
        gap: '0.5rem',
        height: '80px'
    },
    star: {
        fontSize: '3.5rem',
        textShadow: '0px 4px 6px rgba(0,0,0,0.3)',
        transition: 'all 0.3s ease'
    },
    damageBox: {
        backgroundColor: '#34495e',
        color: 'white',
        border: '3px solid #2c3e50',
        borderRadius: '12px',
        padding: '0.8rem 1.5rem',
        width: '80%',
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
        gap: '0.4rem',
        boxShadow: 'inset 0 2px 4px rgba(0,0,0,0.5)'
    },
    damageText: {
        fontSize: '1.3rem',
        textShadow: '1px 1px 0px #000'
    },
    trophyBadge: {
        color: '#f1c40f',
        fontSize: '1.1rem',
        textShadow: '1px 1px 0px #000'
    },
    lootSection: {
        width: '100%'
    },
    lootTitle: {
        margin: '0.5rem 0',
        color: '#2c3e50',
        fontSize: '1.4rem'
    },
    lootGrid: {
        display: 'flex',
        justifyContent: 'center',
        gap: '1rem',
        width: '100%'
    },
    lootCard: {
        flex: 1,
        backgroundColor: '#2c3e50',
        border: '2px solid #1a252f',
        borderRadius: '10px',
        padding: '0.6rem 0.5rem',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        gap: '0.5rem',
        boxShadow: '0 4px 6px rgba(0,0,0,0.15)'
    },
    lootEmoji: {
        fontSize: '1.8rem'
    },
    lootAmount: {
        color: 'white',
        fontSize: '1.2rem',
        textShadow: '1px 1px 0px #000'
    },
    returnBtn: {
        backgroundColor: '#2ecc71',
        color: 'white',
        border: 'none',
        borderRadius: '10px',
        padding: '0.85rem 2rem',
        fontSize: '1.4rem',
        cursor: 'pointer',
        fontFamily: 'inherit',
        textShadow: '2px 2px 0px #27ae60',
        boxShadow: '0 4px 0px #27ae60',
        marginTop: '0.5rem',
        transition: 'transform 0.1s ease',
        width: '80%'
    }
};
