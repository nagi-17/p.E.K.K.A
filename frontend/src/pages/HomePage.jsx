import AppShell from '../components/layout/AppShell';
import VillageCanvas from '../components/village/VillageCanvas';

export default function HomePage() {
    return (
        <AppShell>
            <VillageCanvas />
        </AppShell>
    );
}