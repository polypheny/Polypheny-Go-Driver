<p align="center">
  <a href="https://polypheny.org/">
    <picture>
      <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/polypheny/Admin/master/Logo/logo-white-text_cropped.png">
      <img width="50%" src="https://raw.githubusercontent.com/polypheny/Admin/master/Logo/logo-transparent_cropped.png" alt="Polypheny">
    </picture>
  </a>
</p>

# Polypheny Go Driver

Connect to Polypheny-DB from Go using `database/sql`. The driver supports SQL
queries, positional parameters, and transactions. The raw connection provides
MQL and Cypher queries, and Prism namespace and entity metadata.

## Installation

Use Go 1.23 or newer and a Polypheny server with a compatible Prism 2.x
interface. The driver includes the generated Go bindings from the Prism 2.3 release.
Document collection metadata requires server-side Prism 2.2 support. The Cypher
API targets the graph support in Polypheny's Prism 2.3 implementation.

```sh
go get github.com/polypheny/Polypheny-Go-Driver
```

## Querying with SQL

The driver registers as `polypheny` and accepts
`host:port,username:password` as its connection string. Prism uses port
`20590` by default. The example assumes a local instance configured with
the user `pa` and an empty password. The driver uses plain TCP, which does not
encrypt credentials or query traffic. Connect through a trusted network or an
encrypted tunnel.

```go
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/polypheny/Polypheny-Go-Driver"
)

func main() {
	db, err := sql.Open("polypheny", "localhost:20590,pa:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		log.Fatal(err)
	}
	var value int
	if err := db.QueryRowContext(ctx, "SELECT 1").Scan(&value); err != nil {
		log.Fatal(err)
	}
	fmt.Println(value)
}
```

Use `?` placeholders with `QueryContext`, `QueryRowContext`, and
`ExecContext`, or prepare a statement with `PrepareContext`. For example,
`db.QueryRowContext(ctx, "SELECT name FROM employees WHERE id = ?", 42)`
queries an existing `employees` table. Statements outside an explicit
transaction use auto-commit. Use `db.BeginTx(ctx, nil)` to group statements
and call `Commit` or `Rollback` to finish the transaction.

## MQL and entity metadata

Import the driver as `polypheny` instead of the blank import and use
`(*sql.Conn).Raw` to access driver-specific methods. The following example
queries an existing `employees` collection or table in the default
namespace, then lists the entities in `public`. Call
`SearchEntities` with a pattern pointer to filter the results.
The response types are available in
`github.com/polypheny/Polypheny-Go-Driver/prism`.

```go
conn, err := db.Conn(ctx)
if err != nil {
	return err
}
defer conn.Close()

return conn.Raw(func(raw any) error {
	driverConn := raw.(*polypheny.PolyphenyConn)
	documents, err := driverConn.QueryMongoContext(ctx, "db.employees.find()")
	if err != nil {
		return err
	}
	fmt.Println(*documents)

	entities, err := driverConn.SearchEntities(ctx, "public", nil)
	if err != nil {
		return err
	}
	for _, entity := range entities.GetEntities() {
		fmt.Println(entity)
	}
	return nil
})
```

## Cypher

Use `ExecuteCypherContext(ctx, namespace, query)` inside `(*sql.Conn).Raw`.
It returns a `*prism.StatementResult` with an update count or a fully
collected frame. Read `GetGraphFrame()` for graph-only projections and
`GetRelationalFrame()` for scalar or mixed projections. Mixed rows retain
column order and typed Prism values, including nested graph values.
`SearchNamespaces` can filter to the `GRAPH` model, and `SearchEntities`
returns graph markers for a graph namespace.

See the [Go driver reference](https://docs.polypheny.com/en/latest/drivers/go)
for a complete Cypher example.

## Testing

```sh
go test -race ./...
go vet ./...
```

To run the integration tests, set `POLYPHENY_TEST_DSN` to a test instance
that supports Prism 2.3. The tests create and remove their own fixtures.

```sh
POLYPHENY_TEST_DSN='localhost:20590,pa:' go test -count=1 -run '^TestIntegration' ./...
```

## Updating the Prism bindings

When updating Prism, set the version and SHA-256 of its official Go release
archive in `scripts/update-prism-bindings.sh`, run the script, and commit the
generated `prism/` files. Run the script with `--check` to compare the tracked
files with the verified archive. Builds use the tracked files directly.

## Resources

- [Go package documentation](https://pkg.go.dev/github.com/polypheny/Polypheny-Go-Driver)
- [Prism protocol documentation](https://docs.polypheny.com/en/latest/query_interfaces/prism/protocol)
- [Open Go driver issues](https://github.com/polypheny/Polypheny-DB/labels/A-golang)
- [Polypheny project and contribution guidelines](https://github.com/polypheny/Admin)
- [Code of conduct](https://github.com/polypheny/Admin/blob/master/CODE_OF_CONDUCT.md)

## License

[Apache License 2.0](LICENSE)
