package core

import (
	"github.com/drakedeloz/gator/internal/config"
	"github.com/drakedeloz/gator/internal/database"
)

type State struct {
	Queries *database.Queries
	Config  *config.Config
}
