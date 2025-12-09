package deck

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseDecklistCSV(t *testing.T) {
	t.Run("valid CSV with single card", func(t *testing.T) {
		csvData := `1,"Lightning Bolt",a65e485b-03a2-4634-9218-f5bb7c104d41`
		reader := strings.NewReader(csvData)

		decklist, err := ParseDecklistCSV(reader)
		require.NoError(t, err)
		require.Len(t, decklist.Cards, 1)
		require.Equal(t, 1, decklist.Cards[0].Qty)
		require.Equal(t, "a65e485b-03a2-4634-9218-f5bb7c104d41", decklist.Cards[0].ID)
	})

	t.Run("valid CSV with multiple cards", func(t *testing.T) {
		csvData := `1,"Lightning Bolt",a65e485b-03a2-4634-9218-f5bb7c104d41
2,"Black Lotus",b6a5b3b0-2b4b-4c4b-8b2b-2b2b2b2b2b2b`
		reader := strings.NewReader(csvData)

		decklist, err := ParseDecklistCSV(reader)
		require.NoError(t, err)
		require.Len(t, decklist.Cards, 2)
		require.Equal(t, 1, decklist.Cards[0].Qty)
		require.Equal(t, 2, decklist.Cards[1].Qty)
	})

	t.Run("invalid CSV - missing quantity field", func(t *testing.T) {
		csvData := `"Lightning Bolt",a65e485b-03a2-4634-9218-f5bb7c104d41`
		reader := strings.NewReader(csvData)

		_, err := ParseDecklistCSV(reader)
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid CSV format")
	})

	t.Run("invalid CSV - non-numeric quantity", func(t *testing.T) {
		csvData := `abc,"Lightning Bolt",a65e485b-03a2-4634-9218-f5bb7c104d41`
		reader := strings.NewReader(csvData)

		_, err := ParseDecklistCSV(reader)
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid quantity")
	})

	t.Run("invalid CSV - zero quantity", func(t *testing.T) {
		csvData := `0,"Lightning Bolt",a65e485b-03a2-4634-9218-f5bb7c104d41`
		reader := strings.NewReader(csvData)

		_, err := ParseDecklistCSV(reader)
		require.Error(t, err)
		require.Contains(t, err.Error(), "quantity must be positive")
	})

	t.Run("invalid CSV - invalid Scryfall ID format", func(t *testing.T) {
		csvData := `1,"Lightning Bolt",invalid-id`
		reader := strings.NewReader(csvData)

		_, err := ParseDecklistCSV(reader)
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid Scryfall ID format")
	})
}

func TestCalculateTotalPages(t *testing.T) {
	t.Run("single card with quantity 1", func(t *testing.T) {
		decklist := &Decklist{
			Cards: []CardEntry{
				{Qty: 1, ID: "test-id"},
			},
		}

		pages := CalculateTotalPages(decklist)
		require.Equal(t, 1, pages)
	})

	t.Run("multiple cards with different quantities", func(t *testing.T) {
		decklist := &Decklist{
			Cards: []CardEntry{
				{Qty: 1, ID: "card1"},
				{Qty: 3, ID: "card2"},
				{Qty: 2, ID: "card3"},
			},
		}

		pages := CalculateTotalPages(decklist)
		require.Equal(t, 6, pages) // 1 + 3 + 2
	})

	t.Run("empty decklist", func(t *testing.T) {
		decklist := &Decklist{
			Cards: []CardEntry{},
		}

		pages := CalculateTotalPages(decklist)
		require.Equal(t, 0, pages)
	})
}
