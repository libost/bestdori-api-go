package api

import (
	"context"
	"encoding/json/v2"
	"testing"
	"time"
)

func TestGetRecentLoginBonusForRegionParsesLoginBonuses(t *testing.T) {
	var recent Recent
	err := json.Unmarshal([]byte(`{
		"loginBonuses": {
			"123": {
				"caption": ["JP", "EN"],
				"publishedAt": ["1704067200000", "1704067200000", null, null, null],
				"closedAt": ["4102444800000", "4102444800000", null, null, null]
			},
			"124": {
				"caption": ["JP future", "EN future"],
				"publishedAt": ["4102444800000", "4102444800000", null, null, null],
				"closedAt": ["4102444800000", null, null, null, null]
			}
		}
	}`), &recent)
	if err != nil {
		t.Fatalf("failed to unmarshal recent login bonuses: %v", err)
	}

	results := recent.GetRecentLoginBonusForRegion(RegionEN)
	if len(results) != 1 {
		t.Fatalf("expected 1 active EN login bonus, got %d", len(results))
	}
	if results[0].id != "123" {
		t.Fatalf("expected login bonus ID 123, got %q", results[0].id)
	}
}

func TestGetRecentByRegion(t *testing.T) {
	client := NewApi("", 10*time.Second)
	// Get recent data first
	recent, err := client.GetRecent(context.Background())
	if err != nil {
		t.Fatalf("Expected no error when getting recent data, got error: %v", err)
	}
	if recent == nil {
		t.Fatal("Expected recent data to be non-empty")
	}
	if err = recent.UpdateRecent(context.Background()); err != nil {
		t.Fatalf("Expected no error when updating recent data, got error: %v", err)
	}
	if len(recent.GetRecentEventsForRegion(RegionCN)) == 0 {
		t.Fatal("Expected recent events for region CN to be non-empty")
	}
	if len(recent.GetRecentEventsForRegion(RegionKR)) != 0 {
		t.Fatal("Expected recent events for region KR to be empty")
	}
	if event, err := recent.GetRecentEventsForRegion(RegionCN)[0].GetEvent(context.Background()); err != nil || event == nil {
		t.Fatalf("Expected to get event for recent event in region CN, got event: %v, error: %v", event, err)
	}
	if len(recent.GetRecentGachasForRegion(RegionCN)) == 0 {
		t.Fatal("Expected recent gachas for region CN to be non-empty")
	}
	if len(recent.GetRecentGachasForRegion(RegionKR)) != 0 {
		t.Fatal("Expected recent gachas for region KR to be empty")
	}
	if gacha, err := recent.GetRecentGachasForRegion(RegionCN)[0].GetGacha(context.Background()); err != nil || gacha == nil {
		t.Fatalf("Expected to get gacha for recent gacha in region CN, got gacha: %v, error: %v", gacha, err)
	}
	if len(recent.GetRecentLoginBonusForRegion(RegionKR)) != 0 {
		t.Fatal("Expected recent login bonuses for region KR to be empty")
	}
	if len(recent.GetRecentSongsForRegion(RegionCN)) == 0 {
		t.Fatal("Expected recent songs for region CN to be non-empty")
	}
	if len(recent.GetRecentSongsForRegion(RegionKR)) != 0 {
		t.Fatal("Expected recent songs for region KR to be empty")
	}
	if song, err := recent.GetRecentSongsForRegion(RegionCN)[0].GetSong(context.Background()); err != nil || song == nil {
		t.Fatalf("Expected to get song for recent song in region CN, got song: %v, error: %v", song, err)
	}
	if songd, err := recent.GetRecentSongsForRegion(RegionCN)[0].GetSong(context.Background()); err != nil || songd == nil {
		t.Fatalf("Expected to get song details for recent song in region CN, got song details: %v, error: %v", songd, err)
	}
	lb, err := recent.GetRecentLoginBonusForRegion(RegionEN)[0].GetLoginBonus(context.Background())
	if err != nil || lb == nil {
		t.Fatalf("Expected to get login bonus for recent login bonus in region EN, got login bonus: %v, error: %v", lb, err)
	}
	if recent.GetRecentLoginBonusForRegion(RegionEN)[0].GetCaption(RegionEN) == "" {
		t.Fatalf("Expected recent login bonus caption to be non-empty, got '%s'", recent.GetRecentLoginBonusForRegion(RegionEN)[0].GetCaption(RegionEN))
	}
	if recent.GetRecentLoginBonusForRegion(RegionEN)[0].GetEndTime(RegionEN).IsZero() {
		t.Fatalf("Expected recent login bonus end time to be non-zero, got '%s'", recent.GetRecentLoginBonusForRegion(RegionEN)[0].GetEndTime(RegionEN))
	}
	if recent.GetRecentLoginBonusForRegion(RegionEN)[0].GetStartTime(RegionEN).IsZero() {
		t.Fatalf("Expected recent login bonus start time to be non-zero, got '%s'", recent.GetRecentLoginBonusForRegion(RegionEN)[0].GetStartTime(RegionEN))
	}
	if lb, err := recent.GetRecentLoginBonusForRegion(RegionEN)[0].GetLoginBonus(context.Background()); err != nil || lb == nil {
		t.Fatalf("Expected to get login bonus type for recent login bonus in region EN, got login bonus type: %v, error: %v", lb, err)
	}
}
