package main

import (
	"log"
	"net/http"

	"github.com/tuanpq1998/my-first-go/internal/auth"
	"github.com/tuanpq1998/my-first-go/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (apiCfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		handler(w, r, user)
	}
}
