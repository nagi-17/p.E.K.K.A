import { useState } from 'react';
import { useUiStore } from '../../store/uiStore';
import { useVillageStore } from '../../store/villageStore';
import { startUpgrade, cancelUpgrade, getVillage, getPlayerInfo } from '../../api/village';
import styles from './BuildingInfoPanel.module.css';

export default function BuildingInfoPanel() {
    const selectedId=useUiStore((state) => state.selectedBuildingId);
    const clearSelection=useUiStore((state) => state.clearSelection);

    const startMove=useUiStore((state) => state.startMove);
    
    const buildings=useVillageStore((state) => state.buildings);

    const setGameData=useVillageStore((state) => state.setGameData); 
    
    const [isUpgrading, setIsUpgrading]=useState(false);
    const [error, setError]=useState(null);

    const selectedBuilding=buildings.find(b => b.id === selectedId);

    if (!selectedBuilding) 
        return null;

    const handleUpgrade=async () => {
        setIsUpgrading(true);
        setError(null);
        try {
            await startUpgrade(selectedBuilding.id); 
            const [newVillageData, newStats]=await Promise.all([getVillage(), getPlayerInfo()]);
            setGameData(newVillageData, newStats);
        }
        catch (err) {
            setError(err.message);
        }
        finally {
            setIsUpgrading(false);
        }
    };

    const handleCancel=async () => {
        setIsUpgrading(true);
        setError(null);
        try {
            if (confirm("Are you sure you want to cancel this upgrade? You will get a 50% resource refund.")) {
                await cancelUpgrade(selectedBuilding.id); 
                const [newVillageData, newStats]=await Promise.all([getVillage(), getPlayerInfo()]);
                setGameData(newVillageData, newStats);
                clearSelection();
            }
        }
        catch (err) {
            setError(err.message);
        }
        finally {
            setIsUpgrading(false);
        }
    };

    const isCurrentlyUpgrading=!!selectedBuilding.upgradeCompleteAt;

    return (
        <div className={styles.panelOverlay}>
            <div className={styles.header}>
                <h2>{selectedBuilding.type}</h2>
                <button onClick={clearSelection} className={styles.closeBtn}>X</button>
            </div>
            
            <div className={styles.stats}>
                <p>Location: X: {selectedBuilding.posX}, Y: {selectedBuilding.posY}</p>
                <p>Level: {selectedBuilding.level || 1}</p>
                {error && <p className={styles.errorText}>{error}</p>}
            </div>

            <div className={styles.actions}>
                <button className={styles.moveBtn} onClick={() => startMove(selectedBuilding.id)}>Move</button>
                {isCurrentlyUpgrading ? (
                    <button className={styles.cancelBtn} onClick={handleCancel} disabled={isUpgrading}>
                        {isUpgrading ? 'Canceling...' : 'Cancel'}
                    </button>
                ) : (
                    <button className={styles.upgradeBtn} onClick={handleUpgrade} disabled={isUpgrading}>
                        {isUpgrading ? 'Upgrading...' : 'Upgrade'}
                    </button>
                )}
            </div>
        </div>
    );
}
