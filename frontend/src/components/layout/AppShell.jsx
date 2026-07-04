import TopBar from './TopBar';
import BottomNav from './BottomNav';
import styles from './AppShell.module.css';

export default function AppShell({children}) {
    return (
        <div className={styles.shell}>
            <TopBar /> 
            <div className={styles.mainContent}>
                {children}
            </div>
            <BottomNav />
        </div>
    );
}
