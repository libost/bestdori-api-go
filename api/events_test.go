package api

import (
	"context"
	"encoding/json/v2"
	"slices"
	"testing"
	"time"
)

func TestEventsAttributesUnmarshal(t *testing.T) {
	var events Events
	err := json.Unmarshal([]byte(`{"attributes":[{"attribute":"PURE","percent":10}]}`), &events)
	if err != nil {
		t.Fatalf("Failed to unmarshal events attributes: %v", err)
	}

	attribute, percent := events.GetAttribute()
	if attribute != "PURE" || percent != 10 {
		t.Fatalf("Expected event attribute to be 'PURE' and percent to be 10, got '%s' and %d", attribute, percent)
	}
}

func TestGetEventById(t *testing.T) {
	client := NewApi("", 10*time.Second)
	// Initialize events first
	events, err := client.InitEvents(context.Background())
	if err != nil {
		t.Fatalf("Failed to initialize events: %v", err)
	}
	if events == nil {
		t.Fatal("Expected events, got nil")
	}
	if err = events.UpdateEvents(context.Background()); err != nil {
		t.Fatalf("Failed to update events: %v", err)
	}
	event, err := events.GetEvent("311")
	if err != nil {
		t.Fatalf("Failed to get event by ID: %v", err)
	}
	if event == nil {
		t.Fatal("Expected event, got nil")
	}
	if !event.IsAvailable(RegionCN) {
		t.Fatal("Expected event to be available in RegionCN")
	}
	if event.GetTitle(RegionCN) != "交汇处，双子星相视而笑" {
		t.Fatalf("Expected event title to be '交汇处，双子星相视而笑', got '%s'", event.GetTitle(RegionCN))
	}
	if event.GetBannerURL(RegionCN) != "https://bestdori.com/assets/cn/event/sayohina_cross/images_rip/banner.png" {
		t.Fatalf("Expected event banner URL to be 'https://bestdori.com/assets/cn/event/sayohina_cross/images_rip/banner.png', got '%s'", event.GetBannerURL(RegionCN))
	}
	bannerFile, err := event.GetBannerFile(context.Background(), RegionCN)
	if err != nil {
		t.Fatalf("Failed to get event banner file: %v", err)
	}
	if len(bannerFile) == 0 {
		t.Fatal("Expected event banner file to be non-empty")
	}
	attr, percent := event.GetAttribute()
	if attr != "PURE" || percent != 10 {
		t.Fatalf("Expected event attribute to be 'PURE' and percent to be 10, got '%s' and %d", attr, percent)
	}
	if event.GetStartTime(RegionCN).Format("2006-01-02 15:04:05") != "2026-05-26 02:00:00" {
		t.Fatalf("Expected event start time to be '2026-05-26 02:00:00', got '%s'", event.GetStartTime(RegionCN).Format("2006-01-02 15:04:05"))
	}
	if event.GetEndTime(RegionCN).Format("2006-01-02 15:04:05") != "2026-06-03 14:59:59" {
		t.Fatalf("Expected event end time to be '2026-06-03 14:59:59', got '%s'", event.GetEndTime(RegionCN).Format("2006-01-02 15:04:05"))
	}
	if event.GetType(RegionCN) != "challenge" {
		t.Fatalf("Expected event type to be 'challenge', got '%s'", event.GetType(RegionCN))
	}
	if event.GetEventCards(RegionCN) == nil {
		t.Fatal("Expected event cards to be non-nil")
	}
	if card, percent, err := event.GetEventCards(RegionCN)[0].GetCard(context.Background()); card == nil || err != nil || percent != 20 {
		t.Fatalf("Expected event card to be non-nil, percent to be 20, and no error, got card: %v, percent: %d, error: %v", card, percent, err)
	}
	if char, percent, err := event.GetEventCharacters(RegionCN)[0].GetCharacter(context.Background()); char == nil || percent != 20 || err != nil {
		t.Fatalf("Expected event character to be non-nil, percent to be 20, and no error, got character: %v, percent: %d, error: %v", char, percent, err)
	}
	if card, err := event.GetRewardCards(context.Background()); card == nil || err != nil {
		t.Fatalf("Expected reward cards to be non-nil and no error, got card: %v, error: %v", card, err)
	}
	rate, err := event.GetRates(context.Background(), RegionCN, 50)
	if err != nil {
		t.Fatalf("Expected no error when getting event rates, got error: %v", err)
	}
	eventDetail, err := client.GetEvent(context.Background(), "311")
	if err != nil {
		t.Fatalf("Expected no error when getting event detail, got error: %v", err)
	}
	if !eventDetail.IsAvailable(RegionCN) {
		t.Fatalf("Expected event detail to be available in RegionCN")
	}
	if cutoff, predict, err := rate.GetPredictCutoff(context.Background(), *eventDetail); err != nil || len(cutoff.Cutoffs) == 0 || len(predict.PredictCutoff) == 0 {
		t.Fatalf("Expected no error, and non-nil cutoff and predict, got error: %v, cutoff: %v, predict: %v", err, cutoff, predict)
	}
	if attr, percent := eventDetail.GetAttribute(); attr != "PURE" || percent != 10 {
		t.Fatalf("Expected attribute `PURE` and percent 10, got attribute: %s, percent: %d", attr, percent)
	}
	if eventDetail.GetStartTime(RegionCN).Format("2006-01-02 15:04:05") != "2026-05-26 02:00:00" {
		t.Fatalf("Expected event detail start time to be '2026-05-26 02:00:00', got '%s'", eventDetail.GetStartTime(RegionCN).Format("2006-01-02 15:04:05"))
	}
	if eventDetail.GetEndTime(RegionCN).Format("2006-01-02 15:04:05") != "2026-06-03 14:59:59" {
		t.Fatalf("Expected event detail end time to be '2026-06-03 14:59:59', got '%s'", eventDetail.GetEndTime(RegionCN).Format("2006-01-02 15:04:05"))
	}
	if eventDetail.GetType(RegionCN) != "challenge" {
		t.Fatalf("Expected event detail type to be 'challenge', got '%s'", eventDetail.GetType(RegionCN))
	}
	if eventDetail.GetEventCards() == nil {
		t.Fatal("Expected event detail cards to be non-nil")
	}
	if card, percent, err := eventDetail.GetEventCards()[0].GetCard(context.Background()); card == nil || err != nil || percent != 20 {
		t.Fatalf("Expected event detail card to be non-nil, percent to be 20, and no error, got card: %v, percent: %d, error: %v", card, percent, err)
	}
	if char, percent, err := eventDetail.GetEventCharacters()[0].GetCharacter(context.Background()); char == nil || percent != 20 || err != nil {
		t.Fatalf("Expected event detail character to be non-nil, percent to be 20, and no error, got character: %v, percent: %d, error: %v", char, percent, err)
	}
	if card, err := eventDetail.GetRewardCards(context.Background()); card == nil || err != nil {
		t.Fatalf("Expected reward cards to be non-nil and no error, got card: %v, error: %v", card, err)
	}
	rateDetail, err := eventDetail.GetRates(context.Background(), RegionCN, 50)
	if err != nil {
		t.Fatalf("Expected no error when getting event detail rates, got error: %v", err)
	}
	if cutoff, predict, err := rateDetail.GetPredictCutoff(context.Background(), *eventDetail); err != nil || len(cutoff.Cutoffs) == 0 || len(predict.PredictCutoff) == 0 {
		t.Fatalf("Expected no error, and non-nil cutoff and predict for event detail, got error: %v, cutoff: %v, predict: %v", err, cutoff, predict)
	}
	if cutoffS, predictS, err := rateDetail.GetPredictCutoffWithSummary(context.Background(), *event); err != nil || len(cutoffS.Cutoffs) == 0 || len(predictS.PredictCutoff) == 0 {
		t.Fatalf("Expected no error, and non-nil cutoff and predict for event detail with summary, got error: %v, cutoff: %v, predict: %v", err, cutoffS, predictS)
	}
	if eventDetail.GetBannerURL(RegionCN) != "https://bestdori.com/assets/cn/event/sayohina_cross/images_rip/banner.png" {
		t.Fatalf("Expected event detail banner URL to be 'https://bestdori.com/assets/cn/event/sayohina_cross/images_rip/banner.png', got '%s'", eventDetail.GetBannerURL(RegionCN))
	}
	if bannerFile, err := eventDetail.GetBannerFile(context.Background(), RegionCN); err != nil || len(bannerFile) == 0 {
		t.Fatalf("Expected event detail banner file to be non-empty and no error, got length: %d, error: %v", len(bannerFile), err)
	}
	stories := eventDetail.GetEventStories()
	if len(stories) == 0 {
		t.Fatalf("Expected event detail stories to be non-empty")
	}
	if stories[0].GetCaption(RegionCN) != "序章" {
		t.Fatalf("Expected event detail story caption to be '序章', got '%s'", stories[0].GetCaption(RegionCN))
	}
	if stories[0].GetReleaseConditions(RegionCN) != "序章故事" {
		t.Fatalf("Expected event detail story to be '序章故事', got '%s'", stories[0].GetReleaseConditions(RegionCN))
	}
	if stories[0].GetSynopsis(RegionCN) != "PasPale与Roselia收到了\n某个音乐节的演出邀请。而邀约内容是……" {
		t.Fatalf("Expected event detail story synopsis to be 'PasPale与Roselia收到了\n某个音乐节的演出邀请。而邀约内容是……', got '%s'", stories[0].GetSynopsis(RegionCN))
	}
	if stories[0].GetTitle(RegionCN) != "拼图的两块碎片" {
		t.Fatalf("Expected event detail story title to be '拼图的两块碎片', got '%s'", stories[0].GetTitle(RegionCN))
	}
	if eventDetail.GetTitle(RegionCN) != "交汇处，双子星相视而笑" {
		t.Fatalf("Expected event detail title to be '交汇处，双子星相视而笑', got '%s'", eventDetail.GetTitle(RegionCN))
	}
	if eventDetail.GetPublicStartTime(RegionCN).Format("2006-01-02 15:04:05") != "2026-05-26 02:00:00" {
		t.Fatalf("Expected event detail public start time to be '2026-05-26 02:00:00', got '%s'", eventDetail.GetPublicStartTime(RegionCN).Format("2006-01-02 15:04:05"))
	}
	if eventDetail.GetPublicEndTime(RegionCN).Format("2006-01-02 15:04:05") != "2026-06-05 04:59:59" {
		t.Fatalf("Expected event detail public end time to be '2026-06-05 04:59:59', got '%s'", eventDetail.GetPublicEndTime(RegionCN).Format("2006-01-02 15:04:05"))
	}
	if eventDetail.GetAttributesAndCharatersBonus() != 20 {
		t.Fatalf("Expected event detail attributes and characters bonus to be 20, got %d", eventDetail.GetAttributesAndCharatersBonus())
	}
}

func TestGetEventTitles(t *testing.T) {
	client := NewApi("", 10*time.Second)
	events, err := client.InitEvents(context.Background())
	if err != nil {
		t.Fatalf("Failed to initialize events: %v", err)
	}
	titles := events.GetEventTitleMap()
	if len(titles) == 0 {
		t.Fatal("Expected event titles to be non-empty")
	}
	for i := range titles {
		if titles[i].Id == "" {
			t.Fatal("Expected event ID to be non-empty")
		}
	}
	for i := range titles {
		if titles[i].Id == "311" {
			if slices.Contains(titles[i].Titles, "交汇处，双子星相视而笑") {
				return
			}
			t.Fatal("Expected to find event title '交汇处，双子星相视而笑' for ID '311', but did not")
		}
	}
}
