export const BUILDING_MAP={
    1001: 'Town_Hall1', 1002: 'Town_Hall2', 1003: 'Town_Hall3', 1004: 'Town_Hall4',
    2001: 'Cannon1', 2002: 'Cannon2', 2003: 'Cannon3', 2004: 'Cannon4',
    2101: 'Archer_Tower1', 2102: 'Archer_Tower2', 2103: 'Archer_Tower3', 2104: 'Archer_Tower4',
    2201: 'Mortar1', 2202: 'Mortar2', 2203: 'Mortar3', 2204: 'Mortar4',
    3001: 'Elixir_Collector1', 3002: 'Elixir_Collector2', 3003: 'Elixir_Collector3', 3004: 'Elixir_Collector4',
    3101: 'Pancake_Machine1', 3102: 'Pancake_Machine2', 3103: 'Pancake_Machine3', 3104: 'Pancake_Machine4',
    4001: 'Elixir_Storage1', 4002: 'Elixir_Storage2', 4003: 'Elixir_Storage3', 4004: 'Elixir_Storage4',
    4101: 'Pancake_Stack1', 4102: 'Pancake_Stack2', 4103: 'Pancake_Stack3', 4104: 'Pancake_Stack4',
    5001: 'Laboratory1', 5002: 'Laboratory2', 5003: 'Laboratory3', 5004: 'Laboratory4',
    6001: 'Army_Camp1', 6002: 'Army_Camp2', 6003: 'Army_Camp3', 6004: 'Army_Camp4'
}

export const TROOP_DEFS = {
    'Barbarian': {
        emoji: '🪓',
        space: 1,
        color: 0xf1c40f,
        levels: {
            1: { hp: 45, dps: 8, range: 1, speed: 16, upgradeCost: 0, labReq: 1, upgradeTime: 0 },
            2: { hp: 54, dps: 11, range: 1, speed: 16, upgradeCost: 50000, labReq: 1, upgradeTime: 60 },
            3: { hp: 65, dps: 14, range: 1, speed: 16, upgradeCost: 150000, labReq: 2, upgradeTime: 300 },
            4: { hp: 78, dps: 18, range: 1, speed: 16, upgradeCost: 500000, labReq: 3, upgradeTime: 900 }
        }
    },
    'Archer': {
        emoji: '🏹',
        space: 1,
        color: 0xe91e63,
        levels: {
            1: { hp: 30, dps: 7, range: 3, speed: 24, upgradeCost: 0, labReq: 1, upgradeTime: 0 },
            2: { hp: 36, dps: 9, range: 3, speed: 24, upgradeCost: 50000, labReq: 1, upgradeTime: 60 },
            3: { hp: 44, dps: 12, range: 3, speed: 24, upgradeCost: 250000, labReq: 2, upgradeTime: 600 },
            4: { hp: 54, dps: 16, range: 3, speed: 24, upgradeCost: 750000, labReq: 3, upgradeTime: 1200 }
        }
    },
    'Giant': {
        emoji: '👊',
        space: 5,
        color: 0xd35400,
        levels: {
            1: { hp: 300, dps: 11, range: 1, speed: 12, upgradeCost: 0, labReq: 1, upgradeTime: 0 },
            2: { hp: 360, dps: 14, range: 1, speed: 12, upgradeCost: 100000, labReq: 1, upgradeTime: 300 },
            3: { hp: 430, dps: 19, range: 1, speed: 12, upgradeCost: 250000, labReq: 2, upgradeTime: 900 },
            4: { hp: 520, dps: 24, range: 1, speed: 12, upgradeCost: 750000, labReq: 3, upgradeTime: 1800 }
        }
    },
    'Goblin': {
        emoji: '💰',
        space: 1,
        color: 0x2ecc71,
        levels: {
            1: { hp: 25, dps: 11, range: 1, speed: 32, upgradeCost: 0, labReq: 1, upgradeTime: 0 },
            2: { hp: 30, dps: 14, range: 1, speed: 32, upgradeCost: 50000, labReq: 1, upgradeTime: 60 },
            3: { hp: 36, dps: 19, range: 1, speed: 32, upgradeCost: 250000, labReq: 2, upgradeTime: 600 },
            4: { hp: 43, dps: 24, range: 1, speed: 32, upgradeCost: 750000, labReq: 3, upgradeTime: 1200 }
        }
    },
    'P.E.K.K.A': {
        emoji: '🤖',
        space: 25,
        color: 0x9b59b6,
        levels: {
            1: { hp: 2800, dps: 240, range: 1, speed: 16, upgradeCost: 0, labReq: 1, upgradeTime: 0 },
            2: { hp: 3100, dps: 270, range: 1, speed: 16, upgradeCost: 3000000, labReq: 2, upgradeTime: 600 },
            3: { hp: 3500, dps: 310, range: 1, speed: 16, upgradeCost: 6000000, labReq: 2, upgradeTime: 1800 },
            4: { hp: 4000, dps: 360, range: 1, speed: 16, upgradeCost: 8000000, labReq: 2, upgradeTime: 3600 }
        }
    }
};

export const gameConfig={
    TILE_SIZE: 32,
    GRID_WIDTH: 40,
    GRID_HEIGHT: 40
};