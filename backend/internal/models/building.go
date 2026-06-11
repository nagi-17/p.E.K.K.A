package models

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/nagi-17/p.E.K.K.A/internal/database"
)

func GetOwnedBuildingData(ctx context.Context, playerID uuid.UUID) ([]OwnedBuildingWithData, error) {
	query := `
	SELECT 
		ob.id, ob.player_id, ob.building_data_id, ob.pos_x, ob.pos_y, 
		ob.upgrade_complete_at, ob.last_collected_at,
		bd.building_type, bd.building_level, bd.health, bd.width, bd.height, 
		bd.build_time, bd.upgrade_cost_elixir, bd.upgrade_cost_pancakes, 
		bd.upgrade_time, bd.max_quantity_available, bd.skill_on_upgrade
	FROM owned_buildings ob
	INNER JOIN building_data bd ON ob.building_data_id = bd.id
	WHERE ob.player_id = $1;
	`
	rows, err := database.DB.Query(ctx, query, playerID)
	if err != nil {
		return nil, fmt.Errorf("Error in fetching owned buildings data: %w", err)
	}
	defer rows.Close()

	var buildings []OwnedBuildingWithData
	for rows.Next() {
		var x OwnedBuildingWithData
		err := rows.Scan(
			&x.OwnedBuildingData.ID, &x.OwnedBuildingData.PlayerID, &x.OwnedBuildingData.BuildingDataID,
			&x.OwnedBuildingData.PosX, &x.OwnedBuildingData.PosY, &x.OwnedBuildingData.UpgradeCompleteAt,
			&x.OwnedBuildingData.LastCollectedAt, &x.BuildingType, &x.BuildingLevel, &x.Health, &x.Width,
			&x.Height, &x.BuildTime, &x.UpgradeCostElixir, &x.UpgradeCostPancakes, &x.UpgradeTime,
			&x.MaxQuantityAvailable, &x.SkillOnUpgrade,
		)
		if err != nil {
			return nil, fmt.Errorf("Error in scanning owned building data/rows: %w", err)
		}
		buildings = append(buildings, x)
	}
	return buildings, nil
}
