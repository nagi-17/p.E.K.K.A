import { useUiStore } from '../../store/uiStore';
import { useVillageStore } from '../../store/villageStore';

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
        <div style={styles.overlay}>
            <div style={styles.modal}>
                <div style={styles.header}>
                    <h2>Village Shop</h2>
                    <button onClick={toggleShop} style={styles.closeBtn}>X</button>
                </div>
                
                <div style={styles.grid}>
                    {SHOP_ITEMS.map((item) => {
                        const currentBalance=playerStats[item.currency] || 0;
                        const canAfford=currentBalance >= item.cost;

                        return (
                            <div key={item.id} style={styles.card}>
                                <h3>{item.name}</h3>
                                <p style={{ ...styles.cost, color: canAfford ? '#f1c40f' : '#e74c3c' }}>
                                    Cost: {item.cost} {item.icon}
                                </p>
                                <button style={{...styles.buyBtn, backgroundColor: canAfford ? '#2ecc71' : '#7f8c8d',cursor: canAfford ? 'pointer' : 'not-allowed'}} 
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

const styles={
    overlay: { position: 'absolute', top: 0, left: 0, right: 0, bottom: 0, backgroundColor: 'rgba(0, 0, 0, 0.6)', zIndex: 100, display: 'flex', justifyContent: 'center', alignItems: 'center', fontFamily: '"Luckiest Guy", cursive' },
    modal: { backgroundColor: '#ecf0f1', width: '80%', maxWidth: '600px', borderRadius: '16px', border: '6px solid #f39c12', padding: '1.5rem', boxShadow: '0 10px 25px rgba(0,0,0,0.5)' },
    header: { display: 'flex', justifyContent: 'space-between', alignItems: 'center', borderBottom: '4px solid #bdc3c7', paddingBottom: '1rem', marginBottom: '1rem', color: '#2c3e50', fontSize: '1.5rem' },
    closeBtn: { backgroundColor: '#e74c3c', color: 'white', border: 'none', borderRadius: '8px', padding: '0.5rem 1rem', cursor: 'pointer', fontSize: '1.5rem', fontFamily: 'inherit' },
    grid: { display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(150px, 1fr))', gap: '1rem' },
    card: { backgroundColor: '#34495e', color: 'white', padding: '1rem', borderRadius: '12px', display: 'flex', flexDirection: 'column', alignItems: 'center', textAlign: 'center', border: '3px solid #2c3e50' },
    cost: { margin: '0.5rem 0', fontWeight: 'bold' },
    buyBtn: { width: '100%', color: 'white', border: 'none', borderRadius: '6px', padding: '0.5rem', fontSize: '1.2rem', fontFamily: 'inherit', marginTop: 'auto' }
};