package db

import (
	"fmt"
	"github.com/mattermost/mattermost-server/v6/model"
)

type MattermostClient struct {
	Client *model.Client4
}

func NewMattermostClient(url, token string) *MattermostClient {
	client := model.NewAPIv4Client(url)
	client.SetToken(token)
	return &MattermostClient{
		Client: client,
	}
}

// GetRecentMessages pulls the last N messages from a channel
func (c *MattermostClient) GetRecentMessages(channelID string, limit int) ([]string, error) {
	postList, _, err := c.Client.GetPostsForChannel(channelID, 0, limit, "", false)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch mattermost posts: %w", err)
	}

	var messages []string
	// Order is usually newest first in the list
	for _, postID := range postList.Order {
		post := postList.Posts[postID]
		// Filter out system messages and bots if needed
		if post.Type == "" {
			messages = append(messages, post.Message)
		}
	}

	return messages, nil
}
