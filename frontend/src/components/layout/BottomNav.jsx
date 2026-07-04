import { useUiStore } from "../../store/uiStore";
import { useNavigate } from "react-router-dom";
import styles from './BottomNav.module.css';

export default function BottomNav() {
    const toggleShop=useUiStore((state) => state.toggleShop);
    const toggleArmy=useUiStore((state) => state.toggleArmy);
    const navigate=useNavigate();
    return (
        <div className={styles.bottomNav}>
            <button className={styles.navButton} onClick={toggleShop}>Shop</button>
            <button className={styles.navButton} onClick={toggleArmy}>Train Army</button>
            <button className={styles.navButton} onClick={() => navigate('/battle')}>Attack</button>
        </div>
    );
}
