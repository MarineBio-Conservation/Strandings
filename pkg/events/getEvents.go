package events

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/MarineBio-Conservation/Strandings-Backend/pkg/cors"
	"github.com/jackc/pgx/v4"
)

func GetEventsHandler(w http.ResponseWriter, r *http.Request, db *pgx.Conn) {
	cors.Cors(w, r)

	queryParams := r.URL.Query()
	var err error

	var latMin float64
	latMinString, present := queryParams["latMin"]
	if present {
		latMin, err = strconv.ParseFloat(latMinString[0], 32)
		if err != nil {
			latMin = -90.0
		}
	} else {
		latMin = -90
	}

	var latMax float64
	latMaxString, present := queryParams["latMax"]
	if present {
		latMax, err = strconv.ParseFloat(latMaxString[0], 32)
		if err != nil {
			latMax = 90.0
		}
	} else {
		latMax = 90
	}

	var lngMin float64
	lngMinString, present := queryParams["lngMin"]
	if present {
		lngMin, err = strconv.ParseFloat(lngMinString[0], 32)
		if err != nil {
			lngMin = -180.0
		}
	} else {
		lngMin = -180
	}

	var lngMax float64
	lngMaxString, present := queryParams["lngMax"]
	if present {
		lngMax, err = strconv.ParseFloat(lngMaxString[0], 32)
		if err != nil {
			lngMax = 180.0
		}
	} else {
		lngMax = 180
	}

	var limit int64
	limitString, present := queryParams["limit"]
	if present {
		limit, err = strconv.ParseInt(limitString[0], 10, 32)
		if err != nil {
			limit = 500
		}
	} else {
		limit = 500
	}

	rows, err := db.Query(context.Background(),
		`select event_id, event_date, event_regions, event_animal_type, event_animal_number_died, investigation_type, stranding_causes, event_centroid_lat, event_centroid_long 
			from public.data
			WHERE 
				event_centroid_lat > $1 AND 
				event_centroid_lat < $2 AND
				event_centroid_long > $3 AND 
				event_centroid_long < $4
			ORDER BY event_date DESC
			LIMIT $5;`, latMin, latMax, lngMin, lngMax, limit)
	if err != nil {
		log.Fatalf("conn.Query failed: %v", err)
	}
	defer rows.Close()
	var results []Record
	for rows.Next() {
		var rec Record
		err = rows.Scan(&rec.Id, &rec.Date, &rec.Regions, &rec.AnimalType, &rec.Died, &rec.Investigation, &rec.Causes, &rec.Pos.Lat, &rec.Pos.Lng)
		if err != nil {
			log.Fatalf("conn.Query failed: %v", err)
		}
		results = append(results, rec)
	}

	data, _ := json.Marshal(results)
	w.Write(data)
}
