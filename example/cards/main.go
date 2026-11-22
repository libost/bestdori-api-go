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

	card, err := client.GetCard(context.TODO(), "947")
	if err != nil {
		panic(err)
	}
	if !card.IsAvailable(api.RegionCN) {
		fmt.Println("Card is not available in CN region")
		return
	}
	fmt.Println(card.GetCardName(api.RegionCN))
	char, err := card.GetCharacter(context.TODO())
	if err != nil {
		panic(err)
	}
	if !char.IsAvailable(api.RegionCN) {
		fmt.Println("Character is not available in CN region")
		return
	}
	fmt.Println(char.GetCharacterName(api.RegionCN))
	bandName, err := char.GetBandName(context.TODO(), api.RegionCN)
	if err != nil {
		panic(err)
	}
	fmt.Println(bandName)
	skill, err := card.GetSkill(context.TODO())
	if err != nil {
		panic(err)
	}
	if !skill.IsAvailable(api.RegionCN) {
		fmt.Println("Skill is not available in CN region")
		return
	}
	fmt.Println(skill.GetDescriptionFull(api.RegionCN, 5))
}
