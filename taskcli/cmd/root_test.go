package cmd

import (
	"bytes"
	"testing"
)

func TestRootCmd(t *testing.T) {
	t.Run("execution of taskcli should return the default message", func(t *testing.T) {

		// Arrange
		var capture bytes.Buffer
		rootCmd.SetOut(&capture)

		// Act
		err := rootCmd.Execute()
		if err != nil {
			t.Fatalf("rootCmd.Execute() error = %v, want nil", err)
		}

		// Assert
		expected := `taskcli is your CLI best friend`
		output := capture.String()

		if !contains(output, expected) {
			t.Errorf("Expected output to contain %q, but got %q", expected, output)
		}
	})
}

func contains(output, expected string) bool {
	return bytes.Contains([]byte(output), []byte(expected))
}
