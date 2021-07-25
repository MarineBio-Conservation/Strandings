package events

import (
	"net/http"

	"github.com/MarineBio-Conservation/Strandings-Backend/pkg/cors"
	"github.com/jackc/pgx/v4"
)

func HandleEvent(w http.ResponseWriter, r *http.Request, db *pgx.Conn) {
	cors.Cors(w, r)

	if r.Method == http.MethodPost {
		PostEvent(w, r, db)
	} else if r.Method == http.MethodGet {
		GetEvent(w, r, db)
	}

}
