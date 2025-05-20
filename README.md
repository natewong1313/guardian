# Guardian
a dead simple session authentication library for go web projects. pretty much just a port of lucia auth to golang 

view an example [here](./examples/email-pass/main.go)

## adapters
guardian currently has the following adapters available for installation
- sqlite
- postgres [in progress]

its easy to write your own adapter as well. guardian exposes the following interface you can implement:
```go
type DatabaseAdapter interface {
	CreateSession(session *Session) error
	GetSession(id string) (*Session, error)
	DeleteSession(id string) error
	DeleteAllSessions(userID string) error
	UpdateSession(id string, expiresAt time.Time) error
}
```


## todo
- refactor constants
