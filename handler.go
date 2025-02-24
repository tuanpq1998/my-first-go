package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/tuanpq1998/my-first-go/internal/database"
)

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, 200, struct{}{})
}

func handlerErr(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, 400, "something went wrong!")
}

func (apiCfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Println("handlerCreateUser::JSONDecode::error", err)
		respondWithError(w, 400, "decode JSON failed")
		return
	}

	newUser, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID: pgtype.UUID{
			Bytes: uuid.New(),
			Valid: true,
		},
		Name: params.Name,
		CreatedAt: pgtype.Timestamp{
			Time:  time.Now().UTC(),
			Valid: true,
		},
		UpdatedAt: pgtype.Timestamp{
			Time:  time.Now().UTC(),
			Valid: true,
		},
	})

	if err != nil {
		log.Println("handlerCreateUser::CreateUser::error", err)
		respondWithError(w, 400, "couldnt create user")
		return
	}
	respondWithJSON(w, 201, transformToUserDto(newUser))

}

func (apiCfg apiConfig) handlerGetUserByKey(w http.ResponseWriter, r *http.Request, user database.User) {
	respondWithJSON(w, 200, transformToUserDto(user))
}

func (apiCfg *apiConfig) handlerCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Title string `json:"title"`
		Url   string `json:"url"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Println("handlerCreateFeed::JSONDecode", err)
		respondWithError(w, 400, "decode JSON failed")
		return
	}

	newFeed, err := apiCfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
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
		Title:  params.Title,
		Url:    params.Url,
		UserID: user.ID,
	})
	if err != nil {
		log.Println("handlerCreateFeed::CreateFeed::error", err)
		respondWithError(w, 400, "couldnt create feed")
		return
	}
	respondWithJSON(w, 201, transformToFeedDto(newFeed))
}

func (apiCfg apiConfig) handlerGetFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := apiCfg.DB.GetAllFeeds(r.Context())
	if err != nil {
		log.Println("handlerGetFeeds::GetAllFeeds::error", err)
		respondWithError(w, 400, "couldnt get feeds")
		return
	}
	respondWithJSON(w, 201, transformArrToFeedDto(feeds))
}

func (apiCfg *apiConfig) handlerCreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		// FeedIdStr string `json:"feed_id"`
		FeedId uuid.UUID `json:"feed_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Println("handlerCreateFeedFollow::JSONDecode", err)
		respondWithError(w, 400, "decode JSON failed")
		return
	}

	//parse uuid
	// feedId, err := uuid.Parse(params.FeedIdStr)
	// if err != nil {
	// 	log.Println("handlerCreateFeedFollow::UUIDParse", err)
	// 	respondWithError(w, 400, "parse feed_id failed")
	// 	return
	// }

	// TODO: did this feed_id exist AND did feed follow with this feed_id and user_id has existed?

	feedFollow, err := apiCfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
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
		UserID: user.ID,
		// FeedID: pgtype.UUID{
		// 	Bytes: feedId,
		// 	Valid: true,
		// },
		FeedID: pgtype.UUID{
			Bytes: params.FeedId,
			Valid: true,
		},
	})
	if err != nil {
		log.Println("handlerCreateFeedFollow::CreateFeedFollow::error", err)
		respondWithError(w, 400, "couldnt create feed follow")
		return
	}
	respondWithJSON(w, 201, transformToFeedFollowDto(feedFollow))
}
