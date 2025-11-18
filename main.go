package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/daltonalley/go-scryfall"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v3"
)

type Card struct {
	Qty        int
	Name       string
	ScryfallID string
}

type Decklist struct {
	Cards []Card
}

var Session map[string]Decklist

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	cmd := &cli.Command{
		Name:  "deckforge",
		Usage: "Generate a mtg deck pdf with an exported Archidekt csv file.",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.Args().Len() == 0 {
				return fmt.Errorf("no file provided")
			}

			fileName := cmd.Args().Get(0)
			dir, err := os.Getwd()
			if err != nil {
				return err
			}

			path := filepath.Join(dir, fileName)
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			reader := csv.NewReader(file)
			records, err := reader.ReadAll()
			if err != nil {
				return err
			}

			Session = make(map[string]Decklist)
			decklist := Decklist{Cards: []Card{}}
			for _, record := range records {
				qty, err := strconv.Atoi(record[0])
				if err != nil {
					return err
				}
				card := Card{
					Qty:        qty,
					Name:       record[1],
					ScryfallID: record[2],
				}
				decklist.Cards = append(decklist.Cards, card)
				log.Print(card)

				httpClient := &http.Client{
					Timeout: time.Second * 10,
				}

				c := scryfall.NewClient(httpClient)

				c.FindCardByID(card.ScryfallID)
			}
			Session[uuid.New().String()] = decklist
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal().Err(err).Msg("")
	}
}
