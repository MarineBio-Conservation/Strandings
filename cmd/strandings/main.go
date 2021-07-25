package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/MarineBio-Conservation/Strandings-Backend/pkg/connection"
	"github.com/MarineBio-Conservation/Strandings-Backend/pkg/events"
	"github.com/jackc/pgx/v4"
)

// app struct contains global state.
type app struct {
	// db is the global database connection pool.
	db *pgx.Conn
}

var thisApp *app

func main() {

	thisApp = &app{}

	var err error
	thisApp.db, err = connection.InitTCPConnectionPool()
	if err != nil {
		log.Fatalf("initTCPConnectionPool: unable to connect: %s", err)
	}
	defer thisApp.db.Close(context.Background())

	http.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
		events.GetEventsHandler(w, r, thisApp.db)
	})

	http.HandleFunc("/event", func(w http.ResponseWriter, r *http.Request) {
		events.GetEvent(w, r, thisApp.db)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}

}
