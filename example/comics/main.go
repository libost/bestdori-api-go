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

	comicsMap, err := client.InitComics(context.TODO())
	if err != nil {
		panic(err)
	}
	comic, err := comicsMap.GetComic("1469")
	if err != nil {
		panic(err)
	}
	if !comic.IsAvailable(api.RegionJP) {
		panic("Comic is not available in JP region")
	}
	file, err := comic.GetComicFile(context.TODO(), api.RegionJP)
	if err != nil {
		panic(err)
	}
	os.MkdirAll("./example_output", 0755)
	if err := os.WriteFile("./example_output/comic.png", file, 0644); err != nil {
		panic(err)
	}
	fmt.Println("Comic downloaded successfully")
}
