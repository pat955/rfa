package api

import (
	"context"
	"fmt"
	"sync"
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
				fmt.Println(f.Title)
			}
		}()
	}
	wg.Wait()
	fmt.Println("done with it all")

}

func retrieve(endpoint string) RSS {
	return FetchFeed(endpoint)
}
