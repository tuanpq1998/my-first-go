package main

import (
	"encoding/xml"
	"io"
	"net/http"
	"time"
)

type RssFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RssItem `xml:"item"`
	} `xml:"channel"`
}

type RssItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PublishDate string `xml:"pubDate"`
}

func urlToFeed(url string) (RssFeed, error) {
	res := RssFeed{}
	httpClient := http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := httpClient.Get(url)
	if err != nil {
		return res, err
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return res, err
	}

	err = xml.Unmarshal(raw, &res)
	if err != nil {
		return res, err
	}

	return res, nil
}
