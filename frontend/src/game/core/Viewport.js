import { Viewport } from 'pixi-viewport';
import { gameConfig } from '../gameConfig';

export function createViewport(app, parentDOMElement) {
    const base_width=gameConfig.GRID_WIDTH * gameConfig.TILE_SIZE;
    const base_height=gameConfig.GRID_HEIGHT * gameConfig.TILE_SIZE;

    const viewport=new Viewport({
        screenWidth: parentDOMElement.clientWidth,
        screenHeight: parentDOMElement.clientHeight,
        base_width: base_width,
        base_height: base_height,
        events: app.renderer.events
    });

    viewport
        .drag()
        .pinch()
        .wheel()
        .decelerate()
        .clampZoom({minWidth: 500, maxWidth: base_width * 1.5});

    return viewport;
}