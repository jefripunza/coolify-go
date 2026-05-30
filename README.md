# coolify-go

[![Go Reference](https://pkg.go.dev/badge/github.com/jefripunza/coolify-go.svg)](https://pkg.go.dev/github.com/jefripunza/coolify-go)
[![Go Version](https://img.shields.io/github/go-mod/go-version/jefripunza/coolify-go)](https://golang.org)
[![License](https://img.shields.io/github/license/jefripunza/coolify-go)](LICENSE)

An idiomatic, fully-featured, and zero-dependency Go client library for the [Coolify REST API (v4)](https://coolify.io).

Designed for production environments, `coolify-go` enables robust automation of your self-hosted cloud resources, handling service orchestration, database provisioning, application builds, and real-time deployment monitoring using natural Go paradigms.

---

## Key Features

- **Zero External Dependencies**: Built strictly using the Go Standard Library.
- **Strict Error Handling**: Parsed Laravel-style validation errors (`422 Unprocessable Entity`) and Traefik/Caddy domain conflicts (`409 Conflict`) into structured Go types.
- **Robust Resource Orchestration**: Supports Applications, Databases (PostgreSQL, Redis, MySQL, MongoDB, etc.), Projects, Servers, Private Keys, GitHub Apps, Services, Teams, and System settings.
- **Advanced Service Support**: Build, manage, and deploy complex Docker Compose resources natively using raw YAML payloads.
- **Context-Aware API**: Fully supports Go's native `context.Context` for fine-grained timeouts, cancellations, and request tracing.
- **Pointers Helper Utilities**: Integrated helper functions to make working with optional request payload pointers safe and clean.

---

## Installation

To add `coolify-go` to your Go project, run:

```bash
go get github.com/jefripunza/coolify-go
```

Then, import the package in your source file:

```go
import "github.com/jefripunza/coolify-go"
```

*Note: For local workspace testing, you can alias the package module as needed using Go workspaces or replace directives.*

---

## Quickstart

A minimal working example showing how to initialize the client and fetch the system version:

```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jefripunza/coolify-go"
)

func main() {
	// Initialize client. If custom baseURL is empty, default base URL is used.
	// API version path '/api/v1' is automatically appended if missing.
	client := coolify.NewClient("http://your_coolify_host:8000", "your_bearer_token")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Retrieve Coolify instance version
	version, err := client.System.Version(ctx)
	if err != nil {
		log.Fatalf("Failed to fetch Coolify version: %v", err)
	}

	fmt.Printf("Coolify System Version: %s\n", version)
}
```

---

## Environment Configuration

A standard configuration mapping matches the SDK's integration test suite:

```env
COOLIFY_URL=http://your_coolify_host:8000   # Your Coolify host url
COOLIFY_KEY=your_api_token_here             # Coolify bearer API token
WEB_HOST=your_web_host                      # Target domain for application/service routing
PROJECT_UUID=your_project_uuid              # Target Project UUID
SERVER_UUID=your_server_uuid                # Target Server UUID
ENVIRONMENT_UUID=your_env_uuid              # Target Environment UUID
```

---

## API Reference

The client structure groups endpoints logically into cohesive domains:

| Service Sub-component | Interface/Service Field | Key Capabilities                                                                             |
| :-------------------- | :---------------------- | :------------------------------------------------------------------------------------------- |
| **System**            | `client.System`         | Engine Version, Health Checks, API Key operations.                                           |
| **Applications**      | `client.Applications`   | List, Get, Create (Public/Private), Start, Stop, Restart, Env Variable updates.              |
| **Services**          | `client.Services`       | Custom & One-click Docker Compose creation, Start/Stop lifecycle, Envs management.           |
| **Databases**         | `client.Databases`      | Provision & control PostgreSQL, Redis, MySQL, MariaDB, ClickHouse, and MongoDB.              |
| **Servers**           | `client.Servers`        | Server introspection, Resource metrics, Custom domain validations.                           |
| **Projects**          | `client.Projects`       | Project & sub-environment structures creation and deletion.                                  |
| **Deployments**       | `client.Deployments`    | List deployments, retrieve real-time logs, trigger new deployments, or cancel active builds. |
| **Private Keys**      | `client.PrivateKeys`    | Add and remove server SSH private keys.                                                      |
| **GitHub Apps**       | `client.GitHubApps`     | Register GitHub Integrations and read app configuration.                                     |
| **Teams**             | `client.Teams`          | Manage organizations, check user scopes, retrieve team memberships.                          |

---

## Advanced Examples

### 1. Deploying Custom Docker Compose Services

You can deploy multi-container microservices or customized environments using base64-encoded docker-compose strings. The SDK handles the raw deployment and subdomain mapping dynamically.

```go
package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"

	"github.com/jefripunza/coolify-go"
)

func main() {
	client := coolify.NewClient("http://your_coolify_host:8000", "my-api-key")
	ctx := context.Background()

	// Compose definition
	composeYaml := `
version: '3.8'
services:
  web-app:
    image: nginx:alpine
    ports:
      - "80"
`
	encodedCompose := base64.StdEncoding.EncodeToString([]byte(composeYaml))

	// Provision custom service mapping port 80 to web subdomain
	svcResp, err := client.Services.Create(ctx, coolify.CreateServiceRequest{
		ProjectUUID:     "project_uuid_here",
		EnvironmentUUID: "env_uuid_here",
		ServerUUID:      "server_uuid_here",
		DockerComposeRaw: coolify.String(encodedCompose),
		Urls: []coolify.ServiceUrlMapping{
			{
				Name: "web-app",
				Url:  "https://your-subdomain.example.com:443", // Subdomain routing
			},
		},
		InstantDeploy: coolify.Bool(true),
	})
	if err != nil {
		log.Fatalf("Creation failed: %v", err)
	}

	fmt.Printf("Service successfully created. Service UUID: %s\n", svcResp.UUID)
}
```

### 2. High-Fidelity Validation & Domain Conflict Error Handling

Coolify returns structured JSON responses for schema errors (e.g. Laravel validation) and routing conflicts (e.g. overlapping domains). The SDK handles these by parsing them into custom `coolify.Error` types automatically.

```go
resp, err := client.Applications.CreatePublic(ctx, req)
if err != nil {
	if apiErr, ok := err.(*coolify.Error); ok {
		log.Printf("API Error Response Code: %d", apiErr.StatusCode)
		log.Printf("General Message: %s", apiErr.Message)
		
		// 1. Check for Laravel Validation issues (422)
		if len(apiErr.Errors) > 0 {
			for field, fieldErrs := range apiErr.Errors {
				log.Printf("Field [%s] failed checks: %v", field, fieldErrs)
			}
		}

		// 2. Check for Caddy/Traefik Domain Routing conflicts (409)
		if len(apiErr.Conflicts) > 0 {
			for _, conflict := range apiErr.Conflicts {
				log.Printf("Domain conflict on '%s' by %s (%s)", 
					conflict.Domain, conflict.ResourceName, conflict.ResourceType)
			}
		}
	} else {
		log.Fatalf("Non-API error occurred: %v", err)
	}
}
```

### 3. Modifying Environment Variables on Services

Quickly inject or modify environment variables for target services programmatically:

```go
// Add a new environment variable to a service container
envResp, err := client.Services.CreateEnv(ctx, "service-uuid", coolify.CreateServiceEnvRequest{
	Key:         "DB_PASSWORD",
	Value:       "super-secure-string",
	IsPreview:   false,
	IsLiteral:   true,
	IsMultiline: false,
})
if err != nil {
	log.Fatalf("Failed to register environment variable: %v", err)
}

// Update the variable to a different value
updatedEnv, err := client.Services.UpdateEnv(ctx, "service-uuid", coolify.UpdateServiceEnvRequest{
	Key:   "DB_PASSWORD",
	Value: "an-even-more-secure-string",
})
```

---

## Running the Example Application

The repository features a ready-to-run orchestration pipeline verifying application lifecycle deployment inside `example/`:

1. Copy `.env.example` to `.env` and fill in the connection details:
   ```bash
   cp .env.example .env
   ```
2. Run the integration lifecycle script (automatically executes tests and deploys the SSH demo compose container):
   ```bash
   bash run_example.sh
   ```

---

## License

This project is open-source software licensed under the **MIT License**. Feel free to read the `LICENSE` file for more context.
