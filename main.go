package main

import (
	"exchange-api/client"
	"exchange-api/config"
	"exchange-api/controller"
	"exchange-api/service"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	config := config.NewConfigurations()
	client := client.NewExchangeRateClient()
	service := service.NewExchangeRateService(client, &config.ExchangeRateService)
	controller := controller.NewExchangeRateController(service, &config.ExchangeRateController)
	sm := mux.NewRouter()

	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/currency/{name}", controller.GetExchangeRate)

	ch := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
		handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"}),
	)

	s := http.Server{
		Addr:         ":" + config.Server.Port,              // configure the bind address
		Handler:      ch(sm),                                // set the default handler
		ErrorLog:     log.New(os.Stderr, "", log.LstdFlags), // set the logger for the client
		ReadTimeout:  5 * time.Second,                       // max time to read request from the client
		WriteTimeout: 10 * time.Second,                      // max time to write response to the client
		IdleTimeout:  120 * time.Second,                     // max time for connections using TCP Keep-Alive
	}

	go func() {
		log.Println("Starting server on port: ", s.Addr)

		err := s.ListenAndServe()
		if err != nil {
			log.Printf("Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	sig := <-c
	log.Println("Got signal:", sig)
}
