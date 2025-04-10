package main

import (
	"github.com/sambakker4/gator/internal/config"
	"log"
	"os"
	"errors"
	"database/sql"
	"github.com/sambakker4/gator/internal/database"
)

import _ "github.com/lib/pq"

func main() { 
	cfg, err := config.Read()	
	if err != nil {
		log.Fatal(err)
	}

	currentState := state{config: &cfg}
	commandsMap := commands{cmds: make(map[string]func(*state, command)error)}
	commandsMap.register("login", handlerLogin)
	commandsMap.register("register", handlerRegister)
	commandsMap.register("reset", handlerReset)
	commandsMap.register("users", handlerGetUsers)
	commandsMap.register("agg", handlerAgg)
	commandsMap.register("addfeed", middlewareLoggedIn(handlerNewFeed))
	commandsMap.register("feeds", handlerFeeds)
	commandsMap.register("follow", middlewareLoggedIn(handlerFollow))
	commandsMap.register("following", middlewareLoggedIn(handlerFollowing))
	commandsMap.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	commandsMap.register("browse", middlewareLoggedIn(handlerBrowse))

	args := os.Args

	if len(args) < 2 {
		log.Fatal(errors.New("Must have at least two arguments"))
	}

	db, err := sql.Open("postgres", currentState.config.DBURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	dbQueries := database.New(db)
	currentState.db = dbQueries

	currentCommand := command{name: args[1], args: args[2:]}
	err = commandsMap.run(&currentState, currentCommand)
	if err != nil {
		log.Fatal(err)
	}
}
