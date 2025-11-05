package main

import (
	"embed"
	"fmt"
	"io"
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

// RunHurl executes a hurl file and returns the output
func (a *App) RunHurl(filePath string) (string, error) {
	hurlPath, err := GetHurlPath()
	if err != nil {
		return "", err
	}

	cmd := exec.Command(hurlPath, filePath)
	output, err := cmd.CombinedOutput()

	return string(output), err
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
