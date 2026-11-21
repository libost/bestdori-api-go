package api

type Gachas struct {
	ResourceName          string   `json:"resourceName"`
	BannerAssetBundleName string   `json:"bannerAssetBundleName"`
	GachaName             []string `json:"gachaName"`
	PublishedAt           []string `json:"publishedAt"`
	ClosedAt              []string `json:"closedAt"`
	Type                  string   `json:"type"`
	NewCards              []int    `json:"newCards"`
}

type GachasMap map[string]Gachas

type Gacha struct {
	Details               []GachaDetailsMap     `json:"details"`
	Rates                 []GachaRatesMap       `json:"rates"`
	PaymentMethods        []GachaPaymentMethods `json:"paymentMethods"`
	ResourceName          string                `json:"resourceName"`
	BannerAssetBundleName string                `json:"bannerAssetBundleName"`
	GachaName             []string              `json:"gachaName"`
	PublishedAt           []string              `json:"publishedAt"`
	ClosedAt              []string              `json:"closedAt"`
	Description           []string              `json:"description"`
	Annotation            []string              `json:"annotation"`
	GachaPeriod           []string              `json:"gachaPeriod"`
	GachaType             string                `json:"gachaType"`
	Type                  string                `json:"type"`
	NewCards              []int                 `json:"newCards"`
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
