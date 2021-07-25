package events

import "time"

type position struct {
	Lat float32 `json:"lat"`
	Lng float32 `json:"lng"`
}

type Record struct {
	Id            int       `json:"id"`
	Date          time.Time `json:"date"`
	Regions       *[]string `json:"regions"`
	AnimalType    *[]string `json:"animal_type"`
	Died          int       `json:"died"`
	Investigation *string   `json:"investigation_type"`
	Causes        *[]string `json:"causes"`
	Pos           position  `json:"position"`
}

// Database schema:
// event_id INT GENERATED ALWAYS AS IDENTITY,
// 		event_date date NOT NULL,
// 		event_location_lat real[],
// 		event_location_long real[],
// 		event_centroid_lat real NOT NULL,
// 		event_centroid_long real NOT NULL,
// 		event_regions varchar(63)[],
// 		event_animal_type varchar(63)[],
// 		event_animal_number integer NOT NULL,
// 		event_animal_number_died integer NOT NULL,
// 		investigation_type varchar(31),
// 		investigation_description text,
// 		stranding_causes varchar(63)[],
// 		investigation_results_description text,
// 		investigation_references text,
// 		researcher_id integer NOT NULL,
type FullRecord struct {
	Id                              int         `json:"id"`
	Date                            time.Time   `json:"date"`
	Pos                             position    `json:"position"`
	Bound                           *[]position `json:"bounds"`
	Regions                         *[]string   `json:"regions"`
	AnimalType                      *[]string   `json:"animal_type"`
	Number                          int         `json:"number"`
	Died                            int         `json:"died"`
	Investigation                   *string     `json:"investigation_type"`
	InvestigationDescription        *string     `json:"investigation_description"`
	References                      *string     `json:"references"`
	Causes                          *[]string   `json:"causes"`
	InvestigationResultsDescription *string     `json:"investigation_results_description"`
	ResearcherName                  *string     `json:"researcher_name"`
}
