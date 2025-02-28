package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/TheOTG/gator/internal/database"
	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error {
	if cmd.name != "login" {
		return errors.New("invalid handler")
	}
	if len(cmd.args) == 0 {
		return errors.New("missing argument")
	}

	username := cmd.args[0]

	dbUser, err := s.dbq.GetUser(context.Background(), username)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	err = s.cfg.SetUser(dbUser.Name)
	if err != nil {
		return fmt.Errorf("unable to set user: %s", err)
	}

	fmt.Printf("User has been set to: %s", username)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if cmd.name != "register" {
		return errors.New("invalid handler")
	}
	if len(cmd.args) == 0 {
		return errors.New("missing argument")
	}

	name := cmd.args[0]
	createUserParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
	}
	dbUser, err := s.dbq.CreateUser(context.Background(), createUserParams)
	if err != nil {
		return err
	}

	err = s.cfg.SetUser(name)
	if err != nil {
		return fmt.Errorf("unable to set user: %s", err)
	}

	fmt.Printf("User has been created: %v", dbUser)

	return nil
}

func handlerReset(s *state, cmd command) error {
	if cmd.name != "reset" {
		return errors.New("invalid handler")
	}

	err := s.dbq.DeleteUser(context.Background())
	if err != nil {
		return fmt.Errorf("unable to reset users: %s", err)
	}
	return nil
}

func handlerUsers(s *state, cmd command) error {
	if cmd.name != "users" {
		return errors.New("invalid handler")
	}

	dbUsers, err := s.dbq.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("unable to get users: %s", err)
	}

	for _, user := range dbUsers {
		toPrint := "* " + user.Name
		if user.Name == s.cfg.CurrentUserName {
			toPrint += " (current)"
		}
		fmt.Println(toPrint)
	}

	return nil
}
