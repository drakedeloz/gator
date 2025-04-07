package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/drakedeloz/gator/internal/config"
	"github.com/drakedeloz/gator/internal/core"
	"github.com/drakedeloz/gator/internal/database"
	"github.com/drakedeloz/gator/internal/rss"
	"github.com/drakedeloz/gator/middleware"
	_ "github.com/lib/pq"
)

func main() {
	var state core.State
	var cmds core.Commands
	state.Config = config.Read()

	db, err := sql.Open("postgres", state.Config.DB_URL)
	if err != nil {
		log.Fatalf("could not connect to db: %v", err)
		return
	}

	dbQueries := database.New(db)
	state.Queries = dbQueries

	cmds.Register("login", core.HandlerLogin)
	cmds.Register("register", core.HandlerRegister)
	cmds.Register("users", core.GetUsers)
	cmds.Register("addfeed", middleware.MiddlewareLoggedIn(rss.AddFeed))
	cmds.Register("feeds", rss.Feeds)
	cmds.Register("follow", middleware.MiddlewareLoggedIn(rss.Follow))
	cmds.Register("following", middleware.MiddlewareLoggedIn(rss.Following))
	cmds.Register("unfollow", middleware.MiddlewareLoggedIn(rss.Unfollow))
	cmds.Register("agg", rss.Aggregate)
	cmds.Register("browse", middleware.MiddlewareLoggedIn(rss.Browse))
	cmds.Register("reset", core.HandlerReset)

	args := os.Args
	if len(args) < 2 {
		log.Fatalf("no arguments provided")
		return
	}

	cmd := core.Command{
		Name: args[1],
		Args: args[2:],
	}

	err = cmds.Run(&state, cmd)
	if err != nil {
		log.Fatalf("could not run %v command: %v", cmd.Name, err)
		return
	}
}
