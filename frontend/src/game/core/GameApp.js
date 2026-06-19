import * as PIXI from 'pixi.js';

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