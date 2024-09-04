package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/evgeny-s/go-rss/internal/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT is not found in the ENV params")
	}

	pgPort := os.Getenv("PG_PORT")
	if pgPort == "" {
		log.Fatal("PG_PORT is not found in the ENV params")
	}

	pgHost := os.Getenv("PG_HOST")
	if pgHost == "" {
		log.Fatal("PG_HOST is not found in the ENV params")
	}

	pgUsername := os.Getenv("PG_USERNAME")
	if pgUsername == "" {
		log.Fatal("PG_USERNAME is not found in the ENV params")
	}

	pgPassword := os.Getenv("PG_PASSWORD")
	if pgPassword == "" {
		log.Fatal("PG_PASSWORD is not found in the ENV params")
	}

	pgDb := os.Getenv("PG_DB")
	if pgDb == "" {
		log.Fatal("PG_DB is not found in the ENV params")
	}

	conn, err := sql.Open("postgres", "postgres://"+pgUsername+":"+pgPassword+"@"+pgHost+":"+pgPort+"/"+pgDb+"?sslmode=disable")
	if err != nil {
		log.Fatal("Can't connect to database.")
	}

	db := database.New(conn)
	apiCfg := apiConfig{
		DB: db,
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

	go startScraping(db, 10, time.Minute)

	v1Router := chi.NewRouter()

	v1Router.Get("/healthz", handlerReadiness)
	v1Router.Get("/err", handlerErr)
	v1Router.Post("/users", apiCfg.handlerCreateUser)
	v1Router.Get("/users", apiCfg.middlewareAuth(apiCfg.handlerGetUser))
	v1Router.Post("/feeds", apiCfg.middlewareAuth(apiCfg.handlerCreateFeed))
	v1Router.Get("/feeds", apiCfg.handlerGetGetFeeds)
	v1Router.Post("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerCreateFeedFollow))
	v1Router.Get("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerGetFeedFollows))
	v1Router.Delete("/feed_follows/{feedFollowID}", apiCfg.middlewareAuth(apiCfg.handlerDeleteFeedFollows))
	v1Router.Get("/posts", apiCfg.middlewareAuth(apiCfg.handlerGetPostsForUser))

	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	fmt.Println("Server starting on port: ", portString)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Port: ", portString)
}
