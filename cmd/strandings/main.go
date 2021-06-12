package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/MarineBio-Conservation/Strandings-Backend/pkg/connection"
	"github.com/MarineBio-Conservation/Strandings-Backend/pkg/cors"
	"github.com/jackc/pgx/v4"
)

// app struct contains global state.
type app struct {
	// db is the global database connection pool.
	db *pgx.Conn
}

var thisApp *app

func eventsHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/events" {
		http.NotFound(w, r)
		return
	}

	cors.Cors(w, r)

	rows, err := thisApp.db.Query(context.Background(), "select researcher_name from researchers where researcher_id = 1")
	if err != nil {
		log.Fatalf("conn.Query failed: %v", err)
	}
	defer rows.Close()
	var names []string
	for rows.Next() {
		var name string
		rows.Scan(&name)
		names = append(names, name)
	}

	data, _ := json.Marshal(names)
	w.Write(data)
}

func main() {

	thisApp = &app{}

	var err error
	thisApp.db, err = connection.InitTCPConnectionPool()
	if err != nil {
		log.Fatalf("initTCPConnectionPool: unable to connect: %s", err)
	}
	defer thisApp.db.Close(context.Background())

	http.HandleFunc("/events", eventsHandler)

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
