package events

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/MarineBio-Conservation/Strandings-Backend/pkg/cors"
	"github.com/jackc/pgx/v4"
)

func GetEvent(w http.ResponseWriter, r *http.Request, db *pgx.Conn) {
	cors.Cors(w, r)

	queryParams := r.URL.Query()
	ids, present := queryParams["id"]
	if !present {
		http.Error(w, "id not supplied", 400)
		return
	}
	id, err := strconv.Atoi(ids[0])
	if err != nil {
		http.Error(w, "id has invalid format", 400)
		return
	}
	var record Record
	err = db.QueryRow(context.Background(),
		`select event_id, event_date, event_regions, event_animal_type, event_animal_number_died, investigation_type, stranding_causes, event_centroid_lat, event_centroid_long 
			from public.data
			WHERE 
				event_id = $1;`, id).Scan(&record)
	if err != nil {
		http.Error(w, strconv.Itoa(id)+" not found", 404)
	}

	data, _ := json.Marshal(record)
	w.Write(data)
}
