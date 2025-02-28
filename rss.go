package main

import (
	"context"
	"encoding/xml"
	"html"
	"io"
	"log"
	"net/http"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		log.Printf("Unable to make new request: %s", err)
		return nil, err
	}
	req.Header.Set("User-Agent", "gator")

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Unable to make new request: %s", err)
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Unable to read response body: %s", err)
		return nil, err
	}

	rssFeed := &RSSFeed{}

	err = xml.Unmarshal(body, rssFeed)
	if err != nil {
		log.Printf("Unable to unmarshal xml: %s", err)
		return nil, err
	}

	rssFeed.Channel.Title = html.UnescapeString(rssFeed.Channel.Title)
	rssFeed.Channel.Description = html.UnescapeString(rssFeed.Channel.Description)

	for _, rssItem := range rssFeed.Channel.Item {
		rssItem.Title = html.UnescapeString(rssItem.Title)
		rssItem.Description = html.UnescapeString(rssItem.Description)
	}

	return rssFeed, nil
}
