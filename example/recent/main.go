package main

import (
	"context"
	"fmt"

	bestdoriapi "github.com/libost/bestdori-api-go"
	"github.com/libost/bestdori-api-go/api"
)

func main() {
	client := bestdoriapi.NewClient(nil)
	recent, err := client.GetRecent(context.TODO())
	if err != nil {
		panic(err)
	}
	event := recent.GetRecentEventsForRegion(api.RegionCN)
	if event == nil {
		panic("No recent event found for CN region")
	}
	if len(event) != 1 {
		panic(fmt.Sprintf("Expected 1 recent event for CN region, got %d", len(event)))
	}
	recentEvent := event[0]
	eventDetail, err := recentEvent.GetEvent(context.TODO())
	if err != nil {
		panic(err)
	}
	if !eventDetail.IsAvailable(api.RegionCN) {
		panic("Event is not available in CN region")
	}
	fmt.Println(eventDetail.GetTitle(api.RegionCN))
	gacha := recent.GetRecentGachasForRegion(api.RegionCN)
	if gacha == nil {
		panic("No recent gacha found for CN region")
	}
	gachaDetail, err := gacha[0].GetGacha(context.TODO())
	if err != nil {
		panic(err)
	}
	if !gachaDetail.IsAvailable(api.RegionCN) {
		panic("Gacha is not available in CN region")
	}
	fmt.Println(gachaDetail.GetTitle(api.RegionCN))
	recentsongs := recent.GetRecentSongsForRegion(api.RegionCN)
	if recentsongs == nil {
		panic("No recent songs found for CN region")
	}
	songDetail, err := recentsongs[0].GetSong(context.TODO())
	if err != nil {
		panic(err)
	}
	if !songDetail.IsAvailable(api.RegionCN) {
		panic("Song is not available in CN region")
	}
	fmt.Println(songDetail.GetTitle(api.RegionCN))
}
