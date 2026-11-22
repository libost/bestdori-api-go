package main

import (
	"context"
	"fmt"
	"os"
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
	thumbnail1, err := card.GetThumbnailFile(context.TODO(), api.RegionCN, api.CardStatusNormal)
	if err != nil {
		panic(err)
	}
	thumbnail2, err := card.GetThumbnailFile(context.TODO(), api.RegionCN, api.CardStatusAfterTraining)
	if err != nil {
		panic(err)
	}
	os.Mkdir("./example_output", 0755)
	err = os.WriteFile("./example_output/thumbnail1.png", thumbnail1, 0644)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile("./example_output/thumbnail2.png", thumbnail2, 0644)
	if err != nil {
		panic(err)
	}
	fmt.Println("Thumbnails saved to ./example_output/thumbnail1.png and ./example_output/thumbnail2.png")
}
