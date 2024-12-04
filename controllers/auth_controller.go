package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"rest-services/config"
	"rest-services/utils"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type SignUpRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RevokeRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// SignUp handles user registration
func SignUp(w http.ResponseWriter, r *http.Request) {
	var req SignUpRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, `{"error":"Invalid input"}`, http.StatusBadRequest)
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, `{"error":"Failed to hash password"}`, http.StatusInternalServerError)
		return
	}

	// Insert user into the database
	_, err = config.DB.Exec("INSERT INTO users (email, password) VALUES (?, ?)", req.Email, hashedPassword)
	if err != nil {
		http.Error(w, `{"error":"Email already exists"}`, http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message":"User signed up successfully"}`))
}

// SignIn authenticates a user and returns an access token and refresh token
func SignIn(w http.ResponseWriter, r *http.Request) {
	var req SignInRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, `{"error":"Invalid input"}`, http.StatusBadRequest)
		return
	}

	var storedPassword string
	err = config.DB.QueryRow("SELECT password FROM users WHERE email = ?", req.Email).Scan(&storedPassword)
	if err != nil {
		http.Error(w, `{"error":"Invalid email or password"}`, http.StatusUnauthorized)
		return
	}

	// Compare passwords
	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(req.Password))
	if err != nil {
		http.Error(w, `{"error":"Invalid email or password"}`, http.StatusUnauthorized)
		return
	}

	// Generate access and refresh tokens
	accessToken, refreshToken, err := utils.GenerateToken(req.Email)
	if err != nil {
		http.Error(w, `{"error":"Error generating tokens"}`, http.StatusInternalServerError)
		return
	}

	// Store refresh token in database (optional, if you want server-side tracking)
	_, err = config.DB.Exec("INSERT INTO refresh_tokens (email, token, expiry) VALUES (?, ?, ?)", req.Email, refreshToken, time.Now().Add(7*24*time.Hour))
	if err != nil {
		log.Println("Failed to store refresh token:", err)
	}

	response := map[string]interface{}{
		"token":         accessToken,
		"refresh_token": refreshToken,
		"expires_in":    900, // 15 minutes for access token
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// RefreshToken generates a new access token using the refresh token
func RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req RefreshRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, `{"error":"Invalid input"}`, http.StatusBadRequest)
		return
	}

	// Parse and validate the refresh token
	claims, err := utils.ParseToken(req.RefreshToken)
	if err != nil {
		http.Error(w, `{"error":"Invalid or expired refresh token"}`, http.StatusUnauthorized)
		return
	}

	// Generate new access and refresh tokens
	accessToken, refreshToken, err := utils.GenerateToken(claims.Email)
	if err != nil {
		http.Error(w, `{"error":"Error generating new tokens"}`, http.StatusInternalServerError)
		return
	}

	// Update refresh token in the database (optional)
	_, err = config.DB.Exec("UPDATE refresh_tokens SET token = ?, expiry = ? WHERE email = ?", refreshToken, time.Now().Add(7*24*time.Hour), claims.Email)
	if err != nil {
		log.Println("Failed to update refresh token:", err)
	}

	response := map[string]interface{}{
		"token":         accessToken,
		"refresh_token": refreshToken,
		"expires_in":    900, // 15 minutes for access token
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// RevokeToken invalidates the refresh token
func RevokeToken(w http.ResponseWriter, r *http.Request) {
	var req RevokeRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, `{"error":"Invalid input"}`, http.StatusBadRequest)
		return
	}

	// Invalidate the refresh token (delete it from the database)
	_, err = config.DB.Exec("DELETE FROM refresh_tokens WHERE token = ?", req.RefreshToken)
	if err != nil {
		http.Error(w, `{"error":"Failed to revoke token"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Token revoked successfully"}`))
}

// ProtectedEndpoint is an example of an authenticated endpoint
func ProtectedEndpoint(w http.ResponseWriter, r *http.Request) {
	// Get the user email from the context (set by the AuthMiddleware)
	email := r.Context().Value("email").(string)

	// Respond with a success message
	response := map[string]string{
		"message": "Welcome to the protected endpoint",
		"user":    email,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
