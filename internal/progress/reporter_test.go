package progress

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProgressReporter(t *testing.T) {
	t.Run("progress reporter initialization", func(t *testing.T) {
		reporter := NewProgressReporter().(*ProgressReporter)
		require.NotNil(t, reporter)
		require.Empty(t, reporter.errors)
		require.Equal(t, 0, reporter.currentIndex)
		require.Equal(t, 0, reporter.totalOperations)
	})

	t.Run("start unified progress", func(t *testing.T) {
		reporter := NewProgressReporter().(*ProgressReporter)
		reporter.StartUnified("Test Progress", 10)
		require.Equal(t, 10, reporter.totalOperations)
		require.Equal(t, 0, reporter.currentIndex)
	})

	t.Run("update stage increments index", func(t *testing.T) {
		reporter := NewProgressReporter().(*ProgressReporter)
		reporter.StartUnified("Test", 3)

		reporter.UpdateStage("Operation 1")
		require.Equal(t, 1, reporter.currentIndex)

		reporter.UpdateStage("Operation 2")
		require.Equal(t, 2, reporter.currentIndex)
	})

	t.Run("error tracking", func(t *testing.T) {
		reporter := NewProgressReporter().(*ProgressReporter)
		reporter.StartUnified("Test", 3)

		reporter.AddError("Card1", "Network error")
		require.Len(t, reporter.errors, 1)
		require.Contains(t, reporter.errors[0], "Card1")
		require.Contains(t, reporter.errors[0], "Network error")
	})

	t.Run("finish with success", func(t *testing.T) {
		// Capture stdout to verify output
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		reporter := NewProgressReporter().(*ProgressReporter)
		reporter.StartUnified("Test", 3)

		// Simulate successful completion
		reporter.UpdateStage("Op 1")
		reporter.UpdateStage("Op 2")
		reporter.UpdateStage("Op 3")

		// Create a mock decklist
		mockDecklist := struct{ Cards []interface{} }{Cards: make([]interface{}, 3)}

		reporter.Finish(mockDecklist, "test.pdf")

		w.Close()
		os.Stdout = oldStdout

		output, _ := io.ReadAll(r)
		outputStr := string(output)

		require.Contains(t, outputStr, "✅ Successfully generated PDF")
		require.Contains(t, outputStr, "test.pdf")
	})

	t.Run("finish with errors", func(t *testing.T) {
		oldStderr := os.Stderr
		r, w, _ := os.Pipe()
		os.Stderr = w

		reporter := NewProgressReporter().(*ProgressReporter)
		reporter.StartUnified("Test", 3)

		reporter.AddError("Card1", "Network error")
		reporter.AddError("Card2", "Parse error")

		mockDecklist := struct{ Cards []interface{} }{Cards: make([]interface{}, 3)}
		reporter.Finish(mockDecklist, "test.pdf")

		w.Close()
		os.Stderr = oldStderr

		output, _ := io.ReadAll(r)
		outputStr := string(output)

		require.Contains(t, outputStr, "❌ PDF generation completed with 2 error(s)")
		require.Contains(t, outputStr, "Card1: Network error")
		require.Contains(t, outputStr, "Card2: Parse error")
	})
}

// mockWriter captures output for testing
type mockWriter struct {
	buffer *bytes.Buffer
}

func newMockWriter() *mockWriter {
	return &mockWriter{buffer: &bytes.Buffer{}}
}

func (mw *mockWriter) Write(p []byte) (n int, err error) {
	return mw.buffer.Write(p)
}

func (mw *mockWriter) String() string {
	return mw.buffer.String()
}
