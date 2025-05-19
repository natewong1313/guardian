module github.com/natewong1313/guardian/examples/email-pass

go 1.24.2

require github.com/natewong1313/guardian/pkg/adapters/sqlite v0.0.0

require (
	github.com/natewong1313/guardian v0.0.0
	golang.org/x/crypto v0.38.0
)

require github.com/mattn/go-sqlite3 v1.14.28 // indirect

replace github.com/natewong1313/guardian => ../../../guardian

replace github.com/natewong1313/guardian/pkg/adapters/sqlite => ../../pkg/adapters/sqlite
