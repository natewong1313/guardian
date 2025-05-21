package postgres

import (
	"database/sql"
	"errors"
	"time"

	_ "github.com/lib/pq"
	"github.com/natewong1313/guardian"
)

type PostgresAdapter struct {
	db *sql.DB
}

const (
	CreateTableStmt = `CREATE TABLE IF NOT EXISTS guardian_session (
		id VARCHAR(20) PRIMARY KEY NOT NULL,
		user_id TEXT NOT NULL,
		updated_at TIMESTAMP NOT NULL,
		expires_at TIMESTAMP NOT NULL
	);`
	CreateSessionStmt     = "INSERT INTO guardian_session (id, user_id, updated_at, expires_at) VALUES ($1, $2, $3, $4);"
	GetSessionStmt        = "SELECT * FROM guardian_session WHERE id=$1;"
	DeleteSessionStmt     = "DELETE FROM guardian_session WHERE id=$1;"
	DeleteUserSessionStmt = "DELETE FROM guardian_session WHERE user_id=$1;"
	UpdateSessionStmt     = "UPDATE guardian_session SET expires_at=? WHERE id=$1;"
)

func New(dataSourceName string) (*PostgresAdapter, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	if err := migrate(db); err != nil {
		return nil, err
	}
	return &PostgresAdapter{db}, nil
}

func migrate(db *sql.DB) error {
	_, err := db.Exec(CreateTableStmt)
	return err
}

func (s *PostgresAdapter) CreateSession(session *guardian.Session) error {
	_, err := s.db.Exec(CreateSessionStmt, &session.ID, &session.UserID, &session.UpdatedAt, &session.ExpiresAt)
	return err
}

func (s *PostgresAdapter) GetSession(id string) (*guardian.Session, error) {
	session := &guardian.Session{}
	row := s.db.QueryRow(GetSessionStmt, id)
	err := row.Scan(&session.ID, &session.UserID, &session.UpdatedAt, &session.ExpiresAt)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (s *PostgresAdapter) DeleteSession(id string) error {
	result, err := s.db.Exec(DeleteSessionStmt, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err == nil && rowsAffected == 0 {
		return errors.New("session not found")
	}
	return nil
}

func (s *PostgresAdapter) DeleteAllSessions(userID string) error {
	_, err := s.db.Exec(DeleteUserSessionStmt, userID)
	return err
}

func (s *PostgresAdapter) UpdateSession(id string, expiresAt time.Time) error {
	result, err := s.db.Exec(UpdateSessionStmt, expiresAt, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err == nil && rowsAffected == 0 {
		return errors.New("session not found")
	}
	return nil
}
