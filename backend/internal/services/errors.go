package services

import "errors"

var (
	ErrInvalidBuildingType  = errors.New("Invalid building type")
	ErrTownHallUnderLevel   = errors.New("Town Hall is under levelled")
	ErrNotEnoughPancakes    = errors.New("Not enough pancakes")
	ErrNotEnoughElixir      = errors.New("Not enough elixir")
	ErrBuildingLimitReached = errors.New("All possible buildings of this type have already been placed")
	ErrCellOccupied         = errors.New("Cell is occupied")
)
