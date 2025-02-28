package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/TheOTG/gator/internal/database"
	"github.com/google/uuid"
)

func handlerAggregate(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("missing argument")
	}

	time_between_reqs := cmd.args[0]

	dur, err := time.ParseDuration(time_between_reqs)
	if err != nil {
		return err
	}

	ticker := time.NewTicker(dur)
	fmt.Printf("Collecting feeds every %s\n", time_between_reqs)
	for ; ; <-ticker.C {
		err = scrapeFeeds(s)
		if err != nil {
			log.Print("Unable to scrape feed")
			return err
		}
	}
}

func scrapeFeeds(s *state) error {
	dbFeed, err := s.dbq.GetNextFeedToFetch(context.Background())
	if err != nil {
		log.Print("Unable to get next feed to fetch")
		return err
	}

	err = s.dbq.MarkFeedFetched(context.Background(), dbFeed.ID)
	if err != nil {
		log.Print("Unable to mark feed as fetched")
		return err
	}

	rssFeed, err := fetchFeed(context.Background(), dbFeed.Url)
	if err != nil {
		log.Print("Unable to fetch feed")
		return err
	}

	fmt.Println(rssFeed.Channel.Title)
	for _, rssItem := range rssFeed.Channel.Item {
		pubDate := sql.NullTime{}
		t, err := time.Parse(time.RFC1123Z, rssItem.PubDate)
		if err == nil {
			pubDate = sql.NullTime{
				Time:  t,
				Valid: true,
			}
		}

		createPostParams := database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Title:     rssItem.Title,
			Url:       rssItem.Link,
			Description: sql.NullString{
				String: rssItem.Description,
				Valid:  true,
			},
			PublishedAt: pubDate,
			FeedID:      dbFeed.ID,
		}
		_, err = s.dbq.CreatePost(context.Background(), createPostParams)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				continue
			}
			log.Printf("Couldn't create post: %v", err)
			continue
		}
	}

	return nil
}
