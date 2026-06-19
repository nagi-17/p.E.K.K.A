import * as PIXI from 'pixi.js';
import { createViewport } from './Viewport';
import { VillageScene } from '../scenes/VillageScene';
import { gameConfig } from '../gameConfig';
import { AssetLoader } from './AssetLoader';

export class GameApp {
    constructor() {
        this.app=new PIXI.Application();
    }

    async init(parentContainer) {
        await this.app.init({
            resizeTo: parentContainer,
            backgroundColor: 0x7ec850,
            resolution: window.devicePixelRatio || 1,
            autoDensity: true,
        });

        await AssetLoader.loadAssets();

        this.viewport=createViewport(this.app, parentContainer);
        this.app.stage.addChild(this.viewport);

        this.scene=new VillageScene();
        this.viewport.addChild(this.scene.container);

        const centerX= (gameConfig.GRID_WIDTH*gameConfig.TILE_SIZE)/2;
        const centerY= (gameConfig.GRID_HEIGHT*gameConfig.TILE_SIZE)/2;
        this.viewport.moveCenter(centerX, centerY);
    }

    mount(domElement) {
        domElement.appendChild(this.app.canvas);
    }

    destroy() {
        if (this.app) {
            this.app.destroy(true, {children: true, texture: true, baseTexture: true});
        }
    }
}