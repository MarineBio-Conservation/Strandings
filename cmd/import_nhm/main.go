package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/MarineBio-Conservation/Strandings-Backend/pkg/connection"
	"github.com/jackc/pgx/v4"
)

// app struct contains global state.
type app struct {
	// db is the global database connection pool.
	db *pgx.Conn
}

func main() {

	f, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer f.Close()

	app := &app{}

	app.db, err = connection.InitTCPConnectionPool()
	if err != nil {
		log.Fatalf("initTCPConnectionPool: unable to connect: %s", err)
	}
	defer app.db.Close(context.Background())

	fmt.Println("Creating researchers table")
	// Create the data table if it does not already exist.
	if _, err = app.db.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS researchers
	(
		researcher_id INT GENERATED ALWAYS AS IDENTITY,
		researcher_name VARCHAR(255) UNIQUE NOT NULL,
		researcher_email VARCHAR(127),
		researcher_first_event_date date NOT NULL DEFAULT CURRENT_DATE,
		researcher_latest_event_date date NOT NULL DEFAULT CURRENT_DATE,
		researcher_events_total integer NOT NULL DEFAULT 0,
		PRIMARY KEY(researcher_id)
	 );`); err != nil {
		log.Fatalf("DB.Exec: unable to create table: %s", err)
	}

	fmt.Println("Creating data table")
	if _, err = app.db.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS data
	( 	event_id INT GENERATED ALWAYS AS IDENTITY,
		event_date date NOT NULL,
		event_location_lat real[],
		event_location_long real[],
		event_centroid_lat real NOT NULL,
		event_centroid_long real NOT NULL,
		event_regions varchar(63)[],
		event_animal_type varchar(63)[],
		event_animal_number integer NOT NULL,
		event_animal_number_died integer NOT NULL,
		investigation_type varchar(31),
		investigation_description text,
		stranding_causes varchar(63)[],
		investigation_results_description text,
		investigation_references text,
		researcher_id integer NOT NULL,
		PRIMARY KEY (event_id)	,
		CONSTRAINT fk_researcher
      		FOREIGN KEY(researcher_id) 
	  			REFERENCES researchers(researcher_id)
	);`); err != nil {
		log.Fatalf("DB.Exec: unable to create table: %s", err)
	}

	if _, err = app.db.Exec(context.Background(), `
	CREATE INDEX lat_idx on public.data
	(
		event_centroid_lat ASC
	);
	CREATE INDEX lng_idx on public.data
	(
		event_centroid_long ASC
	);
	CREATE INDEX date_idx on public.data
	(
		event_date DESC
	);
	CREATE INDEX researcher_name_idx on public.researchers
	(
		researcher_name ASC
	);
	`); err != nil {
		log.Fatalf("DB.Exec: unable to create indices: %s", err)
	}

	if _, err = app.db.Exec(context.Background(), `
		INSERT INTO researchers(researcher_name, researcher_email, researcher_first_event_date, researcher_latest_event_date)
		VALUES('Natural History Museum', NULL, '1913-02-13', '1989-12-25' );
	`); err != nil {
		log.Fatalf("DB.Exec: Add researcher: %s", err)
	}

	res, err := app.db.PgConn().CopyFrom(context.Background(), f, "COPY data(event_date,event_centroid_lat,event_centroid_long,event_animal_type,event_animal_number,event_animal_number_died,investigation_description,researcher_id) FROM STDIN DELIMITER ',' CSV HEADER")
	if err != nil {
		panic(err)
	}
	fmt.Print(res.RowsAffected())
}
