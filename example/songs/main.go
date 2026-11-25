package main

import (
	"context"
	"fmt"

	bestdoriapi "github.com/libost/bestdori-api-go"
	"github.com/libost/bestdori-api-go/api"
)

func main() {
	client := bestdoriapi.NewClient(nil)
	songsMap, err := client.InitSongs(context.TODO())
	if err != nil {
		panic(err)
	}
	songs, err := songsMap.GetSong("306")
	if err != nil {
		panic(err)
	}
	if !songs.IsAvailable(api.RegionCN) {
		fmt.Println("Song is not available in CN region")
		return
	}
	fmt.Println(songs.GetTitle(api.RegionCN))
	songMeta, err := client.InitSongMeta(context.TODO())
	if err != nil {
		panic(err)
	}
	m, err := songMeta.GetSongMetaFactors("306", "7", api.DifficultyExpert)
	if err != nil {
		panic(err)
	}
	song, err := client.GetSong(context.TODO(), "306")
	if err != nil {
		panic(err)
	}
	duration := song.GetDuration()
	fmt.Println(m.Efficiency(true, 1.0, 1.1, duration.Seconds()))
	fmt.Println(m.ScoreFactor(true, 1.0, 1.1))
}
