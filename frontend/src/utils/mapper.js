export function mapBuildingData(go_data) {
    //console.log("RAW PLAYER DATA FROM GO:", go_data);
    if (!go_data) return null;
    
    return {
        id: go_data.ID || go_data.id,
        playerId: go_data.PlayerID || go_data.player_id,
        buildingDataId: go_data.BuildingDataID || go_data.building_data_id,

        posX: go_data.PosX ?? go_data.pos_x ?? 0,
        posY: go_data.PosY ?? go_data.pos_y ?? 0,

        upgradeCompleteAt: go_data.UpgradeCompleteAt || go_data.upgrade_complete_at || null,
        lastCollectedAt: go_data.LastCollectedAt || go_data.last_collected_at || null,

        type: go_data.BuildingType || go_data.building_type || 'Unknown',
        level: go_data.BuildingLevel || go_data.building_level || 1,
        health: go_data.Health || go_data.health || 0,
        width: go_data.Width || go_data.width || 3,
        height: go_data.Height || go_data.height || 3,

        buildTime: go_data.BuildTime || go_data.build_time || 0,
        upgradeTime: go_data.UpgradeTime || go_data.upgrade_time || 0,
        costElixir: go_data.UpgradeCostElixir || go_data.upgrade_cost_elixir || 0,
        costPancakes: go_data.UpgradeCostPancakes || go_data.upgrade_cost_pancakes || 0,
        skillGain: go_data.SkillOnUpgrade || go_data.skill_on_upgrade || 0
    };
}

export function mapPlayerStats(go_data) {
    if (!go_data) return null;
    
    return {
        id: go_data.Player_ID || go_data.player_id,
        username: go_data.Username || go_data.username || localStorage.getItem('player_username') || 'Villager',
        trophies: go_data.Trophies || go_data.trophies || 0,
        skillPoints: go_data.Skill_points || go_data.skill_points || 0,
        elixir: go_data.Elixir || go_data.elixir || 0,
        pancakes: go_data.Pancakes || go_data.pancakes || 0,
        shieldEndTime: go_data.Shield_End_Time || go_data.shield_end_time || null,
        maxElixir: go_data.max_elixir || go_data.MaxElixir || 1500,
        maxPancakes: go_data.max_pancakes || go_data.MaxPancakes || 1500
    };
}