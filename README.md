# TMF915 AI Management Suite API

Production-ready REST API implementing the [TM Forum TMF915 AI Management Suite v4.0.0](https://www.tmforum.org/resources/how-to-guide/tmf915-ai-management-api-user-guide/) specification.

## Stack

| Layer | Technology |
|---|---|
| Language | Go 1.22+ |
| HTTP | `net/http` stdlib (Go 1.22 mux with method routing) |
| Database | SQL Server (Windows Integrated Auth) via `sqlx` + `go-mssqldb` |
| Logging | `zerolog` (structured JSON) |
| Auth | `X-API-Key` header |
| Docs | Swagger UI (CDN) at `/docs` |

## Endpoints

Base path: `/tmf-api/AiM/v4`

| Resource | GET (list) | POST | GET (one) | PATCH | DELETE |
|---|---|---|---|---|---|
| `/aiModelSpecification` | ✅ | ✅ | ✅ | ✅ | ✅ |
| `/aiModel` | ✅ | ✅ | ✅ | ✅ | ✅ |
| `/aiContractSpecification` | ✅ | ✅ | ✅ | ✅ | ✅ |
| `/aiContract` | ✅ | ✅ | ✅ | ✅ | ✅ |
| `/aiContractViolation` | ✅ | ✅ | ✅ | — | — |
| `/alarm` | ✅ | ✅ | ✅ | ✅ | ✅ |
| `/rule` | ✅ | ✅ | ✅ | ✅ | ✅ |
| `/monitor` | ✅ | ✅ | ✅ | ✅ | ✅ |
| `/topic` | ✅ | ✅ | ✅ | ✅ | ✅ |
| `/topic/{topicId}/event` | ✅ | — | ✅ | — | — |
| `/hub` | ✅ | ✅ | ✅ | — | ✅ |
| `/listener` | — | ✅ | — | — | ✅ |

All list endpoints support `?fields=`, `?offset=`, and `?limit=` (default 100, max 1000).  
Responses include `X-Result-Count` and `X-Total-Count` headers.

Additional endpoints (no auth required):

| Path | Description |
|---|---|
| `GET /health` | Liveness check — returns `{"status":"ok"}` |
| `GET /docs` | Swagger UI |
| `GET /swagger.json` | OpenAPI spec (host/scheme patched to current server) |

## Prerequisites

- Go 1.22+
- SQL Server (local or remote) with Windows Integrated Authentication
- `sqlcmd` on PATH (comes with SQL Server tools)

## Setup

### 1. Clone

```bash
git clone https://github.com/munteanuser/tmf915-api.git
cd tmf915-api
```

### 2. Configure

Edit `config.yaml`:

```yaml
server:
  port: 8080

database:
  dsn: "server=YOUR_SERVER;database=TMF915;integrated security=true;connection timeout=30;encrypt=false;trustservercertificate=true;app name=TMF915-API"

auth:
  api_key: "your-secret-api-key"

log:
  level: "info"    # debug | info | warn | error
  format: "json"   # json | console
```

The DSN, API key, and port can also be set via environment variables:

| Variable | Overrides |
|---|---|
| `TMF915_DB_DSN` | `database.dsn` |
| `TMF915_API_KEY` | `auth.api_key` |
| `TMF915_PORT` | `server.port` |

### 3. Enable SQL Server TCP/IP

If SQL Server is local and TCP is disabled, run once as Administrator:

```powershell
$tcp = "HKLM:\SOFTWARE\Microsoft\Microsoft SQL Server\MSSQL17.MSSQLSERVER\MSSQLServer\SuperSocketNetLib\Tcp"
Set-ItemProperty $tcp -Name Enabled -Value 1
Set-ItemProperty "$tcp\IPAll" -Name TcpPort -Value "1433"
Set-ItemProperty "$tcp\IPAll" -Name TcpDynamicPorts -Value ""
Restart-Service MSSQLSERVER -Force
```

### 4. Run migrations

```powershell
.\migrate.ps1
```

This creates the `TMF915` database and all 12 tables. Safe to re-run — `CREATE TABLE` statements are idempotent.

### 5. Build and run

```powershell
go build -o tmf915-api.exe ./cmd/server
.\tmf915-api.exe
```

Or with a custom config path:

```powershell
.\tmf915-api.exe path\to\config.yaml
```

Expected startup output:

```json
{"level":"info","server":"...","database":"TMF915","message":"database connected"}
{"level":"info","port":8080,"message":"TMF915 API server starting"}
```

## Usage

### Authentication

All API requests (except `/health`, `/docs`, `/swagger.json`) require the header:

```
X-API-Key: your-secret-api-key
```

Or as a query parameter: `?api_key=your-secret-api-key`

### Example requests

```bash
# List AI models
curl -H "X-API-Key: your-key" http://localhost:8080/tmf-api/AiM/v4/aiModel

# Create an AI model
curl -X POST -H "X-API-Key: your-key" -H "Content-Type: application/json" \
  -d '{"name":"GPT-4","description":"Large language model","state":"active"}' \
  http://localhost:8080/tmf-api/AiM/v4/aiModel

# Patch (partial update)
curl -X PATCH -H "X-API-Key: your-key" -H "Content-Type: application/json" \
  -d '{"state":"inactive"}' \
  http://localhost:8080/tmf-api/AiM/v4/aiModel/{id}

# Delete
curl -X DELETE -H "X-API-Key: your-key" \
  http://localhost:8080/tmf-api/AiM/v4/aiModel/{id}
```

### Swagger UI

Open `http://localhost:8080/docs` in a browser.

Paste your API key in the black bar at the top and click **Save** — it persists in `localStorage` and is injected automatically into every **Try it out** request.

## Event Notifications (Webhooks)

Register a callback URL via `POST /tmf-api/AiM/v4/hub`:

```json
{ "callback": "https://your-server.example.com/notify" }
```

After every create, update, or delete operation the API POSTs a JSON event to all registered callbacks:

```json
{
  "eventType": "AiModelCreateEvent",
  "eventTime": "2026-05-21T01:33:01Z",
  "event": { ... resource payload ... }
}
```

Delivery is async (goroutine per callback), with up to 3 retries and exponential backoff.

## Project Structure

```
tmf915-api/
├── cmd/server/main.go          # Entry point, graceful shutdown
├── internal/
│   ├── config/                 # YAML + env config loader
│   ├── db/                     # sqlx connection pool
│   ├── events/                 # Async webhook dispatcher
│   ├── handlers/               # HTTP handlers (one file per resource)
│   ├── middleware/             # API key auth, request logging, panic recovery
│   ├── models/                 # TMF915 structs with JSON + db tags
│   ├── repository/             # SQL queries (one file per resource)
│   └── router/                 # Route registration + Swagger UI
├── migrations/                 # 12 numbered SQL Server migration scripts
├── config.yaml                 # Runtime configuration
├── migrate.ps1                 # Migration runner (PowerShell)
├── swagger.json                # TMF915 v4.0.0 OpenAPI spec
└── Makefile
```

## Makefile targets

```
make build     # go build → tmf915-api.exe
make run       # build + run
make migrate   # run all SQL migrations via sqlcmd
make vet       # go vet ./...
make clean     # remove binary
```

## License

This implementation is provided as-is. The TMF915 API specification is owned by [TM Forum](https://www.tmforum.org/).
