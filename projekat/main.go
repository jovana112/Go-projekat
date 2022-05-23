package main

import (
	"context"
	"github.com/gorilla/mux"
	cs "github.com/jovana112/Go-Projekat/projekat/configstore"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	router := mux.NewRouter()
	router.StrictSlash(true)

	store, err := cs.CreateNewConfigStore()
	if err != nil {
		log.Fatal(err)
	}

	server := Service{
		store: store,
	}

	router.HandleFunc("/configs", server.createConfigHandler).Methods("POST")
	router.HandleFunc("/configs", server.getAllConfigHandler).Methods("GET")
	router.HandleFunc("/configs/{id}", server.updateConfigWithNewVersionHandler).Methods("PUT")
	router.HandleFunc("/configs/{id}/{version}", server.getConfigHandler).Methods("GET")
	router.HandleFunc("/configs/{id}/{version}", server.deleteConfigHandler).Methods("DELETE")

	router.HandleFunc("/groups", server.createGroupHandler).Methods("POST")
	router.HandleFunc("/groups", server.getAllGroupHandler).Methods("GET")
	router.HandleFunc("/groups/{id}", server.updateGroupWithNewVersionHandler).Methods("PUT")
	router.HandleFunc("/groups/{id}/{version}", server.getGroupHandler).Methods("GET")
	router.HandleFunc("/groups/{id}/{version}", server.deleteGroupHandler).Methods("DELETE")
	router.HandleFunc("/groups/{id}/{version}", server.extendConfigGroupHandler).Methods("PATCH")
	router.HandleFunc("/groups/{id}/{version}/configs", server.getConfigsByLabelsHandler).Methods("GET")

	// start server
	srv := &http.Server{Addr: "0.0.0.0:8000", Handler: router}
	go func() {
		log.Println("server starting")
		if err := srv.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Fatal(err)
			}
		}
	}()

	<-quit

	log.Println("service shutting down ...")

	// gracefully stop server
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
	log.Println("server stopped")
}
