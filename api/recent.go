package api

import (
	"context"
	"encoding/json/v2"
	"net/http"
	"time"
)

type Recent struct {
	api        *Api
	etag       string
	Songs      map[string]RecentSongs      `json:"songs"`
	Events     map[string]RecentEvents     `json:"events"`
	Gacha      map[string]RecentGachas     `json:"gacha"`
	LoginBonus map[string]RecentLoginBonus `json:"loginBonuses"`
}

type RecentSongs struct {
	id          string
	api         *Api
	MusicTitle  []string    `json:"musicTitle"`
	PublishedAt []Timestamp `json:"publishedAt"`
}

type RecentEvents struct {
	id        string
	api       *Api
	EventName []string    `json:"eventName"`
	StartAt   []Timestamp `json:"startAt"`
	EndAt     []Timestamp `json:"endAt"`
}

type RecentGachas struct {
	id          string
	api         *Api
	GachaName   []string    `json:"gachaName"`
	PublishedAt []Timestamp `json:"publishedAt"`
	ClosedAt    []Timestamp `json:"closedAt"`
}

type RecentLoginBonus struct {
	id          string
	api         *Api
	Caption     []string    `json:"caption"`
	PublishedAt []Timestamp `json:"publishedAt"`
	ClosedAt    []Timestamp `json:"closedAt"`
}

// GetRecent retrieves the recent data from the API, including recent songs, events, gachas, and login bonuses. It returns a pointer to the Recent struct and an error if any occurs during the API request or JSON unmarshalling.
func (a *Api) GetRecent(ctx context.Context) (*Recent, error) {
	resp, err := a.do(ctx, http.MethodGet, newsURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var recent Recent
	if err = json.UnmarshalRead(resp.Body, &recent); err != nil {
		return nil, err
	}
	recent.api = a
	return &recent, nil
}

// UpdateRecent updates the Recent struct with the latest data from the API. It uses the ETag header to check if the data has changed since the last update. If the data has not changed, it does nothing. If the data has changed, it updates the Recent struct with the new data and initializes the api field.
func (r *Recent) UpdateRecent(ctx context.Context) error {
	resp, err := r.api.update(ctx, http.MethodGet, newsURL, r.etag)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotModified {
		return nil
	}
	var recent Recent
	if err = json.UnmarshalRead(resp.Body, &recent); err != nil {
		return err
	}
	recent.api = r.api
	recent.etag = resp.Header.Get("ETag")
	*r = recent
	return nil
}

// GetRecentSongsForRegion retrieves the recent songs for a given region. It returns a slice of RecentSongs that are published in the specified region.
//
// Be aware that if region is KR, the function will return an empty slice, as KR has been shut down.
func (r *Recent) GetRecentSongsForRegion(region Region) []RecentSongs {
	results := make([]RecentSongs, 0)
	for id, song := range r.Songs {
		if !song.PublishedAt[region.Idx()].Time.IsZero() {
			if shouldCollectSong(song.PublishedAt, region.Idx()) {
				song.id = id
				song.api = r.api
				results = append(results, song)
			}
		}
	}
	return results
}

func shouldCollectSong(publishedAt []Timestamp, idx int) bool {
	currentPubTime := publishedAt[idx].Time.Unix()

	validValueCount := 0
	for _, value := range publishedAt {
		if value.Time.IsZero() {
			continue
		}
		pubTime := value.Time.Unix()
		validValueCount++
		if pubTime > currentPubTime {
			return false
		}
	}

	return validValueCount > 0
}

// GetRecentEventsForRegion retrieves the recent events for a given region. It returns a slice of RecentEvents that are currently active or upcoming in the specified region.
//
// The function first checks for events that are currently active (i.e., the current time is between the event's start and end times). If no active events are found, it then checks for upcoming events (i.e., events that have not yet started or have already ended).
//
// Be aware that if region is KR, the function will return an empty slice, as KR has been shut down.
func (r *Recent) GetRecentEventsForRegion(region Region) []RecentEvents {
	results := make([]RecentEvents, 0)
	now := time.Now()
	for id, event := range r.Events {
		if !event.StartAt[region.Idx()].Time.IsZero() {
			if event.StartAt[region.Idx()].Time.Before(now) && event.EndAt[region.Idx()].Time.After(now) {
				event.id = id
				event.api = r.api
				results = append(results, event)
			}
		}
		if len(results) > 0 {
			break
		}
	}
	if len(results) == 0 {
		for id, event := range r.Events {
			if !event.StartAt[region.Idx()].Time.IsZero() {
				if event.StartAt[region.Idx()].Time.After(now) || event.EndAt[region.Idx()].Time.Before(now) {
					event.id = id
					event.api = r.api
					results = append(results, event)
				}
			}
		}
	}
	if len(results) > 1 {
		maxVal := 0
		maxIdx := 0
		for r, event := range results {
			if event.StartAt[region.Idx()].Time.Unix() > int64(maxVal) {
				maxVal = int(event.StartAt[region.Idx()].Time.Unix())
				maxIdx = r
			}
		}
		results = []RecentEvents{results[maxIdx]}
	}
	return results
}

// GetRecentGachasForRegion retrieves the recent gachas for a given region. It returns a slice of RecentGachas that are currently active in the specified region.
//
// Be aware that if region is KR, the function will return an empty slice, as KR has been shut down.
func (r *Recent) GetRecentGachasForRegion(region Region) []RecentGachas {
	results := make([]RecentGachas, 0)
	now := time.Now()
	for id, gacha := range r.Gacha {
		if !gacha.PublishedAt[region.Idx()].Time.IsZero() {
			if gacha.PublishedAt[region.Idx()].Time.Before(now) && gacha.ClosedAt[region.Idx()].Time.After(now) {
				gacha.id = id
				gacha.api = r.api
				results = append(results, gacha)
			}
		}
	}
	return results
}

// GetRecentLoginBonusForRegion retrieves the recent login bonuses for a given region. It returns a slice of RecentLoginBonus that are currently active in the specified region.
//
// Be aware that if region is KR, the function will return an empty slice, as KR has been shut down.
func (r *Recent) GetRecentLoginBonusForRegion(region Region) []RecentLoginBonus {
	results := make([]RecentLoginBonus, 0)
	now := time.Now()
	idx := region.Idx()
	if idx < 0 {
		return results
	}
	for id, loginBonus := range r.LoginBonus {
		if idx >= len(loginBonus.PublishedAt) || idx >= len(loginBonus.ClosedAt) {
			continue
		}
		if !loginBonus.PublishedAt[idx].Time.IsZero() {
			if loginBonus.PublishedAt[idx].Time.Before(now) && loginBonus.ClosedAt[idx].Time.After(now) {
				loginBonus.id = id
				loginBonus.api = r.api
				results = append(results, loginBonus)
			}
		}
	}
	return results
}

// GetTitle returns the title of the recent event for the specified region.
func (e *RecentEvents) GetTitle(region Region) string {
	return e.EventName[region.Idx()]
}

// GetStartTime returns the start time of the recent event for the specified region.
func (e *RecentEvents) GetStartTime(region Region) time.Time {
	return e.StartAt[region.Idx()].Time
}

// GetEndTime returns the end time of the recent event for the specified region.
func (e *RecentEvents) GetEndTime(region Region) time.Time {
	return e.EndAt[region.Idx()].Time
}

// GetTitle returns the title of the recent gacha for the specified region.
func (g *RecentGachas) GetTitle(region Region) string {
	return g.GachaName[region.Idx()]
}

// GetPublishTime returns the publish time of the recent gacha for the specified region.
func (g *RecentGachas) GetPublishTime(region Region) time.Time {
	return g.PublishedAt[region.Idx()].Time
}

// GetCloseTime returns the close time of the recent gacha for the specified region.
func (g *RecentGachas) GetCloseTime(region Region) time.Time {
	return g.ClosedAt[region.Idx()].Time
}

// GetTitle returns the title of the recent song for the specified region.
func (s *RecentSongs) GetTitle(region Region) string {
	return s.MusicTitle[region.Idx()]
}

// GetPublishTime returns the publish time of the recent song for the specified region.
func (s *RecentSongs) GetPublishTime(region Region) time.Time {
	return s.PublishedAt[region.Idx()].Time
}
