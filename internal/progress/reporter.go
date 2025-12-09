package progress

import (
	"fmt"
	"os"
	"path/filepath"
)

// Reporter handles simple indexed progress reporting during PDF generation
type Reporter interface {
	StartUnified(description string, total int)
	UpdateStage(description string)
	AddError(cardName, errorMsg string)
	Finish(decklist interface{}, outputPath string)
}

// ProgressReporter implements Reporter with indexed display
type ProgressReporter struct {
	currentIndex    int
	totalOperations int
	errors          []string
}

// NewProgressReporter creates a new progress reporter
func NewProgressReporter() Reporter {
	return &ProgressReporter{
		errors: make([]string, 0),
	}
}

// StartUnified begins unified progress tracking with total operations
func (pr *ProgressReporter) StartUnified(description string, total int) {
	pr.totalOperations = total
	pr.currentIndex = 0
	pr.errors = make([]string, 0)
}

// UpdateStage updates the current stage with indexed display
func (pr *ProgressReporter) UpdateStage(description string) {
	pr.currentIndex++
	fmt.Printf("[%d/%d] %s\n", pr.currentIndex, pr.totalOperations, description)
}

// AddError records an error that occurred during processing
func (pr *ProgressReporter) AddError(cardName, errorMsg string) {
	pr.errors = append(pr.errors, fmt.Sprintf("%s: %s", cardName, errorMsg))
}

// Finish completes the progress reporting and shows comprehensive completion message
func (pr *ProgressReporter) Finish(decklist interface{}, outputPath string) {
	if len(pr.errors) == 0 {
		fmt.Printf("\n✅ Successfully generated PDF '%s'\n", filepath.Base(outputPath))
	} else {
		fmt.Fprintf(os.Stderr, "\n❌ PDF generation completed with %d error(s):\n", len(pr.errors))
		for _, err := range pr.errors {
			fmt.Fprintf(os.Stderr, "   • %s\n", err)
		}
	}
}
