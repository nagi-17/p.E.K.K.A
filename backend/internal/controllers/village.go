package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/nagi-17/p.E.K.K.A/internal/models"
)

type PlaceBuildingReq struct {
	Btype string `json:"building_type"`
	X     int    `json:"pos_x"`
	Y     int    `json:"pos_y"`
}

type MoveBuildingReq struct {
	BuildingID string `json:"building_id"`
	New_x      int    `json:"new_x"`
	New_y      int    `json:"new_y"`
}

type UpgradeBuildingReq struct {
	BuildingID string `json:"building_id"`
}

type Response struct {
	Message string `json:"message"`
}

func LoadVillage(w http.ResponseWriter, request *http.Request) {
	val := request.Context().Value("player_id")
	checkPlayerIDStr, ok := val.(string)
	if !ok {
		http.Error(w, "Player ID missing or is invalid in context", http.StatusInternalServerError)
		return
	}

	playerIDuuid, err := uuid.Parse(checkPlayerIDStr)
	if err != nil {
		http.Error(w, "Player ID conversion failed", http.StatusBadRequest)
		return
	}

	buildings, err := models.GetOwnedBuildingData(request.Context(), playerIDuuid)
	if err != nil {
		http.Error(w, "Couldn't load village(owned buildings)", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(buildings)
}

func PlaceBuilding(w http.ResponseWriter, request *http.Request) {
	val := request.Context().Value("player_id")
	checkPlayerIDStr, ok := val.(string)
	if !ok {
		http.Error(w, "Player ID missing or is invalid in context", http.StatusInternalServerError)
		return
	}
	playerIDuuid, err := uuid.Parse(checkPlayerIDStr)
	if err != nil {
		http.Error(w, "Player ID conversion failed", http.StatusBadRequest)
		return
	}
	var req PlaceBuildingReq
	decoder := json.NewDecoder(request.Body)
	err = decoder.Decode(&req)
	if err != nil {
		http.Error(w, "Bad request-wrong json payload", http.StatusBadRequest)
	}
	if req.X < 0 || req.Y < 0 {
		http.Error(w, "Invalid building coordinates", http.StatusBadRequest)
	}
	err = models.PlaceNewBuilding(request.Context(), playerIDuuid, req.Btype, req.X, req.Y)
	if err != nil {
		switch err.Error() {
		case "Invalid building type", "Town Hall is under levelled", "Not enough pancakes", "Not enough elixir", "All possible buildings of this type have already been placed", "Cell is occupied":
			http.Error(w, err.Error(), http.StatusConflict)
		default:
			http.Error(w, "Failed to place building-internal server error", http.StatusInternalServerError)
		}
	}

	var res Response
	res.Message = "Building placed successfully"

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

func MoveBuildingHandler(w http.ResponseWriter, request *http.Request) {
	val := request.Context().Value("player_id")
	_, ok := val.(string)
	if !ok {
		http.Error(w, "Player ID missing or is invalid in context", http.StatusInternalServerError)
		return
	}

	var req MoveBuildingReq
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&req)
	if err != nil {
		http.Error(w, "Bad request-wrong json payload", http.StatusBadRequest)
	}

	if req.New_x < 0 || req.New_y < 0 {
		http.Error(w, "Invalid building coordinates", http.StatusBadRequest)
	}
	if req.BuildingID == "" {
		http.Error(w, "Buidling id missing", http.StatusBadRequest)
	}

	buildingUUID, err := uuid.Parse(req.BuildingID)
	if err != nil {
		http.Error(w, "Failed to parse building id", http.StatusBadRequest)
	}

	err = models.MoveBuilding(request.Context(), buildingUUID, req.New_x, req.New_y)
	if err != nil {
		if err.Error() == "Cell is occupied" {
			http.Error(w, "Cell is occupied", http.StatusConflict)
		}
		http.Error(w, "Failed to move building-int. server error", http.StatusInternalServerError)
	}

	var res Response
	res.Message = "Building moved successfully"

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func StartUpgradeHandler(w http.ResponseWriter, request *http.Request) {
	val := request.Context().Value("player_id")
	_, ok := val.(string)
	if !ok {
		http.Error(w, "Player ID missing or is invalid in context", http.StatusInternalServerError)
		return
	}

	var req UpgradeBuildingReq
	err := json.NewDecoder(request.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Bad request-invalid json payload", http.StatusBadRequest)
		return
	}

	buildingUUID, err := uuid.Parse(req.BuildingID)
	if err != nil {
		http.Error(w, "Can't parse building id", http.StatusBadRequest)
	}

	err = models.StartUpgrade(request.Context(), buildingUUID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
	}

	var res Response
	res.Message = "Upgrade started successfully"

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func FinishUpgradeHandler(w http.ResponseWriter, request *http.Request) {
	val := request.Context().Value("player_id")
	_, ok := val.(string)
	if !ok {
		http.Error(w, "Player ID missing or is invalid in context", http.StatusInternalServerError)
		return
	}

	var req UpgradeBuildingReq
	err := json.NewDecoder(request.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Bad request-invalid json payload", http.StatusBadRequest)
	}

	buildingUUID, err := uuid.Parse(req.BuildingID)
	if err != nil {
		http.Error(w, "Can't parse building id", http.StatusBadRequest)
	}

	err = models.FinishUpgrade(request.Context(), buildingUUID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
	}

	var res Response
	res.Message = "Building upgraded successfully"

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}
