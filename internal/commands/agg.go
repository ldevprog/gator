package commands

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"

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
	url := "https://www.wagslane.dev/index.xml"

	feed, err := fetchFeed(context.Background(), url)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(feed)

	return nil
}

func fetchFeed(ctx context.Context, feedUrl string) (*RSSFeed, error) {
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
