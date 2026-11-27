package api

import (
	"context"
	"encoding/base64"
	"encoding/binary"
	"encoding/json/v2"
	"errors"
	"fmt"
	"net/http"
	"os"
)

var (
	// ErrNoCookies is returned when the API client does not have any cookies set, which is required for certain API calls.
	ErrNoCookies = errors.New("no cookies found in the API client; please log in first")
)

type UserProfile struct {
	Result   bool      `json:"result"`
	Profiles []Profile `json:"profiles"`
}

// Profile uses RLE.
type Profile struct {
	Data struct {
		Cards CardsProfile `json:"cards"`
		Items ItemsProfile `json:"items"`
	} `json:"data"`
	Name        string `json:"name"`
	Server      int    `json:"server"`
	Compression string `json:"compression"`
}

type CardsProfile struct {
	Eps      []int  `json:"eps"`
	Ids      string `json:"ids"`
	Arts     []int  `json:"arts"`
	Levels   []int  `json:"levels"`
	Skills   []int  `json:"skills"`
	Trains   []int  `json:"trains"`
	Masters  []int  `json:"masters"`
	Excludes []any  `json:"excludes"`
}

type ItemsProfile struct {
	Menu            []int `json:"Menu"`
	MyGo            []int `json:"MyGO"`
	Plaza           []int `json:"Plaza"`
	Roselia         []int `json:"Roselia"`
	Everyone        []int `json:"Everyone"`
	Magazine        []int `json:"Magazine"`
	Afterglow       []int `json:"Afterglow"`
	Morfonica       []int `json:"Morfonica"`
	Potentials      []int `json:"potentials"`
	PoppinParty     []int `json:"PoppinParty"`
	RaiseASuilen    []int `json:"RaiseASuilen"`
	PastelPalettes  []int `json:"PastelPalettes"`
	HelloHappyWorld []int `json:"HelloHappyWorld"`
}

// GetUserProfile fetches the user profile from the Bestdori API.
//
// It requires the API client to have cookies set (i.e., the user must be logged in).
// If the client does not have cookies, it returns an ErrNoCookies error.
func (a *Api) GetUserProfile(ctx context.Context) (*UserProfile, error) {
	if a.client.Jar == nil {
		return nil, ErrNoCookies
	}
	url := "https://bestdori.com/api/user/profiles"
	resp, err := a.do(ctx, http.MethodGet, url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var userProfile UserProfile
	if err = json.UnmarshalRead(resp.Body, &userProfile); err != nil {
		return nil, err
	}
	return &userProfile, nil
}

func decodeUserProfile(ids string) ([]uint16, error) {
	raw, err := base64.StdEncoding.DecodeString(ids)
	if err != nil {
		return nil, err
	}
	if len(raw)%2 != 0 {
		return nil, errors.New("invalid length of decoded data")
	}
	output := make([]uint16, len(raw)/2)
	for i := range output {
		output[i] = binary.LittleEndian.Uint16(raw[i*2 : i*2+2])
	}
	return output, nil
}

func encodeUserProfile(ids []uint16) string {
	raw := make([]byte, len(ids)*2)
	for i, id := range ids {
		binary.LittleEndian.PutUint16(raw[i*2:i*2+2], id)
	}
	return base64.StdEncoding.EncodeToString(raw)
}

func decodeProfileRLE(encoded []int, size int) ([]int, error) {
	if len(encoded)%2 != 0 {
		return nil, errors.New("invalid RLE data: expected count-value pairs")
	}
	if len(encoded) == 0 {
		return make([]int, size), nil
	}
	output := make([]int, 0, size)
	for i := 0; i < len(encoded); i += 2 {
		count, value := encoded[i], encoded[i+1]
		if count <= 0 || count > size-len(output) {
			return nil, errors.New("invalid RLE data: run length does not match expected size")
		}
		for range count {
			output = append(output, value)
		}
	}
	if len(output) != size {
		return nil, errors.New("invalid RLE data: run length does not match expected size")
	}
	return output, nil
}

func encodeProfileRLE(values []int) []int {
	if len(values) == 0 {
		return nil
	}
	encoded := make([]int, 0, len(values)*2)
	for start := 0; start < len(values); {
		end := start + 1
		for end < len(values) && values[end] == values[start] {
			end++
		}
		encoded = append(encoded, end-start, values[start])
		start = end
	}
	return encoded
}

// UploadUserProfile uploads the provided user profile to the Bestdori API.
//
// It requires the API client to have cookies set (i.e., the user must be logged in).
// If the client does not have cookies, it returns an ErrNoCookies error.
//
// If the upload fails, it returns an error with the status code and response body.
func (a *Api) UploadUserProfile(ctx context.Context, userProfile *UserProfile) error {
	if a.client.Jar == nil {
		return ErrNoCookies
	}
	url := "https://bestdori.com/api/user/profiles"
	body, err := json.Marshal(userProfile)
	if err != nil {
		return err
	}
	resp, err := a.doWithBody(ctx, http.MethodPost, url, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New("failed to upload user profile: " + resp.Status)
	}
	if err = json.UnmarshalRead(resp.Body, &userProfile); err != nil {
		return err
	}
	if !userProfile.Result {
		return errors.New("failed to upload user profile: result is false")
	}
	return nil
}

type SanitizedUserProfile struct {
	Name   string     `json:"name"`
	Server Region     `json:"server"`
	Cards  []CardInfo `json:"cards"`
	Items  ItemInfo   `json:"items"`
}

type CardInfo struct {
	ID      int  `json:"id"`
	Level   int  `json:"level"`
	Master  int  `json:"master"`
	Skill   int  `json:"skill"`
	Trained bool `json:"trained"`
	Eps     int  `json:"eps"`
	Arts    int  `json:"arts"`
}

type ItemInfo struct {
	PoppinParty     [7]int  `json:"PoppinParty"`
	Afterglow       [7]int  `json:"Afterglow"`
	HelloHappyWorld [7]int  `json:"HelloHappyWorld"`
	PastelPalettes  [7]int  `json:"PastelPalettes"`
	Roselia         [7]int  `json:"Roselia"`
	RaiseASuilen    [7]int  `json:"RaiseASuilen"`
	Morfonica       [7]int  `json:"Morfonica"`
	MyGO            [7]int  `json:"MyGO"`
	Everyone        [7]int  `json:"Everyone"`
	Magazine        [3]int  `json:"Magazine"`
	Plaza           [4]int  `json:"Plaza"`
	Menu            [4]int  `json:"Menu"`
	Potentials      [40]int `json:"potentials"`
}

// Sanitize processes the Profile data and returns a SanitizedUserProfile with structured card information.
func (p *Profile) Sanitize() (*SanitizedUserProfile, error) {
	ids := p.Data.Cards.Ids
	decodedIds, err := decodeUserProfile(ids)
	if err != nil {
		return nil, err
	}
	id := make([]int, len(decodedIds))
	for i, v := range decodedIds {
		id[i] = int(v)
	}
	result := &SanitizedUserProfile{
		Name:   p.Name,
		Server: regionHelper(p.Server),
		Cards:  make([]CardInfo, len(id)),
		Items:  ItemInfo{},
	}
	attributes := []struct {
		name string
		data []int
	}{
		{"eps", p.Data.Cards.Eps},
		{"arts", p.Data.Cards.Arts},
		{"levels", p.Data.Cards.Levels},
		{"skills", p.Data.Cards.Skills},
		{"trains", p.Data.Cards.Trains},
		{"masters", p.Data.Cards.Masters},
	}
	decoded := make([][]int, len(attributes))
	for i, attribute := range attributes {
		decoded[i], err = decodeProfileRLE(attribute.data, len(id))
		if err != nil {
			return nil, fmt.Errorf("decode card %s: %w", attribute.name, err)
		}
	}
	for i := range id {
		result.Cards[i] = CardInfo{
			ID:      id[i],
			Eps:     decoded[0][i],
			Arts:    decoded[1][i],
			Level:   decoded[2][i],
			Skill:   decoded[3][i],
			Trained: decoded[4][i] != 0,
			Master:  decoded[5][i],
		}
	}
	poppinParty, err := decodeProfileRLE(p.Data.Items.PoppinParty, len(result.Items.PoppinParty))
	if err != nil {
		return nil, fmt.Errorf("decode PoppinParty items: %w", err)
	}
	copy(result.Items.PoppinParty[:], poppinParty)
	afterglow, err := decodeProfileRLE(p.Data.Items.Afterglow, len(result.Items.Afterglow))
	if err != nil {
		return nil, fmt.Errorf("decode Afterglow items: %w", err)
	}
	copy(result.Items.Afterglow[:], afterglow)
	helloHappyWorld, err := decodeProfileRLE(p.Data.Items.HelloHappyWorld, len(result.Items.HelloHappyWorld))
	if err != nil {
		return nil, fmt.Errorf("decode HelloHappyWorld items: %w", err)
	}
	copy(result.Items.HelloHappyWorld[:], helloHappyWorld)
	pastelPalettes, err := decodeProfileRLE(p.Data.Items.PastelPalettes, len(result.Items.PastelPalettes))
	if err != nil {
		return nil, fmt.Errorf("decode PastelPalettes items: %w", err)
	}
	copy(result.Items.PastelPalettes[:], pastelPalettes)
	roselia, err := decodeProfileRLE(p.Data.Items.Roselia, len(result.Items.Roselia))
	if err != nil {
		return nil, fmt.Errorf("decode Roselia items: %w", err)
	}
	copy(result.Items.Roselia[:], roselia)
	raiseASuilen, err := decodeProfileRLE(p.Data.Items.RaiseASuilen, len(result.Items.RaiseASuilen))
	if err != nil {
		return nil, fmt.Errorf("decode RaiseASuilen items: %w", err)
	}
	copy(result.Items.RaiseASuilen[:], raiseASuilen)
	morfonica, err := decodeProfileRLE(p.Data.Items.Morfonica, len(result.Items.Morfonica))
	if err != nil {
		return nil, fmt.Errorf("decode Morfonica items: %w", err)
	}
	copy(result.Items.Morfonica[:], morfonica)
	myGO, err := decodeProfileRLE(p.Data.Items.MyGo, len(result.Items.MyGO))
	if err != nil {
		return nil, fmt.Errorf("decode MyGO items: %w", err)
	}
	copy(result.Items.MyGO[:], myGO)
	everyone, err := decodeProfileRLE(p.Data.Items.Everyone, len(result.Items.Everyone))
	if err != nil {
		return nil, fmt.Errorf("decode Everyone items: %w", err)
	}
	copy(result.Items.Everyone[:], everyone)
	magazine, err := decodeProfileRLE(p.Data.Items.Magazine, len(result.Items.Magazine))
	if err != nil {
		return nil, fmt.Errorf("decode Magazine items: %w", err)
	}
	copy(result.Items.Magazine[:], magazine)
	plaza, err := decodeProfileRLE(p.Data.Items.Plaza, len(result.Items.Plaza))
	if err != nil {
		return nil, fmt.Errorf("decode Plaza items: %w", err)
	}
	copy(result.Items.Plaza[:], plaza)
	menu, err := decodeProfileRLE(p.Data.Items.Menu, len(result.Items.Menu))
	if err != nil {
		return nil, fmt.Errorf("decode Menu items: %w", err)
	}
	copy(result.Items.Menu[:], menu)
	potentials, err := decodeProfileRLE(p.Data.Items.Potentials, len(result.Items.Potentials))
	if err != nil {
		return nil, fmt.Errorf("decode Potentials items: %w", err)
	}
	copy(result.Items.Potentials[:], potentials)
	return result, nil
}

func (s *SanitizedUserProfile) Desanitize() (*Profile, error) {
	ids := make([]uint16, len(s.Cards))
	for i, card := range s.Cards {
		if card.ID < 0 || card.ID > 65535 {
			return nil, fmt.Errorf("card %d has ID out of uint16 range: %d", i, card.ID)
		}
		ids[i] = uint16(card.ID)
	}

	profile := &Profile{
		Name:   s.Name,
		Server: s.Server.Idx(),
	}
	if profile.Server < 0 {
		return nil, fmt.Errorf("invalid profile server: %q", s.Server)
	}
	profile.Data.Cards.Ids = encodeUserProfile(ids)

	eps := make([]int, len(s.Cards))
	arts := make([]int, len(s.Cards))
	levels := make([]int, len(s.Cards))
	skills := make([]int, len(s.Cards))
	trains := make([]int, len(s.Cards))
	masters := make([]int, len(s.Cards))
	for i, card := range s.Cards {
		eps[i] = card.Eps
		arts[i] = card.Arts
		levels[i] = card.Level
		skills[i] = card.Skill
		masters[i] = card.Master
		if card.Trained {
			trains[i] = 1
		}
	}
	profile.Data.Cards.Eps = encodeProfileRLE(eps)
	profile.Data.Cards.Arts = encodeProfileRLE(arts)
	profile.Data.Cards.Levels = encodeProfileRLE(levels)
	profile.Data.Cards.Skills = encodeProfileRLE(skills)
	profile.Data.Cards.Trains = encodeProfileRLE(trains)
	profile.Data.Cards.Masters = encodeProfileRLE(masters)
	profile.Data.Items.PoppinParty = encodeProfileRLE(s.Items.PoppinParty[:])
	profile.Data.Items.Afterglow = encodeProfileRLE(s.Items.Afterglow[:])
	profile.Data.Items.HelloHappyWorld = encodeProfileRLE(s.Items.HelloHappyWorld[:])
	profile.Data.Items.PastelPalettes = encodeProfileRLE(s.Items.PastelPalettes[:])
	profile.Data.Items.Roselia = encodeProfileRLE(s.Items.Roselia[:])
	profile.Data.Items.RaiseASuilen = encodeProfileRLE(s.Items.RaiseASuilen[:])
	profile.Data.Items.Morfonica = encodeProfileRLE(s.Items.Morfonica[:])
	profile.Data.Items.MyGo = encodeProfileRLE(s.Items.MyGO[:])
	profile.Data.Items.Everyone = encodeProfileRLE(s.Items.Everyone[:])
	profile.Data.Items.Magazine = encodeProfileRLE(s.Items.Magazine[:])
	profile.Data.Items.Plaza = encodeProfileRLE(s.Items.Plaza[:])
	profile.Data.Items.Menu = encodeProfileRLE(s.Items.Menu[:])
	profile.Data.Items.Potentials = encodeProfileRLE(s.Items.Potentials[:])

	return profile, nil
}

// saveToFile saves the sanitized user profile to a JSON file with the specified filename.
//
// This function is primarily intended for testing purposes and may contain sensitive information. Ensure that the file is stored securely and not shared publicly.
//
// Therefore, it is possible that this function will be removed in the future without notice. Use it with caution.
func saveToFile(profile *SanitizedUserProfile, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := json.Marshal(profile)
	if err != nil {
		return err
	}
	_, err = file.Write(data)
	return err
}

type EventPlan struct {
	EventID string   `json:"event_id"`
	Region  Region   `json:"region"`
	Team    []string `json:"team"`
}

type Team struct {
	Members       []string `json:"members"`
	Pt            int      `json:"pt"`
	PtBonus       int      `json:"pt_bonus"`
	Ep            int      `json:"ep"`
	EpRoom        int      `json:"ep_room"`
	BandPower     int      `json:"band_power"`
	BandPowerRaw  int      `json:"band_power_raw"`
	BandItem      int      `json:"band_item"`
	AttrItem      string   `json:"attr_item"`
	ParameterItem string   `json:"parameter_item"`
}
