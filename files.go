package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// FileEntry represents a file or directory in the file system
type FileEntry struct {
	Name        string    `json:"name"`
	Path        string    `json:"path"`
	IsDirectory bool      `json:"isDirectory"`
	Size        int64     `json:"size"`
	ModTime     time.Time `json:"modTime"`
}

// CurrentFilesState represents the state after opening a file or directory
type CurrentFilesState struct {
	CurrentDir  *FileEntry `json:"currentDir,omitempty"`
	CurrentFile *FileEntry `json:"currentFile,omitempty"`
}

// ListFiles lists all files and directories in the given path
// Filters to only show directories and .hurl files
func (a *App) ListFiles(path string) ([]FileEntry, error) {
	dirPath := path
	if dirPath == "" {
		dirPath = a.currentDir
	}

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %w", dirPath, err)
	}

	var fileEntries []FileEntry
	for _, entry := range entries {
		// Skip hidden files
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			fmt.Printf("Error getting info for %s: %v\n", entry.Name(), err)
			continue
		}

		// Only include directories, .hurl files, or markdown files
		if !entry.IsDir() {
			ext := strings.ToLower(filepath.Ext(entry.Name()))
			if ext != ".hurl" && ext != ".md" && ext != ".markdown" {
				continue
			}
		}

		fileEntries = append(fileEntries, FileEntry{
			Name:        entry.Name(),
			Path:        filepath.Join(dirPath, entry.Name()),
			IsDirectory: entry.IsDir(),
			Size:        info.Size(),
			ModTime:     info.ModTime(),
		})
	}
	return fileEntries, nil
}

// OpenFile opens a file or changes the current directory
func (a *App) OpenFile(path string) (CurrentFilesState, error) {
	info, err := os.Stat(path)
	if err != nil {
		return CurrentFilesState{}, fmt.Errorf("failed to get file info for %s: %w", path, err)
	}

	if info.IsDir() {
		a.currentDir = path
	} else {
		a.currentFile = path
	}

	return a.GetCurrentFilesState(), nil
}

// GoUp navigates one directory up from the current directory
func (a *App) GoUp() (*CurrentFilesState, error) {
	if a.currentDir == "" {
		return nil, fmt.Errorf("no current directory set")
	}

	// Get the parent directory
	parentDir := filepath.Dir(a.currentDir)

	// Check if we're already at the root
	if parentDir == a.currentDir {
		return nil, fmt.Errorf("already at root directory")
	}

	// Update current directory to parent
	a.currentDir = parentDir

	// Get parent directory info
	dirInfo, err := os.Stat(parentDir)
	if err != nil {
		return nil, fmt.Errorf("failed to get directory info for %s: %w", parentDir, err)
	}

	state := &CurrentFilesState{
		CurrentDir: &FileEntry{
			Name:        filepath.Base(parentDir),
			Path:        parentDir,
			IsDirectory: true,
			Size:        dirInfo.Size(),
			ModTime:     dirInfo.ModTime(),
		},
	}

	// Preserve CurrentFile if it still exists
	if a.currentFile != "" {
		fileInfo, err := os.Stat(a.currentFile)
		if err == nil {
			state.CurrentFile = &FileEntry{
				Name:        filepath.Base(a.currentFile),
				Path:        a.currentFile,
				IsDirectory: false,
				Size:        fileInfo.Size(),
				ModTime:     fileInfo.ModTime(),
			}
		}
	}

	return state, nil
}

// GetFileContent reads and returns the content of the specified file
func (a *App) GetFileContent(path string) (string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return "", fmt.Errorf("failed to get file info for %s: %w", path, err)
	}

	if info.IsDir() {
		return "", fmt.Errorf("path %s is a directory, not a file", path)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read file %s: %w", path, err)
	}

	return string(content), nil
}

// GetCurrentFilesState returns the current directory and file state
func (a *App) GetCurrentFilesState() CurrentFilesState {
	state := CurrentFilesState{}

	// Populate CurrentDir
	if a.currentDir != "" {
		dirInfo, err := os.Stat(a.currentDir)
		if err == nil {
			state.CurrentDir = &FileEntry{
				Name:        filepath.Base(a.currentDir),
				Path:        a.currentDir,
				IsDirectory: true,
				Size:        dirInfo.Size(),
				ModTime:     dirInfo.ModTime(),
			}
		}
	}

	// Populate CurrentFile
	if a.currentFile != "" {
		fileInfo, err := os.Stat(a.currentFile)
		if err == nil {
			state.CurrentFile = &FileEntry{
				Name:        filepath.Base(a.currentFile),
				Path:        a.currentFile,
				IsDirectory: false,
				Size:        fileInfo.Size(),
				ModTime:     fileInfo.ModTime(),
			}
		}
	}

	return state
}

// CreateFile creates a new .hurl file with the given name in the current directory
func (a *App) CreateFile(name string) error {
	if name == "" {
		return fmt.Errorf("file name cannot be empty")
	}

	// Add .hurl extension if not present
	if !strings.HasSuffix(strings.ToLower(name), ".hurl") {
		name = name + ".hurl"
	}

	filePath := filepath.Join(a.currentDir, name)

	// Check if file already exists
	if _, err := os.Stat(filePath); err == nil {
		return fmt.Errorf("file %s already exists", name)
	}

	// Create the file with a basic template
	initialContent := `# New Hurl Request
GET https://example.com

HTTP 200
`
	err := os.WriteFile(filePath, []byte(initialContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", name, err)
	}

	return nil
}

// CreateDir creates a new directory with the given name in the current directory
func (a *App) CreateDir(name string) error {
	if name == "" {
		return fmt.Errorf("directory name cannot be empty")
	}

	dirPath := filepath.Join(a.currentDir, name)

	// Check if directory already exists
	if _, err := os.Stat(dirPath); err == nil {
		return fmt.Errorf("directory %s already exists", name)
	}

	// Create the directory
	err := os.Mkdir(dirPath, 0755)
	if err != nil {
		return fmt.Errorf("failed to create directory %s: %w", name, err)
	}

	return nil
}

// DeleteFile deletes a file or directory
func (a *App) DeleteFile(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("failed to get file info for %s: %w", path, err)
	}

	if info.IsDir() {
		// If deleting current directory, move to parent
		if a.currentDir == path {
			a.currentDir = filepath.Dir(path)
		}
	} else {
		// If deleting current file, clear current file state
		if a.currentFile == path {
			a.currentFile = ""
		}
	}

	err = os.RemoveAll(path) // RemoveAll handles both files and directories recursively
	if err != nil {
		return fmt.Errorf("failed to delete %s: %w", path, err)
	}
	return nil
}

// RenameFile renames a file or directory
func (a *App) RenameFile(oldPath string, newName string) error {
	if newName == "" {
		return fmt.Errorf("new name cannot be empty")
	}

	// Check if old path exists
	info, err := os.Stat(oldPath)
	if err != nil {
		return fmt.Errorf("failed to get file info for %s: %w", oldPath, err)
	}

	// Build the new path (same directory, new name)
	dir := filepath.Dir(oldPath)
	newPath := filepath.Join(dir, newName)

	// Check if new path already exists
	if _, err := os.Stat(newPath); err == nil {
		return fmt.Errorf("a file or directory named %s already exists", newName)
	}

	// Rename the file/directory
	err = os.Rename(oldPath, newPath)
	if err != nil {
		return fmt.Errorf("failed to rename %s to %s: %w", oldPath, newName, err)
	}

	// Update current file/dir state if necessary
	if info.IsDir() && a.currentDir == oldPath {
		a.currentDir = newPath
	} else if !info.IsDir() && a.currentFile == oldPath {
		a.currentFile = newPath
	}

	return nil
}

// SaveFile saves content to the specified file
func (a *App) SaveFile(path string, content string) error {
	if path == "" {
		return fmt.Errorf("file path cannot be empty")
	}

	// Check if file exists
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("failed to get file info for %s: %w", path, err)
	}

	// Don't allow saving to directories
	if info.IsDir() {
		return fmt.Errorf("path %s is a directory, not a file", path)
	}

	// Write content to file
	err = os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to save file %s: %w", path, err)
	}

	return nil
}

// ClearCurrentFile clears the current file state and returns the updated state
func (a *App) ClearCurrentFile() CurrentFilesState {
	a.currentFile = ""
	return a.GetCurrentFilesState()
}

// SaveLastOpenedState saves the current directory and file to a state file
func (a *App) SaveLastOpenedState() error {
	stateFile := filepath.Join(os.TempDir(), "hurlstudio", "last_state.json")

	// Create directory if it doesn't exist
	stateDir := filepath.Dir(stateFile)
	if err := os.MkdirAll(stateDir, 0755); err != nil {
		return fmt.Errorf("failed to create state directory: %w", err)
	}

	state := CurrentFilesState{}
	if a.currentDir != "" {
		dirInfo, err := os.Stat(a.currentDir)
		if err == nil {
			state.CurrentDir = &FileEntry{
				Name:        filepath.Base(a.currentDir),
				Path:        a.currentDir,
				IsDirectory: true,
				Size:        dirInfo.Size(),
				ModTime:     dirInfo.ModTime(),
			}
		}
	}

	if a.currentFile != "" {
		fileInfo, err := os.Stat(a.currentFile)
		if err == nil {
			state.CurrentFile = &FileEntry{
				Name:        filepath.Base(a.currentFile),
				Path:        a.currentFile,
				IsDirectory: false,
				Size:        fileInfo.Size(),
				ModTime:     fileInfo.ModTime(),
			}
		}
	}

	data, err := json.Marshal(state)
	if err != nil {
		return fmt.Errorf("failed to marshal state: %w", err)
	}

	if err := os.WriteFile(stateFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write state file: %w", err)
	}

	return nil
}

// LoadLastOpenedState loads the last opened directory and file from the state file
func (a *App) LoadLastOpenedState() (CurrentFilesState, error) {
	stateFile := filepath.Join(os.TempDir(), "hurlstudio", "last_state.json")

	// Check if state file exists
	if _, err := os.Stat(stateFile); os.IsNotExist(err) {
		// No state file exists, return empty state
		return CurrentFilesState{}, nil
	}

	data, err := os.ReadFile(stateFile)
	if err != nil {
		return CurrentFilesState{}, fmt.Errorf("failed to read state file: %w", err)
	}

	var state CurrentFilesState
	if err := json.Unmarshal(data, &state); err != nil {
		return CurrentFilesState{}, fmt.Errorf("failed to unmarshal state: %w", err)
	}

	// Verify that the paths still exist
	if state.CurrentDir != nil {
		if _, err := os.Stat(state.CurrentDir.Path); err == nil {
			a.currentDir = state.CurrentDir.Path
		} else {
			state.CurrentDir = nil
		}
	}

	if state.CurrentFile != nil {
		if _, err := os.Stat(state.CurrentFile.Path); err == nil {
			a.currentFile = state.CurrentFile.Path
		} else {
			state.CurrentFile = nil
		}
	}

	return state, nil
}
