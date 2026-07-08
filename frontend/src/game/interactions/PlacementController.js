import * as PIXI from 'pixi.js';
import { gameConfig, BUILDING_MAP } from '../gameConfig';
import { AssetLoader } from '../core/AssetLoader';
import { useUiStore } from '../../store/uiStore';
import { placeBuilding, getVillage, getPlayerInfo } from '../../api/village';
import { useVillageStore } from '../../store/villageStore';
import toast from 'react-hot-toast';

const PLACEMENT_DEFS={
    1001: { type: 'Town Hall', w: 4, h: 4 },
    2001: { type: 'Cannon', w: 3, h: 3 },
    2101: { type: 'Archer Tower', w: 3, h: 3 },
    2201: { type: 'Mortar', w: 3, h: 3 },
    3001: { type: 'Elixir Collector', w: 3, h: 3 },
    3101: { type: 'Pancake Machine', w: 3, h: 3 },
    4001: { type: 'Elixir Storage', w: 3, h: 3 },
    4101: { type: 'Pancake Stack', w: 3, h: 3 },
    5001: { type: 'Laboratory', w: 3, h: 3 },
    6001: { type: 'Army Camp', w: 5, h: 5 }
};

export class PlacementController {
    constructor(viewport, scene) {
        this.viewport=viewport;
        this.scene=scene;
        this.ghost=null;
        this.activeData=null;
        this.currentGridX=-1;
        this.currentGridY=-1;

        this.viewport.on('pointermove', this.onPointerMove, this);
        this.viewport.on('clicked', this.onPointerClick, this);
    }

    startPlacement(buildingDataId) {
        this.cancelPlacement();

        const alias=BUILDING_MAP[buildingDataId];
        const def=PLACEMENT_DEFS[buildingDataId] || { type: 'Cannon', w: 3, h: 3 };

        this.activeData={ id: buildingDataId, alias, type: def.type, width: def.w, height: def.h };

        const texture=alias ? AssetLoader.getTexture(alias):null;
        if (texture) {
            this.ghost=new PIXI.Sprite(texture);
        }
        else {
            this.ghost=new PIXI.Graphics();
            this.ghost.rect(0, 0, def.w * gameConfig.TILE_SIZE, def.h * gameConfig.TILE_SIZE).fill(0xffffff).stroke({ width: 2, color: 0x000000 });
        }

        this.ghost.width=def.w * gameConfig.TILE_SIZE;
        this.ghost.height=def.h * gameConfig.TILE_SIZE;
        this.ghost.alpha=0.6;
        this.ghost.zIndex=9999;
        
        this.scene.container.addChild(this.ghost);
    }

    onPointerMove(event) {
        if (!this.ghost)
            return;

        const localPos=this.scene.container.toLocal(event.global);

        const gridX=Math.floor(localPos.x / gameConfig.TILE_SIZE);
        const gridY=Math.floor(localPos.y / gameConfig.TILE_SIZE);

        if (gridX !== this.currentGridX || gridY !== this.currentGridY) {
            this.currentGridX=gridX;
            this.currentGridY=gridY;

            this.ghost.x=gridX * gameConfig.TILE_SIZE;
            this.ghost.y=gridY * gameConfig.TILE_SIZE;

            const isFree=this.scene.isAreaFree(gridX, gridY, this.activeData.width, this.activeData.height);
            this.ghost.tint=isFree ? 0x2ecc71 : 0xe74c3c; 
        }
    }

    async onPointerClick() {
        if (!this.ghost)
            return;

        const isFree=this.scene.isAreaFree(this.currentGridX, this.currentGridY, this.activeData.width, this.activeData.height);
        
        if (isFree) {
            const bType=this.activeData.type;
            const x=this.currentGridX;
            const y=this.currentGridY;

            this.cancelPlacement();
            useUiStore.getState().cancelPlacement();

            try {
                await placeBuilding({ building_type: bType, pos_x: x, pos_y: y });

                const [newVillage, newStats]=await Promise.all([ getVillage(), getPlayerInfo() ]);
                useVillageStore.getState().setGameData(newVillage, newStats);

            }
            catch (err) {
                toast.error(`Placement failed: ${err.message}`);
            }
        }
        else {
            this.ghost.alpha=0.2;
            setTimeout(() => { if(this.ghost) this.ghost.alpha=0.6; }, 100);
        }
    }

    cancelPlacement() {
        if (this.ghost) {
            this.ghost.destroy();
            this.ghost=null;
        }
        this.activeData=null;
    }
}