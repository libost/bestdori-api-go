package api

import (
	"bytes"
	"context"
	"encoding/json/v2"
	"fmt"
	"io"
	"math"
	"math/rand/v2"
	"net/http"
	"time"
)

type Gachas struct {
	id                    string
	api                   *Api
	ResourceName          string      `json:"resourceName"`
	BannerAssetBundleName string      `json:"bannerAssetBundleName"`
	GachaName             []string    `json:"gachaName"`
	PublishedAt           []Timestamp `json:"publishedAt"`
	ClosedAt              []Timestamp `json:"closedAt"`
	Type                  string      `json:"type"`
	NewCards              []NewCards  `json:"newCards"`
}

type NewCards struct {
	a  *Api
	id int
}

func (n *NewCards) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte("null")) {
		n.id = 0
		return nil
	}
	var v int
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	n.id = v
	return nil
}

func (n NewCards) MarshalJSON() ([]byte, error) {
	if n.id == 0 {
		return []byte("null"), nil
	}
	return json.Marshal(n.id)
}

type GachasMap map[string]Gachas

type Gacha struct {
	api                   *Api
	Details               []GachaDetailsMap     `json:"details"`
	Rates                 []GachaRatesMap       `json:"rates"`
	PaymentMethods        []GachaPaymentMethods `json:"paymentMethods"`
	ResourceName          string                `json:"resourceName"`
	BannerAssetBundleName string                `json:"bannerAssetBundleName"`
	GachaName             []string              `json:"gachaName"`
	PublishedAt           []Timestamp           `json:"publishedAt"`
	ClosedAt              []Timestamp           `json:"closedAt"`
	Description           []string              `json:"description"`
	Annotation            []string              `json:"annotation"`
	GachaPeriod           []string              `json:"gachaPeriod"`
	GachaType             string                `json:"gachaType"`
	Type                  string                `json:"type"`
	NewCards              []NewCards            `json:"newCards"`
	Information           struct {
		Description   []string `json:"description"`
		Term          []string `json:"term"`
		NewMemberInfo []string `json:"newMemberInfo"`
		Notice        []string `json:"notice"`
	} `json:"information"`
}

type GachaDetails struct {
	RarityIndex int  `json:"rarityIndex"`
	Weight      int  `json:"weight"`
	Pickup      bool `json:"pickup"`
}

type GachaDetailsMap map[string]GachaDetails

type GachaRates struct {
	Rate        float64 `json:"rate"`
	WeightTotal int     `json:"weightTotal"`
}

type GachaRatesMap map[string]GachaRates

type GachaPaymentMethods struct {
	GachaID          int    `json:"gachaId"`
	PaymentMethod    string `json:"paymentMethod"`
	Quantity         int    `json:"quantity"`
	PaymentMethodID  int    `json:"paymentMethodId"`
	Count            int    `json:"count"`
	Behavior         string `json:"behavior"`
	Pickup           bool   `json:"pickup"`
	CostItemQuantity int    `json:"costItemQuantity"`
	DiscountType     int    `json:"discountType"`
}

func (a *Api) InitGachas(ctx context.Context) (*GachasMap, error) {
	resp, err := a.do(ctx, http.MethodGet, gachaURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var gachas GachasMap
	if err = json.UnmarshalRead(resp.Body, &gachas); err != nil {
		return nil, err
	}
	for id, gacha := range gachas {
		gacha.api = a
		gacha.id = id
		gachas[id] = gacha
	}
	return &gachas, nil
}

func (g *GachasMap) GetGacha(gachaID string) (*Gachas, error) {
	gacha, ok := (*g)[gachaID]
	if !ok {
		return nil, fmt.Errorf("gacha with ID %s not found", gachaID)
	}
	return &gacha, nil
}

func (a *Api) GetGacha(ctx context.Context, gachaID string) (*Gacha, error) {
	resp, err := a.do(ctx, http.MethodGet, "https://bestdori.com/api/gacha/"+gachaID+".json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var gacha Gacha
	if err = json.UnmarshalRead(resp.Body, &gacha); err != nil {
		return nil, err
	}
	gacha.api = a
	return &gacha, nil
}

func (g *Gachas) IsAvailable(region Region) bool {
	if g.GachaName[region.Idx()] == "" {
		return false
	}
	return true
}

func (g *Gacha) IsAvailable(region Region) bool {
	if g.GachaName[region.Idx()] == "" {
		return false
	}
	return true
}

func (r *RecentGachas) GetGacha(ctx context.Context) (*Gacha, error) {
	return r.api.GetGacha(ctx, r.id)
}

func (g *Gachas) GetGacha(ctx context.Context) (*Gacha, error) {
	return g.api.GetGacha(ctx, g.id)
}

func (g *Gachas) GetTitle(region Region) string {
	return g.GachaName[region.Idx()]
}

func (g *Gacha) GetTitle(region Region) string {
	return g.GachaName[region.Idx()]
}

func (g *Gachas) GetPublishTime(region Region) time.Time {
	return g.PublishedAt[region.Idx()].Time
}

func (g *Gacha) GetPublishTime(region Region) time.Time {
	return g.PublishedAt[region.Idx()].Time
}

func (g *Gachas) GetCloseTime(region Region) time.Time {
	return g.ClosedAt[region.Idx()].Time
}

func (g *Gacha) GetCloseTime(region Region) time.Time {
	return g.ClosedAt[region.Idx()].Time
}

func (g *Gachas) GetType() string {
	return g.Type
}

func (g *Gacha) GetType() string {
	return g.Type
}

func (g *Gachas) GetBannerURL(ctx context.Context, region Region) string {
	return getBannerURL(ctx, g.api, region, g.BannerAssetBundleName, g.ResourceName)
}

func (g *Gacha) GetBannerURL(ctx context.Context, region Region) string {
	return getBannerURL(ctx, g.api, region, g.BannerAssetBundleName, g.ResourceName)
}

func getBannerURL(ctx context.Context, api *Api, region Region, bannerAssetBundleName, resourceName string) string {
	url := "https://bestdori.com/assets/" + region.String() + "/homebanner_rip/" + bannerAssetBundleName + ".png"
	if !isValidPNG(ctx, api, url) {
		url = "https://bestdori.com/assets/" + region.String() + "/gacha/screen/" + resourceName + "_rip/logo.png"
		if !isValidPNG(ctx, api, url) {
			url = "" // give up if both URLs are invalid
		}
	}
	return url
}

func (g *Gachas) GetBannerFile(ctx context.Context, region Region) ([]byte, error) {
	url := g.GetBannerURL(ctx, region)
	if url == "" {
		return nil, fmt.Errorf("no valid banner URL found for region %s", region.String())
	}
	return getGachaBannerFile(ctx, g.api, url)
}

func (g *Gacha) GetBannerFile(ctx context.Context, region Region) ([]byte, error) {
	url := g.GetBannerURL(ctx, region)
	if url == "" {
		return nil, fmt.Errorf("no valid banner URL found for region %s", region.String())
	}
	return getGachaBannerFile(ctx, g.api, url)
}

func getGachaBannerFile(ctx context.Context, api *Api, url string) ([]byte, error) {
	resp, err := api.do(ctx, http.MethodGet, url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch banner: %s", resp.Status)
	}
	return io.ReadAll(resp.Body)
}

func (g *Gachas) GetNewCards() []NewCards {
	return g.NewCards
}

func (g *Gacha) GetNewCards() []NewCards {
	var newCards []NewCards
	for _, card := range g.NewCards {
		newCards = append(newCards, NewCards{
			a:  g.api,
			id: card.id,
		})
	}
	return newCards
}

func (g *Gacha) GetDescription(region Region) string {
	return g.Description[region.Idx()]
}

func (g *Gacha) GetAnnotation(region Region) string {
	return g.Annotation[region.Idx()]
}

func (g *Gacha) GetGachaPeriod(region Region) string {
	return g.GachaPeriod[region.Idx()]
}

func (g *Gacha) GetCardSpecifiedWeight() GachaDetailsMap {
	return g.Details[0]
}

func (g *Gacha) GetRaritySpecifiedWeight() GachaRatesMap {
	return g.Rates[0]
}

func (g *Gacha) GetPaymentMethods() []GachaPaymentMethods {
	return g.PaymentMethods
}

type gachaMap struct {
	Id          string
	Weight      int
	WeightTotal int
	Rarity      int
}

func (g *Gacha) EmulateGachaDraw(paymentMethods GachaPaymentMethods) ([]string, error) {
	var gachaMaps []gachaMap
	details := g.GetCardSpecifiedWeight()
	for id, gacha := range details {
		gachaMaps = append(gachaMaps, gachaMap{
			Id:     id,
			Weight: gacha.Weight,
			Rarity: gacha.RarityIndex,
		})
	}
	var rarityRates [5]float64
	rates := g.GetRaritySpecifiedWeight()
	for rarity, rate := range rates {
		var rarityIndex int
		if _, err := fmt.Sscanf(rarity, "%d", &rarityIndex); err != nil {
			return nil, fmt.Errorf("invalid rarity %q: %w", rarity, err)
		}

		if rarityIndex < 1 || rarityIndex > len(rarityRates) {
			return nil, fmt.Errorf("rarity out of range: %d", rarityIndex)
		}

		rarityRates[rarityIndex-1] = rate.Rate
	}
	rarityWeightsInt := make([]int, len(rarityRates))
	for i, rate := range rarityRates {
		rarityWeightsInt[i] = int(math.Round(float64(rate) * 100))
	}
	guaranteedRarity := 3
	switch paymentMethods.Behavior {
	case "fixed_5_star_once":
		guaranteedRarity = 5
	}
	hasFulfilledGuarantee := false
	results := make([]string, 0, paymentMethods.Count)
	for i := range paymentMethods.Count {
		if i == paymentMethods.Count-1 && paymentMethods.Count != 1 && !hasFulfilledGuarantee {

			result, err := generateGachaResultByRates(gachaMaps, rarityWeightsInt, guaranteedRarity)
			if err != nil {
				return nil, err
			}
			results = append(results, result)
		} else {
			result, err := generateGachaResultByRates(gachaMaps, rarityWeightsInt, 1)
			if err != nil {
				return nil, err
			}
			gacha := details[result]
			if gacha.RarityIndex >= guaranteedRarity {
				hasFulfilledGuarantee = true
			}
			results = append(results, result)
		}
	}
	return results, nil
}

func generateGachaResultByRates(gachaMaps []gachaMap, rarityWeights []int, guarantee int) (string, error) {
	if len(gachaMaps) == 0 {
		return "", fmt.Errorf("no gacha cards available")
	}

	rarityTotal := 0
	eligibleRarities := make([]int, 0, len(rarityWeights))
	for rarity, weight := range rarityWeights {
		rarityIndex := rarity + 1
		if rarityIndex < guarantee || weight <= 0 {
			continue
		}
		hasCards := false
		for _, gacha := range gachaMaps {
			if gacha.Rarity == rarityIndex && gacha.Weight > 0 {
				hasCards = true
				break
			}
		}
		if hasCards {
			rarityTotal += weight
			eligibleRarities = append(eligibleRarities, rarityIndex)
		}
	}
	if len(eligibleRarities) == 0 {
		return "", fmt.Errorf("no gacha rarities meet rarity guarantee: %d", guarantee)
	}

	randomNum, err := randomNumGen(1, rarityTotal)
	if err != nil {
		return "", err
	}
	selectedRarity := 0
	for _, rarity := range eligibleRarities {
		randomNum -= rarityWeights[rarity-1]
		if randomNum <= 0 {
			selectedRarity = rarity
			break
		}
	}
	if selectedRarity == 0 {
		return "", fmt.Errorf("failed to generate gacha rarity result")
	}

	cardWeightTotal := 0
	for _, gacha := range gachaMaps {
		if gacha.Rarity == selectedRarity && gacha.Weight > 0 {
			cardWeightTotal += gacha.Weight
		}
	}
	randomNum, err = randomNumGen(1, cardWeightTotal)
	if err != nil {
		return "", err
	}
	for _, gacha := range gachaMaps {
		if gacha.Rarity != selectedRarity || gacha.Weight <= 0 {
			continue
		}
		randomNum -= gacha.Weight
		if randomNum <= 0 {
			return gacha.Id, nil
		}
	}
	return "", fmt.Errorf("failed to generate gacha card result")
}

func randomNumGen(min, max int) (int, error) {
	if min > max {
		return 0, fmt.Errorf("min cannot be greater than max")
	}
	return rand.N(max-min+1) + min, nil
}
