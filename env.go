package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// EnvConfig represents the environment variables configuration
type EnvConfig struct {
	Version             string                       `json:"version"`
	ActiveEnvironment   string                       `json:"activeEnvironment"`
	Global              map[string]string            `json:"global"`
	Environments        map[string]map[string]string `json:"environments"`
}

// getEnvFilePath returns the path to the env.json file
func getEnvFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	hurlStudioDir := filepath.Join(homeDir, ".hurlstudio")

	// Create directory if it doesn't exist
	if err := os.MkdirAll(hurlStudioDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create .hurlstudio directory: %w", err)
	}

	return filepath.Join(hurlStudioDir, "env.json"), nil
}

// LoadEnvVariables loads the environment variables from env.json
// If the file doesn't exist, it creates one with default values
func (a *App) LoadEnvVariables() (string, error) {
	envFilePath, err := getEnvFilePath()
	if err != nil {
		return "", err
	}

	// Check if file exists
	if _, err := os.Stat(envFilePath); os.IsNotExist(err) {
		// Create default config
		defaultConfig := EnvConfig{
			Version:           "1.0",
			ActiveEnvironment: "development",
			Global: map[string]string{
				"timeout":     "30",
				"api_version": "v1",
			},
			Environments: map[string]map[string]string{
				"development": {
					"api_url": "http://localhost:3000",
				},
				"staging": {
					"api_url": "https://staging.example.com",
				},
				"production": {
					"api_url": "https://api.example.com",
				},
			},
		}

		// Write default config to file
		data, err := json.MarshalIndent(defaultConfig, "", "  ")
		if err != nil {
			return "", fmt.Errorf("failed to marshal default config: %w", err)
		}

		if err := os.WriteFile(envFilePath, data, 0644); err != nil {
			return "", fmt.Errorf("failed to write default env.json: %w", err)
		}

		return string(data), nil
	}

	// Read existing file
	data, err := os.ReadFile(envFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to read env.json: %w", err)
	}

	return string(data), nil
}

// loadEnvConfig loads and caches the environment config
func (a *App) loadEnvConfig() (*EnvConfig, error) {
	// Check cache first (read lock)
	a.envConfigMu.RLock()
	if a.envConfig != nil {
		defer a.envConfigMu.RUnlock()
		return a.envConfig, nil
	}
	a.envConfigMu.RUnlock()

	// Load config (write lock)
	a.envConfigMu.Lock()
	defer a.envConfigMu.Unlock()

	// Double-check after acquiring write lock
	if a.envConfig != nil {
		return a.envConfig, nil
	}

	// Load from file
	configJSON, err := a.LoadEnvVariables()
	if err != nil {
		return nil, err
	}

	var config EnvConfig
	if err := json.Unmarshal([]byte(configJSON), &config); err != nil {
		return nil, fmt.Errorf("failed to parse env.json: %w", err)
	}

	a.envConfig = &config
	return a.envConfig, nil
}

// invalidateEnvCache clears the cached environment config
func (a *App) invalidateEnvCache() {
	a.envConfigMu.Lock()
	defer a.envConfigMu.Unlock()
	a.envConfig = nil
}

// SaveEnvVariables saves the environment variables to env.json
func (a *App) SaveEnvVariables(content string) error {
	// Validate JSON structure first
	var config EnvConfig
	if err := json.Unmarshal([]byte(content), &config); err != nil {
		return fmt.Errorf("invalid JSON format: %w", err)
	}

	// Validate required fields
	if config.Version == "" {
		return fmt.Errorf("version field is required")
	}
	if config.ActiveEnvironment == "" {
		return fmt.Errorf("activeEnvironment field is required")
	}
	if config.Global == nil {
		config.Global = make(map[string]string)
	}
	if config.Environments == nil {
		return fmt.Errorf("environments field is required")
	}

	envFilePath, err := getEnvFilePath()
	if err != nil {
		return err
	}

	// Format the JSON nicely before saving
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(envFilePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write env.json: %w", err)
	}

	// Invalidate cache after save
	a.invalidateEnvCache()

	return nil
}

// GetFlattenedVariables returns a flattened map of variables for the given environment
// Priority: environment variables override global variables
func (a *App) GetFlattenedVariables(environment string) (map[string]string, error) {
	// Load config from cache
	config, err := a.loadEnvConfig()
	if err != nil {
		return nil, err
	}

	// Start with global variables
	result := make(map[string]string)
	for k, v := range config.Global {
		result[k] = v
	}

	// Override with environment-specific variables
	if envVars, ok := config.Environments[environment]; ok {
		for k, v := range envVars {
			result[k] = v
		}
	}

	return result, nil
}

// GetActiveEnvironment returns the currently active environment
func (a *App) GetActiveEnvironment() (string, error) {
	// Load config from cache
	config, err := a.loadEnvConfig()
	if err != nil {
		return "", err
	}

	return config.ActiveEnvironment, nil
}
