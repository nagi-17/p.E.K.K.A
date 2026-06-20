import { useState, useEffect } from 'react';
import {useAuthStore} from '../../store/authStore'
import { useVillageStore } from '../../store/villageStore';

export default function TopBar() {

    const logout=useAuthStore(function(state) {
        return state.logout;
    });
    const playerStats=useVillageStore((state) => state.playerStats);
    const [shieldText, setShieldText]=useState('No Shield');

    useEffect(() => {
        if (!playerStats || !playerStats.shieldEndTime) {
            setShieldText('No Shield');
            return;
        }

        const updateShield=() => {
            const endTime=new Date(playerStats.shieldEndTime).getTime();
            const now=Date.now();
            const diff=endTime - now;
            if (diff <= 0) {
                setShieldText('No Shield');
            } else {
                const h=Math.floor(diff / (3600 * 1000));
                const m=Math.floor((diff % (3600 * 1000)) / (60 * 1000));
                setShieldText(`🛡️ ${h}h ${m}m`);
            }
        };

        updateShield();
        const interval=setInterval(updateShield, 60000);
        return () => clearInterval(interval);
    }, [playerStats?.shieldEndTime]);

    if (!playerStats) {
        return <div style={styles.navContainer}>Loading Profile...</div>;
    }

    return (
        <div style={styles.navContainer}>
            <div style={styles.leftSection}>
                <h2 style={styles.username}>{playerStats.username || "VILLAGER"}</h2>
            </div>

            <div style={styles.centerSection}>
                <div style={styles.badge} title="Shield Status">
                    <span>{shieldText}</span>
                </div>

                <div style={styles.badge}>
                    <span role="img" aria-label="Trophies">🏆</span>
                    <span>{playerStats.trophies.toLocaleString()}</span>
                </div>
                
                <div style={styles.badge}>
                    <span role="img" aria-label="Pancakes">🥞</span>
                    <span>{playerStats.pancakes.toLocaleString()} / {(playerStats.maxPancakes || 1500).toLocaleString()}</span>
                </div>
                
                <div style={styles.badge}>
                    <span role="img" aria-label="Elixir">
                        <svg width="24" height="24" viewBox="0 0 24 24" style={{ marginRight: '0.5rem' }}>
                        <path d="M12 2.5C12 2.5 5 10 5 15.5C5 19.09 7.91 22 11.5 22C15.09 22 18 19.09 18 15.5C18 10 12 2.5 12 2.5Z" fill="#d24dff"/>
                    </svg>
                    </span>
                    <span>{playerStats.elixir.toLocaleString()} / {(playerStats.maxElixir || 1500).toLocaleString()}</span>
                </div>

                <div style={styles.badge}>
                    <span role="img" aria-label="Skill">⭐</span>
                    <span>{playerStats.skillPoints.toLocaleString()}</span>
                </div>
            </div>

            <div style={styles.rightSection}>
                <button style={styles.logoutBtn} onClick={logout}>LOGOUT</button>
            </div>
        </div>
    );
}

const styles={
    navContainer: {
        position: 'absolute', top: 0, left: 0, right: 0,
        height: '60px', backgroundColor: '#27ae60',
        display: 'flex', justifyContent: 'space-between', alignItems: 'center',
        padding: '0 1rem', zIndex: 10,
        fontFamily: '"Luckiest Guy", cursive',
        boxShadow: '0px 4px 10px rgba(0,0,0,0.3)',
    },
    leftSection: { flex: 1 },
    username: { color: 'white', margin: 0, textShadow: '2px 2px 0px #000' },
    centerSection: {
        flex: 2, display: 'flex', justifyContent: 'center', gap: '1rem'
    },
    badge: {
        backgroundColor: 'rgba(0,0,0,0.4)', color: 'white',
        padding: '0.25rem 0.75rem', borderRadius: '20px',
        display: 'flex', alignItems: 'center', gap: '0.5rem',
        border: '2px solid rgba(255,255,255,0.2)', fontSize: '1.2rem',
    },
    rightSection: { flex: 1, display: 'flex', justifyContent: 'flex-end' },
    logoutBtn: {
        backgroundColor: '#e74c3c', color: 'white', border: '2px solid #c0392b',
        borderRadius: '8px', padding: '0.5rem 1rem', cursor: 'pointer',
        fontFamily: 'inherit', fontSize: '1rem', textShadow: '1px 1px 0px #000'
    }
};