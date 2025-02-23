package main

import (
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/tuanpq1998/my-first-go/internal/database"
)

type UserDto struct {
	ID        pgtype.UUID      `json:"id"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
	UpdatedAt pgtype.Timestamp `json:"updated_at"`
	Name      string           `json:"name"`
	ApiKey    string           `json:"key"`
}

type FeedDto struct {
	ID        pgtype.UUID      `json:"id"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
	UpdatedAt pgtype.Timestamp `json:"updated_at"`
	Title     string           `json:"title"`
	Url       string           `json:"url"`
	UserId    pgtype.UUID      `json:"user_id"`
}

func transformToUserDto(dbUser database.User) UserDto {
	return UserDto{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Name:      dbUser.Name,
		ApiKey:    dbUser.ApiKey,
	}
}

func transformToFeedDto(db database.Feed) FeedDto {
	return FeedDto{
		ID:        db.ID,
		CreatedAt: db.CreatedAt,
		UpdatedAt: db.UpdatedAt,
		Title:     db.Title,
		Url:       db.Url,
		UserId:    db.UserID,
	}
}

func transformArrToFeedDto(db []database.Feed) []FeedDto {
	feeds := []FeedDto{}
	for _, feed := range db {
		feeds = append(feeds, transformToFeedDto(feed))
	}
	return feeds
}
