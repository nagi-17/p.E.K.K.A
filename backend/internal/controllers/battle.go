package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/nagi-17/p.E.K.K.A/internal/models"
	"github.com/nagi-17/p.E.K.K.A/internal/services"
)

type AttackRequest struct {
	DefenderID string `json:"defender_id"`
}

func MatchMakeHandler(w http.ResponseWriter, request *http.Request) {
	val := request.Context().Value("player_id")
	playerID, ok := val.(string)
	if !ok {
		http.Error(w, "Player ID missing or is invalid in context", http.StatusInternalServerError)
		return
	}

	defender, err := models.FindOpponent(request.Context(), playerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(defender)
}

func AttackHandler(w http.ResponseWriter, request *http.Request) {
	val := request.Context().Value("player_id")
	playerID, ok := val.(string)
	if !ok {
		http.Error(w, "Player ID missing or is invalid in context", http.StatusInternalServerError)
		return
	}

	var req AttackRequest
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&req)
	if err != nil {
		http.Error(w, "Bad request-wrong json payload", http.StatusBadRequest)
		return
	}

	if req.DefenderID == "" {
		http.Error(w, "Who will you attack? Defender ID empty", http.StatusBadRequest)
		return
	}

	battleLog, err := services.Attack(request.Context(), playerID, req.DefenderID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(battleLog)
}
