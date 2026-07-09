package services

import "github.com/nagi-17/p.E.K.K.A/internal/utils"

var (
	ErrInvalidBuildingType  = utils.NewUserError("Invalid building type")
	ErrTownHallUnderLevel   = utils.NewUserError("Town Hall is under levelled")
	ErrNotEnoughPancakes    = utils.NewUserError("Not enough pancakes")
	ErrNotEnoughElixir      = utils.NewUserError("Not enough elixir")
	ErrBuildingLimitReached = utils.NewUserError("All possible buildings of this type have already been placed")
	ErrCellOccupied         = utils.NewUserError("Cell is occupied")
)
