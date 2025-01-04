package login

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/ctroller/chirper/authn/internal/inject"
	"github.com/golang-jwt/jwt/v5"
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
		slog.Debug("Authentication failed", "error", err)
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

var JWT_KEY = []byte("dummy-key") // TODO: change this to a secure key

func authenticateUser(username, password string) (string, error) {
	user, err := inject.App.UserRepository.Find(username)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "chirper",
		"sub": user.ID,
	})
	return token.SignedString(JWT_KEY)
}
