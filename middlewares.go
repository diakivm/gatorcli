package main

import (
	"context"
	"fmt"

	"github.com/diakivm/gatorcli/internal/database"
)

func middlewareLoggedIn(handler func(*state, command, database.User) error) func(*state, command) error {
	return func(st *state, cmd command) error {
		user, err := st.db.GetUser(context.Background(), st.config.CurrentUserName)
		if err != nil {
			return fmt.Errorf("failed to get user: %v", err)
		}
		return handler(st, cmd, user)
	}
}
