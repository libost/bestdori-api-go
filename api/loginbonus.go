package api

import (
	"context"
	"encoding/json/v2"
	"net/http"
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
	api          *Api
	LoginBonusId int    `json:"loginBonusId"`
	Days         int    `json:"days"`
	ResourceType string `json:"resourceType"`
	Quantity     int    `json:"quantity"`
	Seq          int    `json:"seq"`
	GrantType    string `json:"grantType"`
	ResourceId   *int   `json:"resourceId,omitempty"`
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
	*l = loginBonusesMap
	return nil
}
