package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v4/stdlib"
)

// app struct contains global state.
type app struct {
	// db is the global database connection pool.
	db *sql.DB
}

func getRecentEvents(w http.ResponseWriter, r *http.Request, app *app) error {
	return nil
}

// eventListHandler handles requests to the /events route.
func (app *app) eventListHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if err := getRecentEvents(w, r, app); err != nil {
			log.Printf("getRecentEvents: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	default:
		http.Error(w, fmt.Sprintf("HTTP Method %s Not Allowed", r.Method), http.StatusMethodNotAllowed)
	}
}

func main() {
	app := newApp()

	http.HandleFunc("/events", app.eventListHandler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func newApp() *app {
	var err error
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

	// Create the votes table if it does not already exist.
	if _, err = app.db.Exec(`CREATE TABLE IF NOT EXISTS votes
	( id SERIAL NOT NULL, created_at timestamp NOT NULL, updated_at timestamp NOT NULL,
	candidate VARCHAR(6) NOT NULL, PRIMARY KEY (id) );`); err != nil {
		log.Fatalf("DB.Exec: unable to create table: %s", err)
	}

	return app
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
func initSocketConnectionPool() (*sql.DB, error) {
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
	dbPool, err := sql.Open("pgx", dbURI)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %v", err)
	}

	// [START_EXCLUDE]
	configureConnectionPool(dbPool)
	// [END_EXCLUDE]

	return dbPool, nil
	// [END cloud_sql_postgres_databasesql_create_socket]
}

// initTCPConnectionPool initializes a TCP connection pool for a Cloud SQL
// instance of SQL Server.
func initTCPConnectionPool() (*sql.DB, error) {
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
	dbPool, err := sql.Open("pgx", dbURI)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %v", err)
	}

	// [START_EXCLUDE]
	configureConnectionPool(dbPool)
	// [END_EXCLUDE]

	return dbPool, nil
	// [END cloud_sql_postgres_databasesql_create_tcp]
}

// configureConnectionPool sets database connection pool properties.
// For more information, see https://golang.org/pkg/database/sql
func configureConnectionPool(dbPool *sql.DB) {
	// [START cloud_sql_postgres_databasesql_limit]

	// Set maximum number of connections in idle connection pool.
	dbPool.SetMaxIdleConns(5)

	// Set maximum number of open connections to the database.
	dbPool.SetMaxOpenConns(7)

	// [END cloud_sql_postgres_databasesql_limit]

	// [START cloud_sql_postgres_databasesql_lifetime]

	// Set Maximum time (in seconds) that a connection can remain open.
	dbPool.SetConnMaxLifetime(1800)

	// [END cloud_sql_postgres_databasesql_lifetime]
}
