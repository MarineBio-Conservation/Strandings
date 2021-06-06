package main

import (
	"context"
	"fmt"
	"log"
	"os"

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

	// If the optional DB_HOST environment variable is set, it contains
	// the IP address and port number of a TCP connection pool to be created,
	// such as "127.0.0.1:5432". If DB_HOST is not set, a Unix socket
	// connection pool will be created instead.
	if os.Getenv("DB_HOST") != "" {
		app.db, err = initTCPConnectionPool()
		if err != nil {
			log.Fatalf("initTCPConnectionPool: unable to connect: %s", err)
		}
	} else {
		app.db, err = initSocketConnectionPool()
		if err != nil {
			log.Fatalf("initSocketConnectionPool: unable to connect: %s", err)
		}
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

// mustGetEnv is a helper function for getting environment variables.
// Displays a warning if the environment variable is not set.
func mustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("Warning: %s environment variable not set.\n", k)
	}
	return v
}

// initSocketConnectionPool initializes a Unix socket connection pool for
// a Cloud SQL instance of SQL Server.
func initSocketConnectionPool() (*pgx.Conn, error) {
	// [START cloud_sql_postgres_databasesql_create_socket]
	var (
		dbUser                 = mustGetenv("DB_USER")                  // e.g. 'my-db-user'
		dbPwd                  = mustGetenv("DB_PASS")                  // e.g. 'my-db-password'
		instanceConnectionName = mustGetenv("INSTANCE_CONNECTION_NAME") // e.g. 'project:region:instance'
		dbName                 = mustGetenv("DB_NAME")                  // e.g. 'my-database'
	)

	socketDir, isSet := os.LookupEnv("DB_SOCKET_DIR")
	if !isSet {
		socketDir = "/cloudsql"
	}

	var dbURI = fmt.Sprintf("user=%s password=%s database=%s host=%s/%s", dbUser, dbPwd, dbName, socketDir, instanceConnectionName)

	// dbPool is the pool of database connections.
	dbPool, err := pgx.Connect(context.Background(), dbURI)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %v", err)
	}

	return dbPool, nil
	// [END cloud_sql_postgres_databasesql_create_socket]
}

// initTCPConnectionPool initializes a TCP connection pool for a Cloud SQL
// instance of SQL Server.
func initTCPConnectionPool() (*pgx.Conn, error) {
	// [START cloud_sql_postgres_databasesql_create_tcp]
	var (
		dbUser    = mustGetenv("DB_USER") // e.g. 'my-db-user'
		dbPwd     = mustGetenv("DB_PASS") // e.g. 'my-db-password'
		dbTCPHost = mustGetenv("DB_HOST") // e.g. '127.0.0.1' ('172.17.0.1' if deployed to GAE Flex)
		dbPort    = mustGetenv("DB_PORT") // e.g. '5432'
		dbName    = mustGetenv("DB_NAME") // e.g. 'my-database'
	)

	var dbURI = fmt.Sprintf("host=%s user=%s password=%s port=%s database=%s", dbTCPHost, dbUser, dbPwd, dbPort, dbName)

	// dbPool is the pool of database connections.
	dbPool, err := pgx.Connect(context.Background(), dbURI)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %v", err)
	}

	// [START_EXCLUDE]
	// [END_EXCLUDE]

	return dbPool, nil
	// [END cloud_sql_postgres_databasesql_create_tcp]
}
