import { useUiStore } from '../../store/uiStore';
import { useVillageStore } from '../../store/villageStore';
import styles from './Shop.module.css';

const SHOP_ITEMS=[
    { id: 2001, name: 'Cannon', type: 'Cannon', cost: 250, currency: 'pancakes', icon: '🥞' },
    { id: 2101, name: 'Archer Tower', type: 'Archer Tower', cost: 1000, currency: 'pancakes', icon: '🥞' },
    { id: 2201, name: 'Mortar', type: 'Mortar', cost: 8000, currency: 'pancakes', icon: '🥞' },
    { id: 3001, name: 'Elixir Collector', type: 'Elixir Collector', cost: 150, currency: 'pancakes', icon: '🥞' },
    { id: 3101, name: 'Pancake Machine', type: 'Pancake Machine', cost: 150, currency: 'elixir', icon: '💧' },
    { id: 4001, name: 'Elixir Storage', type: 'Elixir Storage', cost: 300, currency: 'pancakes', icon: '🥞' },
    { id: 4101, name: 'Pancake Stack', type: 'Pancake Stack', cost: 300, currency: 'elixir', icon: '💧' },
    { id: 5001, name: 'Laboratory', type: 'Laboratory', cost: 25000, currency: 'elixir', icon: '💧' },
    { id: 6001, name: 'Army Camp', type: 'Army Camp', cost: 250, currency: 'elixir', icon: '💧' },
];

export default function Shop() {
    const isShopOpen=useUiStore((state) => state.isShopOpen);
    const toggleShop=useUiStore((state) => state.toggleShop);
    const startPlacement=useUiStore((state) => state.startPlacement);

    const playerStats=useVillageStore((state) => state.playerStats);

    if (!isShopOpen || !playerStats) return null;

    return (
        <div className={styles.overlay}>
            <div className={styles.modal}>
                <div className={styles.header}>
                    <h2>Village Shop</h2>
                    <button onClick={toggleShop} className={styles.closeBtn}>X</button>
                </div>
                
                <div className={styles.grid}>
                    {SHOP_ITEMS.map((item) => {
                        const currentBalance=playerStats[item.currency] || 0;
                        const canAfford=currentBalance >= item.cost;

                        return (
                            <div key={item.id} className={styles.card}>
                                <h3>{item.name}</h3>
                                <p className={styles.cost} style={{ color: canAfford ? '#f1c40f' : '#e74c3c' }}>
                                    Cost: {item.cost} {item.icon}
                                </p>
                                <button className={`${styles.buyBtn} ${canAfford ? styles.buyBtnAffordable : styles.buyBtnExpensive}`} 
                                    onClick={() => startPlacement(item.id)} disabled={!canAfford}>
                                    {canAfford ? 'Buy' : 'Not Enough!'}
                                </button>
                            </div>
                        );
                    })}
                </div>
            </div>
        </div>
    );
}
