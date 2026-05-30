package coolify

import (
	"context"
	"fmt"
	"net/http"
)

// DatabasesService handles communication with the database-related endpoints of the Coolify API.
type DatabasesService struct {
	client *Client
}

// CommonDatabaseRequest contains common properties shared among all database creations.
type CommonDatabaseRequest struct {
	ServerUUID              string  `json:"server_uuid"`
	ProjectUUID             string  `json:"project_uuid"`
	EnvironmentName         *string `json:"environment_name,omitempty"`
	EnvironmentUUID         *string `json:"environment_uuid,omitempty"`
	DestinationUUID         *string `json:"destination_uuid,omitempty"`
	Name                    *string `json:"name,omitempty"`
	Description             *string `json:"description,omitempty"`
	Image                   *string `json:"image,omitempty"`
	IsPublic                *bool   `json:"is_public,omitempty"`
	PublicPort              *int    `json:"public_port,omitempty"`
	PublicPortTimeout       *int    `json:"public_port_timeout,omitempty"`
	LimitsMemory            *string `json:"limits_memory,omitempty"`
	LimitsMemorySwap        *string `json:"limits_memory_swap,omitempty"`
	LimitsMemorySwappiness  *int    `json:"limits_memory_swappiness,omitempty"`
	LimitsMemoryReservation *string `json:"limits_memory_reservation,omitempty"`
	LimitsCPUs              *string `json:"limits_cpus,omitempty"`
	LimitsCPUSet            *string `json:"limits_cpuset,omitempty"`
	LimitsCPUShares         *int    `json:"limits_cpu_shares,omitempty"`
	InstantDeploy           *bool   `json:"instant_deploy,omitempty"`
}

// CreatePostgreSQLRequest represents the request body for creating a PostgreSQL database.
type CreatePostgreSQLRequest struct {
	CommonDatabaseRequest
	PostgresUser           *string `json:"postgres_user,omitempty"`
	PostgresPassword       *string `json:"postgres_password,omitempty"`
	PostgresDB             *string `json:"postgres_db,omitempty"`
	PostgresInitDBArgs     *string `json:"postgres_initdb_args,omitempty"`
	PostgresHostAuthMethod *string `json:"postgres_host_auth_method,omitempty"`
	PostgresConf           *string `json:"postgres_conf,omitempty"`
}

// CreateRedisRequest represents the request body for creating a Redis database.
type CreateRedisRequest struct {
	CommonDatabaseRequest
	RedisPassword *string `json:"redis_password,omitempty"`
	RedisConf     *string `json:"redis_conf,omitempty"`
}

// CreateClickhouseRequest represents the request body for Clickhouse.
type CreateClickhouseRequest struct {
	CommonDatabaseRequest
	ClickhouseUser     *string `json:"clickhouse_user,omitempty"`
	ClickhousePassword *string `json:"clickhouse_password,omitempty"`
}

// CreateDragonflyRequest represents the request body for Dragonfly.
type CreateDragonflyRequest struct {
	CommonDatabaseRequest
	DragonflyPassword *string `json:"dragonfly_password,omitempty"`
}

// CreateKeyDBRequest represents the request body for KeyDB.
type CreateKeyDBRequest struct {
	CommonDatabaseRequest
	KeyDBPassword *string `json:"keydb_password,omitempty"`
}

// CreateMariaDBRequest represents the request body for MariaDB.
type CreateMariaDBRequest struct {
	CommonDatabaseRequest
	MariadbUser     *string `json:"mariadb_user,omitempty"`
	MariadbPassword *string `json:"mariadb_password,omitempty"`
	MariadbDatabase *string `json:"mariadb_database,omitempty"`
}

// CreateMySQLRequest represents the request body for MySQL.
type CreateMySQLRequest struct {
	CommonDatabaseRequest
	MysqlUser     *string `json:"mysql_user,omitempty"`
	MysqlPassword *string `json:"mysql_password,omitempty"`
	MysqlDatabase *string `json:"mysql_database,omitempty"`
}

// CreateMongoDBRequest represents the request body for MongoDB.
type CreateMongoDBRequest struct {
	CommonDatabaseRequest
	MongodbUser     *string `json:"mongodb_user,omitempty"`
	MongodbPassword *string `json:"mongodb_password,omitempty"`
	MongodbDatabase *string `json:"mongodb_database,omitempty"`
}

// List retrieves all databases.
func (s *DatabasesService) List(ctx context.Context) ([]Database, error) {
	req, err := s.client.newRequest(http.MethodGet, "databases", nil)
	if err != nil {
		return nil, err
	}

	var dbs []Database
	_, err = s.client.do(ctx, req, &dbs)
	return dbs, err
}

// Get retrieves detailed database configurations.
func (s *DatabasesService) Get(ctx context.Context, uuid string) (*Database, error) {
	path := fmt.Sprintf("databases/%s", uuid)
	req, err := s.client.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	db := new(Database)
	_, err = s.client.do(ctx, req, db)
	return db, err
}

// Delete removes a database from Coolify.
func (s *DatabasesService) Delete(ctx context.Context, uuid string) error {
	path := fmt.Sprintf("databases/%s", uuid)
	req, err := s.client.newRequest(http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	_, err = s.client.do(ctx, req, nil)
	return err
}

// CreatePostgreSQL registers and deploys a new PostgreSQL instance.
func (s *DatabasesService) CreatePostgreSQL(ctx context.Context, reqBody CreatePostgreSQLRequest) (*CreateResponse, error) {
	req, err := s.client.newRequest(http.MethodPost, "databases/postgresql", reqBody)
	if err != nil {
		return nil, err
	}

	resp := new(CreateResponse)
	_, err = s.client.do(ctx, req, resp)
	return resp, err
}

// CreateRedis registers and deploys a new Redis instance.
func (s *DatabasesService) CreateRedis(ctx context.Context, reqBody CreateRedisRequest) (*CreateResponse, error) {
	req, err := s.client.newRequest(http.MethodPost, "databases/redis", reqBody)
	if err != nil {
		return nil, err
	}

	resp := new(CreateResponse)
	_, err = s.client.do(ctx, req, resp)
	return resp, err
}

// CreateClickhouse registers and deploys a new Clickhouse instance.
func (s *DatabasesService) CreateClickhouse(ctx context.Context, reqBody CreateClickhouseRequest) (*CreateResponse, error) {
	req, err := s.client.newRequest(http.MethodPost, "databases/clickhouse", reqBody)
	if err != nil {
		return nil, err
	}

	resp := new(CreateResponse)
	_, err = s.client.do(ctx, req, resp)
	return resp, err
}

// CreateDragonfly registers and deploys a new Dragonfly instance.
func (s *DatabasesService) CreateDragonfly(ctx context.Context, reqBody CreateDragonflyRequest) (*CreateResponse, error) {
	req, err := s.client.newRequest(http.MethodPost, "databases/dragonfly", reqBody)
	if err != nil {
		return nil, err
	}

	resp := new(CreateResponse)
	_, err = s.client.do(ctx, req, resp)
	return resp, err
}

// CreateKeyDB registers and deploys a new KeyDB instance.
func (s *DatabasesService) CreateKeyDB(ctx context.Context, reqBody CreateKeyDBRequest) (*CreateResponse, error) {
	req, err := s.client.newRequest(http.MethodPost, "databases/keydb", reqBody)
	if err != nil {
		return nil, err
	}

	resp := new(CreateResponse)
	_, err = s.client.do(ctx, req, resp)
	return resp, err
}

// CreateMariaDB registers and deploys a new MariaDB instance.
func (s *DatabasesService) CreateMariaDB(ctx context.Context, reqBody CreateMariaDBRequest) (*CreateResponse, error) {
	req, err := s.client.newRequest(http.MethodPost, "databases/mariadb", reqBody)
	if err != nil {
		return nil, err
	}

	resp := new(CreateResponse)
	_, err = s.client.do(ctx, req, resp)
	return resp, err
}

// CreateMySQL registers and deploys a new MySQL instance.
func (s *DatabasesService) CreateMySQL(ctx context.Context, reqBody CreateMySQLRequest) (*CreateResponse, error) {
	req, err := s.client.newRequest(http.MethodPost, "databases/mysql", reqBody)
	if err != nil {
		return nil, err
	}

	resp := new(CreateResponse)
	_, err = s.client.do(ctx, req, resp)
	return resp, err
}

// CreateMongoDB registers and deploys a new MongoDB instance.
func (s *DatabasesService) CreateMongoDB(ctx context.Context, reqBody CreateMongoDBRequest) (*CreateResponse, error) {
	req, err := s.client.newRequest(http.MethodPost, "databases/mongodb", reqBody)
	if err != nil {
		return nil, err
	}

	resp := new(CreateResponse)
	_, err = s.client.do(ctx, req, resp)
	return resp, err
}

// Start starts a stopped database container.
func (s *DatabasesService) Start(ctx context.Context, uuid string) error {
	path := fmt.Sprintf("databases/%s/start", uuid)
	req, err := s.client.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return err
	}

	_, err = s.client.do(ctx, req, nil)
	return err
}

// Stop stops an active database container.
func (s *DatabasesService) Stop(ctx context.Context, uuid string) error {
	path := fmt.Sprintf("databases/%s/stop", uuid)
	req, err := s.client.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return err
	}

	_, err = s.client.do(ctx, req, nil)
	return err
}

// Restart restarts the database.
func (s *DatabasesService) Restart(ctx context.Context, uuid string) error {
	path := fmt.Sprintf("databases/%s/restart", uuid)
	req, err := s.client.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return err
	}

	_, err = s.client.do(ctx, req, nil)
	return err
}
