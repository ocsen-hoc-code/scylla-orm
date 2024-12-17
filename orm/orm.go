package orm

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gocql/gocql"
)

// ORM struct manages table operations
type ORM struct {
	Session        *gocql.Session
	Table          string
	whereClauses   []string
	values         []interface{}
	limit          int
	allowFiltering bool
}

// NewORM initializes a new ORM instance
func NewORM(session *gocql.Session, table string) *ORM {
	return &ORM{
		Session:      session,
		Table:        table,
		whereClauses: []string{},
		values:       []interface{}{},
	}
}

// Where adds a WHERE condition
func (o *ORM) Where(condition string, values ...interface{}) *ORM {
	o.whereClauses = append(o.whereClauses, condition)
	o.values = append(o.values, values...)
	return o
}

// And adds an AND condition
func (o *ORM) And(condition string, values ...interface{}) *ORM {
	return o.Where("AND "+condition, values...)
}

// In adds an IN condition
func (o *ORM) In(column string, values []interface{}) *ORM {
	placeholders := make([]string, len(values))
	for i := range values {
		placeholders[i] = "?"
		o.values = append(o.values, values[i])
	}
	return o.Where(fmt.Sprintf("%s IN (%s)", column, strings.Join(placeholders, ", ")))
}

// Limit sets the LIMIT for SELECT
func (o *ORM) Limit(limit int) *ORM {
	o.limit = limit
	return o
}

// AllowFiltering enables ALLOW FILTERING
func (o *ORM) AllowFiltering() *ORM {
	o.allowFiltering = true
	fmt.Println("Warning: ALLOW FILTERING may cause full table scans and impact performance.")
	return o
}

// ClearWhere clears all WHERE conditions
func (o *ORM) ClearWhere() *ORM {
	o.limit = 0
	o.whereClauses = []string{}
	o.values = []interface{}{}
	o.allowFiltering = false
	return o
}

// Insert inserts a new record
func (o *ORM) Insert(data map[string]interface{}, ttl int) error {
	columns, placeholders, values := []string{}, []string{}, []interface{}{}
	for k, v := range data {
		columns = append(columns, k)
		placeholders = append(placeholders, "?")
		values = append(values, v)
	}

	ttlClause := ""
	if ttl > 0 {
		ttlClause = fmt.Sprintf(" USING TTL %d", ttl)
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)%s",
		o.Table, strings.Join(columns, ", "), strings.Join(placeholders, ", "), ttlClause)

	QueryLogger.Log(query, values...)
	return o.Session.Query(query, values...).Exec()
}

// Select retrieves records and binds them into a slice of structs
func (o *ORM) Select(dest interface{}) error {
	query := fmt.Sprintf("SELECT * FROM %s", o.Table)
	if len(o.whereClauses) > 0 {
		query += " WHERE " + strings.Join(o.whereClauses, " ")
	}
	if o.limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", o.limit)
	}
	if o.allowFiltering {
		query += " ALLOW FILTERING"
	}

	QueryLogger.Log(query, o.values...)
	iter := o.Session.Query(query, o.values...).Iter()
	results := reflect.ValueOf(dest).Elem()

	for {
		row := map[string]interface{}{}
		if !iter.MapScan(row) {
			break
		}
		elem := reflect.New(results.Type().Elem()).Elem()
		for k, v := range row {
			fieldName := strings.Title(strings.ReplaceAll(k, "_", ""))
			field := elem.FieldByName(fieldName)
			if field.IsValid() && field.CanSet() {
				field.Set(reflect.ValueOf(v))
			}
		}
		results.Set(reflect.Append(results, elem))
	}
	return iter.Close()
}

// Update updates records based on WHERE conditions with optional TTL
func (o *ORM) Update(data map[string]interface{}, ttl int) error {
	if len(o.whereClauses) == 0 {
		return fmt.Errorf("update requires at least one WHERE condition")
	}

	setClauses, values := []string{}, []interface{}{}
	for k, v := range data {
		setClauses = append(setClauses, fmt.Sprintf("%s = ?", k))
		values = append(values, v)
	}

	ttlClause := ""
	if ttl > 0 {
		ttlClause = fmt.Sprintf(" USING TTL %d", ttl)
	}

	query := fmt.Sprintf("UPDATE %s%s SET %s WHERE %s", o.Table, ttlClause,
		strings.Join(setClauses, ", "), strings.Join(o.whereClauses, " AND "))

	values = append(values, o.values...)
	QueryLogger.Log(query, values...)
	return o.Session.Query(query, values...).Exec()
}

// Delete removes records based on WHERE conditions
func (o *ORM) Delete() error {
	if len(o.whereClauses) == 0 {
		return fmt.Errorf("delete requires at least one WHERE condition")
	}

	query := fmt.Sprintf("DELETE FROM %s WHERE %s", o.Table, strings.Join(o.whereClauses, " AND "))
	QueryLogger.Log(query, o.values...)
	return o.Session.Query(query, o.values...).Exec()
}
