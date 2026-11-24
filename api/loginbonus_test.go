package api

import (
	"context"
	"slices"
	"testing"
	"time"
)

func TestGetLoginBonusById(t *testing.T) {
	client := NewApi("", 10*time.Second)
	// Initialize login bonuses first
	loginBonuses, err := client.InitLoginBonuses(context.Background())
	if err != nil {
		t.Fatalf("Expected no error when initializing login bonuses, got error: %v", err)
	}
	if loginBonuses == nil {
		t.Fatal("Expected login bonuses to be non-empty")
	}
	if err = loginBonuses.UpdateLoginBonuses(context.Background()); err != nil {
		t.Fatalf("Expected no error when updating login bonuses, got error: %v", err)
	}
	lb, err := loginBonuses.GetLoginBonus("1590")
	if err != nil {
		t.Fatalf("Expected no error when getting login bonus by ID, got error: %v", err)
	}
	if lb == nil {
		t.Fatal("Expected login bonus to be non-empty")
	}
	if !lb.IsAvailable(RegionJP) {
		t.Fatal("Expected login bonus to be available in JP region")
	}
	if lb.GetCaption(RegionJP) != "ガルパ×魔法少女ノ魔女裁判コラボ開催記念プレゼントPart1" {
		t.Fatalf("Expected login bonus caption to be 'ガルパ×魔法少女ノ魔女裁判コラボ開催記念プレゼントPart1', got '%s'", lb.GetCaption(RegionJP))
	}
	if lb.GetType() != "event" {
		t.Fatalf("Expected login bonus type to be 'event', got '%s'", lb.GetType())
	}
	if lb.GetStartTime(RegionJP).Format("2006-01-02 15:04:05") != "2026-11-06 06:00:00" {
		t.Fatalf("Expected login bonus start time to be '2026-11-06 06:00:00', got '%s'", lb.GetStartTime(RegionJP).Format("2006-01-02 15:04:05"))
	}
	if lb.GetEndTime(RegionJP).Format("2006-01-02 15:04:05") != "2026-11-18 05:59:59" {
		t.Fatalf("Expected login bonus end time to be '2026-11-18 05:59:59', got '%s'", lb.GetEndTime(RegionJP).Format("2006-01-02 15:04:05"))
	}
	if lb.GetBannerURL(context.Background(), RegionJP) != "https://bestdori.com/assets/jp/event/loginbonus/collabo_manosaba_kinen_present_1_rip/background.png" {
		t.Fatalf("Expected login bonus banner URL to be 'https://bestdori.com/assets/jp/event/loginbonus/collabo_manosaba_kinen_present_1_rip/background.png', got '%s'", lb.GetBannerURL(context.Background(), RegionJP))
	}
	if bannerFile, err := lb.GetBannerFile(context.Background(), RegionJP); bannerFile == nil || err != nil {
		t.Fatalf("Expected login bonus banner file to be non-empty and no error, got bannerFile: %v, error: %v", bannerFile, err)
	}
	if lb.GetAssetBundleNames(RegionJP) != "collabo_manosaba_kinen_present_1" {
		t.Fatalf("Expected login bonus asset bundle names to be 'collabo_manosaba_kinen_present_1', got '%s'", lb.GetAssetBundleNames(RegionJP))
	}
	details, err := lb.GetLoginBonus(context.Background())
	if err != nil {
		t.Fatalf("Expected no error when getting login bonus details, got error: %v", err)
	}
	if details == nil {
		t.Fatal("Expected login bonus details to be non-empty")
	}
	if details.GetType() != "event" {
		t.Fatalf("Expected login bonus details type to be 'event', got '%s'", details.GetType())
	}
	if details.GetCaption(RegionJP) != "ガルパ×魔法少女ノ魔女裁判コラボ開催記念プレゼントPart1" {
		t.Fatalf("Expected login bonus details caption to be 'ガルパ×魔法少女ノ魔女裁判コラボ開催記念プレゼントPart1', got '%s'", details.GetCaption(RegionJP))
	}
	if details.GetStartTime(RegionJP).Format("2006-01-02 15:04:05") != "2026-11-06 06:00:00" {
		t.Fatalf("Expected login bonus details start time to be '2026-11-06 06:00:00', got '%s'", details.GetStartTime(RegionJP).Format("2006-01-02 15:04:05"))
	}
	if details.GetEndTime(RegionJP).Format("2006-01-02 15:04:05") != "2026-11-18 05:59:59" {
		t.Fatalf("Expected login bonus details end time to be '2026-11-18 05:59:59', got '%s'", details.GetEndTime(RegionJP).Format("2006-01-02 15:04:05"))
	}
	if details.GetBannerURL(context.Background(), RegionJP) != "https://bestdori.com/assets/jp/event/loginbonus/collabo_manosaba_kinen_present_1_rip/background.png" {
		t.Fatalf("Expected login bonus details banner URL to be 'https://bestdori.com/assets/jp/event/loginbonus/collabo_manosaba_kinen_present_1_rip/background.png', got '%s'", details.GetBannerURL(context.Background(), RegionJP))
	}
	if bannerFile, err := details.GetBannerFile(context.Background(), RegionJP); bannerFile == nil || err != nil {
		t.Fatalf("Expected login bonus details banner file to be non-empty and no error, got bannerFile: %v, error: %v", bannerFile, err)
	}
	if details.GetAssetBundleNames(RegionJP) != "collabo_manosaba_kinen_present_1" {
		t.Fatalf("Expected login bonus details asset bundle names to be 'collabo_manosaba_kinen_present_1', got '%s'", details.GetAssetBundleNames(RegionJP))
	}
	if len(details.GetDetails(RegionJP)) == 0 {
		t.Fatal("Expected login bonus details to have at least one detail")
	}
	if details.GetDetails(RegionJP)[0].GetGrantType() != "direct" {
		t.Fatalf("Expected login bonus detail grant type to be 'direct', got '%s'", details.GetDetails(RegionJP)[0].GetGrantType())
	}
	if details.GetDetails(RegionJP)[0].GetQuantity() != 1 {
		t.Fatalf("Expected login bonus detail quantity to be 1, got %d", details.GetDetails(RegionJP)[0].GetQuantity())
	}
	if *details.GetDetails(RegionJP)[0].GetResourceId() != 1022 {
		t.Fatalf("Expected login bonus detail resource ID to be 1022, got %d", *details.GetDetails(RegionJP)[0].GetResourceId())
	}
	if details.GetDetails(RegionJP)[0].GetResourceType() != "stamp" {
		t.Fatalf("Expected login bonus detail resource type to be 'stamp', got '%s'", details.GetDetails(RegionJP)[0].GetResourceType())
	}
	if details.GetDetails(RegionJP)[0].GetSeq() != 1 {
		t.Fatalf("Expected login bonus detail sequence to be 1, got %d", details.GetDetails(RegionJP)[0].GetSeq())
	}
}

func TestGetLoginBonusTitles(t *testing.T) {
	client := NewApi("", 10*time.Second)
	loginBonuses, err := client.InitLoginBonuses(context.Background())
	if err != nil {
		t.Fatalf("Expected no error when initializing login bonuses, got error: %v", err)
	}
	titles := loginBonuses.GetLoginBonusesTitleMap()
	if len(titles) == 0 {
		t.Fatal("Expected login bonuses title map to be non-empty")
	}
	for i := range titles {
		if titles[i].Id == "" {
			t.Fatal("Expected login bonus ID to be non-empty")
		}
	}
	for i := range titles {
		if titles[i].Id == "1590" {
			if slices.Contains(titles[i].Titles, "ガルパ×魔法少女ノ魔女裁判コラボ開催記念プレゼントPart1") {
				return
			}
			t.Fatal("Expected to find login bonus title 'ガルパ×魔法少女ノ魔女裁判コラボ開催記念プレゼントPart1' for ID '1590', but did not")
		}
	}
}
