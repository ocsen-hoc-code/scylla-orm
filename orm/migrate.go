package orm

import (
	"log"

	"github.com/gocql/gocql"
)

type Migration struct {
	Session *gocql.Session
}

func NewMigration(session *gocql.Session) *Migration {
	return &Migration{Session: session}
}

func (m *Migration) Run(migrations []string) error {
	for _, migration := range migrations {
		if err := m.Session.Query(migration).Exec(); err != nil {
			log.Printf("Migration failed: %v", err)
			return err
		}
		log.Printf("Migration executed: %v", migration)
	}
	return nil
}
