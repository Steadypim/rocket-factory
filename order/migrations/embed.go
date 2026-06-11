package migrations

import "embed"

// FS contains all SQL migrations required by Order Service.
//
//go:embed *.sql
var FS embed.FS
