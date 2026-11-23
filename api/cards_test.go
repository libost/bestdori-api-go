package api

import (
	"context"
	"testing"
	"time"
)

func TestGetCardById(t *testing.T) {
	client := NewApi("", 10*time.Second)
	// Initialize cards first
	cardsMap, err := client.InitCards(context.Background())
	if err != nil {
		t.Fatalf("Failed to initialize cards: %v", err)
	}
	if cardsMap == nil {
		t.Fatal("Expected cards map, got nil")
	}
	if err = cardsMap.UpdateCards(context.Background()); err != nil {
		t.Fatalf("Failed to update cards: %v", err)
	}
	hina947, err := cardsMap.GetCard("947")
	if err != nil {
		t.Fatalf("Failed to get card by ID: %v", err)
	}
	if hina947 == nil {
		t.Fatal("Expected card, got nil")
	}
	if hina947.id != "947" {
		t.Fatalf("Expected card ID '947', got '%s'", hina947.id)
	}
	if !hina947.IsAvailable(RegionCN) {
		t.Fatal("Expected card to be available in RegionCN")
	}
	if hina947.GetCardName(RegionCN) != "日菜抽中的大奖" {
		t.Fatalf("Expected card name '日菜抽中的大奖', got '%s'", hina947.GetCardName(RegionCN))
	}
	if hina947.GetRarity() != 4 {
		t.Fatalf("Expected card rarity 4 , got '%d'", hina947.GetRarity())
	}
	if hina947.GetAttribute() != "PURE" {
		t.Fatalf("Expected card attribute 'PURE', got '%s'", hina947.GetAttribute())
	}
	if hina947.GetReleaseTime(RegionCN).Format("2006-01-02 15:04:05") != "2021-08-17 05:00:00" {
		t.Fatalf("Expected card release time '2021-08-17 05:00:00', got '%s'", hina947.GetReleaseTime(RegionCN).Format("2006-01-02 15:04:05"))
	}
	if hina947.GetType() != "permanent" {
		t.Fatalf("Expected card type 'permanent', got '%s'", hina947.GetType())
	}
	if hina947.GetLevelLimit() != 50 {
		t.Fatalf("Expected card level limit 50, got '%d'", hina947.GetLevelLimit())
	}
	normal, err := hina947.GetArtsFile(context.Background(), RegionCN, CardStatusNormal)
	if err != nil {
		t.Fatalf("Failed to get card arts file: %v", err)
	}
	if normal == nil {
		t.Fatal("Expected card arts file, got nil")
	}
	if hina947.GetArtsURL(context.Background(), RegionCN, CardStatusNormal) != "https://bestdori.com/assets/cn/characters/resourceset/res017040_rip/card_normal.png" {
		t.Fatalf("Expected card arts URL 'https://bestdori.com/assets/cn/characters/resourceset/res017040_rip/card_normal.png', got '%s'", hina947.GetArtsURL(context.Background(), RegionCN, CardStatusNormal))
	}
	after, err := hina947.GetArtsFile(context.Background(), RegionCN, CardStatusAfterTraining)
	if err != nil {
		t.Fatalf("Failed to get card arts file: %v", err)
	}
	if after == nil {
		t.Fatal("Expected card arts file, got nil")
	}
	if hina947.GetArtsURL(context.Background(), RegionCN, CardStatusAfterTraining) != "https://bestdori.com/assets/cn/characters/resourceset/res017040_rip/card_after_training.png" {
		t.Fatalf("Expected card arts URL 'https://bestdori.com/assets/cn/characters/resourceset/res017040_rip/card_after_training.png', got '%s'", hina947.GetArtsURL(context.Background(), RegionCN, CardStatusAfterTraining))
	}
	attr, err := hina947.GetAttributeIconFile(context.Background())
	if err != nil {
		t.Fatalf("Failed to get card attribute icon file: %v", err)
	}
	if attr == nil {
		t.Fatal("Expected card attribute icon file, got nil")
	}
	if hina947.GetAttributeIconURL() != "https://bestdori.com/res/icon/pure.svg" {
		t.Fatalf("Expected card attribute icon URL 'https://bestdori.com/res/icon/pure.svg', got '%s'", hina947.GetAttributeIconURL())
	}
	thumbNormal, err := hina947.GetThumbnailFile(context.Background(), RegionCN, CardStatusNormal)
	if err != nil {
		t.Fatalf("Failed to get card thumbnail file: %v", err)
	}
	if thumbNormal == nil {
		t.Fatal("Expected card thumbnail file, got nil")
	}
	if hina947.GetThumbnailURL(context.Background(), RegionCN, CardStatusNormal) != "https://bestdori.com/assets/cn/thumb/chara/card00018_rip/res017040_normal.png" {
		t.Fatalf("Expected card thumbnail URL 'https://bestdori.com/assets/cn/thumb/chara/card00018_rip/res017040_normal.png', got '%s'", hina947.GetThumbnailURL(context.Background(), RegionCN, CardStatusNormal))
	}
	thumbAfter, err := hina947.GetThumbnailFile(context.Background(), RegionCN, CardStatusAfterTraining)
	if err != nil {
		t.Fatalf("Failed to get card thumbnail file: %v", err)
	}
	if thumbAfter == nil {
		t.Fatal("Expected card thumbnail file, got nil")
	}
	if hina947.GetThumbnailURL(context.Background(), RegionCN, CardStatusAfterTraining) != "https://bestdori.com/assets/cn/thumb/chara/card00018_rip/res017040_after_training.png" {
		t.Fatalf("Expected card thumbnail URL 'https://bestdori.com/assets/cn/thumb/chara/card00018_rip/res017040_after_training.png', got '%s'", hina947.GetThumbnailURL(context.Background(), RegionCN, CardStatusAfterTraining))
	}
	frame, err := hina947.GetThumbnailFrameFile(context.Background())
	if err != nil {
		t.Fatalf("Failed to get card thumbnail frame file: %v", err)
	}
	if frame == nil {
		t.Fatal("Expected card thumbnail frame file, got nil")
	}
	if hina947.GetThumbnailFrameURL() != "https://bestdori.com/res/image/card-4.png" {
		t.Fatalf("Expected card thumbnail frame URL 'https://bestdori.com/res/image/card-4.png', got '%s'", hina947.GetThumbnailFrameURL())
	}
	char, err := hina947.GetCharacter(context.Background())
	if err != nil {
		t.Fatalf("Failed to get character for card: %v", err)
	}
	if char == nil {
		t.Fatal("Expected character, got nil")
	}
	month, day := char.GetBirthday()
	if month != 3 || day != 20 {
		t.Fatalf("Expected character birthday March 20, got %d-%d", month, day)
	}
	if char.GetVoice(RegionCN) != "小泽亚李" {
		t.Fatalf("Expected character voice actor '小泽亚李', got '%s'", char.GetVoice(RegionCN))
	}
	if char.GetCharacterType() != "unique" {
		t.Fatalf("Expected character type 'unique', got '%s'", char.GetCharacterType())
	}
	if char.GetCharacterColor() != "#55DDEE" {
		t.Fatalf("Expected character color '#55DDEE', got '%s'", char.GetCharacterColor())
	}
	if char.GetCharacterName(RegionCN) != "冰川 日菜" {
		t.Fatalf("Expected character name '冰川 日菜', got '%s'", char.GetCharacterName(RegionCN))
	}
	if char.GetFavoriteFood(RegionCN) != "泡泡糖、垃圾食品" {
		t.Fatalf("Expected character favorite food '泡泡糖、垃圾食品', got '%s'", char.GetFavoriteFood(RegionCN))
	}
	name, err := char.GetBandName(context.Background(), RegionCN)
	if err != nil {
		t.Fatalf("Failed to get character band name: %v", err)
	}
	if name != "Pastel＊Palettes" {
		t.Fatalf("Expected character band name 'Pastel＊Palettes', got '%s'", name)
	}
	if char.GetConstellation() != "pisces" {
		t.Fatalf("Expected character constellation 'pisces', got '%s'", char.GetConstellation())
	}
	if char.GetHeight() != 156 {
		t.Fatalf("Expected character height 156, got '%d'", char.GetHeight())
	}
	if char.GetPart() != "guitar" {
		t.Fatalf("Expected character part 'guitar', got '%s'", char.GetPart())
	}
	if char.GetHobby(RegionCN) != "配制芳香精油" {
		t.Fatalf("Expected character hobby '配制芳香精油', got '%s'", char.GetHobby(RegionCN))
	}
	if char.GetSelfIntroduction(RegionCN) != "担任Pastel＊Palettes的吉他手，是一名天才少女。非常喜欢很噜♪且有趣的事物，虽然其直截了当的措辞经常险些引发麻烦，但她非常珍视PasPale。一直梦想着和姐姐纱夜站上同一个舞台。" {
		t.Fatalf("Expected character self-introduction '担任Pastel＊Palettes的吉他手，是一名天才少女。非常喜欢很噜♪且有趣的事物，虽然其直截了当的措辞经常险些引发麻烦，但她非常珍视PasPale。一直梦想着和姐姐纱夜站上同一个舞台。', got '%s'", char.GetSelfIntroduction(RegionCN))
	}
	if char.GetCharacterFirstName(RegionCN) != "日菜" {
		t.Fatalf("Expected character first name '日菜', got '%s'", char.GetCharacterFirstName(RegionCN))
	}
	if char.GetCharacterLastName(RegionCN) != "冰川" {
		t.Fatalf("Expected character last name '冰川', got '%s'", char.GetCharacterLastName(RegionCN))
	}
	if char.GetCharacterNickname(RegionCN) != "" {
		t.Fatalf("Expected character nickname '', got '%s'", char.GetCharacterNickname(RegionCN))
	}
	if char.GetHatedFood(RegionCN) != "豆腐、腐竹之类的味道淡的食物" {
		t.Fatalf("Expected character hated food '豆腐、腐竹之类的味道淡的食物', got '%s'", char.GetHatedFood(RegionCN))
	}
	if char.GetSchool(RegionCN) != "庆鹏女子大学" {
		t.Fatalf("Expected character school '庆鹏女子大学', got '%s'", char.GetSchool(RegionCN))
	}
	if char.GetSchoolYear(RegionCN) != "大一" {
		t.Fatalf("Expected character school year '大一', got '%s'", char.GetSchoolYear(RegionCN))
	}
	if char.GetSchoolCls(RegionCN) != "B班" {
		t.Fatalf("Expected character school class 'B班', got '%s'", char.GetSchoolCls(RegionCN))
	}
	if char.GetBandLogoURL() != "https://bestdori.com/res/icon/band_4.svg" {
		t.Fatalf("Expected character band logo URL 'https://bestdori.com/res/icon/band_4.svg', got '%s'", char.GetBandLogoURL())
	}
	bandlogo, err := char.GetBandLogoFile(context.Background())
	if err != nil {
		t.Fatalf("Failed to get character band logo file: %v", err)
	}
	if bandlogo == nil {
		t.Fatal("Expected character band logo file, got nil or empty")
	}
	skill, err := hina947.GetSkill(context.Background())
	if err != nil {
		t.Fatalf("Failed to get skill for card: %v", err)
	}
	if skill == nil {
		t.Fatal("Expected skill, got nil")
	}
	if !skill.IsAvailable(RegionCN) {
		t.Fatal("Expected skill to be available in RegionCN")
	}
	if skill.GetSimpleDescription(RegionCN) != "生命600以上则得分提升80%|生命不满600则生命回复500" {
		t.Fatalf("Expected skill simple description '生命600以上则得分提升80%%|生命不满600则生命回复500', got '%s'", skill.GetSimpleDescription(RegionCN))
	}
	if skill.GetDescription(RegionCN) != "发动时若生命在600以上则{1}秒内 得分提升80%。若生命不满600时 生命回复{0}" {
		t.Fatalf("Expected skill description '发动时若生命在600以上则{1}秒内 得分提升80%%。若生命不满600时 生命回复{0}', got '%s'", skill.GetDescription(RegionCN))
	}
	if skill.GetDescriptionFull(RegionCN, 5) != "发动时若生命在600以上则7秒内 得分提升80%。若生命不满600时 生命回复500" {
		t.Fatalf("Expected skill description '发动时若生命在600以上则7秒内 得分提升80%%。若生命不满600时 生命回复500', got '%s'", skill.GetDescriptionFull(RegionCN, 5))
	}
	if skill.GetDuration()[4] != 7*time.Second {
		t.Fatalf("Expected skill duration 7 seconds, got '%s'", skill.GetDuration()[4])
	}
	if skill.GetOnceEffect().OnceEffectType != "life" {
		t.Fatalf("Expected skill once effect type 'life', got '%s'", skill.GetOnceEffect().OnceEffectType)
	}
	cardtest, err := client.GetCard(context.Background(), "947")
	if err != nil || cardtest == nil {
		t.Fatalf("Expected to get card by ID '947', got card: %v, error: %v", cardtest, err)
	}
	if !cardtest.IsAvailable(RegionCN) {
		t.Fatal("Expected card to be available in RegionCN")
	}
	if normalCard, err := cardtest.GetArtsFile(context.Background(), RegionCN, CardStatusNormal); err != nil || normalCard == nil {
		t.Fatalf("Expected to get normal card arts file, got card: %v, error: %v", normalCard, err)
	}
	if cardtest.GetArtsURL(context.Background(), RegionCN, CardStatusNormal) != "https://bestdori.com/assets/cn/characters/resourceset/res017040_rip/card_normal.png" {
		t.Fatalf("Expected card arts URL 'https://bestdori.com/assets/cn/characters/resourceset/res017040_rip/card_normal.png', got '%s'", cardtest.GetArtsURL(context.Background(), RegionCN, CardStatusNormal))
	}
	if afterCard, err := cardtest.GetArtsFile(context.Background(), RegionCN, CardStatusAfterTraining); err != nil || afterCard == nil {
		t.Fatalf("Expected to get after training card arts file, got card: %v, error: %v", afterCard, err)
	}
	if cardtest.GetArtsURL(context.Background(), RegionCN, CardStatusAfterTraining) != "https://bestdori.com/assets/cn/characters/resourceset/res017040_rip/card_after_training.png" {
		t.Fatalf("Expected card arts URL 'https://bestdori.com/assets/cn/characters/resourceset/res017040_rip/card_after_training.png', got '%s'", cardtest.GetArtsURL(context.Background(), RegionCN, CardStatusAfterTraining))
	}
	if cardtest.GetAttribute() != "PURE" {
		t.Fatalf("Expected card attribute 'PURE', got '%s'", cardtest.GetAttribute())
	}
	if attr, err := cardtest.GetAttributeIconFile(context.Background()); err != nil || attr == nil {
		t.Fatalf("Expected to get card attribute icon file, got attribute: %v, error: %v", attr, err)
	}
	if cardtest.GetAttributeIconURL() != "https://bestdori.com/res/icon/pure.svg" {
		t.Fatalf("Expected card attribute icon URL 'https://bestdori.com/res/icon/pure.svg', got '%s'", cardtest.GetAttributeIconURL())
	}
	if cardtest.GetCardName(RegionCN) != "日菜抽中的大奖" {
		t.Fatalf("Expected card name '日菜抽中的大奖', got '%s'", cardtest.GetCardName(RegionCN))
	}
	if chara, err := cardtest.GetCharacter(context.Background()); err != nil || chara == nil {
		t.Fatalf("Expected to get card character, got character: %v, error: %v", chara, err)
	}
	if cardtest.GetGachaText(RegionCN) != "比起这个，去找赠品\n更有噜♪的感觉！" {
		t.Fatalf("Expected card gacha text '比起这个，去找赠品\n更有噜♪的感觉！', got '%s'", cardtest.GetGachaText(RegionCN))
	}
	if cardtest.GetLevelLimit() != 50 {
		t.Fatalf("Expected card level limit 50, got '%d'", cardtest.GetLevelLimit())
	}
	if cardtest.GetRarity() != 4 {
		t.Fatalf("Expected card rarity 4 , got '%d'", cardtest.GetRarity())
	}
	if cardtest.GetReleaseTime(RegionCN).Format("2006-01-02 15:04:05") != "2021-08-17 05:00:00" {
		t.Fatalf("Expected card release time '2021-08-17 05:00:00', got '%s'", cardtest.GetReleaseTime(RegionCN).Format("2006-01-02 15:04:05"))
	}
	if cardtest.GetType() != "permanent" {
		t.Fatalf("Expected card type 'permanent', got '%s'", cardtest.GetType())
	}
	if cardtest.GetSkillName(RegionCN) != "很噜♪的要求" {
		t.Fatalf("Expected card skill name '很噜♪的要求', got '%s'", cardtest.GetSkillName(RegionCN))
	}
	if skill, err := cardtest.GetSkill(context.Background()); err != nil || skill == nil {
		t.Fatalf("Expected to get card skill, got skill: %v, error: %v", skill, err)
	}
	if thumbnail, err := cardtest.GetThumbnailFile(context.Background(), RegionCN, CardStatusNormal); err != nil || thumbnail == nil {
		t.Fatalf("Expected to get card thumbnail file, got thumbnail: %v, error: %v", thumbnail, err)
	}
	if cardtest.GetThumbnailURL(context.Background(), RegionCN, CardStatusNormal) != "https://bestdori.com/assets/cn/thumb/chara/card00018_rip/res017040_normal.png" {
		t.Fatalf("Expected card thumbnail URL 'https://bestdori.com/assets/cn/thumb/chara/card00018_rip/res017040_normal.png', got '%s'", cardtest.GetThumbnailURL(context.Background(), RegionCN, CardStatusNormal))
	}
	if after, err := cardtest.GetThumbnailFile(context.Background(), RegionCN, CardStatusAfterTraining); err != nil || after == nil {
		t.Fatalf("Expected to get card after training thumbnail file, got thumbnail: %v, error: %v", after, err)
	}
	if cardtest.GetThumbnailURL(context.Background(), RegionCN, CardStatusAfterTraining) != "https://bestdori.com/assets/cn/thumb/chara/card00018_rip/res017040_after_training.png" {
		t.Fatalf("Expected card after training thumbnail URL 'https://bestdori.com/assets/cn/thumb/chara/card00018_rip/res017040_after_training.png', got '%s'", cardtest.GetThumbnailURL(context.Background(), RegionCN, CardStatusAfterTraining))
	}
	if frame, err := cardtest.GetThumbnailFrameFile(context.Background()); err != nil || frame == nil {
		t.Fatalf("Expected to get card thumbnail frame file, got frame: %v, error: %v", frame, err)
	}
	if cardtest.GetThumbnailFrameURL() != "https://bestdori.com/res/image/card-4.png" {
		t.Fatalf("Expected card thumbnail frame URL 'https://bestdori.com/res/image/card-4.png', got '%s'", cardtest.GetThumbnailFrameURL())
	}
}
