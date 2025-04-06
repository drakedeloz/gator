package core

import (
	"context"
	"fmt"
	"log"
)

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	AllCommands map[string]func(*State, Command) error
}

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("no username provided")
	}

	dbUser, err := s.Queries.GetUser(context.Background(), cmd.Args[0])
	if err != nil {
		return fmt.Errorf("user not found")
	}

	err = s.Config.SetUser(dbUser.Name)
	if err != nil {
		return err
	}

	fmt.Printf("User has been set to: %v\n", cmd.Args[0])
	return nil
}

func HandlerRegister(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("no username provided")
	}

	dbUser, err := s.Queries.CreateUser(context.Background(), cmd.Args[0])
	if err != nil {
		return err
	}

	err = s.Config.SetUser(dbUser.Name)
	if err != nil {
		return err
	}

	fmt.Printf("Created User: %v\n", dbUser.Name)
	fmt.Println(dbUser)
	return nil
}

func (c *Commands) Register(name string, f func(*State, Command) error) {
	if c.AllCommands == nil {
		c.AllCommands = make(map[string]func(*State, Command) error)
	}
	if _, ok := c.AllCommands[name]; ok {
		log.Fatalf("%v command has already been registered", name)
		return
	}
	c.AllCommands[name] = f
}

func (c *Commands) Run(s *State, cmd Command) error {
	if _, ok := c.AllCommands[cmd.Name]; !ok {
		return fmt.Errorf("%v command does not exist", cmd.Name)
	}
	cmdFunc := c.AllCommands[cmd.Name]
	err := cmdFunc(s, cmd)
	if err != nil {
		return err
	}

	return nil
}
