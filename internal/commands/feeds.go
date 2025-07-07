package commands

import (
	"context"
	"fmt"

	"github.com/levon-dalakyan/gator/internal/state"
)

func HandlerFeeds(s *state.State, cmd state.Command) error {
	feeds, err := s.DB.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	type feedDataOut struct {
		name              string
		url               string
		userAddedFeedName string
	}

	var feedsOut []feedDataOut
	for _, feed := range feeds {
		feedUserAddedName, err := s.DB.GetUserNameById(context.Background(), feed.UserID)
		if err != nil {
			return err
		}

		feedsOut = append(feedsOut, feedDataOut{
			name:              feed.Name,
			url:               feed.Url,
			userAddedFeedName: feedUserAddedName,
		})
	}

	for _, fo := range feedsOut {
		fmt.Println("--Feed--")
		fmt.Println("\tName:", fo.name)
		fmt.Println("\tUrl:", fo.url)
		fmt.Println("\tuserAddedFeedName:", fo.userAddedFeedName)
	}

	return nil
}
