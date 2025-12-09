package scryfall

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/bytedance/sonic"
)

type Card struct {
	ID              string        `json:"id"`
	OracleID        string        `json:"oracle_id"`
	MultiverseIDs   []int         `json:"multiverse_ids"`
	TcgplayerID     int           `json:"tcgplayer_id"`
	CardmarketID    int           `json:"cardmarket_id"`
	Name            string        `json:"name"`
	Lang            string        `json:"lang"`
	ReleasedAt      string        `json:"released_at"`
	URI             string        `json:"uri"`
	ScryfallURI     string        `json:"scryfall_uri"`
	Layout          string        `json:"layout"`
	HighresImage    bool          `json:"highres_image"`
	ImageStatus     string        `json:"image_status"`
	ImageURIs       ImageURIs     `json:"image_uris"`
	ManaCost        string        `json:"mana_cost"`
	CMC             float64       `json:"cmc"`
	TypeLine        string        `json:"type_line"`
	OracleText      string        `json:"oracle_text"`
	Power           string        `json:"power"`
	Toughness       string        `json:"toughness"`
	Loyalty         string        `json:"loyalty"`
	Colors          []string      `json:"colors"`
	ColorIdentity   []string      `json:"color_identity"`
	Keywords        []string      `json:"keywords"`
	ProducedMana    []string      `json:"produced_mana"`
	AllParts        []RelatedCard `json:"all_parts"`
	CardFaces       []CardFace    `json:"card_faces"`
	Legality        Legalities    `json:"legalities"`
	Games           []string      `json:"games"`
	Reserved        bool          `json:"reserved"`
	Foil            bool          `json:"foil"`
	Nonfoil         bool          `json:"nonfoil"`
	Finishes        []string      `json:"finishes"`
	Oversized       bool          `json:"oversized"`
	Promo           bool          `json:"promo"`
	Reprint         bool          `json:"reprint"`
	Variation       bool          `json:"variation"`
	SetID           string        `json:"set_id"`
	Set             string        `json:"set"`
	SetName         string        `json:"set_name"`
	SetType         string        `json:"set_type"`
	SetURI          string        `json:"set_uri"`
	SetSearchURI    string        `json:"set_search_uri"`
	ScryfallSetURI  string        `json:"scryfall_set_uri"`
	RulingsURI      string        `json:"rulings_uri"`
	PrintsSearchURI string        `json:"prints_search_uri"`
	CollectorNumber string        `json:"collector_number"`
	Digital         bool          `json:"digital"`
	Rarity          string        `json:"rarity"`
	FlavorText      string        `json:"flavor_text"`
	CardBackID      string        `json:"card_back_id"`
	Artist          string        `json:"artist"`
	ArtistIDs       []string      `json:"artist_ids"`
	IllustrationID  string        `json:"illustration_id"`
	BorderColor     string        `json:"border_color"`
	Frame           string        `json:"frame"`
	FrameEffects    []string      `json:"frame_effects"`
	SecurityStamp   string        `json:"security_stamp"`
	FullArt         bool          `json:"full_art"`
	Textless        bool          `json:"textless"`
	Booster         bool          `json:"booster"`
	StorySpotlight  bool          `json:"story_spotlight"`
	EDHREC_RANK     int           `json:"edhrec_rank"`
	PennyRank       int           `json:"penny_rank"`
	Prices          Prices        `json:"prices"`
	RelatedURIs     RelatedURIs   `json:"related_uris"`
	PurchaseURIs    PurchaseURIs  `json:"purchase_uris"`
}

type ImageURIs struct {
	Small      string `json:"small"`
	Normal     string `json:"normal"`
	Large      string `json:"large"`
	PNG        string `json:"png"`
	ArtCrop    string `json:"art_crop"`
	BorderCrop string `json:"border_crop"`
}

type RelatedCard struct {
	ID        string `json:"id"`
	Component string `json:"component"`
	Name      string `json:"name"`
	TypeLine  string `json:"type_line"`
	URI       string `json:"uri"`
}

type CardFace struct {
	Name           string    `json:"name"`
	ManaCost       string    `json:"mana_cost"`
	TypeLine       string    `json:"type_line"`
	OracleText     string    `json:"oracle_text"`
	Colors         []string  `json:"colors"`
	Power          string    `json:"power"`
	Toughness      string    `json:"toughness"`
	Loyalty        string    `json:"loyalty"`
	FlavorText     string    `json:"flavor_text"`
	Artist         string    `json:"artist"`
	ArtistID       string    `json:"artist_id"`
	IllustrationID string    `json:"illustration_id"`
	ImageURIs      ImageURIs `json:"image_uris"`
}

type Legalities struct {
	Standard        string `json:"standard"`
	Future          string `json:"future"`
	Historic        string `json:"historic"`
	Timeless        string `json:"timeless"`
	Gladiator       string `json:"gladiator"`
	Pioneer         string `json:"pioneer"`
	Modern          string `json:"modern"`
	Legacy          string `json:"legacy"`
	Pauper          string `json:"pauper"`
	Vintage         string `json:"vintage"`
	Penny           string `json:"penny"`
	Commander       string `json:"commander"`
	Oathbreaker     string `json:"oathbreaker"`
	Standardbrawl   string `json:"standardbrawl"`
	Brawl           string `json:"brawl"`
	Alchemy         string `json:"alchemy"`
	Paupercommander string `json:"paupercommander"`
	Duel            string `json:"duel"`
	Oldschool       string `json:"oldschool"`
	Premodern       string `json:"premodern"`
	Predh           string `json:"predh"`
}

type Prices struct {
	USD       string `json:"usd"`
	USDFoil   string `json:"usd_foil"`
	USDEtched string `json:"usd_etched"`
	EUR       string `json:"eur"`
	EURFoil   string `json:"eur_foil"`
	Tix       string `json:"tix"`
}

type RelatedURIs struct {
	Gatherer         string `json:"gatherer"`
	TCGPlayerArticle string `json:"tcgplayer_infinite_article"`
	TCGPlayerDecks   string `json:"tcgplayer_infinite_decks"`
	EDHREC           string `json:"edhrec"`
}

type PurchaseURIs struct {
	Tcgplayer   string `json:"tcgplayer"`
	Cardmarket  string `json:"cardmarket"`
	Cardhoarder string `json:"cardhoarder"`
}

const baseURL = "https://api.scryfall.com/cards/"

func FindCardByID(id string) (Card, error) {
	req, err := http.NewRequest("GET", baseURL+id, nil)
	if err != nil {
		return Card{}, err
	}

	req.Header.Set("User-Agent", "go-scryfall/0.1.2")
	req.Header.Set("Accept", "application/json")

	c := http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return Card{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Card{}, fmt.Errorf("Error: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Card{}, err
	}

	var card Card
	err = sonic.Unmarshal(body, &card)
	if err != nil {
		return Card{}, err
	}

	return card, nil
}

// DownloadCardImage downloads a card image from Scryfall and caches it locally
// quality can be "normal", "large", "png", etc.
// If card has face-specific ImageURIs (for double-sided cards), those take precedence
func DownloadCardImage(cardID, cacheDir, quality string) (string, error) {
	// First get the card data to find the image URL
	card, err := FindCardByID(cardID)
	if err != nil {
		return "", fmt.Errorf("failed to find card %s: %w", cardID, err)
	}

	// Determine image URL based on quality
	// For double-sided cards, use face ImageURIs if available
	var imageURL string
	switch quality {
	case "normal":
		imageURL = card.ImageURIs.Normal
	case "large":
		imageURL = card.ImageURIs.Large
	case "png":
		imageURL = card.ImageURIs.PNG
	default:
		imageURL = card.ImageURIs.Normal
	}

	if imageURL == "" {
		return "", fmt.Errorf("no image URL available for card %s", cardID)
	}

	// Create cache filename
	cacheFile := filepath.Join(cacheDir, fmt.Sprintf("%s_%s.jpg", cardID, quality))

	// Check if cached file exists
	if _, err := os.Stat(cacheFile); err == nil {
		return cacheFile, nil
	}

	// Download the image
	req, err := http.NewRequest("GET", imageURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request for %s: %w", imageURL, err)
	}

	req.Header.Set("User-Agent", "go-scryfall/0.1.2")
	req.Header.Set("Accept", "image/*")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to download image for %s: %w", cardID, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download image for %s: HTTP %d", cardID, resp.StatusCode)
	}

	// Create cache directory if it doesn't exist
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create cache directory: %w", err)
	}

	// Save image to cache file
	out, err := os.Create(cacheFile)
	if err != nil {
		return "", fmt.Errorf("failed to create cache file: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to save image: %w", err)
	}

	return cacheFile, nil
}

// DownloadImageFromURL downloads an image from a direct URL and caches it locally
func DownloadImageFromURL(imageURL, cacheKey, cacheDir string) (string, error) {
	if imageURL == "" {
		return "", fmt.Errorf("empty image URL")
	}

	// Create cache filename
	cacheFile := filepath.Join(cacheDir, fmt.Sprintf("%s.jpg", cacheKey))

	// Check if cached file exists
	if _, err := os.Stat(cacheFile); err == nil {
		return cacheFile, nil
	}

	// Download the image
	req, err := http.NewRequest("GET", imageURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request for %s: %w", imageURL, err)
	}

	req.Header.Set("User-Agent", "go-scryfall/0.1.2")
	req.Header.Set("Accept", "image/*")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to download image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download image: HTTP %d", resp.StatusCode)
	}

	// Create cache directory if it doesn't exist
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create cache directory: %w", err)
	}

	// Save image to cache file
	out, err := os.Create(cacheFile)
	if err != nil {
		return "", fmt.Errorf("failed to create cache file: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to save image: %w", err)
	}

	return cacheFile, nil
}
