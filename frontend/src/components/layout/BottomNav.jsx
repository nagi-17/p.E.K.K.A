export default function BottomNav() {
    return (
        <div style={styles.bottomNav}>
            <button style={styles.navButton}>Shop</button>
            <button style={styles.navButton}>Train Army</button>
            <button style={styles.navButton}>Attack</button>
        </div>
    );
}

const styles = {
    bottomNav: {
        display: 'flex',
        justifyContent: 'space-around',
        padding: '1rem',
        backgroundColor: 'rgba(0, 0, 0, 0.6)',
        zIndex: 10,
    },
    navButton: {
        fontFamily: '"Luckiest Guy", cursive',
        fontSize: '1.5rem',
        padding: '0.75rem 2rem',
        backgroundColor: '#f1c40f',
        color: '#fff',
        border: '0.25rem solid #e67e22',
        borderRadius: '0.5rem',
        cursor: 'pointer',
        textShadow: '2px 2px 0px #000',
        boxShadow: '0 0.25rem 0px #d35400',
        transition: 'transform 0.1s',
    }
};