package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"coolify"
)

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func main() {
	coolify_url := os.Getenv("COOLIFY_URL")
	coolify_key := os.Getenv("COOLIFY_KEY")
	web_host := os.Getenv("WEB_HOST")
	docker_repo := os.Getenv("DOCKER_REPO")
	project_uuid := os.Getenv("PROJECT_UUID")
	server_uuid := os.Getenv("SERVER_UUID")
	environment_uuid := os.Getenv("ENVIRONMENT_UUID")

	if coolify_url == "" || coolify_key == "" || web_host == "" || docker_repo == "" || project_uuid == "" || server_uuid == "" || environment_uuid == "" {
		log.Fatalf("Error: COOLIFY_URL or COOLIFY_KEY or WEB_HOST or DOCKER_REPO or PROJECT_UUID or SERVER_UUID or ENVIRONMENT_UUID is not configured in .env file or system environment")
	}

	fmt.Printf("Initializing Coolify Client targeting: %s\n\n", coolify_url)

	// Initialize the Coolify API client with credentials from .env.
	client := coolify.NewClient(coolify_url, coolify_key)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	// 1. Query System version and status
	version, err := client.System.Version(ctx)
	if err != nil {
		log.Printf("Could not check version: %v\n", err)
	} else {
		fmt.Printf("1. Connected successfully! Coolify Version: %s\n\n", version)
	}

	// 2. List Applications
	fmt.Println("2. Listing applications...")
	apps, err := client.Applications.List(ctx)
	if err != nil {
		if apiErr, ok := err.(*coolify.Error); ok {
			fmt.Printf("   API Error Listing Apps: Code=%d Message=%q\n", apiErr.StatusCode, apiErr.Message)
		} else {
			log.Printf("   Failed to list applications: %v\n", err)
		}
	} else {
		fmt.Printf("   Retrieved %d applications:\n", len(apps))
		for _, app := range apps {
			fmt.Printf("   - %s (UUID: %s, Status: %s)\n", app.Name, app.UUID, app.Status)
		}
	}
	fmt.Println()

	// 3. List Services
	fmt.Println("3. Listing one-click services...")
	svcs, err := client.Services.List(ctx)
	if err != nil {
		if apiErr, ok := err.(*coolify.Error); ok {
			fmt.Printf("   API Error Listing Services: Code=%d Message=%q\n", apiErr.StatusCode, apiErr.Message)
		} else {
			log.Printf("   Failed to list services: %v\n", err)
		}
	} else {
		fmt.Printf("   Retrieved %d services:\n", len(svcs))
		for _, s := range svcs {
			statusStr := "unknown"
			if s.Status != nil {
				statusStr = *s.Status
			}
			fmt.Printf("   - %s (UUID: %s, Status: %s)\n", s.Name, s.UUID, statusStr)
		}
	}
	fmt.Println()

	// 4. Demonstrating custom Docker Image application creation
	timestamp := time.Now().Unix()
	appName := fmt.Sprintf("tester-%d", timestamp)

	domain_terminal := fmt.Sprintf("https://vps-terminal-%d.%s", timestamp, web_host)
	domain_app := fmt.Sprintf("https://vps-app-%d.%s:8080", timestamp, web_host)

	// join string using comma
	domain := strings.Join([]string{domain_terminal, domain_app}, ",")

	// Query the Docker Hub repository for the latest version tag.
	latestVersion := "latest" // Fallback default version
	fmt.Println("4a. Querying Docker Hub for the latest vps-ubuntu-server tag...")
	fetchedVersion, fetchErr := getLatestDockerHubTag(docker_repo)
	if fetchErr != nil {
		fmt.Printf("   Warning: Could not fetch latest version from Docker Hub: %v. Using fallback: %s\n", fetchErr, latestVersion)
	} else {
		latestVersion = fetchedVersion
		fmt.Printf("   Successfully retrieved latest version tag: %s\n", latestVersion)
	}

	newAppReq := coolify.CreateDockerImageRequest{
		ProjectUUID:             project_uuid,
		ServerUUID:              server_uuid,
		EnvironmentName:         "production",
		EnvironmentUUID:         environment_uuid,
		DockerRegistryImageName: docker_repo,
		DockerRegistryImageTag:  latestVersion,
		Name:                    coolify.String(appName),
		Domains:                 coolify.String(domain),
		PortsExposes:            coolify.String("6080"),
		InstantDeploy:           coolify.Bool(false),
	}

	fmt.Printf("4. Attempting VPS deployment via Docker Image: Name=%q, Domain=%q...\n", appName, domain)
	resp, err := client.Applications.CreateDockerImage(ctx, newAppReq)
	if err != nil {
		if apiErr, ok := err.(*coolify.Error); ok {
			fmt.Printf("   Create Failed (expected for mock parameters): Status Code = %d\n", apiErr.StatusCode)
			fmt.Printf("   Reason: %s\n", apiErr.Message)

			// Inspect validation errors if any (422)
			if len(apiErr.Errors) > 0 {
				fmt.Println("   Validation Errors details:")
				for field, messages := range apiErr.Errors {
					fmt.Printf("     - Field %q: %v\n", field, messages)
				}
			}

			// Inspect domain routing conflicts if any (409)
			if len(apiErr.Conflicts) > 0 {
				fmt.Printf("   Warning: %s\n", apiErr.Warning)
				fmt.Println("   Detected domain conflicts:")
				for _, conflict := range apiErr.Conflicts {
					fmt.Printf("     - Domain %q is already occupied by resource %q (%s)\n",
						conflict.Domain, conflict.ResourceName, conflict.ResourceType)
				}
			}
		} else {
			log.Printf("   Unexpected error: %v\n", err)
		}
	} else {
		fmt.Printf("   Application registered successfully! Assigned UUID: %s\n", resp.UUID)

		ttydUser := getEnv("TTYD_USER", "ubuntu")
		ttydPassword := getEnv("TTYD_PASSWORD", "ubuntu")
		sshUser := getEnv("SSH_USER", "ubuntu")
		sshPassword := getEnv("SSH_PASSWORD", "ubuntu")
		sshHostname := getEnv("SSH_HOSTNAME", "server")
		sawangApiKey := getEnv("SAWANG_CLOUD_API_KEY", "-")
		sawangBaseUrl := getEnv("SAWANG_CLOUD_BASE_URL", "http://localhost:20128/v1")

		// Create the requested environment variables
		fmt.Println("\n   Creating application environment variables...")
		envVariables := []coolify.CreateEnvRequest{
			{Key: "TTYD_USER", Value: ttydUser, IsLiteral: true},
			{Key: "TTYD_PASSWORD", Value: ttydPassword, IsLiteral: true},
			{Key: "SSH_USER", Value: sshUser, IsLiteral: true},
			{Key: "SSH_PASSWORD", Value: sshPassword, IsLiteral: true},
			{Key: "SSH_HOSTNAME", Value: sshHostname, IsLiteral: true},
			{Key: "SAWANG_CLOUD_API_KEY", Value: sawangApiKey, IsLiteral: true},
			{Key: "SAWANG_CLOUD_BASE_URL", Value: sawangBaseUrl, IsLiteral: true},
		}

		for _, envReq := range envVariables {
			envResp, envErr := client.Applications.CreateEnv(ctx, resp.UUID, envReq)
			if envErr != nil {
				fmt.Printf("   Warning: Failed to create env variable %s: %v\n", envReq.Key, envErr)
			} else {
				fmt.Printf("   - Environment variable %s created successfully (UUID: %s)\n", envReq.Key, envResp.UUID)
			}
		}

		// Deploy the application manually
		fmt.Printf("\n5. Deploying application %s...\n", resp.UUID)
		deployResp, deployErr := client.Deployments.Deploy(ctx, true, resp.UUID, nil)
		if deployErr != nil {
			log.Printf("   Failed to trigger deployment: %v\n", deployErr)
		} else {
			fmt.Printf("   Deployment triggered successfully! Message: %s\n", deployResp.Message)
		}

		// // 5. Stop the application after a 20-second delay
		// fmt.Println("\n5. Delaying 20 seconds before stopping the application...")
		// time.Sleep(20 * time.Second)

		// fmt.Printf("   Stopping application %s...\n", resp.UUID)
		// stopCtx, stopCancel := context.WithTimeout(context.Background(), 20*time.Second)
		// defer stopCancel()
		// err = client.Applications.Stop(stopCtx, resp.UUID)
		// if err != nil {
		// 	log.Printf("   Failed to stop application: %v\n", err)
		// } else {
		// 	fmt.Println("   Application stopped successfully!")
		// }

		// // 6. Update SSH_PASSWORD env to "ubuntu-ok"
		// fmt.Println("\n6. Updating SSH_PASSWORD to \"ubuntu-ok\"...")
		// updatedEnv, updateErr := client.Applications.UpdateEnv(ctx, resp.UUID, coolify.UpdateEnvRequest{
		// 	Key:       "SSH_PASSWORD",
		// 	Value:     "ubuntu-ok",
		// 	IsLiteral: true,
		// })
		// if updateErr != nil {
		// 	log.Printf("   Failed to update SSH_PASSWORD: %v\n", updateErr)
		// } else {
		// 	fmt.Printf("   SSH_PASSWORD updated successfully (UUID: %s, new value: ubuntu-ok)\n", updatedEnv.UUID)
		// }
	}
}
