package guardian

import "time"

type DatabaseAdapter interface {
	CreateSession(session *Session) error
	GetSession(id string) (*Session, error)
	DeleteSession(id string) error
	DeleteAllSessions(user_id int) error
	UpdateSession(id string, expires_at time.Time) error
}
