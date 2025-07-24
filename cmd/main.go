package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"plexcache/api"
	red "plexcache/redis"

	"github.com/LukeHagar/plexgo"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func main() {
	log.Println("Starting")

	ctx := context.Background()
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: "",
		DB:       0,
	})

	_, err = rdb.Ping(ctx).Result()

	if err != nil {
		log.Fatalf("Error connecing to redis %e", err)
		return
	}

	subscriber := red.SubscribeToExpired(rdb)
	defer subscriber.Close()

	plexApi := plexgo.New(
		plexgo.WithSecurity(os.Getenv("PLEX_API_KEY")),
		plexgo.WithIP(os.Getenv("PLEX_IP")),
		plexgo.WithProtocol("http"),
	)

	r := mux.NewRouter()
	r.HandleFunc("/", api.WebhookHandler(rdb, plexApi)).Methods("POST")

	if err := http.ListenAndServe(":4001", r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}

}
