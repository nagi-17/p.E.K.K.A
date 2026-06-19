import * as PIXI from 'pixi.js';
import { GridOverlay } from '../entities/GridOverlay';
import { BuildingSprite } from '../entities/BuildingSprite';

export class VillageScene {
    constructor() {
        this.container=new PIXI.Container();
        this.buildings=new Map();
        this.grid=new GridOverlay();
        this.container.addChild(this.grid.container);
    }

    loadVillage(buildingData) {
        this.buildings.forEach(b=>b.container.destroy());
        this.buildings.clear();

        if (buildingData&&buildingData.length>0)
        {
            buildingData.forEach(data=>{
                this.addBuilding(data);
            });
        }
    }

    addBuilding(buildingData) {
        const building=new BuildingSprite(buildingData);
        this.buildings.set(building.id, building);
        this.container.addChild(building.container);
    }
}