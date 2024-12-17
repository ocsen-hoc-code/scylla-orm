# Rewriting the README file after reset
readme_content = """
# **Scylla-ORM**

Scylla-ORM is a lightweight and flexible ORM for ScyllaDB, built using **GoCQL**. It simplifies interaction with ScyllaDB, supporting features like **basic CRUD operations**, **batch queries**, **WHERE conditions**, **ALLOW FILTERING**, and **migrations**.

---

## **Features**

- **Auto-create Keyspace**: Automatically create keyspace during initialization.
- **Migration Support**: Easily run table and schema migrations.
- **CRUD Operations**: Supports `SELECT`, `INSERT`, `UPDATE`, and `DELETE` operations.
- **Dynamic WHERE Conditions**: Chainable API for building flexible WHERE clauses (`AND`, `IN`, `ALLOW FILTERING`).
- **Batch Queries**: Execute multiple queries in a single batch.
- **TTL Support**: Set time-to-live for records in `INSERT` and `UPDATE`.
- **Logging**: Logs queries and parameters for better debugging.

---

## **Installation**

To install Scylla-ORM, use `go get`:


```bash
go get github.com/ocsen-hoc-code/scylla-orm
```
## **Usage**

## **1. Initialize ScyllaDB Session**

To connect to ScyllaDB and initialize a session:

```bash
import "github.com/ocsen-hoc-code/scylla-orm/orm"

session := orm.NewSession(
    []string{"127.0.0.1"}, // ScyllaDB nodes
    "example_keyspace",    // Keyspace name
    "user",                // Username (optional)
    "password",            // Password (optional)
    "SimpleStrategy",      // Replication strategy
    1,                     // Replication factor
)
defer session.Close()
```

---
## **2. Migrations**

To create or modify tables:

```bash
migrations := []string{
    `CREATE TABLE IF NOT EXISTS users (
        id UUID PRIMARY KEY,
        name TEXT,
        age INT,
        created_at TIMESTAMP
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
```

---
## **3. CRUD Operations**

### Insert a Record

```bash
db := orm.NewORM(session.Session, "users")
if err := db.Insert(map[string]interface{}{
    "id":         gocql.TimeUUID(),
    "name":       "Alice",
    "age":        30,
    "created_at": time.Now(),
}, 0); err != nil {
    log.Fatalf("Insert failed: %v", err)
}
```

### Select Records**

Use dynamic WHERE conditions and ALLOW FILTERING:

```bash
db := orm.NewORM(session.Session, "users")
if err := db.Insert(map[string]interface{}{
    "id":         gocql.TimeUUID(),
    "name":       "Alice",
    "age":        30,
    "created_at": time.Now(),
}, 0); err != nil {
    log.Fatalf("Insert failed: %v", err)
}
```

### Update Records
```bash
if err := db.Where("id = ?", someUUID).
    Update(map[string]interface{}{
        "age": 35,
    }, 0); err != nil {
    log.Fatalf("Update failed: %v", err)
}
```

### Delete Records
```bash
if err := db.Where("id = ?", someUUID).Delete(); err != nil {
    log.Fatalf("Delete failed: %v", err)
}
```

---
## **4. Batch Queries**

Execute multiple operations in a single batch:

```bash
batch := db.NewBatch()

// Insert
batch.Insert("users", map[string]interface{}{
    "id":         gocql.TimeUUID(),
    "name":       "Bob",
    "age":        28,
    "created_at": time.Now(),
})

// Update
batch.Update("users", map[string]interface{}{
    "age": 29,
}, "name = ?", "Bob")

// Delete
batch.Delete("users", "name = ?", "Charlie")

// Execute batch
if err := batch.Execute(); err != nil {
    log.Fatalf("Batch execution failed: %v", err)
}
```

---
## **5. TTL Support**

Set TTL (time-to-live) for records:

#### Insert with TTL:

```bash
db.Insert(map[string]interface{}{
    "id":         gocql.TimeUUID(),
    "name":       "Eve",
    "age":        22,
    "created_at": time.Now(),
}, 10) // TTL: 10 seconds
```

#### Update with TTL:

```bash
db.Where("name = ?", "Eve").
    Update(map[string]interface{}{
        "age": 25,
    }, 5) // TTL: 5 seconds
```

---
## **6. Logging**

Scylla-ORM logs queries and parameters for better debugging:

```bash
orm.QueryLogger.Log("SELECT * FROM users WHERE age > ?", 25)
```

---
# Contributing
Contributions are welcome! Please submit a pull request with changes or feature suggestions.

---
# License
This project is licensed under the MIT License.

---
# Support
If you have any questions or need support, feel free to open an issue or contact us at [ocsen.hoc.code@gmail.com].