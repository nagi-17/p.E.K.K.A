import {useAuthStore} from '../../store/authStore'

export default function TopBar() {

    const logout=useAuthStore(function(state) {
        return state.logout;
    });

    return (
        <div style={styles.topBar}>
            <div style={styles.playerInfo}>
                <span>ironman</span>
            </div>
            <div style={styles.resources}>
                <div style={styles.resourceBadge}>
                    <span style={styles.icon}>🏆</span>
                    <span>0</span>
                </div>
                <div style={styles.resourceBadge}>
                    <span style={styles.icon}>🥞</span>
                    <span>1,000</span>
                </div>
                <div style={styles.resourceBadge}>
                    <svg width="24" height="24" viewBox="0 0 24 24" style={{ marginRight: '0.5rem' }}>
                        <path d="M12 2.5C12 2.5 5 10 5 15.5C5 19.09 7.91 22 11.5 22C15.09 22 18 19.09 18 15.5C18 10 12 2.5 12 2.5Z" fill="#d24dff"/>
                    </svg>
                    <span>1000</span>
                </div>
                <div style={styles.resourceBadge}>
                    <span style={styles.icon}>★</span>
                    <span>10</span>
                </div>
            </div>
            <button onClick={logout} style={styles.logoutBtn}>Logout</button>
        </div>
    );
}

const styles = {
    topBar: {
        display: 'flex',
        justifyContent: 'space-between',
        alignItems: 'center',
        padding: '0.5rem 1rem',
        backgroundColor: 'rgba(0, 0, 0, 0.6)',
        color: 'white',
        fontFamily: '"Luckiest Guy", cursive',
        fontSize: '1.25rem',
        boxShadow: '0 4px 6px rgba(0,0,0,0.3)',
        zIndex: 10,
    },
    playerInfo: {
        textShadow: '2px 2px 0px #000',
    },
    resources: {
        display: 'flex',
        gap: '1rem',
    },
    resourceBadge: {
        display: 'flex',
        alignItems: 'center',
        backgroundColor: 'rgba(255, 255, 255, 0.2)',
        padding: '0.25rem 0.75rem',
        borderRadius: '1rem',
        border: '2px solid rgba(255,255,255,0.4)',
    },
    icon: {
        marginRight: '0.5rem',
    },
    logoutBtn: {
        padding: '0.5rem 1rem',
        backgroundColor: '#e74c3c',
        color: 'white',
        border: '2px solid #c0392b',
        borderRadius: '0.5rem',
        fontFamily: '"Luckiest Guy", cursive',
        fontSize: '1rem',
        cursor: 'pointer',
        textShadow: '1px 1px 0px black'
    }
};