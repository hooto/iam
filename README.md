# hooto IAM

Identity and Access Management service built with Go and Svelte.

## Features

- User sign-in/sign-out with session management (access token + HTTP-only cookie)
- User registration with password reset via email verification
- Access Key (AK/SK) management for CLI / programmatic access
- Role-Based Access Control
- Third-party app registration with privilege scoping
- OAuth-style auth-code flow for cross-host app login, plus direct sign-in for same-host apps
- Sign-in rate limiting against brute-force attempts
- Svelte 5 SPA admin UI, embedded into the Go binary via `go:embed`
- Embedded [kvgo](https://github.com/lynkdb/kvgo) storage, no external database required
- Reusable `pkg/iamserver` integration library for host apps

## Tech Stack

Go 1.26 / Svelte 5 / Bootstrap 5 / kvgo / [httpsrv/v2](https://github.com/hooto/httpsrv) (Fiber v3) / inauth JWT

## Quick Start

```bash
make install-deps    # Install frontend & backend dependencies
make all             # Build both frontends + both backend binaries
make run-be          # Build + run IAM server on http://localhost:3000
```

Default admin: `sysadmin` / `changeme`

### Development

```bash
make run-fe          # IAM frontend dev server (HMR, proxies /iam/v2 → :3000)
make run-be          # Build + run IAM backend
make run-iam         # Build + run IAM backend and frontend dev server together
make run-demo-be     # Demo third-party app backend (port 3001)
make run-demo-fe     # Demo app frontend dev server
make run-demoapp     # Build + run demo backend and frontend together
```

## Configuration

Config file: `{prefix}/etc/iam_config.toml` (auto-generated on first run). The `-prefix` flag defaults to the parent of the binary directory; data lives at `{prefix}/var/iam_db/`.

```bash
./bin/iam-server -prefix /opt/hooto/iam
```

| Field           | Default             | Description                                      |
|-----------------|---------------------|--------------------------------------------------|
| `http_port`     | `3000`              | HTTP listen port                                 |
| `service_name`  | `hooto IAM Service` | Service display name                             |
| `instance_id`   | auto-generated      | Unique instance ID                               |
| `access_keys`   | auto-generated      | Service-level AK/SK pairs                        |
| `database`      | embedded local      | Remote kvgo client config; omitted → local replica |

## Integration

Third-party apps delegate auth to IAM via `pkg/iamserver`. The host app calls `AppVerifier.Setup()` with its app credentials, registers the user-auth routes, then uses `Resolve()` as the single entry point that handles both browser sessions (cookie) and CLI access-key tokens (`Authorization: Bearer`, verified at IAM via `/open/app-auth/introspect`).

```go
import (
    "github.com/hooto/httpsrv/v2"
    "github.com/hooto/iam/v2/pkg/iamapi"
    "github.com/hooto/iam/v2/pkg/iamserver"
)

// App credentials issued by the IAM server
cfg := &iamserver.AppAuthConfig{
    AppId:     "<app-id>",
    SecretKey: "<app-secret-key>",
    BaseURL:   "http://localhost:3000",
}

iamserver.AppVerifier.Setup(cfg)
iamserver.AppVerifier.Update(&iamapi.AppInstance{
    ID:     cfg.AppId,
    Status: 1,
    Permissions: []*iamapi.AppPermission{
        {Permission: "read", Roles: []string{iamapi.Role_User, iamapi.Role_Guest}},
        {Permission: "write", Roles: []string{iamapi.Role_User}},
    },
})

// Mount browser session routes (/user-auth/{session,sign-in,callback,sign-out})
app := httpsrv.New()
api := app.Group("/myapp/api")
iamserver.RegisterRoutes(api)
```

For host apps built directly on Fiber v3 instead of httpsrv, use the Fiber adapter: `iamserver/authfiber.RegisterAuthRoutes(fiber.Router)`.

`cmd/demoapp/main.go` is a complete, minimal usage example.

## Project Structure

```
cmd/server/              IAM server entry point (:3000), serves /iam/v2 API + SPA
cmd/demoapp/             Reference third-party app (:3001), integration example
frontend/server/         Admin UI (Svelte 5 + Bootstrap 5)
frontend/demoapp/        Demo app UI
internal/apiserver/      API handlers: auth, user, open, admin (+ util.go)
internal/config/         TOML configuration
internal/data/           Data layer (embedded kvgo)
pkg/iamapi/              Shared types, key namespaces (Ns*), validators, role constants
pkg/iamserver/           Reusable auth-integration library for host apps
pkg/iamserver/authfiber/ Fiber v3 adapter for the user-auth routes
```

## License

[Apache License 2.0](LICENSE)
