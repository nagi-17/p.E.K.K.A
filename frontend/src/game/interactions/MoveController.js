import * as PIXI from 'pixi.js';
import { gameConfig, BUILDING_MAP } from '../gameConfig';
import { AssetLoader } from '../core/AssetLoader';
import { useUiStore } from '../../store/uiStore';
import { useVillageStore } from '../../store/villageStore';
import { moveBuilding, getVillage } from '../../api/village';
import toast from 'react-hot-toast';

export class MoveController {
    constructor(viewport, scene) {
        this.viewport=viewport;
        this.scene=scene;
        this.ghost=null;
        this.activeBuilding=null; 
        
        this.viewport.on('pointermove', this.onPointerMove, this);
        this.viewport.on('clicked', this.onPointerClick, this);
    }

    startMove(instanceId) {
        this.cancelMove();
        const building=useVillageStore.getState().buildings.find(b => b.id === instanceId);
        if (!building) return;

        this.activeBuilding=building;
        const alias=BUILDING_MAP[building.buildingDataId];
        
        const realSprite=this.scene.buildings.get(instanceId);
        if (realSprite) realSprite.container.visible=false;

        const texture=alias ? AssetLoader.getTexture(alias) : null;
        if (texture) {
            this.ghost=new PIXI.Sprite(texture);
        } else {
            this.ghost=new PIXI.Graphics();
            this.ghost.rect(0, 0, building.width * gameConfig.TILE_SIZE, building.height * gameConfig.TILE_SIZE).fill(0xffffff).stroke({ width: 2, color: 0x000000 });
        }

        this.ghost.width=building.width * gameConfig.TILE_SIZE;
        this.ghost.height=building.height * gameConfig.TILE_SIZE;
        this.ghost.alpha=0.6;

        this.ghost.x=building.posX * gameConfig.TILE_SIZE;
        this.ghost.y=building.posY * gameConfig.TILE_SIZE;
        this.currentGridX=building.posX;
        this.currentGridY=building.posY;

        this.scene.container.addChild(this.ghost);
    }

    onPointerMove(event) {
        if (!this.ghost || !this.activeBuilding) return;

        const localPos=this.scene.container.toLocal(event.global);
        const gridX=Math.floor(localPos.x / gameConfig.TILE_SIZE);
        const gridY=Math.floor(localPos.y / gameConfig.TILE_SIZE);

        if (gridX !== this.currentGridX || gridY !== this.currentGridY) {
            this.currentGridX=gridX;
            this.currentGridY=gridY;
            this.ghost.x=gridX * gameConfig.TILE_SIZE;
            this.ghost.y=gridY * gameConfig.TILE_SIZE;

            const isFree=this.scene.isAreaFree(gridX, gridY, this.activeBuilding.width, this.activeBuilding.height, this.activeBuilding.id);
            this.ghost.tint=isFree ? 0x2ecc71 : 0xe74c3c; 
        }
    }

    async onPointerClick() {
        if (!this.ghost || !this.activeBuilding) return;

        const isFree=this.scene.isAreaFree(this.currentGridX, this.currentGridY, this.activeBuilding.width, this.activeBuilding.height, this.activeBuilding.id);
        
        if (isFree) {
            const bId=this.activeBuilding.id;
            const x=this.currentGridX;
            const y=this.currentGridY;

            this.cancelMove();
            useUiStore.getState().cancelMove();

            try {
                await moveBuilding(bId, x, y);
                const newVillage=await getVillage();
                useVillageStore.getState().setVillage(newVillage);
            }
            catch (err) {
                toast.error(`Move failed: ${err.message}`);
                const realSprite=this.scene.buildings.get(bId);
                if (realSprite) realSprite.container.visible=true;
            }
        }
        else {
            this.ghost.alpha=0.2;
            setTimeout(() => { if(this.ghost) this.ghost.alpha=0.6; }, 100);
        }
    }

    cancelMove() {
        if (this.ghost) {
            this.ghost.destroy();
            this.ghost=null;
        }
        if (this.activeBuilding) {
            const realSprite=this.scene.buildings.get(this.activeBuilding.id);
            if (realSprite) realSprite.container.visible=true;
        }
        this.activeBuilding=null;
    }
}