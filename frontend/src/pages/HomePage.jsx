import AppShell from '../components/layout/AppShell';
import Shop from '../components/shop/Shop';
import BuildingInfoPanel from '../components/village/BuildingInfoPanel';
import VillageCanvas from '../components/village/VillageCanvas';
import ArmyPanel from '../components/army/ArmyPanel';

export default function HomePage() {
    return (
        <AppShell>
            <VillageCanvas />
            <BuildingInfoPanel />
            <Shop />
            <ArmyPanel />
        </AppShell>
    );
}