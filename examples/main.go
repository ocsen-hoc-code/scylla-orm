package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gocql/gocql"
	"github.com/ocsen-hoc-code/scylla-orm/orm"
)

type User struct {
	ID        gocql.UUID `db:"id"`
	Name      string     `db:"name"`
	Age       int        `db:"age"`
	CreatedAt time.Time  `db:"created_at"`
}

func main() {
	// Step 1: Initialize ScyllaDB Session
	session := orm.NewSession(
		[]string{"127.0.0.1"}, // ScyllaDB nodes
		"example_keyspace",    // Keyspace name
		"user",                // Username (optional)
		"password",            // Password (optional)
		"SimpleStrategy",      // Replication strategy
		1,                     // Replication factor
	)
	defer session.Close()

	// Step 2: Run Migration to create 'users' table and materialized view
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS users (
			age INT,
			id UUID,
			name TEXT,
			created_at TIMESTAMP,
			PRIMARY KEY (id)
		)`,
		`CREATE MATERIALIZED VIEW IF NOT EXISTS users_by_age AS
			SELECT * FROM users
			WHERE age IS NOT NULL
			PRIMARY KEY (age, id);`,
	}
	migration := orm.NewMigration(session.Session)
	if err := migration.Run(migrations); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	// Step 3: Initialize ORM instance for 'users' table
	db := orm.NewORM(session.Session, "users")

	// ------------------------- Basic CRUD Operations -------------------------

	// INSERT multiple users
	uuid := gocql.TimeUUID()
	usersData := []map[string]interface{}{
		{
			"id":         uuid,
			"name":       "Alice",
			"age":        30,
			"created_at": time.Now(),
		},
		{
			"id":         gocql.TimeUUID(),
			"name":       "Bob",
			"age":        25,
			"created_at": time.Now(),
		},
		{
			"id":         gocql.TimeUUID(),
			"name":       "Charlie",
			"age":        35,
			"created_at": time.Now(),
		},
	}

	// Insert each user
	for _, user := range usersData {
		if err := db.Insert(user, 0); err != nil {
			log.Fatalf("Failed to insert user: %v", err)
		}
	}

	// ------------------------- Pagination (LIMIT, OFFSET) -------------------------

	// Insert more users for pagination example
	for i := 0; i < 10; i++ {
		if err := db.Insert(map[string]interface{}{
			"id":         gocql.TimeUUID(),
			"name":       fmt.Sprintf("User %d", i+1),
			"age":        20 + i,
			"created_at": time.Now(),
		}, 0); err != nil {
			log.Fatalf("Failed to insert user: %v", err)
		}
	}

	// Step 4: Fetch data with LIMIT and OFFSET for pagination
	viewTable := orm.NewORM(session.Session, "users_by_age")

	// Fetch paginated users with `ALLOW FILTERING`
	var paginatedUsers []User
	if err := viewTable.Where("age > ?", 20).Limit(3).AllowFiltering().Select(&paginatedUsers); err != nil {
		log.Fatalf("Failed to fetch paginated users: %v", err)
	}
	for _, user := range paginatedUsers {
		fmt.Printf("User (Paginated): %+v\n", user)
	}

	log.Println("Completed Pagination query.")

	// ------------------------- TTL Support for INSERT and UPDATE -------------------------
	uuid = gocql.TimeUUID()
	// INSERT with TTL
	if err := db.Insert(map[string]interface{}{
		"id":         gocql.TimeUUID(),
		"name":       "Eve",
		"age":        22,
		"created_at": time.Now(),
	}, 10); err != nil { // 10 seconds TTL
		log.Fatalf("Failed to insert user with TTL: %v", err)
	}

	// UPDATE with TTL
	if err := db.Where("id = ?", uuid).Update(map[string]interface{}{
		"age": 30,
	}, 5); err != nil { // 5 seconds TTL
		log.Fatalf("Failed to update user with TTL: %v", err)
	}

	log.Println("Completed TTL example.")

	// ------------------------- Combined Queries -------------------------

	var users []User

	// 1. WHERE + AND
	fmt.Println("\n--- WHERE + AND ---")
	if err := viewTable.ClearWhere().
		Where("age > ?", 25).
		And("name = ?", "Alice").
		AllowFiltering().
		Select(&users); err != nil {
		log.Fatalf("Failed to fetch users: %v", err)
	}
	for _, user := range users {
		fmt.Printf("User (WHERE + AND): %+v\n", user)
	}

	// 2. IN Condition
	fmt.Println("\n--- IN Condition ---")
	if err := viewTable.ClearWhere().
		In("age", []interface{}{25, 30, 35}).
		AllowFiltering().
		Select(&users); err != nil {
		log.Fatalf("Failed to fetch users with IN condition: %v", err)
	}
	for _, user := range users {
		fmt.Printf("User (IN): %+v\n", user)
	}

	// 3. Complex Query with AND + IN
	fmt.Println("\n--- Complex Query (AND + IN) ---")
	if err := viewTable.ClearWhere().
		Where("name = ?", "Alice").
		And("age IN (25, 30, 35)"). // Custom condition
		AllowFiltering().
		Select(&users); err != nil {
		log.Fatalf("Failed to fetch users with complex query: %v", err)
	}
	for _, user := range users {
		fmt.Printf("User (Complex): %+v\n", user)
	}

	log.Println("\nCompleted all combined queries.")

	// ------------------------- Logging -------------------------

	// Log an example query
	orm.QueryLogger.Log("SELECT * FROM users WHERE age > ?", 20)

	log.Println("All operations completed.")
}
