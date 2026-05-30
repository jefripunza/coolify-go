package main

import (
	"os"
	"strings"
)

// parses the .env file manually from the local directory or its parent directory.
func init() {
	// Try loading from the current working directory first
	content, err := os.ReadFile(".env")
	if err != nil {
		// Fallback to parent directory (useful if running from inside the example/ folder)
		content, err = os.ReadFile("../.env")
		if err != nil {
			return
		}
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
			os.Setenv(key, val)
		}
	}
}
