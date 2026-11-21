package api

import (
	"context"
	"encoding/json/v2"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	ErrCardNotFound   = errors.New("card not found")
	ErrNotValidCardId = errors.New("not a valid card id")
)

type Cards struct {
	api             *Api
	CharacterID     int         `json:"characterId"`
	Rarity          int         `json:"rarity"`
	Attribute       string      `json:"attribute"`
	LevelLimit      int         `json:"levelLimit"`
	ResourceSetName string      `json:"resourceSetName"`
	Prefix          []string    `json:"prefix"`
	ReleasedAt      []Timestamp `json:"releasedAt"`
	SkillID         int         `json:"skillId"`
	Type            string      `json:"type"`
	Stat            CardStat    `json:"stat"`
}

type CardStatValue struct {
	Performance int `json:"performance"`
	Technique   int `json:"technique"`
	Visual      int `json:"visual"`
}

type CardStat struct {
	Levels   map[string]CardStatValue `json:"-"`
	Episodes []CardStatValue          `json:"episodes"`
}

type CardsMap map[string]Cards

type Card struct {
	api             *Api
	CharacterID     int    `json:"characterId"`
	Rarity          int    `json:"rarity"`
	Attribute       string `json:"attribute"`
	LevelLimit      int    `json:"levelLimit"`
	ResourceSetName string `json:"resourceSetName"`
	SdResourceName  string `json:"sdResourceName"`
	Episodes        struct {
		Entries []struct {
			EpisodeId         int    `json:"episodeId"`
			EpisodeType       string `json:"episodeType"`
			SituationId       int    `json:"situationId"`
			ScenarioId        string `json:"scenarioId"`
			AppendPerformance int    `json:"appendPerformance"`
			AppendTechnique   int    `json:"appendTechnique"`
			AppendVisual      int    `json:"appendVisual"`
			ReleaseLevel      int    `json:"releaseLevel"`
			Costs             struct {
				ResourceId   int    `json:"resourceId"`
				ResourceType string `json:"resourceType"`
				Quantity     int    `json:"quantity"`
				LbBonus      int    `json:"lbBonus"`
			} `json:"costs"`
			Rewards struct {
				ResourceId   *int   `json:"resourceId,omitempty"`
				ResourceType string `json:"resourceType"`
				Quantity     int    `json:"quantity"`
				LbBonus      int    `json:"lbBonus"`
			} `json:"rewards"`
			Title       []*string `json:"title"`
			CharacterId int       `json:"characterId"`
		} `json:"entries"`
	} `json:"episodes"`
	CostumeID  int         `json:"costumeId"`
	GachaText  []string    `json:"gachaText"`
	Prefix     []string    `json:"prefix"`
	ReleasedAt []Timestamp `json:"releasedAt"`
	SkillName  []string    `json:"skillName"`
	SkillID    int         `json:"skillId"`
	Source     []struct {
		Gacha map[string]struct {
			Probability float64 `json:"probability"`
		} `json:"gacha"`
	} `json:"source"`
	Type  string                   `json:"type"`
	Stats map[string]CardStatValue `json:"stats"`
}

// InitCards returns a pointer to a CardsMap containing all cards from the API. It initializes the api field for each card in the map.
func (a *Api) InitCards(ctx context.Context) (*CardsMap, error) {
	resp, err := a.do(ctx, http.MethodGet, cardsURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var cardsMap CardsMap
	if err = json.UnmarshalRead(resp.Body, &cardsMap); err != nil {
		return nil, err
	}
	for cardID, card := range cardsMap {
		card.api = a
		cardsMap[cardID] = card
	}
	return &cardsMap, nil
}

// GetCard returns a pointer to a Card struct for the given cardId. It returns an error if the cardId is not valid or if the card is not found in the CardsMap.
//
// Be aware that the struct returned by this function is a summary of the card data and does not contain all the information available in the API.
//
// For a more detailed representation of the card, use the GetCard function of the Api struct, which fetches the full card data from the API.
func (c *CardsMap) GetCard(cardId string) (*Cards, error) {
	if card, ok := (*c)[cardId]; ok {
		return &card, nil
	}
	return nil, ErrCardNotFound
}

// IsAvailable checks if the card is available in the given region. It returns true if the card has a non-zero release time for the specified region, indicating that it is available.
//
// It is recommended to use this method before calling GetCardName, GetSkillName, or GetArtsURL to ensure that the card is available in the desired region.
func (c *Cards) IsAvailable(region Region) bool {
	idx := indexHelper(region.String())
	if c.ReleasedAt[idx].Time.IsZero() {
		return false
	}
	return true
}

func getCard(ctx context.Context, api *Api, cardId string) (*Card, error) {
	if _, err := strconv.Atoi(cardId); err != nil {
		return nil, ErrNotValidCardId
	}
	url := "https://bestdori.com/api/cards/" + cardId + ".json"
	resp, err := api.do(ctx, http.MethodGet, url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var card Card
	if err = json.UnmarshalRead(resp.Body, &card); err != nil {
		return nil, err
	}
	card.api = api
	return &card, nil
}

// GetCard returns a pointer to a Card struct for the given cardId. It returns an error if the cardId is not valid or if the card is not found in the API.
//
// This function fetches the full card data from the API, providing a detailed representation of the card, including all its attributes, stats, and related information.
func (a *Api) GetCard(ctx context.Context, cardId string) (*Card, error) {
	return getCard(ctx, a, cardId)
}

// GetCard returns a pointer to a Card struct and the bonus percentage of the card for specified event. It returns an error if the API call fails.
//
// This function fetches the full card data from the API, providing a detailed representation of the card, including all its attributes, stats, and related information.
func (e *EventMembers) GetCard(ctx context.Context) (*Card, int, error) {
	card, err := getCard(ctx, e.api, strconv.Itoa(e.SituationID))
	return card, e.Percent, err
}

// IsAvailable checks if the card is available in the given region. It returns true if the card has a non-zero release time for the specified region, indicating that it is available.
//
// It is recommended to use this method before calling GetCardName, GetSkillName, or GetArtsURL to ensure that the card is available in the desired region.
func (c *Card) IsAvailable(region Region) bool {
	idx := indexHelper(region.String())
	if c.ReleasedAt[idx].Time.IsZero() {
		return false
	}
	return true
}

// GetCardName returns the name of the card for the given region. It returns an empty string if the card is not available in the specified region.
func (c *Card) GetCardName(region Region) string {
	return c.Prefix[region.Idx()]
}

// GetCardName returns the name of the card for the given region. It returns an empty string if the card is not available in the specified region.
func (c *Cards) GetCardName(region Region) string {
	return c.Prefix[region.Idx()]
}

// GetSkillName returns the name of the skill for the given region. It returns an empty string if the card is not available in the specified region.
func (c *Card) GetSkillName(region Region) string {
	return c.SkillName[region.Idx()]
}

// GetRarity returns the rarity of the card.
func (c *Cards) GetRarity() int {
	return c.Rarity
}

// GetRarity returns the rarity of the card.
func (c *Card) GetRarity() int {
	return c.Rarity
}

// GetAttribute returns the attribute of the card.
func (c *Cards) GetAttribute() string {
	return strings.ToUpper(c.Attribute)
}

// GetAttribute returns the attribute of the card.
func (c *Card) GetAttribute() string {
	return strings.ToUpper(c.Attribute)
}

// GetLevelLimit returns the level limit of the card.
func (c *Cards) GetLevelLimit() int {
	return c.LevelLimit
}

// GetLevelLimit returns the level limit of the card.
func (c *Card) GetLevelLimit() int {
	return c.LevelLimit
}

// GetType returns the type of the card.
func (c *Cards) GetType() string {
	return c.Type
}

// GetType returns the type of the card.
func (c *Card) GetType() string {
	return c.Type
}

// GetArtsURL returns the URLs of the card arts for the specified region.
//
// It returns two URLs: the first one is for the normal card art, and the second one is for the after training card art.
//
// If any of the URLs is not available, it will return an empty string for that URL.
func (c *Card) GetArtsURL(ctx context.Context, region Region) (string, string) {
	stats := []string{"normal", "after_training"}
	result := make([]string, 0, len(stats))
	for _, cardStat := range stats {
		url := "https://bestdori.com/assets/" + region.String() + "/characters/resourceset/" + c.ResourceSetName + "_rip/card_" + cardStat + ".png"
		if isValidPNG(ctx, c.api, url) {
			result = append(result, url)
		}
	}
	return result[0], result[1]
}

// GetArtsURL returns the URLs of the card arts for the specified region.
//
// It returns two URLs: the first one is for the normal card art, and the second one is for the after training card art.
//
// If any of the URLs is not available, it will return an empty string for that URL.
func (c *Cards) GetArtsURL(ctx context.Context, region Region) (string, string) {
	stats := []string{"normal", "after_training"}
	result := make([]string, 0, len(stats))
	for _, cardStat := range stats {
		url := "https://bestdori.com/assets/" + region.String() + "/characters/resourceset/" + c.ResourceSetName + "_rip/card_" + cardStat + ".png"
		if isValidPNG(ctx, c.api, url) {
			result = append(result, url)
		}
	}
	return result[0], result[1]
}

// GetArtsFile returns the card arts as byte slices for the specified region.
//
// It returns two byte slices: the first one is for the normal card art, and the second one is for the after training card art.
//
// If any of the arts is not available, it will return nil for that art.
func (c *Card) GetArtsFile(ctx context.Context, region Region) ([]byte, []byte, error) {
	url1, url2 := c.GetArtsURL(ctx, region)
	urls := []string{url1, url2}
	arts := make([][]byte, 0, 2)
	for i, url := range urls {
		if url == "" {
			arts[i] = nil
			continue
		}
		resp, err := c.api.do(ctx, http.MethodGet, url)
		if err != nil {
			return nil, nil, err
		}
		data, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, nil, err
		}
		arts[i] = data
	}
	return arts[0], arts[1], nil
}

// GetArtsFile returns the card arts as byte slices for the specified region.
//
// It returns two byte slices: the first one is for the normal card art, and the second one is for the after training card art.
//
// If any of the arts is not available, it will return nil for that art.
func (c *Cards) GetArtsFile(ctx context.Context, region Region) ([]byte, []byte, error) {
	url1, url2 := c.GetArtsURL(ctx, region)
	urls := []string{url1, url2}
	arts := make([][]byte, 0, 2)
	for i, url := range urls {
		if url == "" {
			arts[i] = nil
			continue
		}
		resp, err := c.api.do(ctx, http.MethodGet, url)
		if err != nil {
			return nil, nil, err
		}
		data, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, nil, err
		}
		arts[i] = data
	}
	return arts[0], arts[1], nil
}

// GetReleaseTime returns the release time of the card for the given region.
func (c *Card) GetReleaseTime(region Region) time.Time {
	return c.ReleasedAt[region.Idx()].Time
}

// GetReleaseTime returns the release time of the card for the given region.
func (c *Cards) GetReleaseTime(region Region) time.Time {
	return c.ReleasedAt[region.Idx()].Time
}

// GetGachaText returns the gacha text of the card for the given region. It returns an empty string if the card is not available in the specified region.
func (c *Card) GetGachaText(region Region) string {
	return c.GachaText[region.Idx()]
}

// GetCard returns the card for the given NewCards. It returns an error if the card is not found or if the API call fails.
func (n *NewCards) GetCard(ctx context.Context) (*Card, error) {
	return getCard(ctx, n.a, strconv.Itoa(n.id))
}
