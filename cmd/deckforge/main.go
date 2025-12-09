package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/daltonalley/deckforge-cli/internal/deck"
	"github.com/daltonalley/deckforge-cli/internal/pdf"
	"github.com/daltonalley/deckforge-cli/internal/progress"
	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:  "deckforge",
		Usage: "Generate printable MTG deck PDFs from Archidekt CSV exports",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Value:   "",
				Usage:   "Output PDF filename (defaults to CSV name)",
			},
			&cli.FloatFlag{
				Name:  "bleed",
				Value: 3.0,
				Usage: "Bleed margin in mm around each card",
			},
			&cli.BoolFlag{
				Name:    "quiet",
				Aliases: []string{"q"},
				Usage:   "Suppress progress output",
			},
		},
		Action: runDeckForge,
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runDeckForge(ctx context.Context, cmd *cli.Command) error {
	// Get arguments
	if cmd.Args().Len() == 0 {
		return fmt.Errorf("no CSV file specified")
	}

	csvPath := cmd.Args().Get(0)

	// Determine output path
	outputPath := cmd.String("output")
	if outputPath == "" {
		baseName := strings.TrimSuffix(filepath.Base(csvPath), filepath.Ext(filepath.Base(csvPath)))
		outputPath = baseName + ".pdf"
	}

	// Get configuration
	bleedAmount := cmd.Float("bleed")
	quiet := cmd.Bool("quiet")

	// Open and parse CSV
	file, err := os.Open(csvPath)
	if err != nil {
		return fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer file.Close()

	decklist, err := deck.ParseDecklistCSV(file)
	if err != nil {
		return fmt.Errorf("failed to parse CSV: %w", err)
	}

	// Convert to pdf.Decklist type
	pdfDecklist := &pdf.Decklist{}
	for _, card := range decklist.Cards {
		pdfDecklist.Cards = append(pdfDecklist.Cards, pdf.CardEntry{
			Qty: card.Qty,
			ID:  card.ID,
		})
	}

	// Create components
	pdfGen := pdf.NewGenerator(bleedAmount)
	pdfGen.SetOutputPath(outputPath)
	var progressReporter progress.Reporter
	if !quiet {
		progressReporter = progress.NewProgressReporter()
	}

	// Calculate total operations for progress tracking
	totalOps := calculateTotalOperations(decklist)

	// Start progress tracking
	if progressReporter != nil {
		progressReporter.StartUnified("Starting PDF generation", totalOps)
	}

	// Generate PDF
	if err := pdfGen.GeneratePDF(pdfDecklist, progressReporter); err != nil {
		return fmt.Errorf("failed to generate PDF: %w", err)
	}

	// Finish progress reporting
	if progressReporter != nil {
		progressReporter.Finish(decklist, outputPath)
	}

	return nil
}

// calculateTotalOperations calculates the total number of operations for unified progress tracking
func calculateTotalOperations(decklist *deck.Decklist) int {
	total := 0
	for range decklist.Cards {
		total += 1 // Card data fetch
		total += 1 // Page generation
	}
	total += 1 // Assembling PDF
	return total
}
