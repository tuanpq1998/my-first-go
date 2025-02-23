package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/tuanpq1998/my-first-go/internal/auth"
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

func (apiCfg apiConfig) handlerGetUserByKey(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.ParseApiKey(r.Header)
	if err != nil {
		log.Println("handlerGetUserByKey::ParseApiKey::error", err, "::apiKey::", apiKey)
		respondWithError(w, 403, "authentication failed")
		return
	}

	user, err := apiCfg.DB.GetUserByApiKey(r.Context(), apiKey)
	if err != nil {
		log.Println("handlerGetUserByKey::GetUserByApiKey::error", err, "::apiKey::", apiKey)
		respondWithError(w, 403, "authentication failed")
		return
	}

	respondWithJSON(w, 200, transformToUserDto(user))
}
