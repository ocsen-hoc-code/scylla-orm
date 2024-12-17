package orm

import (
	"log"
	"time"
)

// QueryLogger manages query logging in the ORM
var QueryLogger = &Logger{
	Enabled: true, // Logging is enabled by default
}

// Logger struct manages log state and methods
type Logger struct {
	Enabled bool
}

// Enable turns on query logging
func (l *Logger) Enable() {
	l.Enabled = true
	log.Println("Query logging enabled.")
}

// Disable turns off query logging
func (l *Logger) Disable() {
	l.Enabled = false
	log.Println("Query logging disabled.")
}

// Log logs the query and parameters if logging is enabled
func (l *Logger) Log(query string, params ...interface{}) {
	if l.Enabled {
		log.Printf("[QUERY] %s\n[PARAMS] %v\n[TIME] %s\n", query, params, time.Now().Format(time.RFC3339))
	}
}
