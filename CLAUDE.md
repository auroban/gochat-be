# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Repository overview

`gochat-be` is a multi-module Go workspace (`go.work`) for a chat backend built as separate services:

- `apis` — shared Protobuf/gRPC contracts (buf-managed), consumed by all services.
- `profile-service` — the only implemented service. Handles user profile creation/lookup over gRPC, backed by Postgres.
- `message-service`, `notification-service` — placeholder modules (`go.mod` only, no code yet).

Modules are wired together via `go.work`, so `profile-service` imports `apis` as `github.com/auroban/gochat-be/apis/...` directly from the sibling directory rather than a published version.

## Build, run, and codegen commands

Run these from the relevant module directory (`profile-service`, `apis`) unless noted.

```bash
# Build profile-service
cd profile-service && go build -o profile-service ./cmd/server

# Run profile-service locally (needs Postgres reachable + env vars below)
cd profile-service && go run ./cmd/server

# Run all tests in a module (no tests exist yet, but this is the convention)
go test ./...
go test ./internal/service/...       # single package
go test -run TestName ./internal/... # single test

# Regenerate protobuf/gRPC Go code from apis/proto/**/*.proto after editing .proto files
cd apis && buf generate

# Lint / breaking-change check for proto contracts
cd apis && buf lint
cd apis && buf breaking --against '.git#branch=master'

# Docker (builds profile-service using the root-level build context so it can COPY go.work + both modules)
docker compose up --build
```

`profile-service` requires these environment variables (see `profile-service/.env`, loaded via docker-compose's `env_file` or exported manually for local runs): `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_SCHEMA`, `DB_SSLMODE`. Startup fails fast (`config.validateEnvVariables`) if any are missing. Non-DB settings (server port, log level) come from `profile-service/configs/config.yaml`.

## Architecture

### Service layering (profile-service)

`cmd/server/main.go` → `internal/server.Start` wires everything together in one place:

1. Load config (`internal/config`, viper: YAML file + env var overrides for DB creds).
2. Connect to Postgres (`internal/db.Connect`) and run migrations on boot (`internal/db.RunMigration`, golang-migrate, `profile-service/migrations/*.sql`).
3. Construct `repository.ProfileRepository` → `service.ProfileService` (injected with a `security.PasswordHasher`, currently `security.BcryptHasher`) → `transport/grpc.ProfileServer`.
4. Register the gRPC server (`google.golang.org/grpc`) and serve on the configured port.

Layers: `transport/grpc` (proto <-> domain translation, request validation, gRPC status codes) → `service` (business logic, e.g. hashing passwords before create) → `repository` (raw SQL via `sqlx`, named queries, Postgres error mapping) → Postgres.

Note: `internal/router` (gorilla/mux) and `internal/models` (validator/v10 struct tags) are leftovers from a prior REST implementation (see commit "Added gRPC in user-service, removed REST endpoints"). The router is no longer wired into `server.Start`; gRPC request validation now happens via `protovalidate` constraints declared directly in the `.proto` file.

### Error handling convention

`internal/domainerror` defines sentinel errors (`ErrProfileAlreadyExists`, `ErrNoProfileFound`). Repository methods wrap driver errors into these (e.g. Postgres unique-violation code `23505` → `ErrProfileAlreadyExists`, `sql.ErrNoRows` → `ErrNoProfileFound`). The gRPC layer (`transport/grpc/profile_server.go`) unwraps with `errors.Is` and maps to gRPC status codes (`AlreadyExists`, `NotFound`, `Internal`); everything else is logged and returned as an opaque `Internal` error rather than leaking internals to clients.

### Proto contracts (`apis`)

- Source of truth: `apis/proto/profile/v1/profile.proto`. Field-level validation is declared inline using `buf.validate` constraints (protovalidate) rather than in application code.
- `apis/buf.gen.yaml` generates Go structs + gRPC stubs into `apis/proto/profile/v1/*.pb.go` via remote buf plugins — regenerate with `buf generate` after any `.proto` change instead of hand-editing generated files.
- Adding a new service/message: define it in a `.proto` file under `apis/proto/<domain>/v1/`, run `buf generate`, then import the generated package from the consuming service module.

### Database

- Postgres, schema-per-service (`profile_service` schema, set via `search_path`/`DB_SCHEMA`).
- Migrations live in `profile-service/migrations/*.up.sql`, applied automatically at service startup via golang-migrate (not a separate CLI step). Add new migrations as sequentially numbered `NNN_description.up.sql` files.
- `repository.BaseEntity` (id, timestamps, created_by/updated_by audit columns) is embedded in entity structs like `repository.Profile`; follow this pattern for new tables/entities.
