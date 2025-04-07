package middleware

import (
	"context"

	"github.com/drakedeloz/gator/internal/core"
	"github.com/drakedeloz/gator/internal/database"
)

func MiddlewareLoggedIn(handler func(s *core.State, cmd core.Command, user database.User) error) func(*core.State, core.Command) error {
	return func(s *core.State, cmd core.Command) error {
		dbUser, err := s.Queries.GetUser(context.Background(), s.Config.CurrentUser)
		if err != nil {
			return err
		}

		return handler(s, cmd, dbUser)
	}
}
