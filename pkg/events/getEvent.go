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
	var record FullRecord
	var bounds_lats []float32
	var bounds_lngs []float32
	err = db.QueryRow(context.Background(),
		`select data.event_id, data.event_date,
		data.event_location_lat, data.event_location_long,
		data.event_centroid_lat, data.event_centroid_long, 
		data.event_regions, data.event_animal_type,
		data.event_animal_number, data.event_animal_number_died,
		data.investigation_type, data.investigation_description,
		data.stranding_causes, data.investigation_results_description, 
		data.investigation_references, researchers.researcher_name
			from data
			INNER JOIN researchers ON data.researcher_id=researchers.researcher_id
			WHERE 
				data.event_id = $1;`, id).Scan(
		&record.Id, &record.Date,
		&bounds_lats, &bounds_lngs,
		&record.Pos.Lat, &record.Pos.Lng,
		&record.Regions, &record.AnimalType,
		&record.Number, &record.Died,
		&record.Investigation, &record.InvestigationDescription,
		&record.Causes, &record.InvestigationResultsDescription,
		&record.References, &record.ResearcherName)
	if err != nil {
		http.Error(w, strconv.Itoa(id)+" not found", 404)
		return
	}

	for i := range bounds_lats {
		*record.Bound = append(*record.Bound, position{Lat: bounds_lats[i], Lng: bounds_lngs[i]})
	}

	data, _ := json.Marshal(record)
	w.Write(data)
}
