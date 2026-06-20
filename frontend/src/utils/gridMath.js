import { gameConfig } from '../game/gameConfig';

export function pointerToGrid(container, globalPoint) {
    const local=container.toLocal(globalPoint);
    return {
        gridX: Math.floor(local.x / gameConfig.TILE_SIZE),
        gridY: Math.floor(local.y / gameConfig.TILE_SIZE)
    };
}