import * as PIXI from 'pixi.js';
import { gameConfig, BUILDING_MAP } from '../gameConfig';
import { AssetLoader } from '../core/AssetLoader';

export class BuildingSprite {
    constructor(buildingData) {
        this.id=buildingData.ID; 
        this.typeAlias=BUILDING_MAP[buildingData.BuildingDataID];
        this.container=new PIXI.Container();
        
        let texture=null;
        const pixelSize=gameConfig.TILE_SIZE * 3;

        if (this.typeAlias) {
            texture = AssetLoader.getTexture(this.typeAlias);
        }

        if (texture) {
            this.sprite=new PIXI.Sprite(texture);
            this.sprite.width=pixelSize;
            this.sprite.height=pixelSize;
            this.container.addChild(this.sprite);
        }
        else {
            const temp=new PIXI.Graphics();
            temp.rect(0, 0, pixelSize, pixelSize).fill(0x3498db).stroke({width: 2, color: 0x2980b9});
            this.container.addChild(temp);
        }

        this.updatePosition(buildingData.PosX, buildingData.PosY);
    }

    updatePosition(gridX, gridY) {
        this.container.x=gridX*gameConfig.TILE_SIZE;
        this.container.y=gridY*gameConfig.TILE_SIZE;
    }
}