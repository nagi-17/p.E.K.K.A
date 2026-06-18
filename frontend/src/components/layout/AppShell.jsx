import TopBar from './TopBar';
import BottomNav from './BottomNav';

export default function AppShell({children}) {
    return (
        <div style={styles.shell}>
            <TopBar /> 
            <div style={styles.mainContent}>
                {children}
            </div>
            <BottomNav />
        </div>
    );
}

const styles = {
    shell: {
        display: 'flex',
        flexDirection: 'column',
        height: '100dvh',
        width: '100vw',
        overflow: 'hidden',
        backgroundColor: '#2ecc71',
    },
    mainContent: {
        flex: 1,
        position: 'relative',
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center',
    }
};