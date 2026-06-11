package coolify

import "time"

// Application represents a Coolify application resource.
type Application struct {
	ID                              int        `json:"id"`
	UUID                            string     `json:"uuid"`
	Name                            string     `json:"name"`
	Description                     *string    `json:"description,omitempty"`
	RepositoryProjectID             *int       `json:"repository_project_id,omitempty"`
	FQDN                            *string    `json:"fqdn,omitempty"`
	ConfigHash                      string     `json:"config_hash"`
	GitRepository                   string     `json:"git_repository"`
	GitBranch                       string     `json:"git_branch"`
	GitCommitSHA                    string     `json:"git_commit_sha"`
	GitFullURL                      string     `json:"git_full_url"`
	DockerRegistryImageName         *string    `json:"docker_registry_image_name,omitempty"`
	DockerRegistryImageTag          *string    `json:"docker_registry_image_tag,omitempty"`
	BuildPack                       string     `json:"build_pack"`
	StaticImage                     *string    `json:"static_image,omitempty"`
	InstallCommand                  *string    `json:"install_command,omitempty"`
	BuildCommand                    *string    `json:"build_command,omitempty"`
	StartCommand                    *string    `json:"start_command,omitempty"`
	PortsExposes                    string     `json:"ports_exposes"`
	PortsMappings                   *string    `json:"ports_mappings,omitempty"`
	CustomNetworkAliases            *string    `json:"custom_network_aliases,omitempty"`
	BaseDirectory                   string     `json:"base_directory"`
	PublishDirectory                *string    `json:"publish_directory,omitempty"`
	HealthCheckEnabled              bool       `json:"health_check_enabled"`
	HealthCheckPath                 string     `json:"health_check_path"`
	HealthCheckPort                 *string    `json:"health_check_port,omitempty"`
	HealthCheckHost                 *string    `json:"health_check_host,omitempty"`
	HealthCheckMethod               string     `json:"health_check_method"`
	HealthCheckReturnCode           int        `json:"health_check_return_code"`
	HealthCheckScheme               string     `json:"health_check_scheme"`
	HealthCheckResponseText         *string    `json:"health_check_response_text,omitempty"`
	HealthCheckInterval             int        `json:"health_check_interval"`
	HealthCheckTimeout              int        `json:"health_check_timeout"`
	HealthCheckRetries              int        `json:"health_check_retries"`
	HealthCheckStartPeriod          int        `json:"health_check_start_period"`
	HealthCheckType                 *string    `json:"health_check_type,omitempty"`
	HealthCheckCommand              *string    `json:"health_check_command,omitempty"`
	LimitsMemory                    string     `json:"limits_memory"`
	LimitsMemorySwap                string     `json:"limits_memory_swap"`
	LimitsMemorySwappiness          int        `json:"limits_memory_swappiness"`
	LimitsMemoryReservation         string     `json:"limits_memory_reservation"`
	LimitsCPUs                      string     `json:"limits_cpus"`
	LimitsCPUSet                    *string    `json:"limits_cpuset,omitempty"`
	LimitsCPUShares                 int        `json:"limits_cpu_shares"`
	Status                          string     `json:"status"`
	PreviewURLTemplate              *string    `json:"preview_url_template,omitempty"`
	DestinationType                 *string    `json:"destination_type,omitempty"`
	DestinationID                   *int       `json:"destination_id,omitempty"`
	SourceID                        *int       `json:"source_id,omitempty"`
	PrivateKeyID                    *int       `json:"private_key_id,omitempty"`
	EnvironmentID                   int        `json:"environment_id"`
	Dockerfile                      *string    `json:"dockerfile,omitempty"`
	DockerfileLocation              string     `json:"dockerfile_location"`
	CustomLabels                    *string    `json:"custom_labels,omitempty"`
	DockerfileTargetBuild           *string    `json:"dockerfile_target_build,omitempty"`
	ManualWebhookSecretGitHub       *string    `json:"manual_webhook_secret_github,omitempty"`
	ManualWebhookSecretGitLab       *string    `json:"manual_webhook_secret_gitlab,omitempty"`
	ManualWebhookSecretBitbucket    *string    `json:"manual_webhook_secret_bitbucket,omitempty"`
	ManualWebhookSecretGitea        *string    `json:"manual_webhook_secret_gitea,omitempty"`
	DockerComposeLocation           string     `json:"docker_compose_location"`
	DockerCompose                   *string    `json:"docker_compose,omitempty"`
	DockerComposeRaw                *string    `json:"docker_compose_raw,omitempty"`
	DockerComposeDomains            *string    `json:"docker_compose_domains,omitempty"`
	DockerComposeCustomStartCommand *string    `json:"docker_compose_custom_start_command,omitempty"`
	DockerComposeCustomBuildCommand *string    `json:"docker_compose_custom_build_command,omitempty"`
	SwarmReplicas                   *int       `json:"swarm_replicas,omitempty"`
	SwarmPlacementConstraints       *string    `json:"swarm_placement_constraints,omitempty"`
	CustomDockerRunOptions          *string    `json:"custom_docker_run_options,omitempty"`
	PostDeploymentCommand           *string    `json:"post_deployment_command,omitempty"`
	PostDeploymentCommandContainer  *string    `json:"post_deployment_command_container,omitempty"`
	PreDeploymentCommand            *string    `json:"pre_deployment_command,omitempty"`
	PreDeploymentCommandContainer   *string    `json:"pre_deployment_command_container,omitempty"`
	WatchPaths                      *string    `json:"watch_paths,omitempty"`
	CustomHealthcheckFound          bool       `json:"custom_healthcheck_found"`
	Redirect                        *string    `json:"redirect,omitempty"`
	CreatedAt                       time.Time  `json:"created_at"`
	UpdatedAt                       time.Time  `json:"updated_at"`
	DeletedAt                       *time.Time `json:"deleted_at,omitempty"`
	ComposeParsingVersion           *string    `json:"compose_parsing_version,omitempty"`
	CustomNginxConfiguration        *string    `json:"custom_nginx_configuration,omitempty"`
	IsHTTPBasicAuthEnabled          bool       `json:"is_http_basic_auth_enabled"`
	HTTPBasicAuthUsername           *string    `json:"http_basic_auth_username,omitempty"`
	HTTPBasicAuthPassword           *string    `json:"http_basic_auth_password,omitempty"`
}

// ApplicationDeploymentQueue represents a deployment queue item.
type ApplicationDeploymentQueue struct {
	ID                     int       `json:"id"`
	ApplicationID          int       `json:"application_id"`
	DeploymentUUID         string    `json:"deployment_uuid"`
	PullRequestID          *int      `json:"pull_request_id,omitempty"`
	DockerRegistryImageTag string    `json:"docker_registry_image_tag"`
	ConfigurationHash      string    `json:"configuration_hash"`
	ConfigurationSnapshot  *string   `json:"configuration_snapshot,omitempty"`
	ConfigurationDiff      *string   `json:"configuration_diff,omitempty"`
	ForceRebuild           bool      `json:"force_rebuild"`
	Commit                 *string   `json:"commit,omitempty"`
	Status                 string    `json:"status"`
	IsWebhook              bool      `json:"is_webhook"`
	IsAPI                  bool      `json:"is_api"`
	CreatedAt              time.Time `json:"created_at"`
	UpdatedAt              time.Time `json:"updated_at"`
	Logs                   *string   `json:"logs,omitempty"`
	CurrentProcessID       *int      `json:"current_process_id,omitempty"`
	RestartOnly            bool      `json:"restart_only"`
	GitType                *string   `json:"git_type,omitempty"`
	ServerID               int       `json:"server_id"`
	ApplicationName        string    `json:"application_name"`
	ServerName             string    `json:"server_name"`
	DeploymentURL          *string   `json:"deployment_url,omitempty"`
	DestinationID          int       `json:"destination_id"`
	OnlyThisServer         bool      `json:"only_this_server"`
	Rollback               bool      `json:"rollback"`
	CommitMessage          *string   `json:"commit_message,omitempty"`
}

// Environment represents a Coolify environment resource.
type Environment struct {
	ID          int       `json:"id"`
	UUID        string    `json:"uuid"`
	Name        string    `json:"name"`
	ProjectID   int       `json:"project_id"`
	Description *string   `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// EnvironmentVariable represents an environment variable.
type EnvironmentVariable struct {
	ID               int       `json:"id"`
	UUID             string    `json:"uuid"`
	ResourceableType string    `json:"resourceable_type"`
	ResourceableID   int       `json:"resourceable_id"`
	IsLiteral        bool      `json:"is_literal"`
	IsMultiline      bool      `json:"is_multiline"`
	IsPreview        bool      `json:"is_preview"`
	IsRuntime        bool      `json:"is_runtime"`
	IsBuildtime      bool      `json:"is_buildtime"`
	IsShared         bool      `json:"is_shared"`
	IsShownOnce      bool      `json:"is_shown_once"`
	Key              string    `json:"key"`
	Value            string    `json:"value"`
	RealValue        *string   `json:"real_value,omitempty"`
	Comment          *string   `json:"comment,omitempty"`
	Version          *string   `json:"version,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// PrivateKey represents an SSH Private Key used in Coolify.
type PrivateKey struct {
	ID           int       `json:"id"`
	UUID         string    `json:"uuid"`
	Name         string    `json:"name"`
	Description  *string   `json:"description,omitempty"`
	PrivateKey   string    `json:"private_key"`
	PublicKey    *string   `json:"public_key,omitempty"`
	Fingerprint  *string   `json:"fingerprint,omitempty"`
	IsGitRelated bool      `json:"is_git_related"`
	TeamID       int       `json:"team_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Project represents a Coolify project.
type Project struct {
	ID           int           `json:"id"`
	UUID         string        `json:"uuid"`
	Name         string        `json:"name"`
	Description  *string       `json:"description,omitempty"`
	Environments []Environment `json:"environments,omitempty"`
}

// ScheduledTask represents a task scheduled on a Coolify service or application.
type ScheduledTask struct {
	ID        int       `json:"id"`
	UUID      string    `json:"uuid"`
	Enabled   bool      `json:"enabled"`
	Name      string    `json:"name"`
	Command   string    `json:"command"`
	Frequency string    `json:"frequency"`
	Container *string   `json:"container,omitempty"`
	Timeout   int       `json:"timeout"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ScheduledTaskExecution represents an execution history item of a scheduled task.
type ScheduledTaskExecution struct {
	UUID       string     `json:"uuid"`
	Status     string     `json:"status"`
	Message    *string    `json:"message,omitempty"`
	RetryCount int        `json:"retry_count"`
	Duration   int        `json:"duration"`
	StartedAt  time.Time  `json:"started_at"`
	FinishedAt *time.Time `json:"finished_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

// ServerProxy represents proxy settings on a server.
type ServerProxy struct {
	RedirectEnabled *bool `json:"redirect_enabled,omitempty"`
}

// Server represents a server resource managed by Coolify.
type Server struct {
	ID                            int            `json:"id"`
	UUID                          string         `json:"uuid"`
	Name                          string         `json:"name"`
	Description                   *string        `json:"description,omitempty"`
	IP                            string         `json:"ip"`
	User                          string         `json:"user"`
	Port                          int            `json:"port"`
	Proxy                         *ServerProxy   `json:"proxy,omitempty"`
	ProxyType                     *string        `json:"proxy_type,omitempty"`
	HighDiskUsageNotificationSent bool           `json:"high_disk_usage_notification_sent"`
	UnreachableNotificationSent   bool           `json:"unreachable_notification_sent"`
	UnreachableCount              int            `json:"unreachable_count"`
	ValidationLogs                *string        `json:"validation_logs,omitempty"`
	LogDrainNotificationSent      bool           `json:"log_drain_notification_sent"`
	SwarmCluster                  bool           `json:"swarm_cluster"`
	Settings                      *ServerSetting `json:"settings,omitempty"`
}

// ServerSetting represents the advanced configuration settings for a server.
type ServerSetting struct {
	ID                                int       `json:"id"`
	ConcurrentBuilds                  int       `json:"concurrent_builds"`
	DeploymentQueueLimit              int       `json:"deployment_queue_limit"`
	DynamicTimeout                    int       `json:"dynamic_timeout"`
	ForceDisabled                     bool      `json:"force_disabled"`
	ForceServerCleanup                bool      `json:"force_server_cleanup"`
	IsBuildServer                     bool      `json:"is_build_server"`
	IsCloudflareTunnel                bool      `json:"is_cloudflare_tunnel"`
	IsJumpServer                      bool      `json:"is_jump_server"`
	IsLogdrainAxiomEnabled            bool      `json:"is_logdrain_axiom_enabled"`
	IsLogdrainCustomEnabled           bool      `json:"is_logdrain_custom_enabled"`
	IsLogdrainHighlightEnabled        bool      `json:"is_logdrain_highlight_enabled"`
	IsLogdrainNewrelicEnabled         bool      `json:"is_logdrain_newrelic_enabled"`
	IsMetricsEnabled                  bool      `json:"is_metrics_enabled"`
	IsReachable                       bool      `json:"is_reachable"`
	IsSentinelEnabled                 bool      `json:"is_sentinel_enabled"`
	IsSwarmManager                    bool      `json:"is_swarm_manager"`
	IsSwarmWorker                     bool      `json:"is_swarm_worker"`
	IsTerminalEnabled                 bool      `json:"is_terminal_enabled"`
	IsUsable                          bool      `json:"is_usable"`
	LogdrainAxiomAPIKey               *string   `json:"logdrain_axiom_api_key,omitempty"`
	LogdrainAxiomDatasetName          *string   `json:"logdrain_axiom_dataset_name,omitempty"`
	LogdrainCustomConfig              *string   `json:"logdrain_custom_config,omitempty"`
	LogdrainCustomConfigParser        *string   `json:"logdrain_custom_config_parser,omitempty"`
	LogdrainHighlightProjectID        *string   `json:"logdrain_highlight_project_id,omitempty"`
	LogdrainNewrelicBaseURI           *string   `json:"logdrain_newrelic_base_uri,omitempty"`
	LogdrainNewrelicLicenseKey        *string   `json:"logdrain_newrelic_license_key,omitempty"`
	SentinelMetricsHistoryDays        int       `json:"sentinel_metrics_history_days"`
	SentinelMetricsRefreshRateSeconds int       `json:"sentinel_metrics_refresh_rate_seconds"`
	SentinelToken                     *string   `json:"sentinel_token,omitempty"`
	DockerCleanupFrequency            string    `json:"docker_cleanup_frequency"`
	DockerCleanupThreshold            int       `json:"docker_cleanup_threshold"`
	ServerID                          int       `json:"server_id"`
	WildcardDomain                    *string   `json:"wildcard_domain,omitempty"`
	CreatedAt                         time.Time `json:"created_at"`
	UpdatedAt                         time.Time `json:"updated_at"`
	DeleteUnusedVolumes               bool      `json:"delete_unused_volumes"`
	DeleteUnusedNetworks              bool      `json:"delete_unused_networks"`
	ConnectionTimeout                 int       `json:"connection_timeout"`
}

// Service represents a Coolify one-click service (like WordPress, Ghost, etc.).
type Service struct {
	ID                              int        `json:"id"`
	UUID                            string     `json:"uuid"`
	Name                            string     `json:"name"`
	EnvironmentID                   int        `json:"environment_id"`
	ServerID                        int        `json:"server_id"`
	Description                     *string    `json:"description,omitempty"`
	DockerComposeRaw                *string    `json:"docker_compose_raw,omitempty"`
	DockerCompose                   *string    `json:"docker_compose,omitempty"`
	DestinationType                 *string    `json:"destination_type,omitempty"`
	DestinationID                   int        `json:"destination_id"`
	ConnectToDockerNetwork          bool       `json:"connect_to_docker_network"`
	IsContainerLabelEscapeEnabled   bool       `json:"is_container_label_escape_enabled"`
	IsContainerLabelReadonlyEnabled bool       `json:"is_container_label_readonly_enabled"`
	ConfigHash                      *string    `json:"config_hash,omitempty"`
	ServiceType                     *string    `json:"service_type,omitempty"`
	ComposeParsingVersion           *string    `json:"compose_parsing_version,omitempty"`
	ServerStatus                    *bool      `json:"server_status,omitempty"`
	Status                          *string    `json:"status,omitempty"`
	CreatedAt                       time.Time  `json:"created_at"`
	UpdatedAt                       time.Time  `json:"updated_at"`
	DeletedAt                       *time.Time `json:"deleted_at,omitempty"`
}

// Team represents a tenant/team in Coolify.
type Team struct {
	ID                int       `json:"id"`
	Name              string    `json:"name"`
	Description       *string   `json:"description,omitempty"`
	PersonalTeam      bool      `json:"personal_team"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	ShowBoarding      bool      `json:"show_boarding"`
	CustomServerLimit *string   `json:"custom_server_limit,omitempty"`
	Members           []User    `json:"members,omitempty"`
}

// User represents a user account belonging to a team.
type User struct {
	ID                   int        `json:"id"`
	Name                 string     `json:"name"`
	Email                string     `json:"email"`
	EmailVerifiedAt      *time.Time `json:"email_verified_at,omitempty"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
	TwoFactorConfirmedAt *time.Time `json:"two_factor_confirmed_at,omitempty"`
	ForcePasswordReset   bool       `json:"force_password_reset"`
	MarketingEmails      bool       `json:"marketing_emails"`
}

// Database represents a Database resource.
type Database struct {
	ID            int        `json:"id"`
	UUID          string     `json:"uuid"`
	Name          string     `json:"name"`
	Description   *string    `json:"description,omitempty"`
	Status        string     `json:"status"`
	Type          string     `json:"type"` // e.g. postgresql, redis, mysql, etc.
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at,omitempty"`
	EnvironmentID int        `json:"environment_id"`
	ServerID      int        `json:"server_id"`
}

// CloudToken represents a cloud provider authentication token.
type CloudToken struct {
	UUID      string    `json:"uuid"`
	Name      string    `json:"name"`
	Provider  string    `json:"provider"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// GitHubApp represents GitHub App integration settings.
type GitHubApp struct {
	ID        int       `json:"id"`
	UUID      string    `json:"uuid"`
	Name      string    `json:"name"`
	APIURL    string    `json:"api_url"`
	HTMLURL   string    `json:"html_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
