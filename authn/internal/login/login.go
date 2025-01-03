package login

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/ctroller/chirper/authn/internal/inject"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token  string `json:"token"`
	Status string `json:"status"`
}

// LoginHandler godoc
//
//	@Summary		Login
//	@Description	login
//	@Tags			login
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	LoginResponse
//	@Router			/api/v1/login [post]
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	token, err := authenticateUser(req.Username, req.Password)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	res := LoginResponse{
		Token:  token,
		Status: "successful",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func authenticateUser(username, password string) (string, error) {
	var passwordHash string
	err := inject.App.DBPool.QueryRow(context.Background(), "SELECT password_hash FROM users WHERE username=$1", username).Scan(&passwordHash)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	if err != nil {
		return "", err
	}

	return "dummy-token", nil
}
