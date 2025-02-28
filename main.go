package main

import (
	"database/sql"
	"errors"
	"log"
	"os"

	"github.com/TheOTG/gator/internal/config"
	"github.com/TheOTG/gator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	db  *sql.DB
	dbq *database.Queries
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	list map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.list[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	fn, ok := c.list[cmd.name]
	if !ok {
		return errors.New("unregistered command")
	}

	return fn(s, cmd)
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Printf("Unable to read config: %s", err)
		os.Exit(1)
	}

	db, err := sql.Open("postgres", cfg.DBUrl)
	if err != nil {
		log.Printf("Unable to connect to database: %s", err)
		os.Exit(1)
	}
	defer db.Close()

	dbQueries := database.New(db)

	s := &state{
		db:  db,
		dbq: dbQueries,
		cfg: cfg,
	}

	cmds := commands{
		list: map[string]func(*state, command) error{},
	}

	args := os.Args
	if len(args) < 2 {
		log.Printf("Not enough arguments were provided")
		os.Exit(1)
	}

	cmd := command{
		name: args[1],
		args: args[2:],
	}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAggregate)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerListFeed)
	cmds.register("follow", middlewareLoggedIn(handlerFollowFeed))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	cmds.register("browse", middlewareLoggedIn(handlerBrowse))

	err = cmds.run(s, cmd)
	if err != nil {
		log.Printf("Error running command: %s", err)
		os.Exit(1)
	}
}
