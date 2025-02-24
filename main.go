package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/tuanpq1998/my-first-go/internal/database"
	"github.com/tuanpq1998/my-first-go/middlewares"

	"github.com/jackc/pgx/v5/pgxpool"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	// feed, err := urlToFeed("https://hnrss.org/frontpage")
	// if err != nil {
	// 	return
	// }
	// log.Println(feed)

	godotenv.Load()

	portStr := os.Getenv("PORT")
	if portStr == "" {
		log.Fatal("PORT isnt found in the environment")
		// log.Println("PORT isnt found in the environment")
		// return
	}

	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatal("DB_URL isnt found in the env")
	}

	// conn, err := sql.Open("pgx", dbUrl)
	// if err != nil {
	// 	log.Fatal("Cant connect to database")
	// }

	// apiCfg := apiConfig{
	// 	DB: database.New(tx),
	// }

	conn, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		log.Fatalf("Cant connect to database: %v", err)
	}

	err = conn.Ping(context.Background())
	if err == nil {
		log.Println("Connected to database")
	} else {
		log.Fatalf("Couldnt ping to database: %v", err)
	}

	db := database.New(conn)
	apiCfg := apiConfig{
		DB: db,
	}

	router := http.NewServeMux()

	v1Router := http.NewServeMux()
	v1Router.HandleFunc("GET /healthz", handlerReadiness)
	v1Router.HandleFunc("GET /error", handlerErr)

	v1Router.HandleFunc("POST /user", apiCfg.handlerCreateUser)
	v1Router.HandleFunc("GET /user", apiCfg.middlewareAuth(apiCfg.handlerGetUserByKey))

	v1Router.HandleFunc("POST /feed", apiCfg.middlewareAuth(apiCfg.handlerCreateFeed))
	v1Router.HandleFunc("GET /feed", (apiCfg.handlerGetFeeds))

	v1Router.HandleFunc("POST /feedFollow", apiCfg.middlewareAuth(apiCfg.handlerCreateFeedFollow))
	v1Router.HandleFunc("GET /feedFollow", apiCfg.middlewareAuth(apiCfg.handlerGetFeedFollows))
	v1Router.HandleFunc("DELETE /feedFollow/{feedFollowId}", apiCfg.middlewareAuth(apiCfg.handlerDeleteFeedFollow))

	v1Router.HandleFunc("GET /posts", apiCfg.middlewareAuth(apiCfg.handlerGetPostsForUser))

	router.Handle("/v1/", http.StripPrefix("/v1", v1Router))

	server := &http.Server{
		Handler: middlewares.AllowCors(middlewares.HttpLogging(router)),
		Addr:    ":" + portStr,
	}

	log.Println("Server started on", server.Addr)
	go startScraping(db, 10, time.Minute)

	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("Server starts failed:%v", err)
	}
}
