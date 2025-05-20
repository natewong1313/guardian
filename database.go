package guardian

import "time"

type DatabaseAdapter interface {
	CreateSession(session *Session) error
	GetSession(id string) (*Session, error)
	DeleteSession(id string) error
	DeleteAllSessions(userID string) error
	UpdateSession(id string, expiresAt time.Time) error
}
