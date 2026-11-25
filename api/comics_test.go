package api

import (
	"context"
	"slices"
	"testing"
)

func TestGetComicById(t *testing.T) {
	client := NewApi(nil)
	// Initialize comics first
	comics, err := client.InitComics(context.Background())
	if err != nil {
		t.Fatalf("Failed to initialize comics: %v", err)
	}
	if comics == nil {
		t.Fatal("Expected comics, got nil")
	}
	if err = comics.UpdateComics(context.Background()); err != nil {
		t.Fatalf("Failed to update comics: %v", err)
	}
	comic, err := comics.GetComic("1469")
	if err != nil {
		t.Fatalf("Failed to get comic by ID: %v", err)
	}
	if comic == nil {
		t.Fatal("Expected comic, got nil")
	}
	if !comic.IsAvailable(RegionJP) {
		t.Fatal("Expected comic to be available in JP region")
	}
	if comic.GetTitle(RegionJP) != "第469話 「氷川式勉強術」" {
		t.Fatalf("Expected comic title to be '第469話 「氷川式勉強術」', got '%s'", comic.GetTitle(RegionJP))
	}
	if comic.GetSubTitle(RegionJP) != "第469話" {
		t.Fatalf("Expected comic subtitle to be '第469話', got '%s'", comic.GetSubTitle(RegionJP))
	}
	if comic.GetPublishTime(RegionJP).Format("2006-01-02 15:04:05") != "2026-09-10 09:00:00" {
		t.Fatalf("Expected comic publish time to be '2026-09-10 09:00:00', got '%s'", comic.GetPublishTime(RegionJP).Format("2006-01-02 15:04:05"))
	}
	if comic.GetComicURL(RegionJP) != "https://bestdori.com/assets/jp/comic/comic_fourframe/comic_fourframe_469_rip/comic_fourframe_469.png" {
		t.Fatalf("Expected comic URL to be 'https://bestdori.com/assets/jp/comic/comic_fourframe/comic_fourframe_469_rip/comic_fourframe_469.png', got '%s'", comic.GetComicURL(RegionJP))
	}
	comicFile, err := comic.GetComicFile(context.Background(), RegionJP)
	if err != nil {
		t.Fatalf("Failed to get comic file: %v", err)
	}
	if len(comicFile) == 0 {
		t.Fatal("Expected comic file to be non-empty")
	}
	chars, err := comic.GetCharacters(context.Background())
	if err != nil {
		t.Fatalf("Failed to get comic characters: %v", err)
	}
	if len(chars) == 0 {
		t.Fatal("Expected comic characters to be non-empty")
	}
}

func TestGetComicTitles(t *testing.T) {
	client := NewApi(nil)
	comics, err := client.InitComics(context.Background())
	if err != nil {
		t.Fatalf("Failed to initialize comics: %v", err)
	}
	comicTitles := comics.GetComicTitleMap()
	if len(comicTitles) == 0 {
		t.Fatal("Expected comic titles to be non-empty")
	}
	for i := range comicTitles {
		if comicTitles[i].Id == "" {
			t.Fatal("Expected comic ID to be non-empty")
		}
	}
	for i := range comicTitles {
		if comicTitles[i].Id == "1469" {
			if slices.Contains(comicTitles[i].Title, "第469話 「氷川式勉強術」") {
				return
			}
			t.Fatal("Expected to find comic title '第469話 「氷川式勉強術」' for ID '1469', but did not")
		}
	}
}
