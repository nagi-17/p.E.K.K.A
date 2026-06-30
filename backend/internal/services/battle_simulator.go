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
const maxBattleTicks = 500

type troopUnit struct {
	health     int
	dps        int
	posX       float64
	posY       float64
	troopRange float64
	speed      float64
}

func distance(x1, y1, x2, y2 float64) float64 {
	return math.Sqrt((x1-x2)*(x1-x2) + (y1-y2)*(y1-y2))
}
func SimulateBattle(army []models.AttackingTroop, defenses []models.DefenseSnapshot, dropX, dropY float64) float32 {
	if len(defenses) == 0 {
		return 100
	}

	var units []troopUnit
	for i := 0; i < len(army); i++ {
		for j := 0; j < army[i].Quantity; j++ {
			var speed float64 = 2.0
			switch army[i].TroopType {
			case "Barbarian":
				speed = 2.0
			case "Archer":
				speed = 3.0
			case "Giant":
				speed = 1.5
			case "Goblin":
				speed = 4.0
			case "P.E.K.K.A":
				speed = 2.0
			}
			units = append(units, troopUnit{
				health:     army[i].Health,
				dps:        army[i].DPS,
				posX:       dropX,
				posY:       dropY,
				troopRange: float64(army[i].TroopRange),
				speed:      speed,
			})
		}
	}
	if len(units) == 0 {
		return 0
	}

	defenseHP := make([]int, len(defenses))
	totalDefenseHP := 0
	for i := 0; i < len(defenses); i++ {
		defenseHP[i] = defenses[i].Health
		totalDefenseHP += defenseHP[i]
	}
	destroyedHP := 0

	for tick := 0; tick < maxBattleTicks; tick++ {
		if allDefensesDestroyed(defenseHP) || allUnitsDead(units) {
			break
		}

		for i := 0; i < len(units); i++ {
			if units[i].health <= 0 {
				continue
			}
			target := nearestLivingDefense(units[i].posX, units[i].posY, defenses, defenseHP)
			if target == -1 {
				break
			}

			dist := distance(units[i].posX, units[i].posY, float64(defenses[target].PosX), float64(defenses[target].PosY))

			if dist > units[i].troopRange {
				dx := float64(defenses[target].PosX) - units[i].posX
				dy := float64(defenses[target].PosY) - units[i].posY

				if dist > 0 {
					temp_x := dx / dist
					temp_y := dy / dist
					units[i].posX += temp_x * units[i].speed
					units[i].posY += temp_y * units[i].speed
				}
			} else {
				before := defenseHP[target]
				defenseHP[target] -= units[i].dps
				if defenseHP[target] < 0 {
					defenseHP[target] = 0
				}
				destroyedHP += before - defenseHP[target]
			}
		}

		for i := 0; i < len(defenses); i++ {
			if defenseHP[i] <= 0 {
				continue
			}

			target := nearestLivingTroop(float64(defenses[i].PosX), float64(defenses[i].PosY), units)
			if target == -1 {
				continue
			}

			dist := distance(float64(defenses[i].PosX), float64(defenses[i].PosY), units[target].posX, units[target].posY)

			if dist <= float64(defenses[i].BuildingRange) {
				units[target].health -= defenses[i].DamagePerShot
			}
		}
	}

	if totalDefenseHP == 0 {
		return 0
	}

	damagePercent := (float32(destroyedHP) / float32(totalDefenseHP)) * 100
	if damagePercent > 100 {
		damagePercent = 100
	}

	return damagePercent
}

func allDefensesDestroyed(defenseHP []int) bool {
	for i := 0; i < len(defenseHP); i++ {
		if defenseHP[i] > 0 {
			return false
		}
	}
	return true
}

func allUnitsDead(units []troopUnit) bool {
	for i := 0; i < len(units); i++ {
		if units[i].health > 0 {
			return false
		}
	}
	return true
}

func nearestLivingDefense(tX, tY float64, defenses []models.DefenseSnapshot, defenseHP []int) int {
	best := -1
	minDist := math.MaxFloat64

	for i := 0; i < len(defenses); i++ {
		if defenseHP[i] > 0 {
			dist := distance(tX, tY, float64(defenses[i].PosX), float64(defenses[i].PosY))
			if dist < minDist {
				minDist = dist
				best = i
			}
		}
	}
	return best
}

func nearestLivingTroop(dX, dY float64, units []troopUnit) int {
	best := -1
	minDist := math.MaxFloat64

	for i := 0; i < len(units); i++ {
		if units[i].health > 0 {
			dist := distance(dX, dY, units[i].posX, units[i].posY)
			if dist < minDist {
				minDist = dist
				best = i
			}
		}
	}
	return best
}

func CalculateLoot(damagePercent float32, defenderElixir int, defenderPancakes int) (int, int) {
	lootFraction := (damagePercent / 100) * maxLootPercent

	elixirLooted := int(float32(defenderElixir) * lootFraction)
	pancakesLooted := int(float32(defenderPancakes) * lootFraction)

	return elixirLooted, pancakesLooted
}

func Attack(ctx context.Context, attackerID string, defenderID string, dropX float64, dropY float64) (*models.BattleLog, error) {
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

	damagePercent := SimulateBattle(army, defenses, dropX, dropY)
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
