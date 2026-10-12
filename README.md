# Chirpy

A boot.dev project: a small JSON HTTP API for posting short messages (chirps),
built with the Go standard library, `database/sql`, and PostgreSQL. SQL queries
are compiled into type-safe Go with `sqlc`.

The server listens on **`:8080`** and serves a static site from the current
directory under `/app/`.

## Getting started

1. Copy `example.env` to `.env` and fill in the values (the server loads `.env`
   at startup with `godotenv`).
2. Apply the database migrations in `sql/schema` (via `goose`)
3. Run the server:

   ```sh
   go run .
   ```

   Or build a binary:

   ```sh
   go build .
   ./chirpy
   ```

Run the tests with `go test ./...`.

## Database

### Migrations (goose)

Schema migrations live in `sql/schema` and are managed with
[goose](https://github.com/pressly/goose). The `GOOSE_DRIVER`, `GOOSE_DBSTRING`,
and `GOOSE_MIGRATION_DIR` environment variables tell goose where the database is
and where to find migrations. Apply them with:

```sh
goose up
```

### Type-safe queries (sqlc)

Hand-written SQL lives in `./sql/queries`. The
[sqlc](https://sqlc.dev) Go binary reads those queries together with the schema
in `sql/schema` (configured in `sqlc.yaml` and generates the
type-safe Go code in `internal/database`. After changing
any file in `sql/queries` or `sql/schema`, regenerate the code with:

```sh
sqlc generate
```

The server uses the generated `database.Queries` for all database access.

## Environment variables

| Variable     | Description                                                        |
| ------------ | ------------------------------------------------------------------ |
| `DB_URL`     | PostgreSQL connection string.                                      |
| `JWT_SECRET` | Secret used to sign and verify access JWTs.                        |
| `POLKA_KEY`  | API key expected on the Polka webhook.                             |

Goose environment variables include `GOOSE_DRIVER`, `GOOSE_DBSTRING`, and
`GOOSE_MIGRATION_DIR` for running migrations.

## Authentication

Most endpoints take a bearer token in the `Authorization` header:

- **Access token**: `Authorization: Bearer <token>`. A JWT issued by
  `POST /api/login`, valid for 1 hour.
- **Refresh token**: `Authorization: Bearer <refresh_token>`. An opaque token
  issued by `POST /api/login`, valid for 60 days.
- **Polka webhook**: `Authorization: ApiKey <POLKA_KEY>`.

## Endpoints

### Users

| Method | Path         | Auth       | Description                                   | Returns                              |
| ------ | ------------ | ---------- | --------------------------------------------- | ------------------------------------ |
| POST   | `/api/users` | none       | Create a user.                                | `201 Created` + user object          |
| PUT    | `/api/users` | Bearer JWT | Update the current user's email and password. | `200 OK` + user object               |
| POST   | `/api/login` | none       | Log in and receive tokens.                    | `200 OK` + user object and tokens    |

**POST `/api/users`** request body:

```json
{ "email": "user@example.com", "password": "secret" }
```

Returns `201 Created` with a user object.

**PUT `/api/users`** takes the same body (email and password) and returns
`200 OK` with the updated user object.

**POST `/api/login`** takes the same body and returns `200 OK`:

```json
{
  "id": "…",
  "created_at": "…",
  "updated_at": "…",
  "email": "user@example.com",
  "is_chirpy_red": false,
  "token": "<access JWT>",
  "refresh_token": "<refresh token>"
}
```

### Tokens

| Method | Path           | Auth           | Description                                     | Returns                                    |
| ------ | -------------- | -------------- | ----------------------------------------------- | ------------------------------------------ |
| POST   | `/api/refresh` | Bearer refresh | Exchange a refresh token for a new access JWT.  | `200 OK` + `{ "token": "<access JWT>" }`   |
| POST   | `/api/revoke`  | Bearer refresh | Revoke a refresh token.                         | `204 No Content` (no body)                 |

### Chirps

| Method | Path               | Auth       | Description            | Returns                          |
| ------ | ------------------ | ---------- | ---------------------- | -------------------------------- |
| GET    | `/api/chirps`      | none       | List chirps.           | `200 OK` + array of chirp objects |
| GET    | `/api/chirps/{id}` | none       | Get one chirp by ID.   | `200 OK` + chirp object          |
| POST   | `/api/chirps`      | Bearer JWT | Create a chirp.        | `201 Created` + chirp object     |
| DELETE | `/api/chirps/{id}` | Bearer JWT | Delete your own chirp. | `204 No Content` (no body)       |

**GET `/api/chirps`** supports optional query parameters:

- `author_id` — only return chirps by this user (UUID).
- `sort` — `asc` (default) or `desc` by creation time.

**POST `/api/chirps`** request body:

```json
{ "body": "Hello, Chirpy!" }
```

Bodies are capped at **140 characters** and a few words (`kerfuffle`,
`sharbert`, `fornax`) are replaced with `****`. Returns `201 Created` with a
chirp object.

**DELETE `/api/chirps/{id}`** returns `204 No Content`. You may only delete your
own chirps; otherwise the server responds `403 Forbidden`.

A chirp object looks like:

```json
{
  "id": "…",
  "created_at": "…",
  "updated_at": "…",
  "body": "Hello, Chirpy!",
  "user_id": "…"
}
```

### Webhooks

| Method | Path                  | Auth                 | Description                   | Returns                    |
| ------ | --------------------- | -------------------- | ----------------------------- | -------------------------- |
| POST   | `/api/polka/webhooks` | `ApiKey <POLKA_KEY>` | Handle Polka upgrade events.  | `204 No Content` (no body) |

On a `user.upgraded` event the referenced user is marked as "Chirpy Red" and the
server returns `204 No Content`:

```json
{ "event": "user.upgraded", "data": { "user_id": "…" } }
```

### Admin & static files

| Method | Path             | Description                                                   | Returns                                    |
| ------ | ---------------- | ------------------------------------------------------------- | ------------------------------------------ |
| GET    | `/api/healthz`   | Health check.                                                 | `200 OK` + plain text `OK`                 |
| GET    | `/admin/metrics` | HTML page showing how many times `/app/` has been visited.    | `200 OK` + HTML page                       |
| POST   | `/admin/reset`   | Clears all data and resets the visit counter.                 | `200 OK` + plain text `Reset`              |
| GET    | `/app/`          | Serves static files from the current directory.               | `200 OK` + requested file                  |

## Errors

Errors are returned as JSON with a non-2xx status code:

```json
{ "error": "description of the problem" }
```

## Notes

* I use AI to generate most of this README documentation. Not any of the code
