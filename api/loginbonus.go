package api

import (
	"context"
	"encoding/json/v2"
	"errors"
	"fmt"
	"io"
	"maps"
	"net/http"
	"time"
)

var (
	// ErrLoginBonusNotFound is returned when a login bonus is not found in the map.
	ErrLoginBonusNotFound = errors.New("login bonus not found")
)

type LoginBonus struct {
	api             *Api
	id              string
	LoginBonusType  string   `json:"loginBonusType"`
	AssetBundleName []string `json:"assetBundleName"`
	AssetMap        struct {
		Entries map[string]struct {
			LoginBonusId int    `json:"loginBonusId"`
			Days         int    `json:"days"`
			ResourceName string `json:"resourceName"`
			VoiceId      string `json:"voiceId"`
		} `json:"entries"`
	} `json:"assetMap"`
	Caption     []string               `json:"caption"`
	PublishedAt []Timestamp            `json:"publishedAt"`
	ClosedAt    []Timestamp            `json:"closedAt"`
	Details     [][]*LoginBonusDetails `json:"details"`
}

type LoginBonusDetails struct {
	LoginBonusId int `json:"loginBonusId"`
	Days         int `json:"days"`
	// ResourceType is the type of the resource granted by the login bonus. It can be "star", "stamp", "practice_ticket", etc.
	ResourceType string `json:"resourceType"`
	// Quantity is the amount of the resource granted by the login bonus.
	Quantity int `json:"quantity"`
	// Seq is the sequence number of the login bonus detail. It is used to determine which day should the bonus be granted.
	Seq int `json:"seq"`
	// GrantType is the way the bonus is granted. It can be "direct"(directly granted to the player), "present"(granted to the present box).
	GrantType string `json:"grantType"`
	// ResourceId is the ID of the resource granted by the login bonus. It is optional and may not be present for all resources. For example, it is used to identify the specific stamp.
	ResourceId *int `json:"resourceId,omitempty"`
}

type LoginBonuses struct {
	api             *Api
	etag            string
	id              string
	LoginBonusType  string      `json:"loginBonusType"`
	AssetBundleName []string    `json:"assetBundleName"`
	Caption         []string    `json:"caption"`
	PublishedAt     []Timestamp `json:"publishedAt"`
	ClosedAt        []Timestamp `json:"closedAt"`
}

type LoginBonusesMap map[string]LoginBonuses

// InitLoginBonuses initializes the login bonuses map by fetching the data from the API. It returns a pointer to the LoginBonusesMap and an error if any occurred during the process.
func (a *Api) InitLoginBonuses(ctx context.Context) (*LoginBonusesMap, error) {
	resp, err := a.do(ctx, http.MethodGet, loginBonusURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var loginBonusesMap LoginBonusesMap
	if err = json.UnmarshalRead(resp.Body, &loginBonusesMap); err != nil {
		return nil, err
	}
	for loginBonusID, loginBonus := range loginBonusesMap {
		loginBonus.etag = resp.Header.Get("ETag")
		loginBonus.id = loginBonusID
		loginBonus.api = a
		loginBonusesMap[loginBonusID] = loginBonus
	}
	return &loginBonusesMap, nil
}

// UpdateLoginBonuses updates the login bonuses map by fetching the data from the API. It returns an error if any occurred during the process.
func (l *LoginBonusesMap) UpdateLoginBonuses(ctx context.Context) error {
	var etag string
	var api *Api
	for _, loginBonus := range *l {
		etag = loginBonus.etag
		if etag == "" {
			continue
		}
		api = loginBonus.api
		if api == nil {
			continue
		}
		break
	}
	resp, err := api.update(ctx, http.MethodGet, loginBonusURL, etag)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotModified {
		return nil
	}
	var loginBonusesMap LoginBonusesMap
	if err = json.UnmarshalRead(resp.Body, &loginBonusesMap); err != nil {
		return err
	}
	for loginBonusID, loginBonus := range loginBonusesMap {
		loginBonus.etag = resp.Header.Get("ETag")
		loginBonus.id = loginBonusID
		loginBonus.api = api
		loginBonusesMap[loginBonusID] = loginBonus
	}
	maps.Copy(*l, loginBonusesMap)
	return nil
}

// GetLoginBonus retrieves a specific login bonus from the map using its ID. It returns a pointer to the LoginBonuses struct and an error if the login bonus is not found.
func (l *LoginBonusesMap) GetLoginBonus(loginBonusID string) (*LoginBonuses, error) {
	loginBonuses, ok := (*l)[loginBonusID]
	if !ok {
		return nil, ErrLoginBonusNotFound
	}
	return &loginBonuses, nil
}

type LoginBonusesTitleMap struct {
	Id     string
	Titles []string
}

// GetLoginBonusesTitleMap returns a slice of LoginBonusesTitleMap containing the ID and title of each login bonus in the map. It is useful for displaying a list of available login bonuses with their titles.
func (l *LoginBonusesMap) GetLoginBonusesTitleMap() []LoginBonusesTitleMap {
	titles := make([]LoginBonusesTitleMap, 0, len(*l))
	for id, loginBonus := range *l {
		titles = append(titles, LoginBonusesTitleMap{
			Id:     id,
			Titles: loginBonus.Caption,
		})
	}
	return titles
}

// IsAvailable checks if the login bonus is available for the specified region. It returns true if the login bonus is available, otherwise false.
func (l *LoginBonuses) IsAvailable(region Region) bool {
	return !l.PublishedAt[region.Idx()].IsZero()
}

// GetType returns the type of the login bonus.
func (l *LoginBonuses) GetType() string {
	return l.LoginBonusType
}

// GetAssetBundleNames returns the asset bundle names for the specified region.
func (l *LoginBonuses) GetAssetBundleNames(region Region) string {
	return l.AssetBundleName[region.Idx()]
}

// GetCaption returns the caption for the specified region.
func (l *LoginBonuses) GetCaption(region Region) string {
	return l.Caption[region.Idx()]
}

// GetStartTime returns the start time of the login bonus for the specified region.
func (l *LoginBonuses) GetStartTime(region Region) time.Time {
	return l.PublishedAt[region.Idx()].Time
}

// GetEndTime returns the end time of the login bonus for the specified region.
func (l *LoginBonuses) GetEndTime(region Region) time.Time {
	return l.ClosedAt[region.Idx()].Time
}

// GetBannerURL returns the URL of the banner image for the specified region. It checks if the image is valid and returns an empty string if not.
//
// If the asset bundle name is empty or the image is not valid, it returns an empty string.
func (l *LoginBonuses) GetBannerURL(ctx context.Context, region Region) string {
	assetBundleName := l.GetAssetBundleNames(region)
	if assetBundleName == "" {
		return ""
	}
	url := fmt.Sprintf("https://bestdori.com/assets/%s/event/loginbonus/%s_rip/background.png", region.String(), assetBundleName)
	if isValidPNG(ctx, l.api, url) {
		return url
	}
	return ""
}

// GetBannerFile returns the file content of the banner image for the specified region. It fetches the image from the URL and returns the raw bytes and an error if any occurred during the process.
//
// If the banner URL is empty or the image cannot be fetched, it returns an error.
//
// The format of the banner image is PNG.
func (l *LoginBonuses) GetBannerFile(ctx context.Context, region Region) ([]byte, error) {
	url := l.GetBannerURL(ctx, region)
	if url == "" {
		return nil, fmt.Errorf("no banner URL available for region %s", region.String())
	}
	resp, err := l.api.do(ctx, http.MethodGet, url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch banner: %s", resp.Status)
	}
	return io.ReadAll(resp.Body)
}

// GetLoginBonus retrieves the login bonus detail associated with the LoginBonuses. It returns a pointer to the LoginBonus struct and an error if any occurred during the process.
func (l *LoginBonuses) GetLoginBonus(ctx context.Context) (*LoginBonus, error) {
	return l.api.GetLoginBonus(ctx, l.id)
}

// GetLoginBonus retrieves a specific login bonus by its ID. It returns a pointer to the LoginBonus struct and an error if any occurred during the process.
func (a *Api) GetLoginBonus(ctx context.Context, loginBonusID string) (*LoginBonus, error) {
	url := fmt.Sprintf("https://bestdori.com/api/loginCampaigns/%s.json", loginBonusID)
	resp, err := a.do(ctx, http.MethodGet, url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var loginBonus LoginBonus
	if err = json.UnmarshalRead(resp.Body, &loginBonus); err != nil {
		return nil, err
	}
	loginBonus.id = loginBonusID
	loginBonus.api = a
	return &loginBonus, nil
}

// IsAvailable checks if the login bonus is available for the specified region. It returns true if the login bonus is available, otherwise false.
func (l *LoginBonus) IsAvailable(region Region) bool {
	return !l.PublishedAt[region.Idx()].IsZero()
}

// GetType returns the type of the login bonus.
func (l *LoginBonus) GetType() string {
	return l.LoginBonusType
}

// GetAssetBundleNames returns the asset bundle names for the specified region.
func (l *LoginBonus) GetAssetBundleNames(region Region) string {
	return l.AssetBundleName[region.Idx()]
}

// GetCaption returns the caption for the specified region.
func (l *LoginBonus) GetCaption(region Region) string {
	return l.Caption[region.Idx()]
}

// GetStartTime returns the start time of the login bonus for the specified region.
func (l *LoginBonus) GetStartTime(region Region) time.Time {
	return l.PublishedAt[region.Idx()].Time
}

// GetEndTime returns the end time of the login bonus for the specified region.
func (l *LoginBonus) GetEndTime(region Region) time.Time {
	return l.ClosedAt[region.Idx()].Time
}

// GetBannerURL returns the URL of the banner image for the specified region. It checks if the image is valid and returns an empty string if not.
//
// If the asset bundle name is empty or the image is not valid, it returns an empty string.
func (l *LoginBonus) GetBannerURL(ctx context.Context, region Region) string {
	assetBundleName := l.GetAssetBundleNames(region)
	if assetBundleName == "" {
		return ""
	}
	url := fmt.Sprintf("https://bestdori.com/assets/%s/event/loginbonus/%s_rip/background.png", region.String(), assetBundleName)
	if isValidPNG(ctx, l.api, url) {
		return url
	}
	return ""
}

// GetBannerFile returns the file content of the banner image for the specified region. It fetches the image from the URL and returns the raw bytes and an error if any occurred during the process.
//
// If the banner URL is empty or the image cannot be fetched, it returns an error.
//
// The format of the banner image is PNG.
func (l *LoginBonus) GetBannerFile(ctx context.Context, region Region) ([]byte, error) {
	url := l.GetBannerURL(ctx, region)
	if url == "" {
		return nil, fmt.Errorf("no banner URL available for region %s", region.String())
	}
	resp, err := l.api.do(ctx, http.MethodGet, url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch banner: %s", resp.Status)
	}
	return io.ReadAll(resp.Body)
}

// GetDetails returns the details of the login bonus for the specified region. It returns a slice of pointers to LoginBonusDetails structs.
func (l *LoginBonus) GetDetails(region Region) []*LoginBonusDetails {
	return l.Details[region.Idx()]
}

// GetResourceType returns the type of the resource granted by the login bonus detail.
func (d *LoginBonusDetails) GetResourceType() string {
	return d.ResourceType
}

// GetQuantity returns the amount of the resource granted by the login bonus detail.
func (d *LoginBonusDetails) GetQuantity() int {
	return d.Quantity
}

// GetSeq returns the sequence number of the login bonus detail. It is used to determine which day should the bonus be granted.
func (d *LoginBonusDetails) GetSeq() int {
	return d.Seq
}

// GetGrantType returns the way the bonus is granted. It can be "direct"(directly granted to the player), "present"(granted to the present box).
func (d *LoginBonusDetails) GetGrantType() string {
	return d.GrantType
}

// GetResourceId returns the ID of the resource granted by the login bonus detail. It is optional and may not be present for all resources. For example, it is used to identify the specific stamp.
func (d *LoginBonusDetails) GetResourceId() *int {
	return d.ResourceId
}

// GetLoginBonus retrieves the login bonus associated with the RecentLoginBonus. It returns a pointer to the LoginBonus struct and an error if any occurred during the process.
func (r *RecentLoginBonus) GetLoginBonus(ctx context.Context) (*LoginBonus, error) {
	return r.api.GetLoginBonus(ctx, r.id)
}

// GetCaption returns the caption of the recent login bonus for the specified region.
func (r *RecentLoginBonus) GetCaption(region Region) string {
	return r.Caption[region.Idx()]
}

// GetStartTime returns the start time of the recent login bonus for the specified region.
func (r *RecentLoginBonus) GetStartTime(region Region) time.Time {
	return r.PublishedAt[region.Idx()].Time
}

// GetEndTime returns the end time of the recent login bonus for the specified region.
func (r *RecentLoginBonus) GetEndTime(region Region) time.Time {
	return r.ClosedAt[region.Idx()].Time
}
