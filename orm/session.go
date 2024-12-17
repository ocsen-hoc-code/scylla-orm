package orm

import (
	"fmt"
	"log"

	"github.com/gocql/gocql"
)

// ScyllaSession manages the ScyllaDB session
type ScyllaSession struct {
	Session *gocql.Session
}

// NewSession creates a new ScyllaDB session and ensures the keyspace exists.
// Supports optional authentication (username/password).
func NewSession(hosts []string, keyspace, username, password, strategy string, replicationFactor int) *ScyllaSession {
	// Cluster configuration
	cluster := gocql.NewCluster(hosts...)
	cluster.Keyspace = "system" // Connect to system keyspace for keyspace management

	// Add authentication only if username and password are provided
	if username != "" && password != "" {
		cluster.Authenticator = gocql.PasswordAuthenticator{
			Username: username,
			Password: password,
		}
		log.Println("Authentication enabled for ScyllaDB connection.")
	} else {
		log.Println("No authentication used for ScyllaDB connection.")
	}

	// Connect to system keyspace
	systemSession, err := cluster.CreateSession()
	if err != nil {
		log.Fatalf("Failed to connect to ScyllaDB: %v", err)
	}
	defer systemSession.Close()

	// Set default strategy to SimpleStrategy if not provided
	if strategy == "" {
		strategy = "SimpleStrategy"
	}

	// Create keyspace if it doesn't exist
	createKeyspace := fmt.Sprintf(`
		CREATE KEYSPACE IF NOT EXISTS %s 
		WITH replication = {'class': '%s', 'replication_factor': %d}`,
		keyspace, strategy, replicationFactor)

	if err := systemSession.Query(createKeyspace).Exec(); err != nil {
		log.Fatalf("Failed to create keyspace '%s': %v", keyspace, err)
	}

	log.Printf("Keyspace '%s' with strategy '%s' ensured.", keyspace, strategy)

	// Reconnect to the specified keyspace
	cluster.Keyspace = keyspace
	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatalf("Failed to connect to keyspace '%s': %v", keyspace, err)
	}

	log.Printf("Connected to keyspace '%s'.", keyspace)
	return &ScyllaSession{Session: session}
}

// Close closes the ScyllaDB session
func (s *ScyllaSession) Close() {
	if s.Session != nil {
		s.Session.Close()
		log.Println("ScyllaDB session closed.")
	}
}
