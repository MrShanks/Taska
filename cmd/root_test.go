package cmd

import (
	"bytes"
	"testing"
)

func TestRootCmd(t *testing.T) {
	t.Run("Run rootCmd should return default message", func(t *testing.T) {
		var capture bytes.Buffer
		rootCmd.SetOut(&capture)

		err := rootCmd.Execute()
		if err != nil {
			t.Fatalf("rootCmd.Execute() error = %v, want nil", err)
		}

<<<<<<< HEAD
		expected := `taskcli is your CLI best friend`
=======
		// expected := `taskcli is your CLI best friend`
		expected := `test build should fail and merge should be blocked`
>>>>>>> 7717e80 (feat: TAS-10 add faulty test)
		output := capture.String()

		if !contains(output, expected) {
			t.Errorf("Expected output to contain %q, but got %q", expected, output)
		}
	})
}

func contains(output, expected string) bool {
	return bytes.Contains([]byte(output), []byte(expected))
}
