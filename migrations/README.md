# Migrations

Migration files are embedded into the Go binary.

Use this naming pattern:

```text
006_short_description.up.sql
006_short_description.down.sql
```

Current project only includes `up` migrations. Add `down` files when rollback is required.

Run migrations:

```bash
APP_ENV=local go run ./cmd/migrate -action up
APP_ENV=local go run ./cmd/migrate -action version
APP_ENV=local go run ./cmd/migrate -action force -version 1
```

`down` requires matching `.down.sql` files.
