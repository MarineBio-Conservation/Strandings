package events

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MarineBio-Conservation/Strandings-Backend/pkg/cors"
	"github.com/jackc/pgx/v4"
)

func PostEvent(w http.ResponseWriter, r *http.Request, db *pgx.Conn) {
	cors.Cors(w, r)

	var newEvent FullRecord
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&newEvent)
	if err != nil {
		http.Error(w, "Unable to decode request", 400)
		return
	}
	fmt.Println("Received a post:")
	fmt.Println(newEvent)
}
