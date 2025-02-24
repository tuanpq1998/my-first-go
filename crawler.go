package main

import (
	"context"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
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
	// logStr := ""
	for _, rssItem := range rssFeed.Channel.Item {
		description := pgtype.Text{}
		if rssItem.Description != "" {
			description.String = rssItem.Description
			description.Valid = true
		}

		publishedAt, err := time.Parse(time.RFC1123Z, rssItem.PublishDate)
		if err != nil {
			// log.Println("Crawler::scrapeFeed::ParseTime::error", err)
			// continue
			publishedAt, _ = time.Parse(time.RFC1123Z, "Tue, 04 Feb 2025 00:00:00 +0000")
		}

		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID: pgtype.UUID{
				Bytes: uuid.New(),
				Valid: true,
			},
			CreatedAt: pgtype.Timestamp{
				Time:  time.Now().UTC(),
				Valid: true,
			},
			UpdatedAt: pgtype.Timestamp{
				Time:  time.Now().UTC(),
				Valid: true,
			},
			Title:       rssItem.Title,
			Description: description,
			Url:         rssItem.Link,
			PublishedAt: pgtype.Timestamp{
				Time:  publishedAt,
				Valid: true,
			},
			FeedID: feed.ID,
		})
		// logStr += fmt.Sprintf("\nFeedItem:%v - Url: %v", rssItem.Title, rssItem.Link)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint \"posts_url_key\"") { //ignore the error of `duplicate key value violates unique constraint "posts_url_key"`
				continue
			}
			log.Println("Crawler::scrapeFeed::CreatePost::error", err)
		}
	}
	// log.Println(logStr)
	log.Println("Crawler::scrapeFeed::Done::Name", feed.Title, "::Len::", len(rssFeed.Channel.Item))

}
