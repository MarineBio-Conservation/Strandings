package connection

import (
	"context"
	"fmt"
	"log"
	"os"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"

	"github.com/jackc/pgx/v4"
)

// mustGetEnv is a helper function for getting environment variables.
// Displays a warning if the environment variable is not set.
func mustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("Warning: %s environment variable not set.\n", k)
	}
	return v
}

func pullDbPassword() string {
	// name := "projects/my-project/secrets/my-secret/versions/5"
	name := "projects/stranding-project/secrets/DB_PASS/versions/latest"
	//name := "projects/676026255512/secrets/DB_PASS/versions/latest"
	fmt.Println("Pulling")
	// Create the client.
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		fmt.Println("ERR 1")
		return fmt.Errorf("failed to create secretmanager client: %v", err).Error()
	}
	defer client.Close()

	// Build the request.
	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: name,
	}

	// Call the API.
	result, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		fmt.Println("Err 2", err)
		return fmt.Errorf("failed to access secret version: %v", err).Error()
	}

	return string(result.Payload.Data)
}

// initTCPConnectionPool initializes a TCP connection pool for a Cloud SQL
// instance of SQL Server.
func InitTCPConnectionPool() (*pgx.Conn, error) {
	// [START cloud_sql_postgres_databasesql_create_tcp]
	var (
		dbUser    = mustGetenv("DB_USER") // e.g. 'my-db-user'
		dbPwd     = os.Getenv("DB_PASS")  // e.g. 'my-db-password'
		dbTCPHost = mustGetenv("DB_HOST") // e.g. '127.0.0.1' ('172.17.0.1' if deployed to GAE Flex)
		dbPort    = mustGetenv("DB_PORT") // e.g. '5432'
		dbName    = mustGetenv("DB_NAME") // e.g. 'my-database'
	)
	if dbPwd == "" {
		dbPwd = pullDbPassword()
	}

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
