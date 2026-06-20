package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/nagi-17/p.E.K.K.A/internal/models"
)

type TrainAmryReq struct {
	TroopType string `json:"troop_type"`
	Quantity  int    `json:"quantity"`
}

func TrainArmy(w http.ResponseWriter, request *http.Request) {
	val := request.Context().Value("player_id")
	playerID, ok := val.(string)
	if !ok {
		http.Error(w, "Player ID missing or is invalid in context", http.StatusInternalServerError)
		return
	}

	var req TrainAmryReq
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&req)
	if err != nil {
		http.Error(w, "Bad request-wrong json payload", http.StatusBadRequest)
		return
	}

	if req.Quantity <= 0 {
		http.Error(w, "Quantity of troops is invalid", http.StatusBadRequest)
		return
	}

	err = models.TrainArmy(request.Context(), playerID, req.TroopType, req.Quantity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	var res Response
	res.Message = "Troop trained successfully"

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func GetArmy(w http.ResponseWriter, request *http.Request) {
	val := request.Context().Value("player_id")
	playerID, ok := val.(string)
	if !ok {
		http.Error(w, "Player ID missing or is invalid in context", http.StatusInternalServerError)
		return
	}

	army, err := models.GetTrainedArmy(request.Context(), playerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(army)
}

type DiscardArmyReq struct {
	TroopType string `json:"troop_type"`
	Quantity  int    `json:"quantity"`
}

func DiscardArmy(w http.ResponseWriter, request *http.Request) {
	val := request.Context().Value("player_id")
	playerID, ok := val.(string)
	if !ok {
		http.Error(w, "Player ID missing or is invalid in context", http.StatusInternalServerError)
		return
	}

	var req DiscardArmyReq
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&req)
	if err != nil {
		http.Error(w, "Bad request-wrong json payload", http.StatusBadRequest)
		return
	}

	if req.Quantity <= 0 {
		http.Error(w, "Quantity of troops is invalid", http.StatusBadRequest)
		return
	}

	err = models.DiscardTroops(request.Context(), playerID, req.TroopType, req.Quantity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	var res Response
	res.Message = "Troops discarded successfully"

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}
