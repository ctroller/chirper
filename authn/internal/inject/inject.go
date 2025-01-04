package inject

import (
	"github.com/ctroller/chirper/authn/lib/user"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Application struct {
	DBPool         *pgxpool.Pool
	UserRepository user.UserRepository
}

var App Application
