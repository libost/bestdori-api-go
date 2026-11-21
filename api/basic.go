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

func (d Difficulty) Idx() int {
	return indexHelper(d.String())
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
