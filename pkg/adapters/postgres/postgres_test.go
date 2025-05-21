package postgres

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/natewong1313/guardian"
)

const (
	DBName   = "postgres"
	Username = "postgres"
	Password = "password"
)

// tests if we can create a session and retrieve it without the contents changing
func TestPostgresSessionCreation(t *testing.T) {
	connectionStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", Username, Password, DBName)
	db, err := New(connectionStr)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	mockSession := &guardian.Session{
		ID:        fmt.Sprintf("session%d", rand.Intn(99999999)),
		UserID:    fmt.Sprintf("user%d", rand.Intn(99999999)),
		UpdatedAt: time.Now().UTC().AddDate(0, 0, -1),
		ExpiresAt: time.Now().UTC().AddDate(0, 0, 1),
	}

	if err := db.CreateSession(mockSession); err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	session, err := db.GetSession(mockSession.ID)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	if session.ID != mockSession.ID {
		t.Errorf("id mismatch")
	}
	if session.UserID != mockSession.UserID {
		t.Errorf("userID mismatch")
	}
	if !session.UpdatedAt.Equal(mockSession.UpdatedAt.Round(time.Microsecond)) {
		t.Errorf("updatedAt mismatch")
	}
	if !session.ExpiresAt.Equal(mockSession.ExpiresAt.Round(time.Microsecond)) {
		t.Errorf("expiresAt mismatch")
	}
}
