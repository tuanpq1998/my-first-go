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
}

func transformToUserDto(dbUser database.User) UserDto {
	return UserDto{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Name:      dbUser.Name,
	}
}
