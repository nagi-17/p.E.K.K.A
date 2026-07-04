package services

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/nagi-17/p.E.K.K.A/internal/database"
	"github.com/nagi-17/p.E.K.K.A/internal/models"
)

const trophyWinGain = 30
const trophyLossCost = 20
const maxLootPercent = 0.20

type troopUnit struct {
	health         int
	dps            int
	posX           float64
	posY           float64
	troopRange     float64
	speed          float64
	attackCooldown float64
}

type buildingUnit struct {
	id             int
	buildingType   string
	health         int
	maxHealth      int
	pixelX         float64
	pixelY         float64
	width          int
	height         int
	isDefense      bool
	rangePx        float64
	dps            int
	damagePerShot  int
	attackCooldown float64
}

type DeploymentEvent struct {
	TroopType string  `json:"troop_type"`
	Quantity  int     `json:"quantity"`
	DropX     float64 `json:"drop_x"`
	DropY     float64 `json:"drop_y"`
	Tick      int     `json:"tick"`
}

func SimulateBattle(events []DeploymentEvent, defenses []models.DefenseSnapshot, army []models.AttackingTroop) float32 {
	if len(defenses) == 0 {
		return 100
	}

	var units []troopUnit

	buildings := make([]buildingUnit, len(defenses))
	totalInitialHP := 0
	for i := 0; i < len(defenses); i++ {
		w := defenses[i].Width
		if w <= 0 {
			w = 3
		}
		h := defenses[i].Height
		if h <= 0 {
			h = 3
		}

		isDefense := defenses[i].BuildingType == "Cannon" || defenses[i].BuildingType == "Archer Tower" || defenses[i].BuildingType == "Mortar"
		
		var dps int
		var rangePx float64
		if defenses[i].BuildingType == "Cannon" {
			dps = 9
			rangePx = 9 * 32.0
		} else if defenses[i].BuildingType == "Archer Tower" {
			dps = 11
			rangePx = 10 * 32.0
		} else if defenses[i].BuildingType == "Mortar" {
			dps = 4
			rangePx = 11 * 32.0
		} else {
			dps = 0
			rangePx = 0
		}

		buildings[i] = buildingUnit{
			id:             i,
			buildingType:   defenses[i].BuildingType,
			health:         defenses[i].Health,
			maxHealth:      defenses[i].Health,
			pixelX:         float64(defenses[i].PosX)*32.0 + (float64(w)*32.0)/2.0,
			pixelY:         float64(defenses[i].PosY)*32.0 + (float64(h)*32.0)/2.0,
			width:          w,
			height:         h,
			isDefense:      isDefense,
			rangePx:        rangePx,
			dps:            dps,
			damagePerShot:  defenses[i].DamagePerShot,
			attackCooldown: 0,
		}
		totalInitialHP += defenses[i].Health
	}

	if totalInitialHP == 0 {
		return 0
	}

	const maxSimulationTicks = 10800
	destroyedHP := 0

	for tick := 0; tick < maxSimulationTicks; tick++ {

		for _, event := range events {
			if event.Tick == tick {
				var baseSpeed float64 = 16.0
				var rangePx float64 = 1.0 * 32.0
				var health int = 100
				var dps int = 20

				var found = false
				for _, troop := range army {
					if troop.TroopType == event.TroopType {
						health = troop.Health
						dps = troop.DPS
						found = true
						break
					}
				}

				if !found {
					switch event.TroopType {
					case "Barbarian":
						health = 45
						dps = 8
					case "Archer":
						health = 30
						dps = 7
					case "Giant":
						health = 300
						dps = 11
					case "Goblin":
						health = 25
						dps = 11
					case "P.E.K.K.A":
						health = 2800
						dps = 240
					}
				}

				switch event.TroopType {
				case "Archer":
					baseSpeed = 24.0
					rangePx = 3.0 * 32.0
				case "Giant":
					baseSpeed = 12.0
					rangePx = 1.0 * 32.0
				case "Goblin":
					baseSpeed = 32.0
					rangePx = 1.0 * 32.0
				case "Barbarian":
					baseSpeed = 16.0
					rangePx = 1.0 * 32.0
				case "P.E.K.K.A":
					baseSpeed = 16.0
					rangePx = 1.0 * 32.0
				}
				speed := baseSpeed / 8.0

				for j := 0; j < event.Quantity; j++ {
					spreadX := (float64(j%5) - 2.0) * 4.0
					spreadY := (float64(j/5) - 2.0) * 4.0

					units = append(units, troopUnit{
						health:         health,
						dps:            dps,
						posX:           event.DropX*32.0 + 16.0 + spreadX,
						posY:           event.DropY*32.0 + 16.0 + spreadY,
						troopRange:     rangePx,
						speed:          speed,
						attackCooldown: 0,
					})
				}
			}
		}

		allBuildingsDestroyed := true
		for i := 0; i < len(buildings); i++ {
			if buildings[i].health > 0 {
				allBuildingsDestroyed = false
				break
			}
		}

		allTroopsDead := true
		for i := 0; i < len(units); i++ {
			if units[i].health > 0 {
				allTroopsDead = false
				break
			}
		}

		hasMoreDeployments := false
		for _, e := range events {
			if e.Tick > tick {
				hasMoreDeployments = true
				break
			}
		}

		if allBuildingsDestroyed || (allTroopsDead && !hasMoreDeployments) {
			break
		}

		for i := 0; i < len(units); i++ {
			if units[i].health <= 0 {
				continue
			}

			var target *buildingUnit = nil
			minDist := math.MaxFloat64

			for j := 0; j < len(buildings); j++ {
				if buildings[j].health <= 0 {
					continue
				}

				dx := buildings[j].pixelX - units[i].posX
				dy := buildings[j].pixelY - units[i].posY
				dist := math.Sqrt(dx*dx + dy*dy)

				if dist < minDist {
					minDist = dist
					target = &buildings[j]
				}
			}

			if target == nil {
				continue
			}

			if minDist > units[i].troopRange {
				dx := target.pixelX - units[i].posX
				dy := target.pixelY - units[i].posY
				if minDist > 0 {
					units[i].posX += (dx / minDist) * units[i].speed
					units[i].posY += (dy / minDist) * units[i].speed
				}
			} else {
				units[i].attackCooldown -= 1.0
				if units[i].attackCooldown <= 0 {
					damage := int(math.Ceil(float64(units[i].dps) * 1.2))
					before := target.health
					target.health -= damage
					if target.health < 0 {
						target.health = 0
					}
					destroyedHP += before - target.health
					units[i].attackCooldown = 60.0
				}
			}
		}

		for i := 0; i < len(buildings); i++ {
			if buildings[i].health <= 0 || !buildings[i].isDefense {
				continue
			}

			var target *troopUnit = nil
			minDist := math.MaxFloat64

			for j := 0; j < len(units); j++ {
				if units[j].health <= 0 {
					continue
				}

				dx := units[j].posX - buildings[i].pixelX
				dy := units[j].posY - buildings[i].pixelY
				dist := math.Sqrt(dx*dx + dy*dy)

				if dist < minDist && dist <= buildings[i].rangePx {
					minDist = dist
					target = &units[j]
				}
			}

			if target != nil {
				buildings[i].attackCooldown -= 1.0
				if buildings[i].attackCooldown <= 0 {
					dmg := int(math.Ceil(float64(buildings[i].dps) * 1.2))
					target.health -= dmg
					buildings[i].attackCooldown = 60.0
				}
			}
		}
	}

	damagePercent := (float32(destroyedHP) / float32(totalInitialHP)) * 100.0
	if damagePercent > 100.0 {
		damagePercent = 100.0
	}
	return damagePercent
}

func CalculateLoot(damagePercent float32, defenderElixir int, defenderPancakes int) (int, int) {
	lootFraction := (damagePercent / 100) * maxLootPercent

	elixirLooted := int(float32(defenderElixir) * lootFraction)
	pancakesLooted := int(float32(defenderPancakes) * lootFraction)

	return elixirLooted, pancakesLooted
}

func Attack(ctx context.Context, attackerID string, defenderID string, events []DeploymentEvent) (*models.BattleLog, error) {
	attackerUUID, err := uuid.Parse(attackerID)
	if err != nil {
		return nil, fmt.Errorf("Error in parsing attacker id: %w", err)
	}

	defenderUUID, err := uuid.Parse(defenderID)
	if err != nil {
		return nil, fmt.Errorf("Error in parsing defender id: %w", err)
	}

	if attackerUUID == defenderUUID {
		return nil, fmt.Errorf("You cannot attack yourself")
	}

	defenderStats, err := models.GetPlayerInfoByID(ctx, defenderUUID)
	if err != nil {
		return nil, fmt.Errorf("Failed to load defender stats: %w", err)
	}

	if defenderStats.Shield_End_Time != nil && defenderStats.Shield_End_Time.After(time.Now()) {
		return nil, fmt.Errorf("Defender is currently under shield")
	}

	attackerStats, err := models.GetPlayerInfoByID(ctx, attackerUUID)
	if err != nil {
		return nil, fmt.Errorf("Failed to load attacker stats: %w", err)
	}

	army, err := models.GetAttackingArmy(ctx, attackerUUID)
	if err != nil {
		return nil, err
	}

	defenses, err := models.GetDefenseSnapshot(ctx, defenderUUID)
	if err != nil {
		return nil, fmt.Errorf("Failed to load defender's village: %w", err)
	}

	damagePercent := SimulateBattle(events, defenses, army)
	elixirLooted, pancakesLooted := CalculateLoot(damagePercent, defenderStats.Elixir, defenderStats.Pancakes)

	attackerWon := damagePercent >= 50
	trophyChangeAttacker := -trophyLossCost
	trophyChangeDefender := trophyWinGain
	if attackerWon {
		trophyChangeAttacker = trophyWinGain
		trophyChangeDefender = -trophyLossCost
	}

	newAttackerTrophies := attackerStats.Trophies + trophyChangeAttacker
	if newAttackerTrophies < 0 {
		newAttackerTrophies = 0
	}
	newDefenderTrophies := defenderStats.Trophies + trophyChangeDefender
	if newDefenderTrophies < 0 {
		newDefenderTrophies = 0
	}

	tx, err := database.DB.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("Failed to begin database transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	newAttackerElixir := attackerStats.Elixir + elixirLooted
	newAttackerPancakes := attackerStats.Pancakes + pancakesLooted
	err = models.UpdateResources(ctx, tx, attackerUUID, newAttackerPancakes, newAttackerElixir)
	if err != nil {
		return nil, err
	}

	newDefenderElixir := defenderStats.Elixir - elixirLooted
	newDefenderPancakes := defenderStats.Pancakes - pancakesLooted
	err = models.UpdateResources(ctx, tx, defenderUUID, newDefenderPancakes, newDefenderElixir)
	if err != nil {
		return nil, err
	}

	err = models.UpdateTrophies(ctx, tx, newAttackerTrophies, newDefenderTrophies, attackerUUID, defenderUUID)
	if err != nil {
		return nil, err
	}

	err = models.UpdateArmyAfterBattle(ctx, tx, attackerUUID)
	if err != nil {
		return nil, err
	}

	shieldEnd := time.Now().Add(8 * time.Hour)
	err = models.UpdateShieldTime(ctx, tx, shieldEnd, defenderUUID)
	if err != nil {
		return nil, err
	}

	err = models.CreateBattleLog(ctx, tx, attackerUUID, defenderUUID, elixirLooted, pancakesLooted, damagePercent)
	if err != nil {
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("Failed to commit to database: %w", err)
	}

	result := &models.BattleLog{
		AttackerID:     attackerUUID,
		DefenderID:     defenderUUID,
		ElixirLooted:   elixirLooted,
		PancakesLooted: pancakesLooted,
		DamagePercent:  damagePercent,
		TimeOfBattle:   time.Now(),
	}

	return result, nil
}
