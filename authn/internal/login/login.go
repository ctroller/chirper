package login

import (
	"encoding/json"
	"log/slog"
	"net/http"
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
		slog.Error("Invalid request payload", slog.Any("error", err))
		return
	}

	// Implement login logic here
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
	return "dummy-token", nil
}
