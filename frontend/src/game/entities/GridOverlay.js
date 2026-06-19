import * as PIXI from 'pixi.js';
import { gameConfig } from '../gameConfig';

export class GridOverlay {
    constructor() {
        this.container=new PIXI.Container();
        this.graphics=new PIXI.Graphics();
        this.container.addChild(this.graphics);
        this.drawGrid();
    }

    drawGrid() {
        const {TILE_SIZE, GRID_WIDTH, GRID_HEIGHT}=gameConfig;
        const base_width=TILE_SIZE * GRID_WIDTH;
        const base_height=TILE_SIZE * GRID_HEIGHT;

        this.graphics.clear();
        const lineStroke={width: 2, color: 0xffffff, alpha: 0.4};

        for (let i=0; i<=GRID_WIDTH; i++) {
            const x=i*TILE_SIZE;
            this.graphics.moveTo(x, 0).lineTo(x, base_height).stroke(lineStroke);
        }
        for (let j=0; j<=GRID_HEIGHT; j++) {
            const y=j*TILE_SIZE;
            this.graphics.moveTo(0, y).lineTo(base_width, y).stroke(lineStroke);
        }

        this.graphics.rect(0, 0, base_width, base_height).stroke({width: 6, color: 0xff0000});
    }
}