// This file contains an algorithm ported from [Tsugu BangDream Bot](https://github.com/Yamamoto-2/tsugu-bangdream-bot).
//
// Original Copyright (c) 2023 Yamamoto-2
// Licensed under the MIT License. See [NOTICE](./../NOTICE) for details.

package api

import (
	"context"
	"encoding/json/v2"
	"errors"
	"fmt"
	"io"
	"maps"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"
)

var (
	ErrEventNotFound     = errors.New("event not found")
	ErrNotValidEventId   = errors.New("not a valid event id")
	ErrTooLessDataPoints = errors.New("too less data points to predict")
)

type Events struct {
	api                   *Api
	id                    string
	etag                  string
	EventType             string             `json:"eventType"`
	EventName             []string           `json:"eventName"`
	AssetBundleName       string             `json:"assetBundleName"`
	BannerAssetBundleName string             `json:"bannerAssetBundleName"`
	StartAt               []Timestamp        `json:"startAt"`
	EndAt                 []Timestamp        `json:"endAt"`
	Attributes            []Attributes       `json:"attributes"`
	Characters            []EventsCharacters `json:"characters"`
	Members               []EventMembers     `json:"members"`
	LimitBreaks           []EventLimitBreaks `json:"limitBreaks"`
	RewardCards           []int              `json:"rewardCards"`
}

type Attributes struct {
	Attribute string `json:"attribute"`
	Percent   int    `json:"percent"`
}

type EventsCharacters struct {
	api *Api
	// CharacterID can be int or string
	CharacterID Any `json:"characterId"`
	Percent     int `json:"percent"`
}

type EventsMap map[string]Events

type Event struct {
	api                   *Api
	id                    string
	EventType             string      `json:"eventType"`
	EventName             []string    `json:"eventName"`
	AssetBundleName       string      `json:"assetBundleName"`
	BannerAssetBundleName string      `json:"bannerAssetBundleName"`
	StartAt               []Timestamp `json:"startAt"`
	EndAt                 []Timestamp `json:"endAt"`
	// We don't know what is `enableFlag` for, it almost always be `null`, so we just use `*bool` to represent it, and it can be `nil` if it's `null`.
	EnableFlag          []*bool     `json:"enableFlag"`
	PublicStartAt       []Timestamp `json:"publicStartAt"`
	PublicEndAt         []Timestamp `json:"publicEndAt"`
	DistributionStartAt []Timestamp `json:"distributionStartAt"`
	DistributionEndAt   []Timestamp `json:"distributionEndAt"`
	BgmAssetBundleName  string      `json:"bgmAssetBundleName"`
	BgmFileName         string      `json:"bgmFileName"`
	AggregateEndAt      []Timestamp `json:"aggregateEndAt"`
	ExchangeEndAt       []Timestamp `json:"exchangeEndAt"`
	PointRewards        [][]struct {
		// Point can be int or string
		Point          Any    `json:"point"`
		RewardType     string `json:"rewardType"`
		RewardId       int    `json:"rewardId"`
		RewardQuantity int    `json:"rewardQuantity"`
	} `json:"pointRewards"`
	RankingRewards [][]struct {
		FromRank       int    `json:"fromRank"`
		ToRank         int    `json:"toRank"`
		RewardType     string `json:"rewardType"`
		RewardId       int    `json:"rewardId"`
		RewardQuantity int    `json:"rewardQuantity"`
	} `json:"rankingRewards"`
	Attributes                      []EventAttribute  `json:"attributes"`
	Characters                      []EventCharacters `json:"characters"`
	EventAttributeAndCharacterBonus struct {
		EventID          int `json:"eventId"`
		PointPercent     int `json:"pointPercent"`
		ParameterPercent int `json:"parameterPercent"`
	} `json:"eventAttributeAndCharacterBonus"`
	Members     []EventMembers     `json:"members"`
	LimitBreaks []EventLimitBreaks `json:"limitBreaks"`
	Stories     []EventStory       `json:"stories"`
	RewardCards []int              `json:"rewardCards"`
}

type EventStory struct {
	api             *Api
	ScenarioId      string `json:"scenarioId"`
	CoverImage      string `json:"coverImage"`
	BackgroundImage string `json:"backgroundImage"`
	// ReleasePt can be int or string
	ReleasePt Any `json:"releasePt"`
	Rewards   []struct {
		RewardType     string `json:"rewardType"`
		RewardId       int    `json:"rewardId"`
		RewardQuantity int    `json:"rewardQuantity"`
	} `json:"rewards"`
	Caption           []string `json:"caption"`
	Title             []string `json:"title"`
	Synopsis          []string `json:"synopsis"`
	ReleaseConditions []string `json:"releaseConditions"`
}

type EventCharacters struct {
	api     *Api
	EventID int `json:"eventId"`
	// CharacterID can be int or string
	CharacterID Any `json:"characterId"`
	Percent     int `json:"percent"`
	Seq         int `json:"seq"`
}

type EventMembers struct {
	api         *Api
	EventID     int `json:"eventId"`
	SituationID int `json:"situationId"`
	Percent     int `json:"percent"`
	Seq         int `json:"seq"`
}

type EventLimitBreaks struct {
	EventID int     `json:"eventId"`
	Rarity  int     `json:"rarity"`
	Rank    int     `json:"rank"`
	Percent float64 `json:"percent"`
}

type EventAttribute struct {
	EventID   int    `json:"eventId"`
	Attribute string `json:"attribute"`
	Percent   int    `json:"percent"`
}

type EventTracker struct {
	Cutoffs []Cutoff `json:"cutoffs"`
	Result  bool     `json:"result"`
}

type Cutoff struct {
	Time int64 `json:"time"`
	Ep   int   `json:"ep"`
}

type Rates struct {
	api    *Api
	Type   string  `json:"type"`
	Server int     `json:"server"`
	Tier   int     `json:"tier"`
	Rate   float64 `json:"rate"`
}

type RatesMap []Rates

// InitEvents initializes the EventsMap by fetching the events data from the API.
//
// It returns a pointer to the EventsMap and an error if any occurs during the API request or JSON unmarshalling.
func (a *Api) InitEvents(ctx context.Context) (*EventsMap, error) {
	resp, err := a.do(ctx, http.MethodGet, eventsURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var eventsMap EventsMap
	if err = json.UnmarshalRead(resp.Body, &eventsMap); err != nil {
		return nil, err
	}
	for eventID, events := range eventsMap {
		events.api = a
		events.id = eventID
		events.etag = resp.Header.Get("ETag")
		eventsMap[eventID] = events
	}
	return &eventsMap, nil
}

// UpdateEvents updates the EventsMap with the latest data from the API.
//
// It uses the ETag header to check if the data has changed since the last update. If the data has not changed, it does nothing. If the data has changed, it updates the EventsMap with the new data and initializes the api field for each event in the map.
func (e *EventsMap) UpdateEvents(ctx context.Context) error {
	var etag string
	var api *Api
	for _, events := range *e {
		etag = events.etag
		if etag == "" {
			continue
		}
		api = events.api
		if api == nil {
			continue
		}
		break
	}
	resp, err := api.update(ctx, http.MethodGet, eventsURL, etag)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotModified {
		return nil
	}
	var eventsMap EventsMap
	if err = json.UnmarshalRead(resp.Body, &eventsMap); err != nil {
		return err
	}
	for eventID, events := range eventsMap {
		events.api = api
		events.id = eventID
		events.etag = resp.Header.Get("ETag")
		eventsMap[eventID] = events
	}
	maps.Copy(*e, eventsMap)
	return nil
}

// GetEvent retrieves an Events struct from the EventsMap based on the provided event ID.
//
// Be aware that the Events struct that is returned is a summary of the event data, and does not include all the details that can be found in the Event struct.
//
// Use api.GetEvent(ctx, eventID) to get the full Event struct.
func (e *EventsMap) GetEvent(eventID string) (*Events, error) {
	if _, err := strconv.Atoi(eventID); err != nil {
		return nil, ErrNotValidEventId
	}
	if event, ok := (*e)[eventID]; ok {
		return &event, nil
	}
	return nil, ErrEventNotFound
}

type EventTitleMap struct {
	Id     string
	Titles []string
}

// GetEventTitleMap returns a slice of EventTitleMap, which contains the event ID and the event title for each event in the EventsMap.
func (e *EventsMap) GetEventTitleMap() []EventTitleMap {
	titles := make([]EventTitleMap, 0, len(*e))
	for id, event := range *e {
		titles = append(titles, EventTitleMap{
			Id:     id,
			Titles: event.EventName,
		})
	}
	return titles
}

// GetEvent retrieves an Event struct from the API based on the provided event ID.
func (a *Api) GetEvent(ctx context.Context, eventID string) (*Event, error) {
	if _, err := strconv.Atoi(eventID); err != nil {
		return nil, ErrNotValidEventId
	}
	return getEvent(ctx, a, eventID)
}

// GetEvent retrieves an Event struct from the API based on the RecentEvents struct.
func (r *RecentEvents) GetEvent(ctx context.Context) (*Event, error) {
	return getEvent(ctx, r.api, r.id)
}

func getEvent(ctx context.Context, api *Api, eventID string) (*Event, error) {
	url := "https://bestdori.com/api/events/" + eventID + ".json"
	resp, err := api.do(ctx, http.MethodGet, url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var event Event
	if err = json.UnmarshalRead(resp.Body, &event); err != nil {
		return nil, err
	}
	event.api = api
	event.id = eventID
	return &event, nil
}

// IsAvailable checks if the event is available in the given region.
//
// It is recommended to use this method before calling other methods that depend on the event's availability, such as GetEventName.
func (e *Events) IsAvailable(region Region) bool {
	idx := indexHelper(region.String())
	return !e.StartAt[idx].IsZero() && !e.EndAt[idx].IsZero()
}

// GetTitle returns the title of the event for the given region.
func (e *Events) GetTitle(region Region) string {
	return e.EventName[indexHelper(region.String())]
}

// GetStartTime returns the start time of the event for the given region.
func (e *Events) GetStartTime(region Region) time.Time {
	return e.StartAt[indexHelper(region.String())].Time
}

// GetEndTime returns the end time of the event for the given region.
func (e *Events) GetEndTime(region Region) time.Time {
	return e.EndAt[indexHelper(region.String())].Time
}

// GetBannerURL returns the URL of the event banner for the given region.
func (e *Events) GetBannerURL(region Region) string {
	return "https://bestdori.com/assets/" + region.String() + "/event/" + e.AssetBundleName + "/images_rip/banner.png"
}

// GetBannerFile retrieves the event banner image file for the given region.
//
// The format of the event banner image is PNG.
func (e *Events) GetBannerFile(ctx context.Context, region Region) ([]byte, error) {
	return getBannerFile(ctx, e.api, e.GetBannerURL(region))
}

// GetType returns the type of the event.
func (e *Events) GetType(region Region) string {
	return e.EventType
}

// GetAttribute returns the attribute and its bonus percentage for the event.
func (e *Events) GetAttribute() (string, int) {
	if len(e.Attributes) == 0 {
		return "", 0
	}
	return strings.ToUpper(e.Attributes[0].Attribute), e.Attributes[0].Percent
}

// GetEventCharacters returns the list of event characters for the given region.
func (e *Events) GetEventCharacters(region Region) []EventsCharacters {
	characters := make([]EventsCharacters, len(e.Characters))
	for i, character := range e.Characters {
		character.api = e.api
		characters[i] = character
	}
	return characters
}

// GetEventCards returns the list of event cards for the given region.
func (e *Events) GetEventCards(region Region) []EventMembers {
	members := make([]EventMembers, len(e.Members))
	for i, member := range e.Members {
		member.api = e.api
		members[i] = member
	}
	return members
}

// GetRewardCards retrieves the reward cards for the event.
func (e *Events) GetRewardCards(ctx context.Context) ([]*Card, error) {
	var rewardCards []*Card
	for _, cardID := range e.RewardCards {
		card, err := getCard(ctx, e.api, strconv.Itoa(cardID))
		if err == nil {
			rewardCards = append(rewardCards, card)
		} else {
			return nil, err
		}
	}
	return rewardCards, nil
}

// IsAvailable checks if the event is available in the given region.
//
// It is recommended to use this method before calling other methods that depend on the event's availability, such as GetEventName.
func (e *Event) IsAvailable(region Region) bool {
	idx := indexHelper(region.String())
	return !e.StartAt[idx].IsZero() && !e.EndAt[idx].IsZero()
}

// GetTitle returns the title of the event for the given region.
func (e *Event) GetTitle(region Region) string {
	return e.EventName[indexHelper(region.String())]
}

// GetStartTime returns the start time of the event for the given region.
func (e *Event) GetStartTime(region Region) time.Time {
	return e.StartAt[indexHelper(region.String())].Time
}

// GetEndTime returns the end time of the event for the given region.
func (e *Event) GetEndTime(region Region) time.Time {
	return e.EndAt[indexHelper(region.String())].Time
}

// GetPublicStartTime returns the public start time of the event for the given region.
func (e *Event) GetPublicStartTime(region Region) time.Time {
	return e.PublicStartAt[indexHelper(region.String())].Time
}

// GetPublicEndTime returns the public end time of the event for the given region.
func (e *Event) GetPublicEndTime(region Region) time.Time {
	return e.PublicEndAt[indexHelper(region.String())].Time
}

// GetBannerURL returns the URL of the event banner for the given region.
func (e *Event) GetBannerURL(region Region) string {
	return "https://bestdori.com/assets/" + region.String() + "/event/" + e.AssetBundleName + "/images_rip/banner.png"
}

// GetBannerFile retrieves the event banner image file for the given region.
//
// The format of the event banner image is PNG.
func (e *Event) GetBannerFile(ctx context.Context, region Region) ([]byte, error) {
	return getBannerFile(ctx, e.api, e.GetBannerURL(region))
}

// GetType returns the type of the event.
func (e *Event) GetType(region Region) string {
	return e.EventType
}

// GetAttribute returns the attribute and its bonus percentage for the event.
func (e *Event) GetAttribute() (string, int) {
	if len(e.Attributes) == 0 {
		return "", 0
	}
	return strings.ToUpper(e.Attributes[0].Attribute), e.Attributes[0].Percent
}

// GetEventCharacters returns the list of event characters.
func (e *Event) GetEventCharacters() []EventCharacters {
	characters := make([]EventCharacters, len(e.Characters))
	for i, character := range e.Characters {
		character.api = e.api
		characters[i] = character
	}
	return characters
}

// GetEventCards returns the list of event cards.
func (e *Event) GetEventCards() []EventMembers {
	members := make([]EventMembers, len(e.Members))
	for i, member := range e.Members {
		member.api = e.api
		members[i] = member
	}
	return members
}

// GetRewardCards retrieves the reward cards for the event.
func (e *Event) GetRewardCards(ctx context.Context) ([]*Card, error) {
	var rewardCards []*Card
	for _, cardID := range e.RewardCards {
		card, err := getCard(ctx, e.api, strconv.Itoa(cardID))
		if err == nil {
			rewardCards = append(rewardCards, card)
		} else {
			return nil, err
		}
	}
	return rewardCards, nil
}

// GetEventStories returns the list of event stories.
func (e *Event) GetEventStories() []EventStory {
	stories := make([]EventStory, len(e.Stories))
	for i, story := range e.Stories {
		story.api = e.api
		stories[i] = story
	}
	return stories
}

// GetTitle returns the title of the event story for the given region.
func (s *EventStory) GetTitle(region Region) string {
	return s.Title[indexHelper(region.String())]
}

// GetSynopsis returns the synopsis of the event story for the given region.
func (s *EventStory) GetSynopsis(region Region) string {
	return s.Synopsis[indexHelper(region.String())]
}

// GetCaption returns the caption of the event story for the given region.
func (s *EventStory) GetCaption(region Region) string {
	return s.Caption[indexHelper(region.String())]
}

// GetReleaseConditions returns the release conditions string of the event story for the given region.
func (s *EventStory) GetReleaseConditions(region Region) string {
	return s.ReleaseConditions[indexHelper(region.String())]
}

// GetAttributesAndCharatersBonus returns the attribute and character bonus percentage for the event.
func (e *Event) GetAttributesAndCharatersBonus() int {
	return e.EventAttributeAndCharacterBonus.PointPercent
}

func getBannerFile(ctx context.Context, api *Api, url string) ([]byte, error) {
	resp, err := api.do(ctx, http.MethodGet, url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get event banner: %s", resp.Status)
	}
	if contentType := resp.Header.Get("Content-Type"); contentType != "image/png" {
		return nil, fmt.Errorf("unexpected content type: %s", contentType)
	}
	return io.ReadAll(resp.Body)
}

func getRates(ctx context.Context, api *Api) (RatesMap, error) {
	resp, err := api.do(ctx, http.MethodGet, ratesURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var rates RatesMap
	if err = json.UnmarshalRead(resp.Body, &rates); err != nil {
		return nil, err
	}
	return rates, nil
}

// GetRates retrieves the rates for the event based on the provided region and optional tier.
//
// Be aware that the tier parameter is specific to the region, and the function will return an error if the tier is not valid for the given region.
//
// See api.*TierList for the valid tiers for each region(* stands for region, "cn", "tw" etc). You can also use the api.Region.Tiers() method to get the valid tiers for a specific region.
func (e *Event) GetRates(ctx context.Context, region Region, tier int) (Rates, error) {
	if !slices.Contains(region.Tiers(), tier) {
		return Rates{}, fmt.Errorf("tier %d is not valid for region %s", tier, region.String())
	}
	rates, err := getRates(ctx, e.api)
	if err != nil {
		return Rates{}, err
	}

	for _, rate := range rates {
		if rate.Type == e.EventType &&
			rate.Server == region.Idx() &&
			rate.Tier == tier {
			return rate, nil
		}
	}

	return Rates{}, fmt.Errorf("no rates found for event type: %s", e.EventType)
}

// GetRates retrieves the rates for the event based on the provided region and optional tier.
//
// Be aware that the tier parameter is specific to the region, and the function will return an error if the tier is not valid for the given region.
//
// See api.*TierList for the valid tiers for each region(* stands for region, "cn", "tw" etc). You can also use the api.Region.Tiers() method to get the valid tiers for a specific region.
func (e *Events) GetRates(ctx context.Context, region Region, tier int) (Rates, error) {
	if !slices.Contains(region.Tiers(), tier) {
		return Rates{}, fmt.Errorf("tier %d is not valid for region %s", tier, region.String())
	}
	rates, err := getRates(ctx, e.api)
	if err != nil {
		return Rates{}, err
	}

	for _, rate := range rates {
		if rate.Type == e.EventType &&
			rate.Server == region.Idx() && rate.Tier == tier {
			rate.api = e.api
			return rate, nil
		}
	}

	return Rates{}, fmt.Errorf("no rates found for event type: %s", e.EventType)
}

type EventPredictedTracker struct {
	PredictCutoff []PredictCutoff
}

type PredictCutoff struct {
	Time    int64
	Ep      int
	Percent float64
}

// GetPredictCutoff fetches the event tracker data and predicts the cutoff points based on the provided rates.
//
// The function returns the actual event tracker data, the predicted tracker data, and an error if any occurred during the process.
//
// predicted tracker data will NOT contain the data points that are within 1 day of the event start time, as they are considered too close to the start time to make accurate predictions.
//
// predicted tracker data will contain a data point at when the event ends, with the latest predicted ep value at that time.
//
// The prediction algorithm is a port of the one used in Tsugu, another BanG Dream! Bot.
//
// Be aware that this function refuse to predict if the latest data point is less than 1 day from the event start time, as it may lead to inaccurate predictions. If this happens, the function will return an ErrTooLessDataPoints error.
func (r *Rates) GetPredictCutoff(ctx context.Context, event Event) (EventTracker, EventPredictedTracker, error) {
	region := regionHelper(r.Server)
	url := fmt.Sprintf("https://bestdori.com/api/tracker/data?server=%d&event=%s&tier=%d", r.Server, event.id, r.Tier)
	return getPredictCutoff(ctx, r.api, r.Rate, url, event.GetStartTime(region).UnixMilli(), event.GetEndTime(region).UnixMilli())
}

// GetPredictCutoffWithSummary is the same as GetPredictCutoff, but it takes an Events struct instead of an Event struct.
func (r *Rates) GetPredictCutoffWithSummary(ctx context.Context, events Events) (EventTracker, EventPredictedTracker, error) {
	region := regionHelper(r.Server)
	url := fmt.Sprintf("https://bestdori.com/api/tracker/data?server=%d&event=%s&tier=%d", r.Server, events.id, r.Tier)
	return getPredictCutoff(ctx, r.api, r.Rate, url, events.GetStartTime(region).UnixMilli(), events.GetEndTime(region).UnixMilli())
}

func getPredictCutoff(ctx context.Context, api *Api, rate float64, url string, timeStart, timeEnd int64) (EventTracker, EventPredictedTracker, error) {
	resp, err := api.do(ctx, http.MethodGet, url)
	if err != nil {
		return EventTracker{}, EventPredictedTracker{}, err
	}
	defer resp.Body.Close()
	var tracker EventTracker
	if err = json.UnmarshalRead(resp.Body, &tracker); err != nil {
		return EventTracker{}, EventPredictedTracker{}, err
	}
	predicted := make([]PredictCutoff, len(tracker.Cutoffs))
	var avgPercent, sumPercent float64
	var avgEp, sumEp int
	if tracker.Cutoffs[len(tracker.Cutoffs)-1].Time < timeStart+86400000 {
		return EventTracker{}, EventPredictedTracker{}, ErrTooLessDataPoints
	}
	for i, cutoff := range tracker.Cutoffs {
		percent := float64(cutoff.Time-timeStart) / float64(timeEnd-timeStart)
		predicted[i].Percent = percent
		sumPercent += predicted[i].Percent
		sumEp += tracker.Cutoffs[i].Ep
	}
	if len(tracker.Cutoffs) > 0 {
		avgPercent = sumPercent / float64(len(tracker.Cutoffs))
		avgEp = sumEp / len(tracker.Cutoffs)
	}
	for i, cutoff := range tracker.Cutoffs {
		predicted[i].Time = cutoff.Time
		if int64(cutoff.Time)-timeStart < 86400000 { // 如果最新时间距离开始时间不足一天，则不进行预测
			predicted[i].Ep = cutoff.Ep
			continue
		}
		var z, w float64
		cutoffAltered := tracker.Cutoffs[:i+1]
		for k, cutoff := range cutoffAltered {
			z += (predicted[k].Percent - avgPercent) * (float64(cutoff.Ep) - float64(avgEp))
			w += (predicted[k].Percent - avgPercent) * (predicted[k].Percent - avgPercent)
		}
		b := z / w
		a := float64(avgEp) - b*avgPercent
		predictedEp := a + b*(float64(1)+rate) // 预测的最终ep值
		predicted[i].Ep = int(predictedEp)
	}
	var lastEqualIndex int
	for i, cutoff := range tracker.Cutoffs {
		if cutoff.Ep != predicted[i].Ep {
			lastEqualIndex = i
			break // 如果预测值与实际值不同，则停止截断
		}
	}
	predicted = predicted[lastEqualIndex:] // 截断预测值，去掉与实际值相同的部分
	if predicted[len(predicted)-1].Time < timeEnd {
		predicted = append(
			predicted,
			PredictCutoff{
				Time:    timeEnd,
				Ep:      predicted[len(predicted)-1].Ep,
				Percent: 1.0,
			},
		)
	}
	return tracker, EventPredictedTracker{PredictCutoff: predicted}, nil
}
