package main

import (
	"fmt"
	"log"
	"os"

	"github.com/csrrmrvll/gator/internal/command"
	"github.com/csrrmrvll/gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	fmt.Printf("Read config: %+v\n", cfg)

	state := command.State{Config: &cfg}
	commands := command.Commands{}
	commands.Register("login", command.HandlerLogin)

	args := os.Args[1:]
	if len(args) < 1 {
		log.Fatalf("no command provided")
	}

	cmd := command.Command{Name: args[0], Args: args[1:]}
	err = commands.Run(&state, cmd)
	if err != nil {
		log.Fatalf("error running command: %v", err)
	}
}
