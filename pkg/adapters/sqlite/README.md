## Sqlite Adapter
a premade adapter for Sqlite that implements the [DatabaseAdapter](../../../database.go) interface.

### Installation
```bash
go get github.com/natewong1313/guardian/pkg/adapters/sqlite
```

### Usage
First import the adapter
```go
import (
    "github.com/natewong1313/guardian/pkg/adapters/sqlite"
)
```
You can then initialize the adapter and pass it to guardian
```go
adapter, err := sqlite.New("<your connection string here>")
sessionToken := guardian.GenerateSessionToken()
session, err := guardian.CreateSession(sessionToken, "<user id>", adapter)
```
