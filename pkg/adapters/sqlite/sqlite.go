package sqlite

import (
	"database/sql"
	"errors"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/natewong1313/guardian"
)

type SQLiteAdapter struct {
	db *sql.DB
}

const (
	CreateTableStmt = `CREATE TABLE IF NOT EXISTS guardian_session (
		id VARCHAR(20) PRIMARY KEY NOT NULL,
		user_id TEXT NOT NULL,
		updated_at DATETIME NOT NULL,
		expires_at DATETIME NOT NULL
	)`
	CreateSessionStmt     = "INSERT INTO guardian_session (id, user_id, updated_at, expires_at) VALUES (?, ?, ?, ?)"
	GetSessionStmt        = "SELECT * FROM guardian_session WHERE id=?"
	DeleteSessionStmt     = "DELETE FROM guardian_session WHERE id=?"
	DeleteUserSessionStmt = "DELETE FROM guardian_session WHERE user_id=?"
	UpdateSessionStmt     = "UPDATE guardian_session SET expires_at=? WHERE id=?"
)

func New(dataSourceName string) (*SQLiteAdapter, error) {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	if err := migrate(db); err != nil {
		return nil, err
	}
	return &SQLiteAdapter{db}, nil
}

func migrate(db *sql.DB) error {
	_, err := db.Exec(CreateTableStmt)
	return err
}

func (s *SQLiteAdapter) CreateSession(session *guardian.Session) error {
	_, err := s.db.Exec(CreateSessionStmt, &session.ID, &session.UserID, &session.UpdatedAt, &session.ExpiresAt)
	return err
}

func (s *SQLiteAdapter) GetSession(id string) (*guardian.Session, error) {
	session := &guardian.Session{}
	row := s.db.QueryRow(GetSessionStmt, id)
	err := row.Scan(&session.ID, &session.UserID, &session.UpdatedAt, &session.ExpiresAt)
	return session, err
}

func (s *SQLiteAdapter) DeleteSession(id string) error {
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

func (s *SQLiteAdapter) DeleteAllSessions(userID string) error {
	_, err := s.db.Exec(DeleteUserSessionStmt, userID)
	return err
}

func (s *SQLiteAdapter) UpdateSession(id string, expiresAt time.Time) error {
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
