import { useEffect } from 'react';
import { useUiStore } from '../../store/uiStore';
import { useArmyStore } from '../../store/armyStore';
import { useVillageStore } from '../../store/villageStore';
import TrainTroopCard from './TrainTroopCard';
import { TROOP_DEFS } from '../../game/gameConfig';

const CAMP_SPACE={ 1: 20, 2: 30, 3: 40, 4: 50 };

export default function ArmyPanel() {
    const isArmyOpen=useUiStore((state) => state.isArmyOpen);
    const toggleArmy=useUiStore((state) => state.toggleArmy);

    const { troops, isLoading, error, fetchArmy, train, startUpgrade, finishUpgrade, discard }=useArmyStore();
    const { buildings, playerStats }=useVillageStore();

    useEffect(() => {
        if (isArmyOpen) {
            fetchArmy();
        }
    }, [isArmyOpen]);

    if (!isArmyOpen || !playerStats) return null;

    let maxSpace=0;
    const now=Date.now();
    buildings.forEach((b) => {
        if (b.type === 'Army Camp') {
            const isUpgrading=b.upgradeCompleteAt && new Date(b.upgradeCompleteAt).getTime() > now;
            if (!isUpgrading) {
                maxSpace += CAMP_SPACE[b.level || 1] || 20;
            }
        }
    });

    let currentSpace=0;
    troops.forEach((t) => {
        const def=TROOP_DEFS[t.troop_type];
        if (def) {
            currentSpace += t.quantity * def.space;
        }
    });

    let labLevel=0;
    let isLabUpgrading=false;
    const labBuilding=buildings.find((b) => b.type === 'Laboratory');
    if (labBuilding) {
        const isUpgrading=labBuilding.upgradeCompleteAt && new Date(labBuilding.upgradeCompleteAt).getTime() > now;
        if (isUpgrading) {
            isLabUpgrading=true;
        } else {
            labLevel=labBuilding.level || 1;
        }
    }

    const handleTrain=async (troopType) => {
        const def=TROOP_DEFS[troopType];
        if (!def) return;
        if (currentSpace + def.space > maxSpace) {
            alert('Not enough housing space! Upgrade or build more Army Camps.');
            return;
        }
        try {
            await train(troopType, 1);
        } catch (err) {
            console.error(err);
        }
    };

    const handleUpgradeStart=async (troopType) => {
        try {
            await startUpgrade(troopType);
        } catch (err) {
            alert(`Upgrade failed: ${err.message}`);
        }
    };

    const handleUpgradeFinish=async (troopType) => {
        try {
            await finishUpgrade(troopType);
        } catch (err) {
            console.error("Failed to finish upgrade:", err);
        }
    };

    const handleDiscard=async (troopType) => {
        try {
            await discard(troopType, 1);
        } catch (err) {
            alert(`Discard failed: ${err.message}`);
        }
    };

    const isCampFull=currentSpace >= maxSpace;

    return (
        <div style={styles.overlay}>
            <div style={styles.modal}>
                <div style={styles.header}>
                    <div style={styles.titleArea}>
                        <h2 style={styles.title}>Barracks & Laboratory</h2>
                        <div style={styles.spaceBadge}>
                            ⛺ Space: {currentSpace} / {maxSpace}
                        </div>
                    </div>
                    <button onClick={toggleArmy} style={styles.closeBtn}>X</button>
                </div>

                {isLabUpgrading && (
                    <div style={styles.alert}>
                        🧪 Laboratory is currently upgrading. Troop upgrades are paused.
                    </div>
                )}

                {maxSpace === 0 && !isLoading && (
                    <div style={styles.alertWarning}>
                        ⚠️ Build or complete upgrades on your Army Camps to house troops!
                    </div>
                )}

                {error && <div style={styles.error}>{error}</div>}

                {isLoading && troops.length === 0 ? (
                    <div style={styles.loading}>Loading Army...</div>
                ) : (
                    <div style={styles.grid}>
                        {troops.map((troop) => (
                            <TrainTroopCard
                                key={troop.troop_type}
                                troop={troop}
                                onTrain={handleTrain}
                                onDiscard={handleDiscard}
                                onUpgradeStart={handleUpgradeStart}
                                onUpgradeFinish={handleUpgradeFinish}
                                labLevel={labLevel}
                                elixirBalance={playerStats.elixir}
                                isCampFull={isCampFull || (currentSpace + (TROOP_DEFS[troop.troop_type]?.space || 1) > maxSpace)}
                            />
                        ))}
                    </div>
                )}
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
        backgroundColor: 'rgba(0, 0, 0, 0.65)',
        zIndex: 100,
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center',
        fontFamily: '"Luckiest Guy", cursive'
    },
    modal: {
        backgroundColor: '#ecf0f1',
        width: '90%',
        maxWidth: '750px',
        maxHeight: '85vh',
        borderRadius: '16px',
        border: '6px solid #e67e22',
        padding: '1.5rem',
        boxShadow: '0 10px 25px rgba(0,0,0,0.5)',
        display: 'flex',
        flexDirection: 'column',
        overflow: 'hidden'
    },
    header: {
        display: 'flex',
        justifyContent: 'space-between',
        alignItems: 'flex-start',
        borderBottom: '4px solid #bdc3c7',
        paddingBottom: '1rem',
        marginBottom: '1rem',
        color: '#2c3e50',
    },
    titleArea: {
        display: 'flex',
        flexDirection: 'column',
        gap: '0.4rem'
    },
    title: {
        margin: 0,
        fontSize: '1.8rem',
        textShadow: '1px 1px 0px rgba(255,255,255,0.8)'
    },
    spaceBadge: {
        alignSelf: 'flex-start',
        backgroundColor: '#34495e',
        color: '#f1c40f',
        padding: '0.3rem 0.8rem',
        borderRadius: '20px',
        fontSize: '1.1rem',
        border: '2px solid #2c3e50',
        textShadow: '1px 1px 0px #000'
    },
    closeBtn: {
        backgroundColor: '#e74c3c',
        color: 'white',
        border: 'none',
        borderRadius: '8px',
        padding: '0.5rem 1rem',
        cursor: 'pointer',
        fontSize: '1.5rem',
        fontFamily: 'inherit',
        boxShadow: '0 3px 0px #c0392b'
    },
    alert: {
        backgroundColor: '#d24dff',
        color: 'white',
        padding: '0.6rem 1rem',
        borderRadius: '8px',
        marginBottom: '1rem',
        fontSize: '0.95rem',
        border: '2px solid #b300b3',
        textShadow: '1px 1px 0px #000'
    },
    alertWarning: {
        backgroundColor: '#e67e22',
        color: 'white',
        padding: '0.6rem 1rem',
        borderRadius: '8px',
        marginBottom: '1rem',
        fontSize: '0.95rem',
        border: '2px solid #d35400',
        textShadow: '1px 1px 0px #000'
    },
    error: {
        backgroundColor: '#e74c3c',
        color: 'white',
        padding: '0.6rem 1rem',
        borderRadius: '8px',
        marginBottom: '1rem',
        fontSize: '0.95rem'
    },
    loading: {
        color: '#2c3e50',
        fontSize: '1.8rem',
        textAlign: 'center',
        padding: '3rem'
    },
    grid: {
        display: 'grid',
        gridTemplateColumns: 'repeat(auto-fit, minmax(200px, 1fr))',
        gap: '1rem',
        overflowY: 'auto',
        padding: '0.2rem'
    }
};
