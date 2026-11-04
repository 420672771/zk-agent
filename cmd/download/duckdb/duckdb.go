package main

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/duckdb/duckdb-go/v2"
)

// duckdbExtensions is the list of DuckDB extensions required by WeKnora's
// data analysis tool. `spatial` is used for layer metadata (st_read_meta)
// so we can enumerate sheet names from Excel files, while `excel` provides
// the dedicated read_xlsx reader with proper type inference.
var duckdbExtensions = []string{"spatial", "excel"}

// main verifies that each required extension can be loaded from the local
// extension directory. The extension binaries are pre-downloaded by
// Dockerfile.app (curl + gunzip) into /root/.duckdb/extensions, because the
// container's DuckDB INSTALL path is unusable here: installing any extension
// triggers autoloading of the `httpfs` dependency via a URL that omits the
// `.gz` suffix, which the extension server rejects with HTTP 404. Loading from
// the pre-seeded directory avoids that network round-trip entirely.
func main() {
	ctx := context.Background()

	sqlDB, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		panic(err)
	}
	defer sqlDB.Close()

	for _, ext := range duckdbExtensions {
		if _, err := sqlDB.ExecContext(ctx, fmt.Sprintf("LOAD %s;", ext)); err != nil {
			panic(fmt.Errorf("failed to load %s extension: %w", ext, err))
		}
	}
}
