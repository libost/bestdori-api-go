package api

import (
	"context"
	"encoding/json/v2"
	"errors"
	"maps"
	"net/http"
	"strconv"
	"time"
)

type Characters struct {
	api           *Api
	etag          string
	CharacterType string   `json:"characterType"`
	CharacterName []string `json:"characterName"`
	FirstName     []string `json:"firstName"`
	LastName      []string `json:"lastName"`
	Nickname      []string `json:"nickname"`
	BandID        int      `json:"bandId"`
	ColorCode     string   `json:"colorCode"`
}

type CharactersMap map[string]Characters

type Character struct {
	api                  *Api
	CharacterType        string   `json:"characterType"`
	CharacterName        []string `json:"characterName"`
	FirstName            []string `json:"firstName"`
	LastName             []string `json:"lastName"`
	Nickname             []string `json:"nickname"`
	BandID               int      `json:"bandId"`
	ColorCode            string   `json:"colorCode"`
	SdAssetBundleName    string   `json:"sdAssetBundleName"`
	DefaultCostumeId     int      `json:"defaultCostumeId"`
	SeasonCostumeListMap struct {
		Entries struct {
			Season_3 struct {
				Entries []CharacterSeasonCostume `json:"entries"`
			} `json:"season_3"`
			Season_2 struct {
				Entries []CharacterSeasonCostume `json:"entries"`
			} `json:"season_2"`
			Season_1 struct {
				Entries []CharacterSeasonCostume `json:"entries"`
			} `json:"season_1"`
		} `json:"entries"`
	} `json:"seasonCostumeListMap"`
	Ruby    []string `json:"ruby"`
	Profile struct {
		CharacterVoice   []string `json:"characterVoice"`
		FavoriteFood     []string `json:"favoriteFood"`
		HatedFood        []string `json:"hatedFood"`
		Hobby            []string `json:"hobby"`
		SelfIntroduction []string `json:"selfIntroduction"`
		School           []string `json:"school"`
		SchoolCls        []string `json:"schoolCls"`
		// SchoolYear can be int or string
		SchoolYear    []Any     `json:"schoolYear"`
		Part          string    `json:"part"`
		Birthday      Timestamp `json:"birthday"`
		Constellation string    `json:"constellation"`
		Height        int       `json:"height"`
	} `json:"profile"`
}

type CharacterSeasonCostume struct {
	CharacterId           int    `json:"characterId"`
	BasicSeasonId         int    `json:"basicSeasonId"`
	CostumeType           string `json:"costumeType"`
	SeasonCostumeType     string `json:"seasonCostumeType"`
	SdAssetBundleName     string `json:"sdAssetBundleName"`
	Live2dAssetBundleName string `json:"live2dAssetBundleName"`
	SeasonType            string `json:"seasonType"`
}

// InitCharacters initializes the CharactersMap by fetching character data from the API.
//
// It returns a pointer to the CharactersMap.
//
// If there is an error during the API request or JSON unmarshalling, it returns an error.
func (a *Api) InitCharacters(ctx context.Context) (*CharactersMap, error) {
	resp, err := a.do(ctx, http.MethodGet, charactersURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var charactersMap CharactersMap
	if err = json.UnmarshalRead(resp.Body, &charactersMap); err != nil {
		return nil, err
	}
	for characterID, character := range charactersMap {
		character.api = a
		character.etag = resp.Header.Get("ETag")
		charactersMap[characterID] = character
	}
	return &charactersMap, nil
}

// UpdateCharacters updates the CharactersMap with the latest data from the API.
//
// It uses the ETag header to check if the data has changed since the last update. If the data has not changed, it does nothing. If the data has changed, it updates the CharactersMap with the new data and initializes the api field for each character in the map.
func (c *CharactersMap) UpdateCharacters(ctx context.Context) error {
	var api *Api
	var etag string
	for _, character := range *c {
		api = character.api
		if api == nil {
			continue
		}
		etag = character.etag
		if etag == "" {
			continue
		}
		break
	}
	if api == nil {
		return errors.New("api is nil")
	}
	resp, err := api.update(ctx, http.MethodGet, charactersURL, etag)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotModified {
		return nil
	}
	var charactersMap CharactersMap
	if err = json.UnmarshalRead(resp.Body, &charactersMap); err != nil {
		return err
	}
	for characterID, character := range charactersMap {
		character.api = api
		charactersMap[characterID] = character
	}
	maps.Copy(*c, charactersMap)
	return nil
}

// GetCharacter retrieves a Character from the CharactersMap based on the provided character ID.
//
// Be aware that the Characters struct that is returned is a summary of the character data, and does not include all the details that can be found in the Character struct.
//
// If the character ID is not found in the CharactersMap, it returns an error.
func (c *CharactersMap) GetCharacter(characterID int) (*Characters, error) {
	character, exists := (*c)[strconv.Itoa(characterID)]
	if !exists {
		return nil, errors.New("character not found")
	}
	return &character, nil
}

func (a *Api) GetCharacter(ctx context.Context, characterID int) (*Character, error) {
	return getCharacter(ctx, a, characterID)
}

// GetCharacter retrieves a Character struct from the API based on the Card.
func (c *Cards) GetCharacter(ctx context.Context) (*Character, error) {
	return getCharacter(ctx, c.api, c.CharacterID)
}

// GetCharacter retrieves a Character struct from the API based on the Card.
func (c *Card) GetCharacter(ctx context.Context) (*Character, error) {
	return getCharacter(ctx, c.api, c.CharacterID)
}

// GetCharacter retrieves a Character struct from the API based on the EventCharacters.
//
// It also returns the bonus percentage associated with the EventCharacters.
func (e *EventCharacters) GetCharacter(ctx context.Context) (*Character, int, error) {
	charId, _ := e.CharacterID.Int()
	char, err := getCharacter(ctx, e.api, charId)
	return char, e.Percent, err
}

// GetCharacter retrieves a Character struct from the API based on the EventsCharacters.
//
// It also returns the bonus percentage associated with the EventsCharacters.
func (e *EventsCharacters) GetCharacter(ctx context.Context) (*Character, int, error) {
	charId, _ := e.CharacterID.Int()
	char, err := getCharacter(ctx, e.api, charId)
	return char, e.Percent, err
}

func getCharacter(ctx context.Context, api *Api, characterID int) (*Character, error) {
	url := "https://bestdori.com/api/characters/" + strconv.Itoa(characterID) + ".json"
	resp, err := api.do(ctx, http.MethodGet, url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var character Character
	if err = json.UnmarshalRead(resp.Body, &character); err != nil {
		return nil, err
	}
	character.api = api
	return &character, nil
}

// IsAvailable checks if the character is available in the given region.
//
// It is recommended to use this method before calling any other methods that require the character to be available in the specified region.
func (c *Characters) IsAvailable(region Region) bool {
	idx := indexHelper(region.String())
	if c.CharacterName[idx] == "" {
		return false
	}
	return true
}

// GetCharacterName returns the character name for the given region.
func (c *Characters) GetCharacterName(region Region) string {
	return c.CharacterName[region.Idx()]
}

// GetCharacterFirstName returns the character's first name for the given region.
func (c *Characters) GetCharacterFirstName(region Region) string {
	return c.FirstName[region.Idx()]
}

// GetCharacterLastName returns the character's last name for the given region.
func (c *Characters) GetCharacterLastName(region Region) string {
	return c.LastName[region.Idx()]
}

// GetCharacterType returns the character type.
func (c *Characters) GetCharacterType() string {
	return c.CharacterType
}

// GetCharacterColor returns the character color code in hex format (e.g., "#FFFFFF").
func (c *Characters) GetCharacterColor() string {
	return c.ColorCode
}

// GetCharacterNickname returns the character nickname for the given region.
//
// Be aware that at this time, only member in the band "RAISE A SUILEN" has a nickname, and the rest of the characters will return an empty string.
func (c *Characters) GetCharacterNickname(region Region) string {
	return c.Nickname[region.Idx()]
}

// IsAvailable checks if the character is available in the given region.
//
// It is recommended to use this method before calling any other methods that require the character to be available in the specified region.
func (c *Character) IsAvailable(region Region) bool {
	if c.CharacterName[region.Idx()] == "" {
		return false
	}
	return true
}

// GetCharacterName returns the character name for the given region.
func (c *Character) GetCharacterName(region Region) string {
	return c.CharacterName[region.Idx()]
}

// GetCharacterFirstName returns the character's first name for the given region.
func (c *Character) GetCharacterFirstName(region Region) string {
	return c.FirstName[region.Idx()]
}

// GetCharacterLastName returns the character's last name for the given region.
func (c *Character) GetCharacterLastName(region Region) string {
	return c.LastName[region.Idx()]
}

// GetCharacterType returns the character type.
func (c *Character) GetCharacterType() string {
	return c.CharacterType
}

// GetCharacterColor returns the character color code in hex format (e.g., "#FFFFFF").
func (c *Character) GetCharacterColor() string {
	return c.ColorCode
}

// GetCharacterNickname returns the character nickname for the given region.
//
// Be aware that at this time, only member in the band "RAISE A SUILEN" has a nickname, and the rest of the characters will return an empty string.
func (c *Character) GetCharacterNickname(region Region) string {
	return c.Nickname[region.Idx()]
}

// GetCharacters returns a slice of Character structs for the given comics.
//
// It fetches the character data from the API based on the character IDs associated with the comics. If any error occurs during the API request or JSON unmarshalling, it returns an error.
func (c *Comics) GetCharacters(ctx context.Context) ([]*Character, error) {
	characters := make([]*Character, 0, len(c.CharacterId))
	for _, charId := range c.CharacterId {
		character, err := getCharacter(ctx, c.api, charId)
		if err != nil {
			return nil, err
		}
		characters = append(characters, character)
	}
	return characters, nil
}

// GetHeight returns the character's height in centimeters.
func (c *Character) GetHeight() int {
	return c.Profile.Height
}

// GetVoice returns the character's voice actor's name localized for the given region.
func (c *Character) GetVoice(region Region) string {
	return c.Profile.CharacterVoice[region.Idx()]
}

// GetFavoriteFood returns the character's favorite food localized for the given region.
func (c *Character) GetFavoriteFood(region Region) string {
	return c.Profile.FavoriteFood[region.Idx()]
}

// GetHatedFood returns the character's hated food localized for the given region.
func (c *Character) GetHatedFood(region Region) string {
	return c.Profile.HatedFood[region.Idx()]
}

// GetHobby returns the character's hobby localized for the given region.
func (c *Character) GetHobby(region Region) string {
	return c.Profile.Hobby[region.Idx()]
}

// GetSelfIntroduction returns the character's self-introduction localized for the given region.
func (c *Character) GetSelfIntroduction(region Region) string {
	return c.Profile.SelfIntroduction[region.Idx()]
}

// GetSchool returns the character's school localized for the given region.
func (c *Character) GetSchool(region Region) string {
	return c.Profile.School[region.Idx()]
}

// GetSchoolCls returns the character's school class localized for the given region.
func (c *Character) GetSchoolCls(region Region) string {
	return c.Profile.SchoolCls[region.Idx()]
}

// GetSchoolYear returns the character's school year localized for the given region.
func (c *Character) GetSchoolYear(region Region) string {
	return c.Profile.SchoolYear[region.Idx()].String()
}

// GetPart returns the character's part (e.g., vocal, guitar, base, drum, keyboard, dj and violin).
//
// It can return values like "guitar_vocal", "base_vocal"
func (c *Character) GetPart() string {
	return c.Profile.Part
}

// GetBirthday returns the character's birthday as a time.Month and day (int).
//
// The birthday is localized to the Asia/Tokyo timezone.
func (c *Character) GetBirthday() (time.Month, int) {
	loc, _ := time.LoadLocation("Asia/Tokyo")
	month := c.Profile.Birthday.Time.In(loc).Month()
	day := c.Profile.Birthday.Time.In(loc).Day()
	return month, day
}

// GetConstellation returns the character's constellation.
func (c *Character) GetConstellation() string {
	return c.Profile.Constellation
}
