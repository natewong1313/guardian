module github.com/natewong1313/guardian/examples/email-pass

go 1.24.2

require github.com/natewong1313/guardian/pkg/adapters/sqlite v0.0.0

require github.com/natewong1313/guardian v0.0.0

require (
	golang.org/x/crypto v0.38.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
)

replace github.com/natewong1313/guardian => ../../../guardian

replace github.com/natewong1313/guardian/pkg/adapters/sqlite => ../../pkg/adapters/sqlite
