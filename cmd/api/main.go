package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ishanshre/GoRestAPIMongoDB/internals/database"
	"github.com/ishanshre/GoRestAPIMongoDB/internals/handlers"
	"github.com/ishanshre/GoRestAPIMongoDB/internals/routers"
	"github.com/joho/godotenv"
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

	// connect to handler interface
	h := handlers.NewHandler(client)

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
