package orm

import (
	"log"
	"strings"

	"github.com/gocql/gocql"
)

// Batch provides a way to execute multiple queries as a single batch
type Batch struct {
	Session *gocql.Session
	batch   *gocql.Batch
	queries []string
	params  [][]interface{}
}

// NewBatch initializes a new Batch instance
func (o *ORM) NewBatch() *Batch {
	return &Batch{
		Session: o.Session,
		batch:   o.Session.NewBatch(gocql.LoggedBatch), // Logged batch ensures atomicity
		queries: []string{},
		params:  [][]interface{}{},
	}
}

// Insert adds an INSERT query to the batch
func (b *Batch) Insert(table string, data map[string]interface{}) *Batch {
	columns, placeholders, values := []string{}, []string{}, []interface{}{}
	for k, v := range data {
		columns = append(columns, k)
		placeholders = append(placeholders, "?")
		values = append(values, v)
	}

	query := "INSERT INTO " + table + " (" + join(columns, ", ") + ") VALUES (" + join(placeholders, ", ") + ")"
	b.addToBatch(query, values)
	return b
}

// Update adds an UPDATE query to the batch
func (b *Batch) Update(table string, data map[string]interface{}, where string, whereValues ...interface{}) *Batch {
	setClauses, values := []string{}, []interface{}{}
	for k, v := range data {
		setClauses = append(setClauses, k+" = ?")
		values = append(values, v)
	}

	query := "UPDATE " + table + " SET " + join(setClauses, ", ") + " WHERE " + where
	values = append(values, whereValues...)
	b.addToBatch(query, values)
	return b
}

// Delete adds a DELETE query to the batch
func (b *Batch) Delete(table string, where string, whereValues ...interface{}) *Batch {
	query := "DELETE FROM " + table + " WHERE " + where
	b.addToBatch(query, whereValues)
	return b
}

// Execute runs the batch queries
func (b *Batch) Execute() error {
	for i, query := range b.queries {
		b.batch.Query(query, b.params[i]...)
		log.Printf("Batch Query: %s | Params: %v", query, b.params[i])
	}
	if err := b.Session.ExecuteBatch(b.batch); err != nil {
		log.Printf("Failed to execute batch: %v", err)
		return err
	}
	log.Println("Batch executed successfully.")
	return nil
}

// addToBatch adds a query and its parameters to the batch
func (b *Batch) addToBatch(query string, values []interface{}) {
	b.queries = append(b.queries, query)
	b.params = append(b.params, values)
}

// join is a helper function to concatenate strings with a separator
func join(elements []string, separator string) string {
	return strings.Join(elements, separator)
}
