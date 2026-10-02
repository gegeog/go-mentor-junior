# Lesson 3: PostgreSQL persistence

> Draft for mentor review. This is not yet an active student assignment.

## Goal

Replace the in-memory stores from Lesson 2 with PostgreSQL implementations without changing the service's API or business behavior. The required behavior remains defined by the [project specification](../../docs/project.md).

## Support material

These optional pages provide starting points for the decisions introduced in this lesson:

- ...

## What you will practice

- deriving a relational schema from queries, ownership, and invariants;
- expressing appropriate invariants with PostgreSQL constraints;
- choosing indexes for known access patterns;
- evolving a schema with versioned migrations;
- making related database changes atomic with transactions;
- handling concurrent changes to persistent state;
- testing a store against a real PostgreSQL instance;
- managing a database connection pool as part of the service lifecycle.

## Common tasks

1. Design and document a PostgreSQL schema for all data owned by the assigned service.
2. Add ordered, versioned migrations that create the schema from an empty database.
3. Replace the in-memory stores with PostgreSQL implementations behind the existing application interfaces.
4. Use transactions where one business operation changes multiple database records, and prevent concurrent commands from violating the project rules.
5. Configure and manage a database connection pool. Keep `/livez` process-only and make `/readyz` report whether the database is available.
6. Add integration tests for the PostgreSQL stores and transactional behavior using a real PostgreSQL instance. The tests must apply the migrations themselves and start from isolated state.
7. Keep the Lesson 2 unit tests and all earlier service requirements working. Document how to configure the database, apply migrations, run the service, and run the integration tests.

Students choose the database driver or toolkit, migration tool, schema, SQL organization, transaction strategy, and integration-test setup. Schema creation must use committed migrations rather than ORM auto-migration or startup code.

## Track A: Order Service

1. Persist Carts, Cart Items, Orders, and Order Items, including the snapshots required by the project specification.
2. Make creating an Order and clearing its Cart atomic.
3. Preserve valid Order transitions when concurrent commands target the same Order.

## Track B: Restaurant Service

1. Persist Restaurants and Menu Items.
2. Preserve idempotent replacement and currency consistency when requests execute concurrently.

## Acceptance criteria

1. The assigned service implements the same API and business behavior as Lesson 2.
2. All service-owned application data is stored in PostgreSQL and survives a service restart.
3. The schema has justified keys, constraints, relationships, and indexes; it does not create foreign keys to data owned by another service.
4. Committed, ordered migrations can create the complete schema in an empty database.
5. Applying all forward migrations to an up-to-date database succeeds without changing existing migration files.
6. The service does not apply migrations automatically during startup.
7. Multi-record business operations are atomic, and concurrent commands cannot produce invalid state.
8. PostgreSQL store and transaction integration tests run against a real PostgreSQL instance with isolated test state.
9. Lesson 2 unit tests remain independent from PostgreSQL and continue to pass.
10. `GET /livez` remains independent from the database; `GET /readyz` returns a non-success status while the database is unavailable.
11. `go test ./...`, the documented integration-test command, `go test -race ./...`, `go fmt ./...`, `go vet ./...`, and `go build ./...` succeed in their documented environments.
12. Database configuration, migration, startup, integration-test, and production migration-policy instructions are documented in the service README.

## Not required

- changing the HTTP API or adding business features;
- a shared database or cross-service foreign keys;
- real communication between the Order and Restaurant Services;
- ORM-generated schema or a specific PostgreSQL library;
- production deployment automation;
- zero-downtime migrations, online backfills, or partitioning;
- an outbox, message broker, saga, or distributed transaction;
- automated production rollback or required down migrations;
- query-performance targets, load tests, replication, or backups.
