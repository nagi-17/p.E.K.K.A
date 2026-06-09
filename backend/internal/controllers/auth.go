package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nagi-17/p.E.K.K.A/internal/config"
	"github.com/nagi-17/p.E.K.K.A/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type RegisterResponse struct {
	Message   string `json:"message"`
	Player_ID string `json:"player_id"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Message   string `json:"message"`
	Token     string `json:"token"`
	Player_ID string `json:"player_id"`
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func generateJWT(player_ID string) (string, error) {
	jwtConfig := config.LoadConfig()

	mapClaims := jwt.MapClaims{
		"player_id":  player_ID,
		"expiration": time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, mapClaims)
	return token.SignedString([]byte(jwtConfig.JWTSecret))
}

func Register(w http.ResponseWriter, request *http.Request) {

	var reg_req RegisterRequest
	err := json.NewDecoder(request.Body).Decode(&reg_req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if reg_req.Username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}
	if reg_req.Email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}
	if reg_req.Password == "" {
		http.Error(w, "Password is required", http.StatusBadRequest)
		return
	}

	hashedpassword, err := hashPassword(reg_req.Password)
	if err != nil {
		http.Error(w, "Failed to hash password due to internal server error", http.StatusInternalServerError)
		return
	}

	player_ID, err := models.RegisterNewPlayer(request.Context(), reg_req.Username, reg_req.Email, hashedpassword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var reg_res RegisterResponse
	reg_res.Message = "User registered succesfully"
	reg_res.Player_ID = player_ID

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(reg_res)
}

func Login(w http.ResponseWriter, request *http.Request) {

	var login_req LoginRequest
	err := json.NewDecoder(request.Body).Decode(&login_req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if login_req.Username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}
	if login_req.Password == "" {
		http.Error(w, "Password is required", http.StatusBadRequest)
		return
	}

	playerInfo, err := models.GetLoginInfoUsingUsername(request.Context(), login_req.Username)
	if err != nil {
		http.Error(w, "Wrong username or password", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(playerInfo.Password_Hash), []byte(login_req.Password))
	if err != nil {
		http.Error(w, "Wrong password: password and hashed password do not match", http.StatusUnauthorized)
		return
	}

	tokenString, err := generateJWT(playerInfo.ID.String())
	if err != nil {
		http.Error(w, "Failed to generate jwt-token due to internal server error", http.StatusInternalServerError)
		return
	}

	var login_res LoginResponse
	login_res.Message = "Login Successful"
	login_res.Player_ID = playerInfo.ID.String()
	login_res.Token = tokenString

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(login_res)
}
