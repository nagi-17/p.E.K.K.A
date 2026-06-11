package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/nagi-17/p.E.K.K.A/internal/models"
)

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
