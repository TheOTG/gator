package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/TheOTG/gator/internal/database"
	"github.com/google/uuid"
)

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if cmd.name != "addfeed" {
		return errors.New("invalid handler")
	}

	if len(cmd.args) < 2 {
		return errors.New("missing argument")
	}

	name := cmd.args[0]
	url := cmd.args[1]
	createFeedParams := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	qtx := s.dbq.WithTx(tx)

	dbFeed, err := qtx.CreateFeed(context.Background(), createFeedParams)
	if err != nil {
		return err
	}

	createFeedFollowParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    dbFeed.ID,
	}

	_, err = qtx.CreateFeedFollow(context.Background(), createFeedFollowParams)
	if err != nil {
		return err
	}

	fmt.Printf("Feed '%s' has been created and you are now following it\n", dbFeed.Name)

	return tx.Commit()
}

func handlerListFeed(s *state, cmd command) error {
	if cmd.name != "feeds" {
		return errors.New("invalid handler")
	}

	dbFeeds, err := s.dbq.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	for _, dbFeed := range dbFeeds {
		fmt.Printf("Name: %s, URL: %s, by: %s\n", dbFeed.Name, dbFeed.Url, dbFeed.UserName)
	}

	return nil
}

func handlerFollowFeed(s *state, cmd command, user database.User) error {
	if cmd.name != "follow" {
		return errors.New("invalid handler")
	}

	if len(cmd.args) == 0 {
		return errors.New("missing argument")
	}

	url := cmd.args[0]

	dbFeed, err := s.dbq.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return err
	}

	createFeedFollowParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    dbFeed.ID,
	}

	dbFeedFollow, err := s.dbq.CreateFeedFollow(context.Background(), createFeedFollowParams)
	if err != nil {
		return err
	}

	fmt.Printf("%s have followed: %s", dbFeedFollow.UserName, dbFeedFollow.FeedName)

	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	if cmd.name != "following" {
		return errors.New("invalid handler")
	}

	dbFeedFollows, err := s.dbq.GetFollowing(context.Background(), user.ID)
	if err != nil {
		return err
	}

	for _, dbFeedFollow := range dbFeedFollows {
		fmt.Println(dbFeedFollow.FeedName)
	}

	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if cmd.name != "unfollow" {
		return errors.New("invalid handler")
	}

	if len(cmd.args) == 0 {
		return errors.New("missing argument")
	}

	url := cmd.args[0]

	dbFeed, err := s.dbq.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return err
	}

	unfollowParams := database.UnfollowParams{
		UserID: user.ID,
		FeedID: dbFeed.ID,
	}

	err = s.dbq.Unfollow(context.Background(), unfollowParams)
	if err != nil {
		return err
	}

	return nil
}
