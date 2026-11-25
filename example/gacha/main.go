package main

import (
	"context"
	"fmt"

	bestdoriapi "github.com/libost/bestdori-api-go"
	"github.com/libost/bestdori-api-go/api"
)

func main() {
	client := bestdoriapi.NewClient(nil)
	gachasMap, err := client.InitGachas(context.TODO())
	if err != nil {
		panic(err)
	}
	gachas, err := gachasMap.GetGacha("1")
	if err != nil {
		panic(err)
	}
	if !gachas.IsAvailable(api.RegionCN) {
		fmt.Println("Gacha is not available in CN region")
		return
	}
	fmt.Println(gachas.GetTitle(api.RegionCN))
	gacha, err := client.GetGacha(context.TODO(), "1")
	if err != nil {
		panic(err)
	}
	if !gacha.IsAvailable(api.RegionCN) {
		fmt.Println("Gacha is not available in CN region")
		return
	}
	fmt.Println(gacha.GetTitle(api.RegionCN))
	paymentMethods := gacha.GetPaymentMethods()
	results, err := gacha.EmulateGachaDraw(paymentMethods[2])
	if err != nil {
		panic(err)
	}
	fmt.Println(results)
}
