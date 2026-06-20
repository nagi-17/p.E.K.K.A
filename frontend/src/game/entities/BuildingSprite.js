import * as PIXI from 'pixi.js';
import { gameConfig, BUILDING_MAP } from '../gameConfig';
import { AssetLoader } from '../core/AssetLoader';
import { useUiStore } from '../../store/uiStore';
import { collectResource,getVillage,getPlayerInfo, finishUpgrade } from '../../api/village';
import {useVillageStore} from '../../store/villageStore'

export class BuildingSprite {
    constructor(buildingData) {
        this.id=buildingData.id;
        this.type=buildingData.type;
        this.typeAlias=BUILDING_MAP[buildingData.buildingDataId];
        this.container=new PIXI.Container();

        this.gridWidth=buildingData.width;
        this.gridHeight=buildingData.height;
        this.pixelWidth=this.gridWidth * gameConfig.TILE_SIZE;
        this.pixelHeight=this.gridHeight * gameConfig.TILE_SIZE;
        
        let texture=null;
        if (this.typeAlias) {
            texture=AssetLoader.getTexture(this.typeAlias);
        }

        if (texture) {
            this.sprite=new PIXI.Sprite(texture);
        } else {
            this.sprite=new PIXI.Graphics();
            this.sprite.rect(0, 0, this.pixelWidth, this.pixelHeight).fill(0x3498db).stroke({width: 2, color: 0x2980b9});
        }
        this.sprite.width=this.pixelWidth;
        this.sprite.height=this.pixelHeight;
        this.container.addChild(this.sprite);

        this.upgradeCompleteAt=buildingData.upgradeCompleteAt ? new Date(buildingData.upgradeCompleteAt).getTime() : null;
        this.totalUpgradeTimeMs=(buildingData.upgradeTime || 60) * 1000;
        
        this.progressBar=null;
        this.ticker=null;

        if (this.upgradeCompleteAt && this.upgradeCompleteAt>Date.now()) {
            this.setupConstructionVisuals();
        } else {
            this.lastCollectedAt=buildingData.lastCollectedAt ? new Date(buildingData.lastCollectedAt).getTime() : Date.now();
            this.setupResourceBubble();
        }

        this.updatePosition(buildingData.posX, buildingData.posY);
        this.setupInteraction();
    }

    setupConstructionVisuals() {
        this.sprite.tint=0x888888; 

        const label=new PIXI.Text({
            text: 'Upgrading...',
            style: { fontFamily: '"Luckiest Guy", cursive', fontSize: 16, fill: 0xffffff, stroke: { color: 0x000000, width: 3 } }
        });
        label.anchor.set(0.5);
        label.position.set(this.pixelWidth/2, -20);
        this.container.addChild(label);

        this.progressBar=new PIXI.Graphics();
        this.progressBar.position.set(0, -10);
        this.container.addChild(this.progressBar);

        this.ticker=PIXI.Ticker.shared.add(this.updateProgress, this);
    }

    updateProgress() {
        const now=Date.now();
        const timeLeft=this.upgradeCompleteAt-now;

        if (timeLeft <= 0) {
            PIXI.Ticker.shared.remove(this.updateProgress, this);
            this.sprite.tint=0xFFFFFF;
            if (this.progressBar) {
                this.container.removeChild(this.progressBar);
                this.progressBar=null;
            }

            this.upgradeCompleteAt=null;
            this.completeUpgradeOnBackend();
            return;
        }

        let percent=1-(timeLeft / this.totalUpgradeTimeMs);
        if (percent<0)
            percent=0;

        this.progressBar.clear();
        this.progressBar.rect(0, 0, this.pixelWidth, 8).fill(0xe74c3c).stroke({ width: 2, color: 0x000000 });
        this.progressBar.rect(0, 0, this.pixelWidth * percent, 8).fill(0x2ecc71);
    }

    async completeUpgradeOnBackend() {
        try {
            await finishUpgrade(this.id);
            const [newVillage, newStats]=await Promise.all([getVillage(), getPlayerInfo()]);
            useVillageStore.getState().setGameData(newVillage, newStats);
            console.log(`Upgrade complete for building ${this.id}!`);
        }
        catch (err) {
            console.error("Failed to finish upgrade on backend:", err);
        }
    }

    setupInteraction() {
        this.container.eventMode='static';
        this.container.cursor='pointer';

        this.container.on('pointerdown', (event) => {
            event.stopPropagation(); 
            useUiStore.getState().selectBuilding(this.id);
            this.container.alpha=0.8;
            setTimeout(() => {this.container.alpha=1;}, 100);
        });
    }

    updatePosition(gridX, gridY) {
        this.gridX=gridX;
        this.gridY=gridY;
        this.container.x=gridX*gameConfig.TILE_SIZE;
        this.container.y=gridY*gameConfig.TILE_SIZE;
    }

    setupResourceBubble() {
        const isGenerator=this.type === 'Elixir Collector' || this.type === 'Pancake Machine';
        if (!isGenerator) return;

        const msSinceLastCollect=Date.now() - this.lastCollectedAt;
        if (msSinceLastCollect < 60000) return; 

        this.bubbleContainer=new PIXI.Container();

        const circle=new PIXI.Graphics();
        circle.circle(0, 0, 20).fill(0xffffff).stroke({ width: 3, color: 0x000000 });

        const iconText=this.type === 'Pancake Machine' ? '🥞' : '💧';
        const icon=new PIXI.Text({
            text: iconText,
            style: { fontSize: 24 }
        });
        icon.anchor.set(0.5);

        this.bubbleContainer.addChild(circle, icon);

        this.bubbleContainer.position.set(this.pixelWidth / 2, -15);

        this.bubbleContainer.eventMode='static';
        this.bubbleContainer.cursor='pointer';

        this.bubbleTime=0;
        this.animateBubble=() => {
            if (!this.bubbleContainer || this.bubbleContainer.destroyed) return;
            this.bubbleTime+=0.05;
            this.bubbleContainer.y=-15+Math.sin(this.bubbleTime)*5;
        };

        PIXI.Ticker.shared.add(this.animateBubble);

        this.bubbleContainer.on('pointerdown', async (event) => {
            event.stopPropagation();
            this.bubbleContainer.visible=false; 

            try {
                await collectResource(this.id);
                const [newVillage, newStats]=await Promise.all([getVillage(), getPlayerInfo()]);
                useVillageStore.getState().setGameData(newVillage, newStats);
                
            }
            catch (err) {
                console.error("Collect failed:", err);
                this.bubbleContainer.visible=true;
            }
        });

        this.container.addChild(this.bubbleContainer);
    }

    destroy() {
        if (this.ticker) {
            PIXI.Ticker.shared.remove(this.updateProgress, this);
        }
        if (this.animateBubble) {
            PIXI.Ticker.shared.remove(this.animateBubble);
        }
        this.container.destroy({ children: true });
    }
}