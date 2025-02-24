package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/tuanpq1998/my-first-go/internal/database"

	"github.com/jackc/pgx/v5"
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

	apiCfg := apiConfig{
		DB: database.New(conn),
	}

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", handlerReadiness)
	v1Router.Get("/error", handlerErr)

	v1Router.Post("/user", apiCfg.handlerCreateUser)
	v1Router.Get("/user", apiCfg.middlewareAuth(apiCfg.handlerGetUserByKey))

	v1Router.Post("/feed", apiCfg.middlewareAuth(apiCfg.handlerCreateFeed))
	v1Router.Get("/feed", (apiCfg.handlerGetFeeds))

	v1Router.Post("/feedFollow", apiCfg.middlewareAuth(apiCfg.handlerCreateFeedFollow))
	v1Router.Get("/feedFollow", apiCfg.middlewareAuth(apiCfg.handlerGetFeedFollows))
	v1Router.Delete("/feedFollow/{feedFollowId}", apiCfg.middlewareAuth(apiCfg.handlerDeleteFeedFollow))

	router.Mount("/v1", v1Router)

	server := &http.Server{
		Handler: router,
		Addr:    ":" + portStr,
	}

	log.Println("Server started on", server.Addr)

	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("Server starts failed:%v", err)
	}
}
