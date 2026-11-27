package api

import (
	"bytes"
	"context"
	"encoding/json/v2"
	"fmt"
	"io"
	"maps"
	"math"
	"math/rand/v2"
	"net/http"
	"time"
)

type Gachas struct {
	id                    string
	api                   *Api
	etag                  string
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

// InitGachas initializes the GachasMap by fetching data from the API and populating the GachasMap. It returns an error if the request fails or if the response cannot be unmarshaled.
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
		gacha.etag = resp.Header.Get("ETag")
		gachas[id] = gacha
	}
	return &gachas, nil
}

// UpdateGachas updates the GachasMap with the latest data from the API.
//
// It uses the ETag header to check if the data has changed since the last update. If the data has not changed, it does nothing. If the data has changed, it updates the GachasMap with the new data and initializes the api field for each gacha in the map.
func (g *GachasMap) UpdateGachas(ctx context.Context) error {
	var etag string
	var api *Api
	for _, gacha := range *g {
		api = gacha.api
		if api == nil {
			continue
		}
		etag = gacha.etag
		if etag == "" {
			continue
		}
		break
	}
	resp, err := api.update(ctx, http.MethodGet, gachaURL, etag)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotModified {
		return nil
	}
	var gachas GachasMap
	if err = json.UnmarshalRead(resp.Body, &gachas); err != nil {
		return err
	}
	for gachaID, gacha := range gachas {
		gacha.api = api
		gacha.etag = resp.Header.Get("ETag")
		gachas[gachaID] = gacha
	}
	maps.Copy(*g, gachas)
	return nil
}

// GetGacha retrieves a specific Gacha from the GachasMap by its ID. It returns an error if the Gacha is not found.
//
// Be aware that the struct returned by this function is a summary of the gacha data and does not contain all the information available in the API.
//
// Use (api.Api).GetGacha(ctx, gachaID) to get the full gacha data.
func (g *GachasMap) GetGacha(gachaID string) (*Gachas, error) {
	gacha, ok := (*g)[gachaID]
	if !ok {
		return nil, fmt.Errorf("gacha with ID %s not found", gachaID)
	}
	return &gacha, nil
}

type GachaTitleMap struct {
	Id     string
	Titles []string
}

// GetGachaTitleMap retrieves a slice of GachaTitleMap containing the IDs and titles of all gachas in the GachasMap.
func (g *GachasMap) GetGachaTitleMap() []GachaTitleMap {
	titles := make([]GachaTitleMap, 0, len(*g))
	for id, gacha := range *g {
		titles = append(titles, GachaTitleMap{
			Id:     id,
			Titles: gacha.GachaName,
		})
	}
	return titles
}

// GetGacha retrieves a specific Gacha by its ID. It returns an error if the Gacha is not found or if the API call fails.
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

// IsAvailable checks if the Gachas is available in the specified region.
//
// It is recommended to use this method before calling functions like GetTitle.
func (g *Gachas) IsAvailable(region Region) bool {
	return g.GachaName[region.Idx()] != ""
}

// IsAvailable checks if the Gacha is available in the specified region.
//
// It is recommended to use this method before calling functions like GetTitle.
func (g *Gacha) IsAvailable(region Region) bool {
	return g.GachaName[region.Idx()] != ""
}

// GetGacha retrieves the Gacha associated with the RecentGachas. It returns an error if the Gacha is not found or if the API call fails.
func (r *RecentGachas) GetGacha(ctx context.Context) (*Gacha, error) {
	return r.api.GetGacha(ctx, r.id)
}

// GetGacha retrieves the detailed information associated with the Gachas. It returns an error if the Gacha is not found or if the API call fails.
func (g *Gachas) GetGacha(ctx context.Context) (*Gacha, error) {
	return g.api.GetGacha(ctx, g.id)
}

// GetTitle retrieves the title of the Gachas in the specified region. It returns an empty string if the Gachas is not available in that region.
func (g *Gachas) GetTitle(region Region) string {
	return g.GachaName[region.Idx()]
}

// GetTitle retrieves the title of the Gacha in the specified region. It returns an empty string if the Gacha is not available in that region.
func (g *Gacha) GetTitle(region Region) string {
	return g.GachaName[region.Idx()]
}

// GetPublishTime retrieves the publish time of the Gachas in the specified region. It returns a zero time if the Gachas is not available in that region.
func (g *Gachas) GetPublishTime(region Region) time.Time {
	return g.PublishedAt[region.Idx()].Time
}

// GetPublishTime retrieves the publish time of the Gacha in the specified region. It returns a zero time if the Gacha is not available in that region.
func (g *Gacha) GetPublishTime(region Region) time.Time {
	return g.PublishedAt[region.Idx()].Time
}

// GetCloseTime retrieves the close time of the Gachas in the specified region. It returns a zero time if the Gachas is not available in that region.
func (g *Gachas) GetCloseTime(region Region) time.Time {
	return g.ClosedAt[region.Idx()].Time
}

// GetCloseTime retrieves the close time of the Gacha in the specified region. It returns a zero time if the Gacha is not available in that region.
func (g *Gacha) GetCloseTime(region Region) time.Time {
	return g.ClosedAt[region.Idx()].Time
}

// GetType retrieves the type of the Gachas.
func (g *Gachas) GetType() string {
	return g.Type
}

// GetType retrieves the type of the Gacha.
func (g *Gacha) GetType() string {
	return g.Type
}

// GetBannerURL retrieves the banner URL of the Gachas in the specified region. It returns an empty string if the Gachas is not available in that region.
//
// If the banner URL is not valid, it will attempt to use a fallback URL. If both URLs are invalid, it will return an empty string.
//
// Be aware that for some gacha, the banner URL for certain regions may not be valid, even if the gacha is available in that region. In such cases, the function will return an empty string.
func (g *Gachas) GetBannerURL(ctx context.Context, region Region) string {
	return getBannerURL(ctx, g.api, region, g.BannerAssetBundleName, g.ResourceName)
}

// GetBannerURL retrieves the banner URL of the Gacha in the specified region. It returns an empty string if the Gacha is not available in that region.
//
// If the banner URL is not valid, it will attempt to use a fallback URL. If both URLs are invalid, it will return an empty string.
//
// Be aware that for some gacha, the banner URL for certain regions may not be valid, even if the gacha is available in that region. In such cases, the function will return an empty string.
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

// GetBannerFile retrieves the banner file of the Gachas in the specified region. It returns an error if the banner file cannot be fetched or if the banner URL is not valid.
//
// If no valid banner URL is found for the specified region, it will return an error indicating that the banner URL is not available.
//
// Be aware that for some gacha, the banner file for certain regions may not be valid, even if the gacha is available in that region. In such cases, the function will return an error.
func (g *Gachas) GetBannerFile(ctx context.Context, region Region) ([]byte, error) {
	url := g.GetBannerURL(ctx, region)
	if url == "" {
		return nil, fmt.Errorf("no valid banner URL found for region %s", region.String())
	}
	return getGachaBannerFile(ctx, g.api, url)
}

// GetBannerFile retrieves the banner file of the Gacha in the specified region. It returns an error if the banner file cannot be fetched or if the banner URL is not valid.
//
// If no valid banner URL is found for the specified region, it will return an error indicating that the banner URL is not available.
//
// Be aware that for some gacha, the banner file for certain regions may not be valid, even if the gacha is available in that region. In such cases, the function will return an error.
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

// GetNewCards retrieves the new cards associated with the Gachas. It returns a slice of NewCards.
func (g *Gachas) GetNewCards() []NewCards {
	var newCards []NewCards
	for _, card := range g.NewCards {
		newCards = append(newCards, NewCards{
			a:  g.api,
			id: card.id,
		})
	}
	return newCards
}

// GetNewCards retrieves the new cards associated with the Gacha. It returns a slice of NewCards.
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

// GetDescription retrieves the description of the Gacha localized in the specified region. It returns an empty string if the Gacha is not available in that region.
func (g *Gacha) GetDescription(region Region) string {
	return g.Description[region.Idx()]
}

// GetAnnotation retrieves the annotation of the Gacha localized in the specified region. It returns an empty string if the Gacha is not available in that region.
func (g *Gacha) GetAnnotation(region Region) string {
	return g.Annotation[region.Idx()]
}

// GetGachaPeriod retrieves the gacha period string of the Gacha localized in the specified region. It returns an empty string if the Gacha is not available in that region.
func (g *Gacha) GetGachaPeriod(region Region) string {
	return g.GachaPeriod[region.Idx()]
}

// GetCardSpecifiedWeight retrieves the card-specific weight details of the Gacha. It returns a GachaDetailsMap containing the weight information for each card in the Gacha.
//
// It can be used to determine the probability of obtaining specific cards during a gacha draw.
func (g *Gacha) GetCardSpecifiedWeight() GachaDetailsMap {
	return g.Details[0]
}

// GetRaritySpecifiedWeight retrieves the rarity-specific weight details of the Gacha. It returns a GachaRatesMap containing the weight information for each rarity level in the Gacha.
//
// It can be used to determine the probability of obtaining cards of specific rarities during a gacha draw.
func (g *Gacha) GetRaritySpecifiedWeight() GachaRatesMap {
	return g.Rates[0]
}

// GetPaymentMethods retrieves the payment methods associated with the Gacha. It returns a slice of GachaPaymentMethods, each representing a different payment option available for the gacha draw.
//
// It can be used to determine the cost and behavior of different gacha draw options.
//
// It can also be used to emulate gacha draws by providing the selected payment method to the EmulateGachaDraw function.
func (g *Gacha) GetPaymentMethods() []GachaPaymentMethods {
	return g.PaymentMethods
}

type gachaMap struct {
	Id          string
	Weight      int
	WeightTotal int
	Rarity      int
}

// EmulateGachaDraw simulates a gacha draw based on the specified payment methods. It returns a slice of card IDs representing the results of the gacha draw.
//
// The paymentMethods parameter specifies the payment method and count for the gacha draw. The function uses the card-specific and rarity-specific weights to determine the outcome of the draw.
//
// It can be used to simulate gacha draws without making actual API calls, allowing users to test different payment methods and understand the probabilities of obtaining specific cards or rarities.
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

// GetMethod retrieves the payment method of the GachaPaymentMethods.
//
// It can be used to determine the type of payment method used for the gacha draw.
func (p *GachaPaymentMethods) GetMethod() string {
	return p.PaymentMethod
}

// GetCost retrieves the cost of the GachaPaymentMethods.
//
// Be aware that the cost returned is the average cost per draw, use GetTotalCost to get the total cost.
func (p *GachaPaymentMethods) GetCost() int {
	return p.Quantity
}

// GetCount retrieves the count of draws for the GachaPaymentMethods.
//
// It can be used to determine how many draws are performed with the specified payment method.
func (p *GachaPaymentMethods) GetDraws() int {
	return p.Count
}

// GetDiscountType retrieves the guaranteed type of the GachaPaymentMethods.
func (p *GachaPaymentMethods) GetBehavior() string {
	return p.Behavior
}

// IsPickup checks if the GachaPaymentMethods is a pickup gacha. It returns true if the gacha is a pickup gacha, and false otherwise.
func (p *GachaPaymentMethods) IsPickup() bool {
	return p.Pickup
}

// GetTotalCost retrieves the total cost of the GachaPaymentMethods.
func (p *GachaPaymentMethods) GetTotalCost() int {
	return p.CostItemQuantity
}
