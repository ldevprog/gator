package commands

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/levon-dalakyan/gator/internal/database"
	"github.com/levon-dalakyan/gator/internal/state"
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

func HandlerAgg(s *state.State, cmd state.Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("The argument for time between requests is required")
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return err
	}
	fmt.Printf("Collecting feeds every %v\n", timeBetweenRequests)

	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func scrapeFeeds(s *state.State) error {
	nextFeedToFetch, err := s.DB.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}
	err = s.DB.MarkFeedFetched(context.Background(), nextFeedToFetch.ID)
	if err != nil {
		return err
	}

	feed, err := fetchFeed(context.Background(), nextFeedToFetch.Url)
	if err != nil {
		return err
	}

	layout := time.RFC822
	for _, item := range feed.Channel.Item {
		pubDate, err := time.Parse(layout, item.PubDate)
		if err != nil {
			fmt.Println("Error parsing time:", err)
		}
		s.DB.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Url:         item.Link,
			Description: item.Description,
			PublishedAt: sql.NullTime{Time: pubDate, Valid: err == nil},
			FeedID:      nextFeedToFetch.ID,
		})
		if err != nil {
			fmt.Println(err)
		}
	}

	return err
}

func fetchFeed(ctx context.Context, feedUrl string) (*RSSFeed, error) {
	fmt.Println("** Fetching feed **")

	req, err := http.NewRequestWithContext(ctx, "GET", feedUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "gator")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(res.Body)
	defer res.Body.Close()

	var feed *RSSFeed
	err = xml.Unmarshal(data, &feed)
	if err != nil {
		return nil, err
	}

	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	for i := 0; i < len(feed.Channel.Item); i++ {
		feed.Channel.Item[i].Title =
			html.UnescapeString(feed.Channel.Item[i].Title)
		feed.Channel.Item[i].Description =
			html.UnescapeString(feed.Channel.Item[i].Description)
	}

	return feed, nil
}
