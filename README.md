# hooto IAM

Identity and Access Management service built with Go and Svelte.

## Features

- User sign-in/sign-out with session management (access token + HTTP-only cookie)
- User registration with password reset via email verification
- Access Key (AK/SK) management for programmatic access
- Role-Based Access Control (`sysadmin`, `user`, `developer`, `guest`)
- Third-party app registration with privilege scoping
- `pkg/apptokenhandler` middleware for third-party token verification
- Svelte 5 SPA admin UI, embedded into Go binary
- Embedded [kvgo](https://github.com/lynkdb/kvgo) storage, no external database required

## Tech Stack

Go 1.26 / Svelte 5 / Bootstrap 5 / kvgo / inauth JWT

## Quick Start

```bash
make install-deps    # Install frontend & backend dependencies
make all             # Build frontend and backend
make run-be          # Start server on http://localhost:3000
```

Default admin: `sysadmin` / `changeme`

### Development

```bash
make run-fe          # Frontend dev server (HMR)
make run-be          # Build & run backend
make run-demo-fe     # Demo app frontend
make run-demo-be     # Demo app backend (port 3001)
```

## Configuration

Config file: `{prefix}/etc/iam_config.toml` (auto-generated on first run).

```bash
./bin/iam-server -prefix /opt/hooto/iam
```

| Field           | Default             | Description                    |
|-----------------|---------------------|--------------------------------|
| `http_port`     | `3000`              | HTTP listen port               |
| `service_name`  | `hooto IAM Service` | Service display name           |
| `instance_id`   | auto-generated      | Unique instance ID             |
| `access_keys`   | auto-generated      | Service-level AK/SK pairs      |

## Integration

Use `pkg/apptokenhandler` to verify IAM tokens in your app:

```go
import "github.com/hooto/iam/v2/pkg/apptokenhandler"

var appAuthConfig = &apptokenhandler.AppAuthConfig{
    AppId:     "<app-id>",
    Endpoint:  "http://localhost:3000",
    AccessKey: "<access-key>",
}

// Register as controller on your httpsrv module
mod.RegisterController(apptokenhandler.NewAppTokenHandler(appAuthConfig))
```

## Project Structure

```
cmd/server/          IAM server entry point
cmd/demoapp/         Demo third-party app
frontend/server/     Admin UI (Svelte 5 + Bootstrap 5)
frontend/demoapp/    Demo app UI
internal/apiserver/  Core API handlers
internal/config/     TOML configuration
internal/data/       Data layer (kvgo)
pkg/iamapi/          Shared types, constants, validators
pkg/apptokenhandler/ Reusable auth middleware for third-party apps
```

## License

[Apache License 2.0](LICENSE)