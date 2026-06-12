package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/nagi-17/p.E.K.K.A/internal/models"
)

func CollectResourceHandler(w http.ResponseWriter, request *http.Request) {
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

	err = models.CollectResource(request.Context(), buildingUUID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
	}

	var res Response
	res.Message = "Resources collected"

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}
