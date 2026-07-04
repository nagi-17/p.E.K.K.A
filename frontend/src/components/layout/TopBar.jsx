import { useState, useEffect } from 'react';
import {useAuthStore} from '../../store/authStore'
import { useVillageStore } from '../../store/villageStore';
import styles from './TopBar.module.css';

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
        return <div className={styles.navContainer}>Loading Profile...</div>;
    }

    return (
        <div className={styles.navContainer}>
            <div className={styles.leftSection}>
                <h2 className={styles.username}>{playerStats.username || "VILLAGER"}</h2>
            </div>

            <div className={styles.centerSection}>
                <div className={styles.badge} title="Shield Status">
                    <span>{shieldText}</span>
                </div>

                <div className={styles.badge}>
                    <span role="img" aria-label="Trophies">🏆</span>
                    <span>{playerStats.trophies.toLocaleString()}</span>
                </div>
                
                <div className={styles.badge}>
                    <span role="img" aria-label="Pancakes">🥞</span>
                    <span>{playerStats.pancakes.toLocaleString()} / {(playerStats.maxPancakes || 1500).toLocaleString()}</span>
                </div>
                
                <div className={styles.badge}>
                    <span role="img" aria-label="Elixir">
                        <svg width="24" height="24" viewBox="0 0 24 24" style={{ marginRight: '0.5rem' }}>
                        <path d="M12 2.5C12 2.5 5 10 5 15.5C5 19.09 7.91 22 11.5 22C15.09 22 18 19.09 18 15.5C18 10 12 2.5 12 2.5Z" fill="#d24dff"/>
                    </svg>
                    </span>
                    <span>{playerStats.elixir.toLocaleString()} / {(playerStats.maxElixir || 1500).toLocaleString()}</span>
                </div>

                <div className={styles.badge}>
                    <span role="img" aria-label="Skill">⭐</span>
                    <span>{playerStats.skillPoints.toLocaleString()}</span>
                </div>
            </div>

            <div className={styles.rightSection}>
                <button className={styles.logoutBtn} onClick={logout}>LOGOUT</button>
            </div>
        </div>
    );
}
