package api

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/pat955/rss_feed_aggregator/internal/database"
)

func RetrieveGroup() {
	fmt.Println("stating retrieving group")
	var n int32 = 5
	db := connect().DB

	var wg sync.WaitGroup

	feeds, _ := db.GetNextFeedsToFetch(context.Background(), n)
	for i, feed := range feeds {
		if feed.Url == "something.com" {
			continue
		}
		i++
		wg.Add(i)
		go func() {
			defer wg.Done()
			for _, f := range retrieve(feed.Url).Channel.Items {
				date, err := time.Parse("Mon, 02 Jan 2006 15:04:05 -0700", f.PubDate)
				if err != nil {
					panic(err)
				}
				err = db.CreatePost(context.Background(), database.CreatePostParams{
					ID:          uuid.New(),
					CreatedAt:   time.Now().UTC(),
					UpdatedAt:   time.Now().UTC(),
					Title:       f.Title,
					Url:         f.Link,
					Description: f.Description,
					PublishedAt: date,
					FeedID:      feed.ID,
				})
				if err != nil && strings.Contains(err.Error(), "duplicate key") {
					continue
				} else {
					log.Fatal(err)
				}
			}
		}()
	}
	wg.Wait()
	fmt.Println("done with it all")

}

func retrieve(endpoint string) RSS {
	return FetchFeed(endpoint)
}
