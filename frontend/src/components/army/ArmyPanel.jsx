import { useEffect } from 'react';
import { useUiStore } from '../../store/uiStore';
import { useArmyStore } from '../../store/armyStore';
import { useVillageStore } from '../../store/villageStore';
import TrainTroopCard from './TrainTroopCard';
import { TROOP_DEFS } from '../../game/gameConfig';
import styles from './ArmyPanel.module.css';

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
        <div className={styles.overlay}>
            <div className={styles.modal}>
                <div className={styles.header}>
                    <div className={styles.titleArea}>
                        <h2 className={styles.title}>Barracks & Laboratory</h2>
                        <div className={styles.spaceBadge}>
                            ⛺ Space: {currentSpace} / {maxSpace}
                        </div>
                    </div>
                    <button onClick={toggleArmy} className={styles.closeBtn}>X</button>
                </div>

                {isLabUpgrading && (
                    <div className={styles.alert}>
                        🧪 Laboratory is currently upgrading. Troop upgrades are paused.
                    </div>
                )}

                {maxSpace === 0 && !isLoading && (
                    <div className={styles.alertWarning}>
                        ⚠️ Build or complete upgrades on your Army Camps to house troops!
                    </div>
                )}

                {error && <div className={styles.error}>{error}</div>}

                {isLoading && troops.length === 0 ? (
                    <div className={styles.loading}>Loading Army...</div>
                ) : (
                    <div className={styles.grid}>
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

