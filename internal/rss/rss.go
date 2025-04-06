package rss

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"

	"github.com/drakedeloz/gator/internal/core"
	"github.com/drakedeloz/gator/internal/database"
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

func Aggregate(s *core.State, cmd core.Command) error {
	feed, err := FetchFeed(context.Background(), "https://www.wowhead.com/news/rss/all")
	if err != nil {
		return err
	}

	fmt.Println(feed)
	return nil
}

func AddFeed(s *core.State, cmd core.Command) error {
	if len(cmd.Args) < 2 {
		return fmt.Errorf("%v command usage: gator addfeed name url\n", cmd.Name)
	}

	dbUser, err := s.Queries.GetUser(context.Background(), s.Config.CurrentUser)
	if err != nil {
		return fmt.Errorf("could not get current user %v from db: %v", s.Config.CurrentUser, err)
	}

	newFeed, err := s.Queries.CreateFeed(context.Background(), database.CreateFeedParams{
		UserID: dbUser.ID,
		Name:   cmd.Args[0],
		Url:    cmd.Args[1],
	})
	if err != nil {
		return fmt.Errorf("failed to create new feed: %v", err)
	}

	fmt.Println(newFeed)
	return nil
}

func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, nil)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("User-Agent", "gator")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("could not read request body: %v", err)
	}

	var feed RSSFeed
	err = xml.Unmarshal(body, &feed)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("could not unmarshal request body: %v", err)
	}

	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)

	for i, item := range feed.Channel.Item {
		feed.Channel.Item[i].Title = html.UnescapeString(item.Title)
		feed.Channel.Item[i].Description = html.UnescapeString(item.Description)
	}

	return &feed, nil
}
