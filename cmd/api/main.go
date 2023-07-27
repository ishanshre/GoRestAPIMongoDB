package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ishanshre/GoRestAPIMongoDB/internals/database"
	"github.com/ishanshre/GoRestAPIMongoDB/internals/handlers"
	"github.com/ishanshre/GoRestAPIMongoDB/internals/routers"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func main() {
	// get default port if no port is assiged by user
	port := flag.Int("port", 8000, "Port that server listens")
	addr := fmt.Sprintf(":%d", *port)

	// load the environment files
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("error in loading environment files: %v", err)
	}

	//connect to mongo db
	client, err := database.ConnectToNoSql(os.Getenv("dsn"))
	if err != nil {
		log.Fatalf("error in connecting to mongo db: %v", err)
	}

	redisPool := redis.NewClient(
		&redis.Options{
			Addr:         os.Getenv("REDIS_URL"),
			Password:     "",
			DB:           0,
			MaxIdleConns: 10,
			PoolSize:     10,
			MinIdleConns: 0,
		},
	)

	if err := redisPool.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("error in connecting to redis: %s", err.Error())
	}
	// connect to handler interface
	h := handlers.NewHandler(client, redisPool)

	// connect to router
	router := routers.Router(h)

	log.Printf("Starting Server at port : %d", *port)

	// configure the serevr
	srv := http.Server{
		Addr:    addr,
		Handler: router,
	}

	// start the server
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("error in listening to server: %s", err.Error())
	}
}
