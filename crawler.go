package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/tuanpq1998/my-first-go/internal/database"
)

func startScraping(db *database.Queries, concurrency int, gapTime time.Duration) {
	log.Println("Crawler::startScraping::starting::concurrency", concurrency, "::gapTime::", gapTime)
	ticker := time.NewTicker(gapTime)
	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(context.Background(), int32(concurrency))
		if err != nil {
			log.Println("Crawler::startScraping::GetNextFeedsToFetch::error", err)
			continue
		}

		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)

			go scrapeFeed(db, feed, wg)
		}
		wg.Wait()
	}
}

func scrapeFeed(db *database.Queries, feed database.Feed, wg *sync.WaitGroup) {
	defer wg.Done()

	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Println("Crawler::scrapeFeed::MarkFeedAsFetched::error", err)
		return
	}

	rssFeed, err := urlToFeed(feed.Url)
	if err != nil {
		log.Println("Crawler::scrapeFeed::urlToFeed::error", err)
		return
	}

	//show result
	log.Println("Crawler::scrapeFeed::Done::Name", feed.Title, "::Len::", len(rssFeed.Channel.Item))
	logStr := ""
	for _, rssItem := range rssFeed.Channel.Item {
		logStr += fmt.Sprintf("FeedItem:%v - Url: %v\n", rssItem.Title, rssItem.Link)
	}
	log.Println(logStr)

}
