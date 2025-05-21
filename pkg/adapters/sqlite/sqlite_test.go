package sqlite

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/natewong1313/guardian"
)

const (
	DbName = "test.db"
)

func newMockSession(userID string) *guardian.Session {
	return &guardian.Session{
		ID:        fmt.Sprintf("session%d", rand.Intn(99999999)),
		UserID:    userID,
		UpdatedAt: time.Now().UTC().AddDate(0, 0, -1),
		ExpiresAt: time.Now().UTC().AddDate(0, 0, 1),
	}
}

func TestSqliteSessionCreation(t *testing.T) {
	os.Remove(DbName)
	db, err := New(DbName)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	mockSession := newMockSession("testUser")
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
	if !session.UpdatedAt.Equal(mockSession.UpdatedAt) {
		t.Errorf("updatedAt mismatch")
	}
	if !session.ExpiresAt.Equal(mockSession.ExpiresAt) {
		t.Errorf("expiresAt mismatch")
	}
}

func TestSessionDeletion(t *testing.T) {
	db, err := New(DbName)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	mockSession := newMockSession("testUser")
	if err := db.CreateSession(mockSession); err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	db.DeleteSession(mockSession.ID)

	session, err := db.GetSession(mockSession.ID)
	if err == nil || session != nil {
		t.Errorf("session not deleted: %v", session)
	} else if err.Error() != "sql: no rows in result set" {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestUserSessionDeletion(t *testing.T) {
	db, err := New(DbName)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	mockSession := newMockSession("testUser")
	mockSession2 := newMockSession("testUser")
	if err := db.CreateSession(mockSession); err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if err := db.CreateSession(mockSession2); err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	db.DeleteAllSessions(mockSession.UserID)

	session, err := db.GetSession(mockSession.ID)
	if err == nil || session != nil {
		t.Errorf("session not deleted: %v", session)
	} else if err.Error() != "sql: no rows in result set" {
		t.Errorf("unexpected error: %v", err)
	}
	session, err = db.GetSession(mockSession2.ID)
	if err == nil || session != nil {
		t.Errorf("session not deleted: %v", session)
	} else if err.Error() != "sql: no rows in result set" {
		t.Errorf("unexpected error: %v", err)
	}
}
