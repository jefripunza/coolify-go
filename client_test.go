package coolify

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// loadEnv parses the local .env file manually.
func loadEnv() {
	content, err := os.ReadFile(".env")
	if err != nil {
		return
	}
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])
			val = strings.Trim(val, `"'`)
			os.Setenv(key, val)
		}
	}
}

func TestNewClient(t *testing.T) {
	token := "test-bearer-token"
	baseURL := "https://my-coolify.com/api/v1/"
	client := NewClient(baseURL, token)

	if client.token != token {
		t.Errorf("Expected token %s, got %s", token, client.token)
	}

	if client.baseURL.String() != baseURL {
		t.Errorf("Expected BaseURL %s, got %s", baseURL, client.baseURL.String())
	}
}

func TestClient_Headers(t *testing.T) {
	token := "secret-auth-key"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify standard headers
		authHeader := r.Header.Get("Authorization")
		expectedAuth := "Bearer " + token
		if authHeader != expectedAuth {
			t.Errorf("Expected Authorization header %q, got %q", expectedAuth, authHeader)
		}

		accept := r.Header.Get("Accept")
		if accept != "application/json" {
			t.Errorf("Expected Accept header to be application/json, got %q", accept)
		}

		userAgent := r.Header.Get("User-Agent")
		if userAgent != "go-coolify/v0.1.0" {
			t.Errorf("Expected User-Agent to be go-coolify/v0.1.0, got %q", userAgent)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[]`))
	}))
	defer server.Close()

	client := NewClient(server.URL, token)

	_, err := client.Applications.List(context.Background())
	if err != nil {
		t.Fatalf("Unexpected list error: %v", err)
	}
}

func TestClient_ErrorParsing(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
		wantErr    func(*testing.T, error)
	}{
		{
			name:       "Plain text error response",
			statusCode: http.StatusBadRequest,
			body:       "Bad Request Header",
			wantErr: func(t *testing.T, err error) {
				apiErr, ok := err.(*Error)
				if !ok {
					t.Fatalf("Expected error type *Error, got %T", err)
				}
				if apiErr.StatusCode != http.StatusBadRequest {
					t.Errorf("Expected status code 400, got %d", apiErr.StatusCode)
				}
				if apiErr.Message != "Bad Request Header" {
					t.Errorf("Expected message 'Bad Request Header', got %q", apiErr.Message)
				}
			},
		},
		{
			name:       "Structured Laravel 422 validation response",
			statusCode: http.StatusUnprocessableEntity,
			body:       `{"message":"The given data was invalid.","errors":{"name":["The name field is required."],"git_repository":["The git repository must be a valid URL."]}}`,
			wantErr: func(t *testing.T, err error) {
				apiErr, ok := err.(*Error)
				if !ok {
					t.Fatalf("Expected error type *Error, got %T", err)
				}
				if apiErr.StatusCode != http.StatusUnprocessableEntity {
					t.Errorf("Expected status code 422, got %d", apiErr.StatusCode)
				}
				if apiErr.Message != "The given data was invalid." {
					t.Errorf("Expected message 'The given data was invalid.', got %q", apiErr.Message)
				}
				if len(apiErr.Errors) != 2 {
					t.Errorf("Expected 2 validation field errors, got %d", len(apiErr.Errors))
				}
				if apiErr.Errors["name"][0] != "The name field is required." {
					t.Errorf("Expected validation text for field name, got %q", apiErr.Errors["name"][0])
				}
			},
		},
		{
			name:       "Structured 409 domain conflicts response",
			statusCode: http.StatusConflict,
			body:       `{"message":"Domain conflicts detected.","warning":"Routing conflicts possible.","conflicts":[{"domain":"test.com","resource_name":"My Web App","resource_uuid":"uuid-123","resource_type":"application","message":"Already in use"}]}`,
			wantErr: func(t *testing.T, err error) {
				apiErr, ok := err.(*Error)
				if !ok {
					t.Fatalf("Expected error type *Error, got %T", err)
				}
				if apiErr.StatusCode != http.StatusConflict {
					t.Errorf("Expected status code 409, got %d", apiErr.StatusCode)
				}
				if apiErr.Warning != "Routing conflicts possible." {
					t.Errorf("Expected warning message, got %q", apiErr.Warning)
				}
				if len(apiErr.Conflicts) != 1 {
					t.Errorf("Expected 1 conflict record, got %d", len(apiErr.Conflicts))
				}
				conflict := apiErr.Conflicts[0]
				if conflict.Domain != "test.com" || conflict.ResourceName != "My Web App" || conflict.ResourceType != "application" {
					t.Errorf("Unexpected conflict details mapping: %+v", conflict)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.body))
			}))
			defer server.Close()

			client := NewClient(server.URL, "")
			_, err := client.Applications.List(context.Background())
			if err == nil {
				t.Fatal("Expected request to fail, but it succeeded")
			}
			tt.wantErr(t, err)
		})
	}
}

func TestApplications_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[
			{"id":1,"uuid":"app-uuid-1","name":"Main API","status":"running","build_pack":"nixpacks","created_at":"2026-05-30T00:00:00Z","updated_at":"2026-05-30T00:00:00Z"},
			{"id":2,"uuid":"app-uuid-2","name":"Admin Dashboard","status":"exited","build_pack":"static","created_at":"2026-05-30T00:00:00Z","updated_at":"2026-05-30T00:00:00Z"}
		]`))
	}))
	defer server.Close()

	client := NewClient(server.URL, "")
	apps, err := client.Applications.List(context.Background())
	if err != nil {
		t.Fatalf("Unexpected listing error: %v", err)
	}

	if len(apps) != 2 {
		t.Fatalf("Expected 2 applications, got %d", len(apps))
	}

	if apps[0].UUID != "app-uuid-1" || apps[0].Name != "Main API" || apps[0].Status != "running" {
		t.Errorf("Unexpected values in application 0: %+v", apps[0])
	}
}

func TestApplications_CreatePublic(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check incoming JSON body
		var body CreatePublicRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("Failed to decode JSON body: %v", err)
		}

		if body.ProjectUUID != "proj-1" || body.ServerUUID != "server-1" || body.GitRepository != "https://github.com/test/repo" {
			t.Errorf("Unexpected values in body: %+v", body)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"uuid":"app-created-123"}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, "")
	reqBody := CreatePublicRequest{
		ProjectUUID:     "proj-1",
		ServerUUID:      "server-1",
		EnvironmentName: "production",
		EnvironmentUUID: "env-uuid-123",
		GitRepository:   "https://github.com/test/repo",
		GitBranch:       "main",
		BuildPack:       "nixpacks",
		PortsExposes:    "3000",
	}

	resp, err := client.Applications.CreatePublic(context.Background(), reqBody)
	if err != nil {
		t.Fatalf("Unexpected create error: %v", err)
	}

	if resp.UUID != "app-created-123" {
		t.Errorf("Expected UUID 'app-created-123', got %q", resp.UUID)
	}
}

func TestServers_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[
			{"id":1,"uuid":"server-uuid-1","name":"Production VPS","ip":"1.2.3.4","user":"root","port":22,"swarm_cluster":false}
		]`))
	}))
	defer server.Close()

	client := NewClient(server.URL, "")
	servers, err := client.Servers.List(context.Background())
	if err != nil {
		t.Fatalf("Unexpected servers listing error: %v", err)
	}

	if len(servers) != 1 {
		t.Fatalf("Expected 1 server, got %d", len(servers))
	}

	if servers[0].UUID != "server-uuid-1" || servers[0].IP != "1.2.3.4" || servers[0].Port != 22 {
		t.Errorf("Unexpected server properties: %+v", servers[0])
	}
}

func TestDatabases_CreatePostgreSQL(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"uuid":"db-created-123"}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, "")
	reqBody := CreatePostgreSQLRequest{
		CommonDatabaseRequest: CommonDatabaseRequest{
			ServerUUID:  "srv-123",
			ProjectUUID: "proj-123",
		},
		PostgresUser: String("postgres"),
	}

	resp, err := client.Databases.CreatePostgreSQL(context.Background(), reqBody)
	if err != nil {
		t.Fatalf("Unexpected database creation error: %v", err)
	}

	if resp.UUID != "db-created-123" {
		t.Errorf("Expected db UUID 'db-created-123', got %q", resp.UUID)
	}
}

func TestSystem_VersionAndHealth(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/version") {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`"v4.0.0-beta.2"`))
		} else if strings.HasSuffix(r.URL.Path, "/health") {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`"healthy"`))
		}
	}))
	defer server.Close()

	client := NewClient(server.URL, "")

	version, err := client.System.Version(context.Background())
	if err != nil {
		t.Fatalf("Unexpected version query error: %v", err)
	}
	if version != "v4.0.0-beta.2" {
		t.Errorf("Expected version 'v4.0.0-beta.2', got %q", version)
	}

	health, err := client.System.Health(context.Background())
	if err != nil {
		t.Fatalf("Unexpected health query error: %v", err)
	}
	if health != "healthy" {
		t.Errorf("Expected health status 'healthy', got %q", health)
	}
}

func TestApplications_Storage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/storages") {
			var req CreateStorageRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				t.Errorf("Failed to decode CreateStorageRequest: %v", err)
			}
			if req.MountPath != "/" || *req.Name != "test-volume" {
				t.Errorf("Unexpected CreateStorageRequest: %+v", req)
			}
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(`{"uuid":"storage-123","name":"test-volume","mount_path":"/"}`))
		} else if r.Method == http.MethodPatch && strings.HasSuffix(r.URL.Path, "/storages") {
			var req UpdateStorageRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				t.Errorf("Failed to decode UpdateStorageRequest: %v", err)
			}
			if *req.MountPath != "/new" {
				t.Errorf("Unexpected UpdateStorageRequest: %+v", req)
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"uuid":"storage-123","name":"test-volume","mount_path":"/new"}`))
		} else if r.Method == http.MethodDelete && strings.HasSuffix(r.URL.Path, "/storages/storage-123") {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"message":"Storage deleted."}`))
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	client := NewClient(server.URL, "")
	ctx := context.Background()

	// 1. Test CreateStorage
	createResp, err := client.Applications.CreateStorage(ctx, "app-123", CreateStorageRequest{
		Type:      "persistent",
		Name:      String("test-volume"),
		MountPath: "/",
	})
	if err != nil {
		t.Fatalf("Unexpected CreateStorage error: %v", err)
	}
	if createResp["uuid"] != "storage-123" {
		t.Errorf("Expected UUID 'storage-123', got %v", createResp["uuid"])
	}

	// 2. Test UpdateStorage
	updateResp, err := client.Applications.UpdateStorage(ctx, "app-123", UpdateStorageRequest{
		Type:      "persistent",
		MountPath: String("/new"),
	})
	if err != nil {
		t.Fatalf("Unexpected UpdateStorage error: %v", err)
	}
	if updateResp["mount_path"] != "/new" {
		t.Errorf("Expected mount_path '/new', got %v", updateResp["mount_path"])
	}

	// 3. Test DeleteStorage
	msg, err := client.Applications.DeleteStorage(ctx, "app-123", "storage-123")
	if err != nil {
		t.Fatalf("Unexpected DeleteStorage error: %v", err)
	}
	if msg != "Storage deleted." {
		t.Errorf("Expected message 'Storage deleted.', got %q", msg)
	}
}

func TestIntegration_RealCoolify(t *testing.T) {
	loadEnv()

	url := os.Getenv("COOLIFY_URL")
	key := os.Getenv("COOLIFY_KEY")

	if url == "" || key == "" {
		t.Skip("Skipping integration test: COOLIFY_URL and COOLIFY_KEY are not configured in .env file")
	}

	t.Logf("Running real integration test against: %s", url)

	client := NewClient(url, key)
	t.Logf("Client baseURL: %s", client.baseURL.String())
	t.Logf("Client token length: %d", len(client.token))
	ctx := context.Background()

	// 1. Query health
	health, err := client.System.Health(ctx)
	if err != nil {
		t.Fatalf("Failed to query health: %v", err)
	}
	t.Logf("Instance Health Status: %s", health)

	// 2. Query version
	version, err := client.System.Version(ctx)
	if err != nil {
		t.Fatalf("Failed to query version: %v", err)
	}
	t.Logf("Instance Version: %s", version)

	// 3. List Servers
	servers, err := client.Servers.List(ctx)
	if err != nil {
		t.Fatalf("Failed to list servers: %v", err)
	}
	t.Logf("Registered Servers count: %d", len(servers))
	for _, s := range servers {
		t.Logf(" - Server: %s (IP: %s, User: %s)", s.Name, s.IP, s.User)
	}

	// 4. List Applications
	apps, err := client.Applications.List(ctx)
	if err != nil {
		t.Fatalf("Failed to list applications: %v", err)
	}
	t.Logf("Applications count: %d", len(apps))
	for _, app := range apps {
		t.Logf(" - Application: %s (UUID: %s, Status: %s)", app.Name, app.UUID, app.Status)
	}

	// 5. List Services
	svcs, err := client.Services.List(ctx)
	if err != nil {
		t.Fatalf("Failed to list services: %v", err)
	}
	t.Logf("Services count: %d", len(svcs))
	for _, s := range svcs {
		statusStr := "unknown"
		if s.Status != nil {
			statusStr = *s.Status
		}
		t.Logf(" - Service: %s (UUID: %s, Status: %s)", s.Name, s.UUID, statusStr)
	}
}
