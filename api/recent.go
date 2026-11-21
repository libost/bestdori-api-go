package api

import (
	"context"
	"encoding/json/v2"
	"net/http"
	"time"
)

type Recent struct {
	api        *Api
	Songs      map[string]RecentSongs      `json:"songs"`
	Events     map[string]RecentEvents     `json:"events"`
	Gacha      map[string]RecentGachas     `json:"gacha"`
	LoginBonus map[string]RecentLoginBonus `json:"loginBonus"`
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
	for id, loginBonus := range r.LoginBonus {
		if !loginBonus.PublishedAt[region.Idx()].Time.IsZero() {
			if loginBonus.PublishedAt[region.Idx()].Time.Before(now) && loginBonus.ClosedAt[region.Idx()].Time.After(now) {
				loginBonus.id = id
				loginBonus.api = r.api
				results = append(results, loginBonus)
			}
		}
	}
	return results
}
