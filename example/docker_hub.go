package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"time"
)

// Queries the Docker Hub registry for semantic version tags and returns the latest version.
func getLatestDockerHubTag(repo string) (string, error) {
	url := fmt.Sprintf("https://hub.docker.com/v2/repositories/%s/tags/?page_size=100", repo)
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP error: %s", resp.Status)
	}

	var data struct {
		Results []struct {
			Name string `json:"name"`
		} `json:"results"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}

	re := regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+$`)
	var versions []string
	for _, result := range data.Results {
		if re.MatchString(result.Name) {
			versions = append(versions, result.Name)
		}
	}

	if len(versions) == 0 {
		return "", fmt.Errorf("no semantic version tags found")
	}

	// Sort version strings semantically
	sort.Slice(versions, func(i, j int) bool {
		return compareVersions(versions[i], versions[j]) < 0
	})

	// The last one is the latest version
	return versions[len(versions)-1], nil
}

// Helper to compare two 3-part semantic version strings
func compareVersions(v1, v2 string) int {
	var maj1, min1, pat1 int
	var maj2, min2, pat2 int
	fmt.Sscanf(v1, "%d.%d.%d", &maj1, &min1, &pat1)
	fmt.Sscanf(v2, "%d.%d.%d", &maj2, &min2, &pat2)

	if maj1 != maj2 {
		if maj1 < maj2 {
			return -1
		}
		return 1
	}
	if min1 != min2 {
		if min1 < min2 {
			return -1
		}
		return 1
	}
	if pat1 != pat2 {
		if pat1 < pat2 {
			return -1
		}
		return 1
	}
	return 0
}
