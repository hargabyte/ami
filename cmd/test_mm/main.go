package main

import (
	"fmt"
	"github.com/hargabyte/ami/internal/db"
	"os"
)

func main() {
	url := "https://chat.hargabyte.com"
	token := "u196aooqppr6fxx43cqpphwydo"
	channelID := "55f4kwfcbjyipemg13ooua5j4a"

	client := db.NewMattermostClient(url, token)
	messages, err := client.GetRecentMessages(channelID, 5)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Fetched %d messages:\n", len(messages))
	for i, msg := range messages {
		fmt.Printf("%d: %s\n", i+1, msg)
	}
}
