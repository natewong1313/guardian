package guardian

import (
	"math"
	"testing"
	"time"
)

// mock the db interface for testing use
type mockDB struct {
	createSession     func(session *Session) error
	getSession        func(id string) (*Session, error)
	deleteSession     func(id string) error
	deleteAllSessions func(userID string) error
	updateSession     func(id string, expiresAt time.Time) error
}

func (m *mockDB) CreateSession(session *Session) error {
	return nil
}

func (m *mockDB) GetSession(id string) (*Session, error) {
	return m.getSession(id)
}

func (m *mockDB) DeleteSession(id string) error {
	return nil
}

func (m *mockDB) DeleteAllSessions(userID string) error {
	return nil
}

func (m *mockDB) UpdateSession(id string, expiresAt time.Time) error {
	return nil
}

func TestSessionTokenExpiryRefresh(t *testing.T) {
	expiryTests := [][]int{
		// test expiry refresh
		{
			-16, 14, 30,
		},
		// shouldn't refresh
		{
			-1, 30, 30,
		},
	}

	for _, testData := range expiryTests {
		updatedAt := time.Now().AddDate(0, 0, testData[0])
		expiresAt := time.Now().AddDate(0, 0, testData[1])
		expectedDaysUntilExpiry := testData[2]
		db := &mockDB{
			getSession: func(id string) (*Session, error) {
				return &Session{
					ID:        id,
					UserID:    "test",
					UpdatedAt: updatedAt,
					ExpiresAt: expiresAt,
				}, nil
			},
		}
		session, err := ValidateSessionToken("token", db)
		if err != nil {
			t.Errorf("unexpected error: %v\n", err)
			continue
		}
		daysUntilExpiry := int(math.Round(session.ExpiresAt.Sub(time.Now()).Hours() / 24))
		if daysUntilExpiry != expectedDaysUntilExpiry {
			t.Errorf("expected %d days until expiry, got %d\n", expectedDaysUntilExpiry, daysUntilExpiry)
		}
	}
}
