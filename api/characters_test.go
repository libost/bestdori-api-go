package api

import (
	"context"
	"slices"
	"testing"
)

func TestGetCharacterById(t *testing.T) {
	client := NewApi(nil)
	// Initialize characters first
	characters, err := client.InitCharacters(context.Background())
	if err != nil {
		t.Fatalf("Failed to initialize characters: %v", err)
	}
	if characters == nil {
		t.Fatal("Expected characters, got nil")
	}
	if err = characters.UpdateCharacters(context.Background()); err != nil {
		t.Fatalf("Failed to update characters: %v", err)
	}
	sayo, err := characters.GetCharacter("22")
	if err != nil {
		t.Fatalf("Failed to get character by ID: %v", err)
	}
	if sayo == nil {
		t.Fatal("Expected character, got nil")
	}
	if !sayo.IsAvailable(RegionCN) {
		t.Fatal("Expected character to be available in CN region")
	}
	if sayo.GetCharacterName(RegionCN) != "冰川 纱夜" {
		t.Fatalf("Expected character name '冰川 纱夜', got '%s'", sayo.GetCharacterName(RegionCN))
	}
	if sayo.GetCharacterColor() != "#00AABB" {
		t.Fatalf("Expected character color '#00AABB', got '%s'", sayo.GetCharacterColor())
	}
	if sayo.GetCharacterType() != "unique" {
		t.Fatalf("Expected character type 'unique', got '%s'", sayo.GetCharacterType())
	}
	bandlogo, err := sayo.GetBandLogoFile(context.Background())
	if err != nil {
		t.Fatalf("Failed to get character band logo file: %v", err)
	}
	if bandlogo == nil {
		t.Fatal("Expected character band logo file, got nil or empty")
	}
	bandname, err := sayo.GetBandName(context.Background(), RegionCN)
	if err != nil {
		t.Fatalf("Failed to get character band name: %v", err)
	}
	if bandname != "Roselia" {
		t.Fatal("Expected character band name 'Roselia', got empty string")
	}
	if sayo.GetCharacterFirstName(RegionCN) != "纱夜" {
		t.Fatalf("Expected character first name '纱夜', got '%s'", sayo.GetCharacterFirstName(RegionCN))
	}
	if sayo.GetCharacterLastName(RegionCN) != "冰川" {
		t.Fatalf("Expected character last name '冰川', got '%s'", sayo.GetCharacterLastName(RegionCN))
	}
	tomori, err := characters.GetCharacter("36")
	if err != nil {
		t.Fatalf("Failed to get character by ID: %v", err)
	}
	if tomori == nil {
		t.Fatal("Expected character, got nil")
	}
	if tomori.IsAvailable(RegionKR) {
		t.Fatal("Expected character to be unavailable in KR region")
	}
	chu2, err := characters.GetCharacter("35")
	if err != nil {
		t.Fatalf("Failed to get character by ID: %v", err)
	}
	if chu2 == nil {
		t.Fatal("Expected character, got nil")
	}
	if chu2.GetCharacterNickname(RegionJP) != "チュチュ" {
		t.Fatalf("Expected character nickname 'チュチュ', got '%s'", chu2.GetCharacterNickname(RegionJP))
	}
	sayod, err := client.GetCharacter(context.Background(), "22")
	if err != nil {
		t.Fatalf("Failed to get character by ID: %v", err)
	}
	if sayod == nil {
		t.Fatal("Expected character, got nil")
	}
	if !sayod.IsAvailable(RegionCN) {
		t.Fatal("Expected character to be available in CN region")
	}
}

func TestGetCharacterNames(t *testing.T) {
	client := NewApi(nil)

	charMap, err := client.InitCharacters(context.Background())
	if err != nil {
		t.Fatalf("Failed to initialize characters: %v", err)
	}
	nameMap := charMap.GetCharacterNameMap()
	if len(nameMap) == 0 {
		t.Fatal("Expected character name map to be non-empty")
	}
	for i := range nameMap {
		if nameMap[i].Id == "" {
			t.Fatal("Expected character ID to be non-empty")
		}
	}
	for i := range nameMap {
		if nameMap[i].Id == "22" {
			if slices.Contains(nameMap[i].Names, "冰川 纱夜") {
				return
			}
			t.Fatal("Expected to find character name '冰川 纱夜' for ID '22', but did not")
		}
	}
	for i := range nameMap {
		if nameMap[i].Id == "31" {
			if slices.Contains(nameMap[i].Names, "Layer") {
				return
			}
			t.Fatal("Expected to find character nickname 'Layer' for ID '31', but did not")
		}
	}
}
