package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/nagi-17/p.E.K.K.A/internal/models"
)

type TroopUpgradeReq struct {
	TroopType string `json:"troop_type"`
}

func StartTroopUpgrade(w http.ResponseWriter, request *http.Request) {
	val := request.Context().Value("player_id")
	playerID, ok := val.(string)
	if !ok {
		http.Error(w, "Player ID missing or is invalid in context", http.StatusInternalServerError)
		return
	}

	var req TroopUpgradeReq
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&req)
	if err != nil {
		http.Error(w, "Bad request-wrong json payload", http.StatusBadRequest)
		return
	}

	err = models.StartTroopUpgrade(request.Context(), playerID, req.TroopType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	var res Response
	res.Message = "Troop upgrade started successfully"

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func FinishTroopUpgrade(w http.ResponseWriter, request *http.Request) {
	val := request.Context().Value("player_id")
	playerID, ok := val.(string)
	if !ok {
		http.Error(w, "Player ID missing or is invalid in context", http.StatusInternalServerError)
		return
	}

	var req TroopUpgradeReq
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&req)
	if err != nil {
		http.Error(w, "Bad request-wrong json payload", http.StatusBadRequest)
		return
	}

	err = models.FinishTroopUpgrade(request.Context(), playerID, req.TroopType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	var res Response
	res.Message = "Troop upgraded successfully"

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}
