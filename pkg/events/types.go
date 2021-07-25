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
