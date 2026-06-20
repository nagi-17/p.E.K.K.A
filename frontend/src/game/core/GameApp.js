import * as PIXI from 'pixi.js';
import { createViewport } from './Viewport';
import { VillageScene } from '../scenes/VillageScene';
import { gameConfig } from '../gameConfig';
import { AssetLoader } from './AssetLoader';
import { useUiStore } from '../../store/uiStore';
import { PlacementController } from '../interactions/PlacementController';
import { MoveController } from '../interactions/MoveController';

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

        this.viewport.on('clicked', () => {
            useUiStore.getState().clearSelection();
        });

        this.placementController=new PlacementController(this.viewport, this.scene);
        this.moveController=new MoveController(this.viewport, this.scene);

        this.unsubscribeUi=useUiStore.subscribe((state, prevState) => {
            if (state.placementModeId && state.placementModeId !== prevState.placementModeId) {
                this.placementController.startPlacement(state.placementModeId);
            } 
            else if (!state.placementModeId && prevState.placementModeId) {
                this.placementController.cancelPlacement();
            }

            if (state.moveModeId && state.moveModeId !== prevState.moveModeId) {
                this.moveController.startMove(state.moveModeId);
            } else if (!state.moveModeId && prevState.moveModeId) {
                this.moveController.cancelMove();
            }
        });

        const centerX= (gameConfig.GRID_WIDTH*gameConfig.TILE_SIZE)/2;
        const centerY= (gameConfig.GRID_HEIGHT*gameConfig.TILE_SIZE)/2;
        this.viewport.moveCenter(centerX, centerY);

    }

    mount(domElement) {
        domElement.appendChild(this.app.canvas);
    }

    destroy() {
        if(this.unsubscribeUi)
            this.unsubscribeUi();
        if (this.app) {
            this.app.destroy(true, {children: true});
        }
    }
}