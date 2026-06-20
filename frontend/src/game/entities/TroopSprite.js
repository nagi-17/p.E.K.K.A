import * as PIXI from 'pixi.js';
import { TROOP_DEFS } from '../gameConfig';

export class TroopSprite {
    constructor(type, level, x, y, maxHealth) {
        this.type=type;
        this.level=level;
        this.maxHealth=maxHealth;
        this.health=maxHealth;

        this.container=new PIXI.Container();

        const def=TROOP_DEFS[type] || { emoji: '👤', color: 0x95a5a6 };

        this.base=new PIXI.Graphics();
        this.base.circle(0, 0, 14)
            .fill({ color: def.color, alpha: 0.85 })
            .stroke({ width: 2, color: 0x000000 });
        this.container.addChild(this.base);

        this.emojiText=new PIXI.Text({
            text: def.emoji,
            style: {
                fontSize: 16,
            }
        });
        this.emojiText.anchor.set(0.5);
        this.container.addChild(this.emojiText);

        this.healthBar=new PIXI.Graphics();
        this.healthBar.position.set(-15, -24);
        this.container.addChild(this.healthBar);
        this.updateHealthBar();

        this.updatePosition(x, y);
    }

    updatePosition(pixelX, pixelY) {
        this.x=pixelX;
        this.y=pixelY;
        this.container.x=pixelX;
        this.container.y=pixelY;
        this.container.zIndex=Math.floor(pixelY);
    }

    setHealth(currentHealth) {
        this.health=Math.max(0, currentHealth);
        this.updateHealthBar();
    }

    updateHealthBar() {
        this.healthBar.clear();
        if (this.health <= 0) return;

        const width=30;
        const height=4;
        const percent=this.health / this.maxHealth;

        this.healthBar.rect(0, 0, width, height)
            .fill(0xc0392b)
            .stroke({ width: 1, color: 0x000000 });

        this.healthBar.rect(0, 0, width * percent, height)
            .fill(0x2ecc71);
    }

    destroy() {
        this.container.destroy({ children: true });
    }
}
