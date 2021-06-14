package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/MarineBio-Conservation/Strandings-Backend/pkg/connection"
	"github.com/MarineBio-Conservation/Strandings-Backend/pkg/cors"
	"github.com/jackc/pgx/v4"
)

// app struct contains global state.
type app struct {
	// db is the global database connection pool.
	db *pgx.Conn
}

type position struct {
	Lat float32 `json:"lat"`
	Lng float32 `json:"lng"`
}

type record struct {
	Id   int       `json:"id"`
	Date time.Time `json:"date"`
	Pos  position  `json:"position"`
}

var thisApp *app

func eventsHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/events" {
		http.NotFound(w, r)
		return
	}

	cors.Cors(w, r)

	rows, err := thisApp.db.Query(context.Background(),
		"select event_id, event_date, event_centroid_lat, event_centroid_long from public.data ORDER BY event_date DESC LIMIT 1000;")
	if err != nil {
		log.Fatalf("conn.Query failed: %v", err)
	}
	defer rows.Close()
	var results []record
	for rows.Next() {
		var rec record
		err = rows.Scan(&rec.Id, &rec.Date, &rec.Pos.Lat, &rec.Pos.Lng)
		if err != nil {
			log.Fatalf("conn.Query failed: %v", err)
		}
		results = append(results, rec)
	}

	data, _ := json.Marshal(results)
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
