package api

import (
	"context"
	"slices"
	"testing"
	"time"
)

func TestGetGachaById(t *testing.T) {
	client := NewApi("", 10*time.Second)
	// Initialize gachas first
	gachas, err := client.InitGachas(context.Background())
	if err != nil {
		t.Fatalf("Expected no error when initializing gachas, got error: %v", err)
	}
	if gachas == nil {
		t.Fatal("Expected gachas to be non-empty")
	}
	if err = gachas.UpdateGachas(context.Background()); err != nil {
		t.Fatalf("Failed to update gachas: %v", err)
	}
	gacha, err := gachas.GetGacha("1166")
	if err != nil {
		t.Fatalf("Expected no error when getting gacha by ID, got error: %v", err)
	}
	if gacha == nil {
		t.Fatal("Expected gacha to be non-empty")
	}
	if !gacha.IsAvailable(RegionJP) {
		t.Fatal("Expected gacha to be available in JP region")
	}
	if bannerFile, err := gacha.GetBannerFile(context.Background(), RegionJP); bannerFile == nil || err != nil {
		t.Fatalf("Expected gacha banner file to be non-empty and no error, got bannerFile: %v, error: %v", bannerFile, err)
	}
	if gacha.GetBannerURL(context.Background(), RegionJP) != "https://bestdori.com/assets/jp/homebanner_rip/banner_gacha1166.png" {
		t.Fatalf("Expected gacha banner URL to be 'https://bestdori.com/assets/jp/homebanner_rip/banner_gacha1166.png', got '%s'", gacha.GetBannerURL(context.Background(), RegionJP))
	}
	if gacha.GetTitle(RegionJP) != "落涙のPouring Rainガチャ" {
		t.Fatalf("Expected gacha title to be '落涙のPouring Rainガチャ', got '%s'", gacha.GetTitle(RegionJP))
	}
	if gacha.GetCloseTime(RegionJP).Format("2006-01-02 15:04:05") != "2023-07-29 05:59:59" {
		t.Fatalf("Expected gacha close time to be '2023-07-29 05:59:59', got '%s'", gacha.GetCloseTime(RegionJP).Format("2006-01-02 15:04:05"))
	}
	if len(gacha.GetNewCards()) == 0 {
		t.Fatal("Expected gacha new cards to be non-empty")
	}
	if card, err := gacha.GetNewCards()[0].GetCard(context.Background()); card == nil || err != nil {
		t.Fatalf("Expected gacha new card to be non-empty and no error, got card: %v, error: %v", card, err)
	}
	if gacha.GetPublishTime(RegionJP).Format("2006-01-02 15:04:05") != "2023-07-20 06:00:00" {
		t.Fatalf("Expected gacha publish time to be '2023-07-20 06:00:00', got '%s'", gacha.GetPublishTime(RegionJP).Format("2006-01-02 15:04:05"))
	}
	if gacha.GetType() != "permanent" {
		t.Fatalf("Expected gacha type to be 'permanent', got '%s'", gacha.GetType())
	}
	gachaDetail, err := gacha.GetGacha(context.Background())
	if err != nil {
		t.Fatalf("Expected no error when getting gacha detail, got error: %v", err)
	}
	if gachaDetail == nil {
		t.Fatal("Expected gacha detail to be non-empty")
	}
	if !gachaDetail.IsAvailable(RegionJP) {
		t.Fatal("Expected gacha detail to be available in JP region")
	}
	if gachaDetail.GetTitle(RegionJP) != "落涙のPouring Rainガチャ" {
		t.Fatalf("Expected gacha detail title to be '落涙のPouring Rainガチャ', got '%s'", gachaDetail.GetTitle(RegionJP))
	}
	if gachaDetail.GetCloseTime(RegionJP).Format("2006-01-02 15:04:05") != "2023-07-29 05:59:59" {
		t.Fatalf("Expected gacha detail close time to be '2023-07-29 05:59:59', got '%s'", gachaDetail.GetCloseTime(RegionJP).Format("2006-01-02 15:04:05"))
	}
	if gachaDetail.GetPublishTime(RegionJP).Format("2006-01-02 15:04:05") != "2023-07-20 06:00:00" {
		t.Fatalf("Expected gacha detail publish time to be '2023-07-20 06:00:00', got '%s'", gachaDetail.GetPublishTime(RegionJP).Format("2006-01-02 15:04:05"))
	}
	if gachaDetail.GetType() != "permanent" {
		t.Fatalf("Expected gacha detail type to be 'permanent', got '%s'", gachaDetail.GetType())
	}
	if gachaDetail.GetAnnotation(RegionJP) != "・下記提供割合でメンバーが出現します。\n・出現メンバーはすべて特訓前の状態になります。\n・｢★3以上1人確定｣は｢確定枠提供割合｣をご確認ください。\n・メンバー毎に表示している提供割合は、\n表示上小数点第4位にて四捨五入しているため、\n全てを合算しても100%にならない場合があります。" {
		t.Fatalf("Expected gacha detail annotation to be '・下記提供割合でメンバーが出現します。\n・出現メンバーはすべて特訓前の状態になります。\n・｢★3以上1人確定｣は｢確定枠提供割合｣をご確認ください。\n・メンバー毎に表示している提供割合は、\n表示上小数点第4位にて四捨五入しているため、\n全てを合算しても100%%にならない場合があります。', got '%s'", gachaDetail.GetAnnotation(RegionJP))
	}
	if gachaDetail.GetBannerURL(context.Background(), RegionJP) != "https://bestdori.com/assets/jp/homebanner_rip/banner_gacha1166.png" {
		t.Fatalf("Expected gacha detail banner URL to be 'https://bestdori.com/assets/jp/homebanner_rip/banner_gacha1166.png', got '%s'", gachaDetail.GetBannerURL(context.Background(), RegionJP))
	}
	if bannerFile, err := gachaDetail.GetBannerFile(context.Background(), RegionJP); bannerFile == nil || err != nil {
		t.Fatalf("Expected gacha detail banner file to be non-empty and no error, got bannerFile: %v, error: %v", bannerFile, err)
	}
	if len(gachaDetail.GetNewCards()) == 0 {
		t.Fatal("Expected gacha detail new cards to be non-empty")
	}
	if card, err := gachaDetail.GetNewCards()[0].GetCard(context.Background()); card == nil || err != nil {
		t.Fatalf("Expected gacha detail new card to be non-empty and no error, got card: %v, error: %v", card, err)
	}
	if len(gachaDetail.GetCardSpecifiedWeight()) == 0 {
		t.Fatal("Expected gacha detail card specified weight to be non-empty")
	}
	if len(gachaDetail.GetRaritySpecifiedWeight()) == 0 {
		t.Fatal("Expected gacha detail rarity specified weight to be non-empty")
	}
	if gachaDetail.GetDescription(RegionJP) != "★★以上のメンバーが出現！ピックアップメンバーは出現率ＵＰ中！一部の★★以上のメンバーはライブ衣装付き！" {
		t.Fatalf("Expected gacha detail description to be '★★以上のメンバーが出現！ピックアップメンバーは出現率ＵＰ中！一部の★★以上のメンバーはライブ衣装付き！', got '%s'", gachaDetail.GetDescription(RegionJP))
	}
	if gachaDetail.GetGachaPeriod(RegionJP) != "07/29 14:59 まで" {
		t.Fatalf("Expected gacha detail period to be '07/29 14:59 まで', got '%s'", gachaDetail.GetGachaPeriod(RegionJP))
	}
	if len(gachaDetail.GetPaymentMethods()) == 0 {
		t.Fatal("Expected gacha detail payment methods to be non-empty")
	}
	if gachaDetail.GetPaymentMethods()[0].GetBehavior() != "over_the_3_star_once" {
		t.Fatalf("Expected gacha detail payment method behavior to be 'over_the_3_star_once', got '%s'", gachaDetail.GetPaymentMethods()[0].GetBehavior())
	}
	if gachaDetail.GetPaymentMethods()[0].GetCost() != 250 {
		t.Fatalf("Expected gacha detail payment method cost to be 250, got '%d'", gachaDetail.GetPaymentMethods()[0].GetCost())
	}
	if gachaDetail.GetPaymentMethods()[0].GetDraws() != 10 {
		t.Fatalf("Expected gacha detail payment method draws to be 10, got '%d'", gachaDetail.GetPaymentMethods()[0].GetDraws())
	}
	if gachaDetail.GetPaymentMethods()[0].GetMethod() != "free_star" {
		t.Fatalf("Expected gacha detail payment method method to be 'free_star', got '%s'", gachaDetail.GetPaymentMethods()[0].GetMethod())
	}
	if gachaDetail.GetPaymentMethods()[0].GetTotalCost() != 2500 {
		t.Fatalf("Expected gacha detail payment method total cost to be 2500, got '%d'", gachaDetail.GetPaymentMethods()[0].GetTotalCost())
	}
	if gachaDetail.GetPaymentMethods()[0].IsPickup() == true {
		t.Fatal("Expected gacha detail payment method pickup to be false")
	}
	if result, err := gachaDetail.EmulateGachaDraw(gachaDetail.GetPaymentMethods()[0]); err != nil || len(result) != 10 {
		t.Fatalf("Expected gacha detail emulated gacha draw to be 10, got '%d'", len(result))
	}
}

func TestGetGachaTitles(t *testing.T) {
	client := NewApi("", 10*time.Second)
	gachas, err := client.InitGachas(context.Background())
	if err != nil {
		t.Fatalf("Failed to initialize gachas: %v", err)
	}
	titles := gachas.GetGachaTitleMap()
	if len(titles) == 0 {
		t.Fatal("Expected gacha titles to be non-empty")
	}
	for i := range titles {
		if titles[i].Id == "" {
			t.Fatal("Expected gacha ID to be non-empty")
		}
	}
	for i := range titles {
		if titles[i].Id == "1166" {
			if slices.Contains(titles[i].Titles, "落涙のPouring Rainガチャ") {
				return
			}
			t.Fatal("Expected to find gacha title '落涙のPouring Rainガチャ' for ID '1166', but did not")
		}
	}
}
