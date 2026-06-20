import * as PIXI from 'pixi.js';
import { GridOverlay } from '../entities/GridOverlay';
import { BuildingSprite } from '../entities/BuildingSprite';
import { TroopSprite } from '../entities/TroopSprite';
import { gameConfig, TROOP_DEFS } from '../gameConfig';

export class BattleScene {
    constructor(onSimulationEnd, onDamageUpdate) {
        this.onSimulationEnd=onSimulationEnd;
        this.onDamageUpdate=onDamageUpdate;

        this.container=new PIXI.Container();
        this.container.sortableChildren=true;

        this.grid=new GridOverlay();
        this.grid.container.zIndex=-100;
        this.container.addChild(this.grid.container);

        this.buildings=new Map();
        this.buildingStats=new Map();

        this.troops=[];
        this.projectiles=[];
        
        this.isBattleActive=false;
        this.isBattleEnded=false;
        this.deployableTroopType=null;
        this.deployedCountByType={};
        this.onTroopDeployedCallback=null;
        this.deployRemaining=null;

        this.totalInitialHP=0;
        this.currentDestroyedHP=0;

        this.setupDeploymentInteraction();

        PIXI.Ticker.shared.add(this.update, this);
    }

    setupDeploymentInteraction() {
        this.grid.container.eventMode='static';
        this.grid.container.cursor='crosshair';
        this.grid.container.on('pointerdown', (event) => {
            if (!this.deployableTroopType || this.isBattleEnded) return;

            const localPos=event.getLocalPosition(this.container);
            const gridX=Math.floor(localPos.x / gameConfig.TILE_SIZE);
            const gridY=Math.floor(localPos.y / gameConfig.TILE_SIZE);

            if (gridX >= 0 && gridX < gameConfig.GRID_WIDTH && gridY >= 0 && gridY < gameConfig.GRID_HEIGHT) {
                if (this.isCellFree(gridX, gridY)) {
                    this.deployTroop(this.deployableTroopType, localPos.x, localPos.y, gridX, gridY);
                }
            }
        });
    }

    isCellFree(gridX, gridY) {
        for (const [id, b] of this.buildingStats) {
            if (gridX >= b.x && gridX < b.x + b.width && gridY >= b.y && gridY < b.y + b.height) {
                return false;
            }
        }
        return true;
    }

    setDeployableTroop(troopType, onTroopDeployed) {
        this.deployableTroopType=troopType;
        this.onTroopDeployedCallback=onTroopDeployed;
    }

    deployTroop(type, pixelX, pixelY, gridX, gridY) {
        const def=TROOP_DEFS[type];
        if (!def) return;

        const level=1;
        const stats=def.levels[level];
        const troop=new TroopSprite(type, level, pixelX, pixelY, stats.hp);

        troop.currentHP=stats.hp;
        troop.maxHP=stats.hp;
        troop.dps=stats.dps;
        troop.range=def.space > 10 ? 1 : (type === 'Archer' ? 3 * gameConfig.TILE_SIZE : 1 * gameConfig.TILE_SIZE);
        troop.speed=(stats.speed || 16) / 8;
        troop.attackCooldown=0;

        this.troops.push(troop);
        this.container.addChild(troop.container);

        if (this.onTroopDeployedCallback) {
            this.onTroopDeployedCallback(type, gridX, gridY);
        }
    }

    loadDefenderVillage(buildingData) {
        this.buildings.forEach(b => b.destroy());
        this.buildings.clear();
        this.buildingStats.clear();
        this.totalInitialHP=0;

        if (buildingData && buildingData.length > 0) {
            buildingData.forEach(data => {
                
                const extendedData={ 
                    ...data, 
                    lastCollectedAt: new Date().toISOString(), 
                    upgradeCompleteAt: null 
                };
                const building=new BuildingSprite(extendedData);
                building.container.eventMode='none';

                this.buildings.set(data.id, building);
                this.container.addChild(building.container);

                const maxHP=data.health || 500;
                this.totalInitialHP += maxHP;

                const isDefense=['Cannon', 'Archer Tower', 'Mortar'].includes(data.type);

                this.buildingStats.set(data.id, {
                    id: data.id,
                    type: data.type,
                    currentHP: maxHP,
                    maxHP: maxHP,
                    x: data.posX,
                    y: data.posY,
                    width: data.width,
                    height: data.height,
                    pixelX: data.posX * gameConfig.TILE_SIZE + (data.width * gameConfig.TILE_SIZE) / 2,
                    pixelY: data.posY * gameConfig.TILE_SIZE + (data.height * gameConfig.TILE_SIZE) / 2,
                    isDefense: isDefense,
                    attackCooldown: 0,
                    range: (data.type === 'Cannon' ? 9 : data.type === 'Archer Tower' ? 10 : 11) * gameConfig.TILE_SIZE,
                    dps: data.type === 'Cannon' ? 9 : data.type === 'Archer Tower' ? 11 : 4
                });

                building.healthBar=new PIXI.Graphics();
                building.healthBar.position.set(0, -6);
                building.container.addChild(building.healthBar);
            });
        }
    }

    startBattle() {
        if (this.troops.length === 0) {
            alert('Deploy at least one troop to attack!');
            return;
        }
        this.isBattleActive=true;
    }

    update(ticker) {
        if (!this.isBattleActive) return;

        const deltaTime=ticker.deltaTime;

        this.updateProjectiles(deltaTime);
        this.updateTroops(deltaTime);

        this.updateDefenses(deltaTime);
        this.checkBattleEnd();
    }

    updateProjectiles(deltaTime) {
        for (let i=this.projectiles.length - 1; i >= 0; i--) {
            const p=this.projectiles[i];
            const dx=p.target.x - p.sprite.x;
            const dy=p.target.y - p.sprite.y;
            const dist=Math.sqrt(dx * dx + dy * dy);

            if (dist < 10) {
                p.target.setHealth(p.target.currentHP - p.damage);
                this.container.removeChild(p.sprite);
                p.sprite.destroy();
                this.projectiles.splice(i, 1);
            } else {
                const speed=10 * deltaTime;
                p.sprite.x += (dx / dist) * speed;
                p.sprite.y += (dy / dist) * speed;
            }
        }
    }

    updateTroops(deltaTime) {
        this.troops.forEach((troop) => {
            if (troop.currentHP <= 0) return;

            let targetBuilding=null;
            let minDist=Infinity;

            for (const [id, stats] of this.buildingStats) {
                if (stats.currentHP <= 0) continue;

                const dx=stats.pixelX - troop.x;
                const dy=stats.pixelY - troop.y;
                const dist=Math.sqrt(dx * dx + dy * dy);

                if (dist < minDist) {
                    minDist=dist;
                    targetBuilding=stats;
                }
            }

            if (!targetBuilding) return;

            if (minDist > troop.range) {
                const dx=targetBuilding.pixelX - troop.x;
                const dy=targetBuilding.pixelY - troop.y;
                const moveX=(dx / minDist) * troop.speed * deltaTime;
                const moveY=(dy / minDist) * troop.speed * deltaTime;
                troop.updatePosition(troop.x + moveX, troop.y + moveY);
            } else {
                troop.attackCooldown -= deltaTime;
                if (troop.attackCooldown <= 0) {
                    const damage=Math.ceil(troop.dps * 1.5);
                    targetBuilding.currentHP=Math.max(0, targetBuilding.currentHP - damage);
                    troop.attackCooldown=45;

                    this.showHitFlash(targetBuilding.pixelX, targetBuilding.pixelY);

                    const bSprite=this.buildings.get(targetBuilding.id);
                    if (bSprite) {
                        this.drawBuildingHealthBar(bSprite, targetBuilding.currentHP, targetBuilding.maxHP);
                    }

                    if (targetBuilding.currentHP <= 0) {
                        this.destroyBuilding(targetBuilding.id);
                    }
                }
            }
        });
    }

    updateDefenses(deltaTime) {
        for (const [id, def] of this.buildingStats) {
            if (def.currentHP <= 0 || !def.isDefense) continue;

            let targetTroop=null;
            let minDist=Infinity;

            this.troops.forEach((troop) => {
                if (troop.currentHP <= 0) return;

                const dx=troop.x - def.pixelX;
                const dy=troop.y - def.pixelY;
                const dist=Math.sqrt(dx * dx + dy * dy);

                if (dist < minDist && dist <= def.range) {
                    minDist=dist;
                    targetTroop=troop;
                }
            });

            if (targetTroop) {
                def.attackCooldown -= deltaTime;
                if (def.attackCooldown <= 0) {
                    const dmg=Math.ceil(def.dps * 1.8);
                    this.spawnProjectile(def.pixelX, def.pixelY, targetTroop, dmg, def.type);
                    def.attackCooldown=50;
                }
            }
        }
    }

    spawnProjectile(startX, startY, targetTroop, damage, defenseType) {
        const pGraphics=new PIXI.Graphics();
        
        if (defenseType === 'Cannon') {
            pGraphics.circle(0, 0, 4).fill(0x34495e).stroke({ width: 1, color: 0x000000 });
        } else if (defenseType === 'Mortar') {
            pGraphics.circle(0, 0, 6).fill(0xd35400).stroke({ width: 1, color: 0x000000 });
        } else {
            pGraphics.rect(-6, -1, 12, 2).fill(0xe91e63);
        }

        pGraphics.x=startX;
        pGraphics.y=startY;
        pGraphics.zIndex=1000;
        this.container.addChild(pGraphics);

        this.projectiles.push({
            sprite: pGraphics,
            target: targetTroop,
            damage: damage
        });
    }

    drawBuildingHealthBar(bSprite, current, max) {
        bSprite.healthBar.clear();
        if (current <= 0 || current === max) return;

        const w=bSprite.pixelWidth;
        const h=4;
        const percent=current / max;

        bSprite.healthBar.rect(0, 0, w, h).fill(0xc0392b);
        bSprite.healthBar.rect(0, 0, w * percent, h).fill(0x2ecc71);
    }

    showHitFlash(x, y) {
        const flash=new PIXI.Graphics();
        flash.circle(0, 0, 8).fill({ color: 0xffffff, alpha: 0.8 });
        flash.x=x;
        flash.y=y;
        this.container.addChild(flash);

        setTimeout(() => {
            if (flash && !flash.destroyed) {
                this.container.removeChild(flash);
                flash.destroy();
            }
        }, 80);
    }

    destroyBuilding(id) {
        const bSprite=this.buildings.get(id);
        if (bSprite) {
            bSprite.sprite.tint=0x333333;
            bSprite.healthBar.clear();

            const xCross=new PIXI.Graphics();
            xCross.lineStyle({ width: 4, color: 0xc0392b });
            xCross.moveTo(10, 10).lineTo(bSprite.pixelWidth - 10, bSprite.pixelHeight - 10);
            xCross.moveTo(bSprite.pixelWidth - 10, 10).lineTo(10, bSprite.pixelHeight - 10);
            bSprite.container.addChild(xCross);
        }

        let destroyedHP=0;
        for (const [bid, stats] of this.buildingStats) {
            destroyedHP += (stats.maxHP - stats.currentHP);
        }
        
        const pct=(destroyedHP / this.totalInitialHP) * 100;
        this.onDamageUpdate(Math.min(100, pct));
    }

    checkBattleEnd() {
        let buildingsLeft=false;
        for (const [id, stats] of this.buildingStats) {
            if (stats.currentHP > 0) {
                buildingsLeft=true;
                break;
            }
        }

        let troopsAlive=false;
        this.troops.forEach((troop) => {
            if (troop.currentHP > 0) {
                troopsAlive=true;
            }
        });

        let undeployedLeft=false;
        if (this.deployRemaining) {
            undeployedLeft=Object.values(this.deployRemaining).some(qty => qty > 0);
        }

        if (!buildingsLeft || (!troopsAlive && !undeployedLeft)) {
            this.endBattle();
        }
    }

    endBattle() {
        if (!this.isBattleActive) return;
        this.isBattleActive=false;
        this.isBattleEnded=true;
        
        let destroyedHP=0;
        for (const [bid, stats] of this.buildingStats) {
            destroyedHP += (stats.maxHP - stats.currentHP);
        }
        const finalDamagePercent=(destroyedHP / this.totalInitialHP) * 100;

        if (this.onSimulationEnd) {
            this.onSimulationEnd(finalDamagePercent);
        }
    }

    destroy() {
        PIXI.Ticker.shared.remove(this.update, this);

        this.projectiles.forEach((p) => {
            this.container.removeChild(p.sprite);
            p.sprite.destroy();
        });
        this.projectiles=[];

        this.troops.forEach((t) => t.destroy());
        this.troops=[];

        this.buildings.forEach((b) => b.destroy());
        this.buildings.clear();

        this.container.destroy({ children: true });
    }
}
