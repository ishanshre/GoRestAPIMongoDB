package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/ishanshre/GoRestAPIMongoDB/internals/handlers"
	"github.com/ishanshre/GoRestAPIMongoDB/internals/routers"
)

func main() {
	port := flag.Int("port", 8000, "Port that server listens")
	addr := fmt.Sprintf(":%d", *port)
	h := handlers.NewHandler()
	router := routers.Router(h)
	log.Printf("Starting Server at port : %d", *port)
	srv := http.Server{
		Addr:    addr,
		Handler: router,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("error in listening to server: %s", err.Error())
	}
}
