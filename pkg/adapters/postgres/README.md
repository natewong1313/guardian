## Postgres Adapter
a premade adapter for Postgres that implements the `DatabaseAdapter` interface.

### Installation
```bash
go get github.com/natewong1313/guardian/pkg/adapters/postgres
```

### Usage
First import the adapter
```go
import (
    "github.com/natewong1313/guardian/pkg/adapters/postgres"
)
```
You can then initialize the adapter
```go
db, err := postgres.New("<your connection string here>")
sessionToken := guardian.GenerateSessionToken()
session, err := guardian.CreateSession(sessionToken, "<user id>", db)
```


