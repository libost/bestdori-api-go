package api

import (
	"bytes"
	"context"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"time"
)

const (
	bandsURL      = "https://bestdori.com/api/bands/all.1.json"
	cardsURL      = "https://bestdori.com/api/cards/all.5.json"
	charactersURL = "https://bestdori.com/api/characters/main.3.json"
	eventsURL     = "https://bestdori.com/api/events/all.5.json"
	newsURL       = "https://bestdori.com/api/news/dynamic/recent.json"
	skillsURL     = "https://bestdori.com/api/skills/all.10.json"
	gachaURL      = "https://bestdori.com/api/gacha/all.5.json"
	ratesURL      = "https://bestdori.com/api/tracker/rates.json"
	songsURL      = "https://bestdori.com/api/songs/all.5.json"
	comicsURL     = "https://bestdori.com/api/comics/all.5.json"
	songsMetaURL  = "https://bestdori.com/api/songs/meta/all.5.json"
	loginBonusURL = "https://bestdori.com/api/loginCampaigns/all.5.json"
)

// Tiers for China
var CNTierList = [...]int{
	20, 30, 40, 50, 100, 200, 300, 400, 500, 1000, 1500, 2000, 3000, 4000, 5000, 10000, 15000, 20000, 30000, 40000, 50000, 100000, 150000, 200000, 300000,
}

// Tiers for Japan
var JPTierList = [...]int{
	20, 30, 40, 50, 100, 200, 300, 400, 500, 1000, 2000, 5000, 10000, 20000, 30000, 50000, 70000, 100000, 300000,
}

// Tiers for Taiwan
var TWTierList = [...]int{
	20, 30, 40, 50, 100, 500, 1000, 3000, 5000,
}

// Tiers for International(English)
var ENTierList = [...]int{
	20, 30, 40, 50, 100, 200, 300, 400, 500, 1000, 2000, 2500, 3000, 4000, 5000, 10000,
}

// Tiers for Korea
var KRTierList = [...]int{
	100,
}

type Tier interface {
	Region() Region
	Tiers() []int
}

// RegionTier identifies a ranking tier in a game region.
type RegionTier struct {
	Region Region
	Tier   int
}

// NewRegionTier creates a region/tier parameter and validates that the tier
// is supported by the region.
func NewRegionTier(region Region, tier int) (RegionTier, error) {
	regionTier := RegionTier{Region: region, Tier: tier}
	if err := regionTier.Validate(); err != nil {
		return RegionTier{}, err
	}
	return regionTier, nil
}

// Validate checks that the region and tier combination is supported.
func (rt RegionTier) Validate() error {
	if slices.Contains(rt.Region.Tiers(), rt.Tier) {
		return nil
	}
	return fmt.Errorf("invalid tier %d for region %s", rt.Tier, rt.Region)
}

type Difficulty string

const (
	DiffcultyEasy    Difficulty = "easy"
	DiffcultyNormal  Difficulty = "normal"
	DiffcultyHard    Difficulty = "hard"
	DiffcultyExpert  Difficulty = "expert"
	DiffcultySpecial Difficulty = "special"
)

func (d Difficulty) String() string {
	return string(d)
}

// Color returns the color associated with the difficulty level.
//
// The color is represented as a hexadecimal string.
func (d Difficulty) Color() string {
	return difficultyColors[d.Idx()]
}

func (d Difficulty) Idx() int {
	return indexHelper(d.String())
}

var difficultyColors = [...]string{
	"#3076F8", // Easy
	"#188A28", // Normal
	"#D1861A", // Hard
	"#FD2A30", // Expert
	"#F12199", // Special
}

type CardStatus string

const (
	CardStatusNormal        CardStatus = "normal"
	CardStatusAfterTraining CardStatus = "after_training"
)

func (cs CardStatus) String() string {
	return string(cs)
}

// Tiers returns the supported ranking tiers for a region.
func (r Region) Tiers() []int {
	switch r {
	case RegionJP:
		return JPTierList[:]
	case RegionEN:
		return ENTierList[:]
	case RegionKR:
		return KRTierList[:]
	case RegionTW:
		return TWTierList[:]
	case RegionCN:
		return CNTierList[:]
	default:
		return nil
	}
}

// Region represents the region of the game.
type Region string

const (
	RegionJP Region = "jp"
	RegionEN Region = "en"
	RegionKR Region = "kr"
	RegionTW Region = "tw"
	RegionCN Region = "cn"
)

func (r Region) String() string {
	return string(r)
}

func (r Region) Idx() int {
	return indexHelper(r.String())
}

// Any is a type that can hold any JSON value, including strings and numbers. It is used to handle fields that can have different types in the JSON data.
type Any string

func (a *Any) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	if bytes.Equal(data, []byte("null")) {
		*a = ""
		return nil
	}

	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		*a = Any(str)
		return nil
	}

	var number jsontext.Value
	if err := json.Unmarshal(data, &number); err != nil {
		return err
	}
	if number.Kind() != jsontext.KindNumber {
		return fmt.Errorf("api.Any: expected JSON string or number, got %s", number.Kind())
	}
	*a = Any(number.String())
	return nil
}

// String returns the string representation of the Any value.
func (a Any) String() string {
	return string(a)
}

// Int attempts to convert the Any value to an int. It returns an error if the conversion fails.
func (a Any) Int() (int, error) {
	return strconv.Atoi(string(a))
}

type Timestamp struct {
	time.Time
}

func (t *Timestamp) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	if bytes.Equal(data, []byte("null")) {
		t.Time = time.Time{}
		return nil
	}

	var ms int64
	if len(data) > 0 && data[0] == '"' {
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		if s == "" {
			t.Time = time.Time{}
			return nil
		}
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return err
		}
		ms = v
	} else {
		if err := json.Unmarshal(data, &ms); err != nil {
			return err
		}
	}
	if ms == 0 {
		t.Time = time.Time{}
		return nil
	}
	t.Time = time.UnixMilli(ms).UTC()
	return nil
}

func (t Timestamp) MarshalJSON() ([]byte, error) {
	if t.Time.IsZero() {
		return []byte("null"), nil
	}
	return strconv.AppendInt(nil, t.UnixMilli(), 10), nil
}

type Timeduration struct {
	time.Duration
}

func (d *Timeduration) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	if bytes.Equal(data, []byte("null")) {
		d.Duration = 0
		return nil
	}

	var seconds float64
	if err := json.Unmarshal(data, &seconds); err != nil {
		return err
	}
	d.Duration = time.Duration(seconds * float64(time.Second))
	return nil
}

func (d Timeduration) MarshalJSON() ([]byte, error) {
	if d.Duration == 0 {
		return []byte("null"), nil
	}
	seconds := d.Duration.Seconds()
	return fmt.Appendf(nil, "%.3f", seconds), nil
}

type BestdoriAPIError struct {
	StatusCode int
	Message    string
}

type Api struct {
	client  *http.Client
	proxy   string
	timeout time.Duration
}

func NewApi(proxy string, timeout time.Duration) *Api {
	if timeout <= 0 {
		timeout = 10 * time.Second
	}

	var transport *http.Transport
	if defaultTransport, ok := http.DefaultTransport.(*http.Transport); ok {
		transport = defaultTransport.Clone()
	} else {
		transport = &http.Transport{Proxy: http.ProxyFromEnvironment}
	}
	if proxy != "" {
		proxyURL, err := url.Parse(proxy)
		if err == nil && proxyURL.Scheme != "" && proxyURL.Host != "" {
			transport.Proxy = http.ProxyURL(proxyURL)
		}
	}

	return &Api{
		client:  &http.Client{Transport: transport, Timeout: timeout},
		proxy:   proxy,
		timeout: timeout,
	}
}

var defaultApi = NewApi("", 10*time.Second)

func (a *Api) do(ctx context.Context, method, requestURL string) (*http.Response, error) {
	if a == nil || a.client == nil {
		a = defaultApi
	}

	req, err := http.NewRequestWithContext(ctx, method, requestURL, nil)
	if err != nil {
		return nil, err
	}
	return a.client.Do(req)
}

func (a *Api) update(ctx context.Context, method, requestURL, etag string) (*http.Response, error) {
	if a == nil || a.client == nil {
		a = defaultApi
	}

	req, err := http.NewRequestWithContext(ctx, method, requestURL, nil)
	if err != nil {
		return nil, err
	}
	if etag != "" {
		req.Header.Set("If-None-Match", etag)
	}
	return a.client.Do(req)
}

func isValidPNG(ctx context.Context, api *Api, url string) bool {
	resp, err := api.do(ctx, http.MethodHead, url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return false
	}
	if contentType := resp.Header.Get("Content-Type"); contentType != "image/png" {
		return false
	}
	return true
}

const (
	card1PowerfulFrame = "https://bestdori.com/res/image/card-1-powerful.png"
	card1CoolFrame     = "https://bestdori.com/res/image/card-1-cool.png"
	card1HappyFrame    = "https://bestdori.com/res/image/card-1-happy.png"
	card1PureFrame     = "https://bestdori.com/res/image/card-1-pure.png"
	card2Frame         = "https://bestdori.com/res/image/card-2.png"
	card3Frame         = "https://bestdori.com/res/image/card-3.png"
	card4Frame         = "https://bestdori.com/res/image/card-4.png"
	card5Frame         = "https://bestdori.com/res/image/card-5.png"
)

var cardFrame1 = []string{
	card1PowerfulFrame,
	card1CoolFrame,
	card1HappyFrame,
	card1PureFrame,
}

var cardFrame = []string{
	card2Frame,
	card3Frame,
	card4Frame,
	card5Frame,
}

var bandFrames = [...]string{
	"https://bestdori.com/res/icon/band_1.svg",
	"https://bestdori.com/res/icon/band_2.svg",
	"https://bestdori.com/res/icon/band_3.svg",
	"https://bestdori.com/res/icon/band_4.svg",
	"https://bestdori.com/res/icon/band_5.svg", // 250*250
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"https://bestdori.com/res/icon/band_18.svg", // 18 44*44
	"",
	"",
	"https://bestdori.com/res/icon/band_21.svg", // 21 44*44
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"https://bestdori.com/res/icon/band_45.svg", // 45
}

const (
	powerfulFrame = "https://bestdori.com/res/icon/powerful.svg"
	coolFrame     = "https://bestdori.com/res/icon/cool.svg"
	happyFrame    = "https://bestdori.com/res/icon/happy.svg"
	pureFrame     = "https://bestdori.com/res/icon/pure.svg"
)

var attributeFrames = [...]string{
	powerfulFrame,
	coolFrame,
	happyFrame,
	pureFrame,
}
