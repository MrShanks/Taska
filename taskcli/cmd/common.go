package cmd

import (
	"os"
	"os/exec"
	"strings"
)

// openEditor opens vim editor and returns the edited content.
// First line is interpreted as title and from the second on will be description
func openEditor() (string, string, error) {
	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "task_edit_*.txt")
	if err != nil {
		return "", "", err
	}
	defer os.Remove(tmpFile.Name()) // Clean up after editing

	// Open the editor
	cmd := exec.Command("vim", tmpFile.Name()) //nolint:gosec // input is passed directly
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Wait for the user to edit and save
	if err := cmd.Run(); err != nil {
		return "", "", err
	}

	// Read the edited file content
	content, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		return "", "", err
	}

	lines := strings.SplitN(string(content), "\n", 2)

	title := strings.TrimSpace(lines[0]) // First line as title

	description := ""
	if len(lines) > 1 {
		description = strings.TrimSpace(lines[1]) // Rest as description
	}

	return title, description, nil
}
