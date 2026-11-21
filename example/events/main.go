package main

import (
	"context"
	"fmt"
	"time"

	bestdoriapi "github.com/libost/bestdori-api-go"
	"github.com/libost/bestdori-api-go/api"
)

func main() {
	client := bestdoriapi.NewClient("", 10*time.Second)

	events, err := client.GetEvent(context.TODO(), "311")
	if err != nil {
		panic(err)
	}
	if !events.IsAvailable(api.RegionCN) {
		println("Event is not available in CN region")
		return
	}
	fmt.Println(events.GetEventTitle(api.RegionCN))
	cards := events.GetEventCards()
	type CardDetails struct {
		CardName      string
		Bonus         string
		CharacterName string
		BandName      string
		Rarity        int
	}
	cardDetails := make([]CardDetails, 0, len(cards))
	for _, card := range cards {
		cardPtr, per, err := card.GetCard(context.TODO())
		if err != nil {
			panic(err)
		}
		if !cardPtr.IsAvailable(api.RegionCN) {
			fmt.Printf("Card %s is not available in CN region\n", cardPtr.GetCardName(api.RegionCN))
			continue
		}
		character, err := cardPtr.GetCharacter(context.TODO())
		if err != nil {
			panic(err)
		}
		if !character.IsAvailable(api.RegionCN) {
			fmt.Printf("Character %s is not available in CN region\n", character.GetCharacterName(api.RegionCN))
			continue
		}
		bandName, err := character.GetBandName(context.TODO(), api.RegionCN)
		if err != nil {
			panic(err)
		}
		cardDetails = append(cardDetails, CardDetails{
			CardName:      cardPtr.GetCardName(api.RegionCN),
			Bonus:         fmt.Sprintf("%d%%", per),
			CharacterName: character.GetCharacterName(api.RegionCN),
			BandName:      bandName,
			Rarity:        cardPtr.Rarity,
		})
	}
	fmt.Println("Event Cards:")
	for _, detail := range cardDetails {
		fmt.Printf("Card Name: %s, Bonus: %s, Character Name: %s, Band Name: %s, Rarity: %d\n",
			detail.CardName, detail.Bonus, detail.CharacterName, detail.BandName, detail.Rarity)
	}
	rates, err := events.GetRates(context.TODO(), api.RegionCN, 50)
	if err != nil {
		panic(err)
	}
	_, predictTracker, err := rates.GetPredictCutoff(context.TODO(), *events)
	if err != nil {
		panic(err)
	}
	//fmt.Println(eventTracker)
	fmt.Println(predictTracker)
}
