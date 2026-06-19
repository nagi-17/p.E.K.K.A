import { Assets } from 'pixi.js';
import { BUILDING_MAP } from '../interactions/gameConfig';

export class AssetLoader {
    static async loadAssets() {
        const aliasesToLoad=[];
        const uniqueAliases=[...new Set(Object.values(BUILDING_MAP))];

        uniqueAliases.forEach(alias=>{
            Assets.add({alias: alias, src: `${alias}.png`});
            aliasesToLoad.push(alias);
        });
        await Assets.load(aliasesToLoad); 
    }

    static getTexture(name) {
        return Assets.get(name);
    }
}