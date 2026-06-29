import { useState } from 'react';
import { useUiStore } from '../../store/uiStore';
import { useVillageStore } from '../../store/villageStore';
import { startUpgrade, cancelUpgrade, getVillage, getPlayerInfo } from '../../api/village';

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
        <div style={styles.panelOverlay}>
            <div style={styles.header}>
                <h2>{selectedBuilding.type}</h2>
                <button onClick={clearSelection} style={styles.closeBtn}>X</button>
            </div>
            
            <div style={styles.stats}>
                <p>Location: X: {selectedBuilding.posX}, Y: {selectedBuilding.posY}</p>
                <p>Level: {selectedBuilding.level || 1}</p>
                {error && <p style={styles.errorText}>{error}</p>}
            </div>

            <div style={styles.actions}>
                <button style={styles.moveBtn} onClick={() => startMove(selectedBuilding.id)}>Move</button>
                {isCurrentlyUpgrading ? (
                    <button style={styles.cancelBtn} onClick={handleCancel} disabled={isUpgrading}>
                        {isUpgrading ? 'Canceling...' : 'Cancel'}
                    </button>
                ) : (
                    <button style={styles.upgradeBtn} onClick={handleUpgrade} disabled={isUpgrading}>
                        {isUpgrading ? 'Upgrading...' : 'Upgrade'}
                    </button>
                )}
            </div>
        </div>
    );
}

const styles={
    panelOverlay: {
        position: 'absolute', bottom: '5.625rem', left: '50%', transform: 'translateX(-50%)',
        width: '20rem', backgroundColor: 'rgba(0, 0, 0, 0.85)', color: '#fff',
        borderRadius: '0.75rem', border: '3px solid #f1c40f', padding: '1rem',
        zIndex: 50, fontFamily: '"Luckiest Guy", cursive',
    },
    header: {
        display: 'flex', justifyContent: 'space-between', alignItems: 'center',
        borderBottom: '2px solid #333', paddingBottom: '0.5rem', marginBottom: '0.5rem',
    },
    closeBtn: {
        backgroundColor: '#e74c3c', color: 'white', border: 'none', borderRadius: '6px',
        padding: '0.25rem 0.75rem', cursor: 'pointer', fontSize: '1.2rem', fontFamily: 'inherit',
    },
    stats: { fontSize: '1.1rem', marginBottom: '1rem' },
    actions: { display: 'flex', justifyContent: 'center', gap: '0.5rem' },
    errorText: { color: '#e74c3c', fontSize: '0.9rem', marginTop: '0.5rem' },

    moveBtn: {
        flex: 1, backgroundColor: '#3498db', color: 'white', border: '2px solid #2980b9',
        borderRadius: '8px', padding: '0.75rem', fontSize: '1.2rem', cursor: 'pointer', fontFamily: 'inherit',
    },
    upgradeBtn: {
        flex: 1, backgroundColor: '#2ecc71', color: 'white', border: '2px solid #27ae60',
        borderRadius: '8px', padding: '0.75rem', fontSize: '1.2rem', cursor: 'pointer', fontFamily: 'inherit',
    },
    cancelBtn: {
        flex: 1, backgroundColor: '#e74c3c', color: 'white', border: '2px solid #c0392b',
        borderRadius: '8px', padding: '0.75rem', fontSize: '1.2rem', cursor: 'pointer', fontFamily: 'inherit',
    }
};