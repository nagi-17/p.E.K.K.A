import * as PIXI from 'pixi.js';
import { GridOverlay } from '../entities/GridOverlay';
import { BuildingSprite } from '../entities/BuildingSprite';
import { gameConfig } from '../gameConfig';

export class VillageScene {
    constructor() {
        this.container=new PIXI.Container();
        this.container.sortableChildren=true;
        this.buildings=new Map();
        this.grid=new GridOverlay();
        this.grid.container.zIndex=-100;
        this.container.addChild(this.grid.container);
    }

    loadVillage(buildingData) {
        this.buildings.forEach(b => {
            if (typeof b.destroy === 'function') {
                b.destroy();
            }
            else {
                b.container.destroy({ children: true });
            }
        });
        this.buildings.clear();

        if (buildingData&&buildingData.length>0) {
            buildingData.forEach(data=>{
                this.addBuilding(data);
            });
        }
    }

    addBuilding(buildingData) {
        const building=new BuildingSprite(buildingData);

        building.container.zIndex=buildingData.posY;
        this.buildings.set(building.id, building);
        this.container.addChild(building.container);
    }

    isAreaFree(targetX, targetY, width, height, ignoreBuildingId=null) {
        if (targetX < 0 || targetY < 0 || targetX+width > gameConfig.GRID_WIDTH || targetY+height > gameConfig.GRID_HEIGHT) {
            return false;
        }

        for (const [id, building] of this.buildings) {
            if (id===ignoreBuildingId) 
                continue;

            const bX=building.gridX;
            const bY=building.gridY;
            const bW=building.gridWidth;
            const bH=building.gridHeight;

            const overlapX=targetX < (bX+bW) && (targetX+width) > bX;
            const overlapY=targetY < (bY+bH) && (targetY+height) > bY;

            if (overlapX && overlapY) 
                return false;
        }
        return true;
    }
}