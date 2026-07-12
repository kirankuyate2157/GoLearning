package ingestion

import (
	"app/db"
	"app/db/repo"
	"app/models"
	"context"
	"fmt"
	"os"
	"os/exec"
	"sync"
	"time"
)

// IngestionProcess represents a running ingestion subprocess
type IngestionProcess struct {
	FeedUID    string
	ScriptType string
	Cmd        *exec.Cmd
	Cancel     context.CancelFunc
	StartedAt  time.Time
}

// IngestionManager manages all active ingestion processes
type IngestionManager struct {
	mu        sync.RWMutex
	processes map[string]*IngestionProcess // key is feed UID
}

var Manager = &IngestionManager{
	processes: make(map[string]*IngestionProcess),
}

// StartProcess starts an ingestion subprocess for a feed
func (im *IngestionManager) StartProcess(feedUID, connectionType string, args []string) error {
	im.mu.Lock()
	defer im.mu.Unlock()

	// Check if process already running
	if _, exists := im.processes[feedUID]; exists {
		return fmt.Errorf("ingestion process for feed %s is already running", feedUID)
	}

	// Get script path from config
	scriptPath, ok := pathConf[connectionType]
	if !ok {
		return fmt.Errorf("no script configured for connection type: %s", connectionType)
	}

	// Read the embedded script
	scriptContent, err := Scripts.ReadFile(scriptPath)
	if err != nil {
		return fmt.Errorf("failed to read script %s: %w", scriptPath, err)
	}

	// Create a temporary file for the script
	tmpFile, err := os.CreateTemp("", fmt.Sprintf("ingestion_%s_*.py", connectionType))
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}

	// Write script content to temp file
	if _, err := tmpFile.Write(scriptContent); err != nil {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
		return fmt.Errorf("failed to write script to temp file: %w", err)
	}
	tmpFile.Close()

	// Make the script executable (for Unix-like systems)
	os.Chmod(tmpFile.Name(), 0755)

	// Create context for process cancellation
	ctx, cancel := context.WithCancel(context.Background())
	allArgs := append([]string{tmpFile.Name()}, args...)
	fmt.Println(allArgs)
	// Start the Python script as subprocess
	cmd := exec.CommandContext(ctx, "python", allArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		cancel()
		os.Remove(tmpFile.Name())
		return fmt.Errorf("failed to start ingestion process: %w", err)
	}

	// Store process info
	im.processes[feedUID] = &IngestionProcess{
		FeedUID:    feedUID,
		ScriptType: connectionType,
		Cmd:        cmd,
		Cancel:     cancel,
		StartedAt:  time.Now(),
	}

	// Monitor process in goroutine
	go im.monitorProcess(feedUID, tmpFile.Name())

	return nil
}

// monitorProcess monitors the subprocess and cleans up when it exits
func (im *IngestionManager) monitorProcess(feedUID, tmpFilePath string) {
	im.mu.RLock()
	process, exists := im.processes[feedUID]
	im.mu.RUnlock()

	if !exists {
		return
	}

	// Wait for process to complete
	err := process.Cmd.Wait()

	// Clean up
	im.mu.Lock()
	delete(im.processes, feedUID)
	im.mu.Unlock()

	// Remove temporary script file
	os.Remove(tmpFilePath)
	pool := db.DBConfig()
	defer pool.Close()

	running := false
	repo.UpdateFeed(pool, feedUID, models.FeedUpdateRequest{IsRunning: &running})

	if err != nil {
		fmt.Printf("Ingestion process for feed %s exited with error: %v\n", feedUID, err)
	} else {
		fmt.Printf("Ingestion process for feed %s completed successfully\n", feedUID)
	}
}

// StopProcess stops an ingestion subprocess for a feed
func (im *IngestionManager) StopProcess(feedUID string) error {
	im.mu.Lock()
	defer im.mu.Unlock()

	process, exists := im.processes[feedUID]
	if !exists {
		return fmt.Errorf("no active ingestion process found for feed: %s", feedUID)
	}

	// Send SIGINT to the process for graceful shutdown
	if process.Cmd.Process != nil {
		process.Cmd.Process.Signal(os.Interrupt)
	}

	// Cancel the context to stop the process
	process.Cancel()

	// Remove from active processes
	delete(im.processes, feedUID)

	return nil
}

// ListActiveProcesses returns information about all active ingestion processes
func (im *IngestionManager) ListActiveProcesses() []map[string]any {
	im.mu.RLock()
	defer im.mu.RUnlock()

	result := make([]map[string]any, 0, len(im.processes))
	for _, process := range im.processes {
		result = append(result, map[string]any{
			"feed_uid":    process.FeedUID,
			"script_type": process.ScriptType,
			"started_at":  process.StartedAt,
			"uptime":      time.Since(process.StartedAt).String(),
			"pid":         process.Cmd.Process.Pid,
		})
	}

	return result
}

// IsProcessActive checks if a process is running for a given feed
func (im *IngestionManager) IsProcessActive(feedUID string) bool {
	im.mu.RLock()
	defer im.mu.RUnlock()

	_, exists := im.processes[feedUID]
	return exists
}
