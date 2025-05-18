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

func New(table string) (*SQLiteAdapter, error) {
	db, err := sql.Open("sqlite3", table)
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
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS guardian_session (
		id VARCHAR(20) PRIMARY KEY,
		user_id TEXT,
		expires_at DATETIME
	)`)
	return err
}

func (s *SQLiteAdapter) CreateSession(session *guardian.Session) error {
	_, err := s.db.Exec("INSERT INTO guardian_session (id, user_id, expires_at) VALUES (?, ?, ?)", &session.ID, &session.UserID, &session.ExpiresAt)
	return err
}

func (s *SQLiteAdapter) GetSession(id string) (*guardian.Session, error) {
	session := &guardian.Session{}
	row := s.db.QueryRow("SELECT * FROM guardian_session WHERE id=?", id)
	err := row.Scan(&session.ID, &session.UserID, &session.ExpiresAt)
	return session, err
}

func (s *SQLiteAdapter) DeleteSession(id string) error {
	result, err := s.db.Exec("DELETE FROM guardian_session WHERE id=?", id)
	if err != nil {
		return err
	}
	rows_affected, err := result.RowsAffected()
	if err == nil && rows_affected == 0 {
		return errors.New("session not found")
	}
	return nil
}

func (s *SQLiteAdapter) DeleteAllSessions(user_id string) error {
	_, err := s.db.Exec("DELETE FROM guardian_session WHERE user_id=?", user_id)
	return err
}

func (s *SQLiteAdapter) UpdateSession(id string, expires_at time.Time) error {
	result, err := s.db.Exec("UPDATE guardian_session SET expires_at=? WHERE id=?", expires_at, id)
	if err != nil {
		return err
	}
	rows_affected, err := result.RowsAffected()
	if err == nil && rows_affected == 0 {
		return errors.New("session not found")
	}
	return nil
}
