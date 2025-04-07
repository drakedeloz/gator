package rss

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"

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
	if len(cmd.Args) < 1 {
		return fmt.Errorf("time_between_reqs not provided")
	}

	duration, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("could not parse duration: %v", err)
	}

	fmt.Printf("Collecting feeds every %v\n", duration)
	ticker := time.NewTicker(duration)
	for ; ; <-ticker.C {
		ScrapeFeeds(s)
	}

}

func AddFeed(s *core.State, cmd core.Command, user database.User) error {
	if len(cmd.Args) < 2 {
		return fmt.Errorf("%v command usage: gator addfeed name url", cmd.Name)
	}

	newFeed, err := s.Queries.CreateFeed(context.Background(), database.CreateFeedParams{
		UserID: user.ID,
		Name:   cmd.Args[0],
		Url:    cmd.Args[1],
	})
	if err != nil {
		return fmt.Errorf("failed to create new feed: %v", err)
	}

	_, err = s.Queries.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		UserID: user.ID,
		FeedID: newFeed.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to follow feed: %v", err)
	}

	fmt.Println(newFeed)
	return nil
}

func Feeds(s *core.State, cmd core.Command) error {
	dbFeeds, err := s.Queries.GetAllFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("could not get feeds: %v", err)
	}

	if len(dbFeeds) == 0 {
		fmt.Println("No feeds found in database")
		return nil
	}

	for _, feed := range dbFeeds {
		fmt.Println("----------")
		fmt.Println(feed.FeedName)
		fmt.Println(feed.Url)
		fmt.Println(feed.UserName)
	}
	fmt.Println("----------")
	return nil
}

func Follow(s *core.State, cmd core.Command, user database.User) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("no url provided")
	}

	dbFeed, err := s.Queries.GetFeedByUrl(context.Background(), cmd.Args[0])
	if err != nil {
		return fmt.Errorf("feed not found")
	}

	feedFollow, err := s.Queries.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		UserID: user.ID,
		FeedID: dbFeed.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to follow feed: %v", err)
	}

	fmt.Println(feedFollow.FeedName)
	fmt.Printf("* %v\n", user.Name)
	return nil
}

func Following(s *core.State, cmd core.Command, user database.User) error {
	feeds, err := s.Queries.GetFeedFollowsByUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("could not get followed feeds: %v", err)
	}

	if len(feeds) == 0 {
		fmt.Println("You are not following any feeds!")
		return nil
	}

	for _, feed := range feeds {
		fmt.Printf("* %v\n", feed.FeedName)
	}
	return nil
}

func Unfollow(s *core.State, cmd core.Command, user database.User) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("no url provided")
	}

	dbFeed, err := s.Queries.GetFeedByUrl(context.Background(), cmd.Args[0])
	if err != nil {
		return fmt.Errorf("feed not found")
	}

	err = s.Queries.UnfollowFeed(context.Background(), database.UnfollowFeedParams{
		UserID: user.ID,
		FeedID: dbFeed.ID,
	})
	if err != nil {
		return fmt.Errorf("could not unfollow feed: %v", err)
	}

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

func ScrapeFeeds(s *core.State) error {
	feed, err := s.Queries.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("could not get next feed to fetch: %v", err)
	}

	err = s.Queries.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		return fmt.Errorf("could not mark feed as fetched: %v", err)
	}

	rssFeed, err := FetchFeed(context.Background(), feed.Url)
	if err != nil {
		return fmt.Errorf("could not fetch feed: %v", err)
	}

	fmt.Println("----------")
	fmt.Println(rssFeed.Channel.Title)
	fmt.Println(rssFeed.Channel.Description)
	fmt.Println("----------")

	for _, item := range rssFeed.Channel.Item {
		fmt.Printf("* %v\n", item.Title)
	}

	fmt.Println("----------")
	return nil
}
