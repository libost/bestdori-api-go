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
	MyGo            []int `json:"MyGo"`
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
	Name   string         `json:"name"`
	Server Region         `json:"server"`
	Cards  []CardInfo     `json:"cards"`
	Items  map[string]int `json:"items"`
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
		Items:  make(map[string]int),
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
	return result, nil
}

func (s *SanitizedUserProfile) Desanitize() (*Profile, error) {
	if len(s.Items) > 0 {
		return nil, errors.New("cannot desanitize items: no item-to-profile mapping is defined")
	}

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

// func (s *SanitizedUserProfile) EventPlanner(event *Event) {}

// "items":{"PoppinParty":[1,1,2,4,1,1,1,4,1,1,1,4],"Afterglow":[1,4,1,1,1,4,1,1,1,2,2,4],"HelloHappyWorld":[1,0,1,6,1,4,1,2,3,4],"PastelPalettes":[1,1,1,6,1,4,1,1,2,4,1,2],"Roselia":[1,4,1,2,1,4,1,2,1,4,1,6,1,4],"RaiseASuilen":[2,4,1,1,1,4,1,2,2,4],"Morfonica":[3,4,1,2,1,4,2,2],"MyGO":[1,2,2,4,1,2,3,4],"Everyone":[2,4,1,1,1,4,1,1,2,4],"Magazine":[1,4,1,1,1,4],"Plaza":[1,4,1,2,2,4],"Menu":[1,4,1,2,2,4],"potentials":[40,1]}}}
