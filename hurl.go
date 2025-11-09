package main

import (
	"embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

//go:embed binaries/*
var binaries embed.FS

// GetHurlPath returns the path to the hurl binary
// It extracts the embedded binary to a temp location if needed
func GetHurlPath() (string, error) {
	// Determine the binary name based on OS
	binaryName := "hurl"
	if runtime.GOOS == "windows" {
		binaryName = "hurl.exe"
	}

	// Check if hurl is already in PATH
	if path, err := exec.LookPath("hurl"); err == nil {
		return path, nil
	}

	// Determine the embedded binary path based on OS and architecture
	embeddedPath := fmt.Sprintf("binaries/%s-%s/%s", runtime.GOOS, runtime.GOARCH, binaryName)

	// Read the embedded binary
	data, err := binaries.ReadFile(embeddedPath)
	if err != nil {
		return "", fmt.Errorf("hurl binary not found for %s-%s: %w", runtime.GOOS, runtime.GOARCH, err)
	}

	// Create a temp directory for the binary
	tempDir := filepath.Join(os.TempDir(), "hurlstudio")
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create temp directory: %w", err)
	}

	// Write the binary to temp location
	binaryPath := filepath.Join(tempDir, binaryName)

	// Only write if file doesn't exist or is different
	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		if err := os.WriteFile(binaryPath, data, 0755); err != nil {
			return "", fmt.Errorf("failed to write binary: %w", err)
		}
	}

	return binaryPath, nil
}

// RunHurl executes a hurl file and returns the JSON report
// Generates a JSON report in /tmp/hurlstudio/<full-file-path>/
func (a *App) RunHurl(filePath string) (string, error) {
	hurlPath, err := GetHurlPath()
	if err != nil {
		return "", err
	}

	// Create report directory preserving the full file path
	// e.g., /tmp/hurlstudio/Users/sumeet/Projects/myproject/api.hurl/
	reportDir := filepath.Join(os.TempDir(), "hurlstudio", filePath)
	fmt.Println("Report directory:", reportDir)

	// Remove and recreate this specific file's report directory to ensure clean state
	os.RemoveAll(reportDir)
	if err := os.MkdirAll(reportDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create report directory: %w", err)
	}

	// Run hurl with JSON report generation
	cmd := exec.Command(hurlPath, "--report-json", reportDir, filePath)
	output, _ := cmd.CombinedOutput()

	// Read the generated JSON report
	// Hurl generates report with the filename: report-<timestamp>.json
	// We'll read the most recent file
	reportFiles, err := filepath.Glob(filepath.Join(reportDir, "*.json"))
	if err != nil || len(reportFiles) == 0 {
		// If no JSON report found, return the stdout output
		return string(output), nil
	}

	// Read the most recent report file (last one in sorted list)
	reportFile := reportFiles[len(reportFiles)-1]
	jsonContent, err := os.ReadFile(reportFile)
	if err != nil {
		return string(output), nil
	}

	return string(jsonContent), nil
}

// RunHurlWithOptions executes a hurl file with custom options
func (a *App) RunHurlWithOptions(filePath string, options []string) (string, error) {
	hurlPath, err := GetHurlPath()
	if err != nil {
		return "", err
	}

	args := append(options, filePath)
	cmd := exec.Command(hurlPath, args...)
	output, err := cmd.CombinedOutput()

	return string(output), err
}

// GetExistingReport checks if a report already exists for the given file path
// Returns the report JSON content if found, empty string if not found
func (a *App) GetExistingReport(filePath string) (string, error) {
	reportDir := filepath.Join(os.TempDir(), "hurlstudio", filePath)

	// Check if report directory exists
	if _, err := os.Stat(reportDir); os.IsNotExist(err) {
		return "", nil // No report exists yet
	}

	// Look for JSON report files
	reportFiles, err := filepath.Glob(filepath.Join(reportDir, "*.json"))
	if err != nil || len(reportFiles) == 0 {
		return "", nil // No report found
	}

	// Read the most recent report file
	reportFile := reportFiles[len(reportFiles)-1]
	jsonContent, err := os.ReadFile(reportFile)
	if err != nil {
		return "", nil // Could not read report
	}

	return string(jsonContent), nil
}

// GetResponseBody reads the response body from the report directory
// bodyPath is the relative path from the report (e.g., "store/response_1.txt")
func (a *App) GetResponseBody(hurlFilePath string, bodyPath string) (string, error) {
	reportDir := filepath.Join(os.TempDir(), "hurlstudio", hurlFilePath)
	bodyFilePath := filepath.Join(reportDir, bodyPath)

	// Read the body file
	content, err := os.ReadFile(bodyFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to read body file: %w", err)
	}

	return string(content), nil
}
