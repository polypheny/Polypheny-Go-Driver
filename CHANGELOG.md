# Changelog

## v0.3.0 - 2026-11-23

- Bundle the official Prism 2.4 Go bindings.
- Add `SQLState()` and `ErrorCode()` to `*ClientError` for structured server errors. Use `errors.As` to access the metadata. An explicit error code of `0` is reported as present.
- Keep the original server message in `Error()`, including for older message-only responses. A nonterminal error response closes the connection and returns the original message with empty metadata.

## [0.2.0] - 2026-11-21

- Bundle the official Prism 2.3 Go bindings.
- Run Cypher directly with `PolyphenyConn.ExecuteCypherContext`. It supports graph, scalar, and mixed results, buffered pagination, and transactions.
- Preserve typed nodes, edges, paths, and graphs in result values.
- Discover namespaces with `SearchNamespaces` and find graph entities through `SearchEntities`.

## [0.1.0] - 2026-11-15

- First tagged release. Requires Go 1.23 or newer and bundles the official Prism 2.2 Go bindings.
- Run SQL through `database/sql` with positional parameters, paginated results, and transactions.
- Run MQL and look up entity metadata through `sql.Conn.Raw`. Document collection metadata requires Prism 2.2 server support.

[0.2.0]: https://github.com/polypheny/Polypheny-Go-Driver/releases/tag/v0.2.0
[0.1.0]: https://github.com/polypheny/Polypheny-Go-Driver/releases/tag/v0.1.0
