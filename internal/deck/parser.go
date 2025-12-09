package deck

import (
	"encoding/csv"
	"fmt"
	"io"
	"regexp"
	"strconv"

	"github.com/daltonalley/deckforge-cli/scryfall"
)

// CardEntry represents a single card in the decklist
type CardEntry struct {
	Qty  int
	ID   string
	Card scryfall.Card // Card data from Scryfall API
}

// Decklist represents a parsed decklist from CSV
type Decklist struct {
	Cards []CardEntry
}

// ParseDecklistCSV parses and validates a CSV reader containing decklist data
// Expected format: quantity,"card name",scryfall_id
func ParseDecklistCSV(reader io.Reader) (*Decklist, error) {
	csvReader := csv.NewReader(reader)
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV: %w", err)
	}

	var decklist Decklist
	scryfallIDRegex := regexp.MustCompile(`^[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}$`)

	for i, record := range records {
		if len(record) < 3 {
			return nil, fmt.Errorf("invalid CSV format at line %d: expected 3 fields, got %d", i+1, len(record))
		}

		// Parse quantity
		qty, err := strconv.Atoi(record[0])
		if err != nil {
			return nil, fmt.Errorf("invalid quantity '%s' at line %d: %w", record[0], i+1, err)
		}
		if qty <= 0 {
			return nil, fmt.Errorf("quantity must be positive at line %d, got %d", i+1, qty)
		}

		// Validate Scryfall ID format
		scryfallID := record[2]
		if !scryfallIDRegex.MatchString(scryfallID) {
			return nil, fmt.Errorf("invalid Scryfall ID format '%s' at line %d", scryfallID, i+1)
		}

		// Create card entry
		entry := CardEntry{
			Qty: qty,
			ID:  scryfallID,
		}
		decklist.Cards = append(decklist.Cards, entry)
	}

	return &decklist, nil
}

// CalculateTotalPages calculates the total number of pages that will be generated
func CalculateTotalPages(decklist *Decklist) int {
	total := 0
	for _, card := range decklist.Cards {
		if len(card.Card.CardFaces) > 0 {
			total += len(card.Card.CardFaces) * card.Qty
		} else {
			total += card.Qty
		}
	}
	return total
}
