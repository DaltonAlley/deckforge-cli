package pdf

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/daltonalley/deckforge-cli/scryfall"
	"github.com/rs/zerolog/log"
	"github.com/signintech/gopdf"
)

// PDFGenerator interface for generating PDFs
type PDFGenerator interface {
	SetOutputPath(path string)
	GeneratePDF(decklist *Decklist, progress Reporter) error
}

// Generator handles PDF creation with bleed margins
type Generator struct {
	bleedAmount float64
	outputPath  string
}

// Decklist represents a parsed decklist from CSV
type Decklist struct {
	Cards []CardEntry
}

// CardEntry represents a single card in the decklist
type CardEntry struct {
	Qty  int
	ID   string
	Card scryfall.Card
}

// Reporter interface for progress reporting
type Reporter interface {
	StartUnified(description string, total int)
	UpdateStage(description string)
	AddError(cardName, errorMsg string)
	Finish(decklist interface{}, outputPath string)
}

// MTG card dimensions in mm
const (
	CardWidth  = 63.0 // 63mm
	CardHeight = 88.0 // 88mm
)

// NewGenerator creates a new PDF generator with the specified bleed amount
func NewGenerator(bleedAmount float64) PDFGenerator {
	return &Generator{
		bleedAmount: bleedAmount,
	}
}

// SetOutputPath sets the output path for the PDF
func (g *Generator) SetOutputPath(path string) {
	g.outputPath = path
}

// TotalWidth returns the total page width including bleed
func (g *Generator) TotalWidth() float64 {
	return CardWidth + (2 * g.bleedAmount)
}

// TotalHeight returns the total page height including bleed
func (g *Generator) TotalHeight() float64 {
	return CardHeight + (2 * g.bleedAmount)
}

// ImagePosition returns the X,Y coordinates where the card image should be positioned
func (g *Generator) ImagePosition() (float64, float64) {
	return g.bleedAmount, g.bleedAmount
}

// GeneratePDF creates a PDF from the decklist
func (g *Generator) GeneratePDF(decklist *Decklist, progress Reporter) error {
	// Create cache directory for images
	cacheDir := ".card_cache"
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}

	pages := [][]byte{}

	// Initialize progress reporting
	if progress != nil {
		totalOps := calculateTotalOperations(decklist)
		progress.StartUnified("Starting PDF generation", totalOps)
	}

	// Process each card entry
	for _, cardEntry := range decklist.Cards {
		// Update progress for fetching card
		if progress != nil {
			progress.UpdateStage(fmt.Sprintf("Fetching card: %s", cardEntry.ID))
		}

		// Fetch card data from Scryfall
		card, err := scryfall.FindCardByID(cardEntry.ID)
		if err != nil {
			log.Error().Err(err).Str("cardID", cardEntry.ID).Msg("Failed to fetch card data")
			if progress != nil {
				progress.AddError(cardEntry.ID, err.Error())
			}
			// Skip error cards for now - just log the error
			continue
		}

		// Update card entry with fetched data
		entryWithCard := cardEntry
		entryWithCard.Card = card

		// Handle double-sided cards
		if len(card.CardFaces) > 0 {
			// For double-sided cards, create separate pages for each face
			for _, face := range card.CardFaces {
				faceEntry := CardEntry{
					Qty: 1, // Each face gets its own page
					ID:  cardEntry.ID,
					Card: scryfall.Card{
						Name:      face.Name,
						ImageURIs: face.ImageURIs,
					},
				}

				// Update progress for generating face page
				if progress != nil {
					progress.UpdateStage(fmt.Sprintf("Generating page: %s", face.Name))
				}

				facePage, err := g.generatePage(faceEntry, cacheDir)
				if err != nil {
					log.Error().Err(err).Str("cardID", cardEntry.ID).Str("face", face.Name).Msg("Failed to generate face page")
					if progress != nil {
						progress.AddError(fmt.Sprintf("%s (%s)", face.Name, cardEntry.ID), err.Error())
					}
					continue
				}

				// Only add non-empty pages
				if len(facePage) > 0 {
					// Add face page for each quantity
					for i := 0; i < cardEntry.Qty; i++ {
						pages = append(pages, facePage)
					}
				}
			}
		} else {
			// Single-sided card
			if progress != nil {
				progress.UpdateStage(fmt.Sprintf("Generating page: %s", card.Name))
			}

			page, err := g.generatePage(entryWithCard, cacheDir)
			if err != nil {
				log.Error().Err(err).Str("cardID", cardEntry.ID).Msg("Failed to generate card page")
				if progress != nil {
					progress.AddError(card.Name, err.Error())
				}
				continue
			}

			// Only add non-empty pages
			if len(page) > 0 {
				// Add page for each quantity
				for i := 0; i < cardEntry.Qty; i++ {
					pages = append(pages, page)
				}
			}
		}
	}

	if len(pages) == 0 {
		return fmt.Errorf("no pages generated")
	}

	// Create final PDF by combining all pages
	if progress != nil {
		progress.UpdateStage("Assembling PDF")
	}

	size := gopdf.Rect{W: g.TotalWidth(), H: g.TotalHeight()}
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: size})
	for _, b := range pages {
		if len(b) == 0 {
			continue
		}

		reader := bytes.NewReader(b)
		if err := pdf.ImportPagesFromSource(reader, "/MediaBox"); err != nil {
			return err
		}
	}
	if err := pdf.WritePdf(g.outputPath); err != nil {
		return err
	}

	return nil
}

// generatePage creates a single PDF page for a card entry
func (g *Generator) generatePage(ce CardEntry, cacheDir string) ([]byte, error) {
	// Create PDF with dimensions including bleed
	size := gopdf.Rect{W: g.TotalWidth(), H: g.TotalHeight()}
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: size})
	pdf.AddPage()

	// Fill bleed area with black background (same color as card border)
	if g.bleedAmount > 0 {
		pdf.SetFillColor(0, 0, 0) // Black
		pdf.RectFromUpperLeftWithStyle(0, 0, g.TotalWidth(), g.TotalHeight(), "F")
	}

	// Determine image URL to use
	var imageURL string
	var imagePath string
	var err error

	if ce.Card.ImageURIs.Normal != "" {
		imageURL = ce.Card.ImageURIs.Normal
		// Download image from URL
		cacheKey := fmt.Sprintf("%s_normal", ce.ID)
		if ce.Card.Name != "" {
			cacheKey = fmt.Sprintf("%s_%s_normal", ce.ID, sanitizeFilename(ce.Card.Name))
		}
		imagePath, err = scryfall.DownloadImageFromURL(imageURL, cacheKey, cacheDir)
	} else {
		// Fallback: try to get from Scryfall API
		imagePath, err = scryfall.DownloadCardImage(ce.ID, cacheDir, "normal")
	}

	if err != nil {
		// Log error (text on PDF causes font issues, so just log for now)
		log.Error().Err(err).Str("cardID", ce.ID).Str("cardName", ce.Card.Name).Msg("Failed to load card image")

		// Return empty page for now (avoids PDF corruption)
		page, err := pdf.GetBytesPdfReturnErr()
		if err != nil {
			return nil, err
		}
		return page, nil
	}

	// Embed the image at bleed offset, scaling to fit card dimensions
	imageX, imageY := g.ImagePosition()
	err = pdf.Image(imagePath, imageX, imageY, &gopdf.Rect{W: CardWidth, H: CardHeight})
	if err != nil {
		log.Warn().Err(err).Str("cardID", ce.ID).Msg("Failed to embed image")
		// Return empty page
		page, err := pdf.GetBytesPdfReturnErr()
		if err != nil {
			return nil, err
		}
		return page, nil
	}

	page, err := pdf.GetBytesPdfReturnErr()
	if err != nil {
		return nil, err
	}
	return page, nil
}

// sanitizeFilename creates a safe filename from card name
func sanitizeFilename(name string) string {
	// Simple sanitization - replace spaces and special chars
	return strings.ReplaceAll(strings.ReplaceAll(name, " ", "_"), "/", "_")
}

// calculateTotalOperations calculates the total number of operations for unified progress tracking
func calculateTotalOperations(decklist *Decklist) int {
	total := 0
	for range decklist.Cards {
		total += 1 // Card data fetch
		total += 1 // Page generation
	}
	total += 1 // Assembling PDF
	return total
}
