package login

import (
	"fmt"
	"testing"

	"github.com/ctroller/chirper/authn/internal/inject"
	"github.com/ctroller/chirper/authn/lib/user"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func setupApp() {
	inject.App = inject.Application{
		UserRepository: &UserRepositoryMock{},
	}
}

type UserRepositoryMock struct{}

func (r *UserRepositoryMock) Find(name string) (*user.User, error) {
	return &user.User{
		ID:           1,
		Username:     "general_kenobi",
		PasswordHash: "$2a$10$.6bM20IjWVkojnVSewjQ0uNTQOMlBlfU7EUKzsJYUJolsb.DSmuW6",
	}, nil
}

func TestLogin(t *testing.T) {
	setupApp()

	token, err := authenticateUser("general_kenobi", "test")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	tk, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return JWT_KEY, nil
	})
	assert.NoError(t, err)

	userId, err := tk.Claims.GetSubject()
	assert.NoError(t, err)
	assert.Equal(t, "1", userId)
}
