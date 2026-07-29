package command

import (
	"fmt"

	"github.com/csrrmrvll/gator/internal/config"
)

type State struct {
	Config *config.Config
}

type Command struct {
	Name string
	Args []string
}

func HandlerLogin(state *State, cmd Command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("username is required")
	}
	username := cmd.Args[0]
	err := state.Config.SetUser(username)
	if err != nil {
		return err
	}
	fmt.Printf("Logged in as %s\n", username)
	return nil
}

type Commands struct {
	commands map[string]func(*State, Command) error
}

func (c *Commands) Register(name string, handler func(*State, Command) error) {
	if c.commands == nil {
		c.commands = make(map[string]func(*State, Command) error)
	}
	c.commands[name] = handler
}

func (c *Commands) Run(state *State, cmd Command) error {
	handler, exists := c.commands[cmd.Name]
	if !exists {
		return fmt.Errorf("unknown command: %s", cmd.Name)
	}
	return handler(state, cmd)
}
