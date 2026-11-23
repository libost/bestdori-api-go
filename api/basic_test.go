package api

import (
	"testing"
)

func TestDifficuty(t *testing.T) {
	testCases := []struct {
		difficulty Difficulty
		str        string
		color      string
		idx        int
	}{
		{DifficultyEasy, "easy", "#3076F8", 0},
		{DifficultyNormal, "normal", "#188A28", 1},
		{DifficultyHard, "hard", "#D1861A", 2},
		{DifficultyExpert, "expert", "#FD2A30", 3},
		{DifficultySpecial, "special", "#F12199", 4},
	}
	for _, tc := range testCases {
		if tc.difficulty.String() != tc.str {
			t.Errorf("Expected difficulty string '%s', got '%s'", tc.str, tc.difficulty.String())
		}
		if tc.difficulty.Color() != tc.color {
			t.Errorf("Expected difficulty color '%s', got '%s'", tc.color, tc.difficulty.Color())
		}
		if tc.difficulty.Idx() != tc.idx {
			t.Errorf("Expected difficulty index %d, got %d", tc.idx, tc.difficulty.Idx())
		}
	}
}

func TestCardStatus(t *testing.T) {
	testCases := []struct {
		cardStatus CardStatus
		str        string
	}{
		{CardStatusNormal, "normal"},
		{CardStatusAfterTraining, "after_training"},
	}
	for _, tc := range testCases {
		if tc.cardStatus.String() != tc.str {
			t.Errorf("Expected card status string '%s', got '%s'", tc.str, tc.cardStatus.String())
		}
	}
}

func TestSongThumbnailType(t *testing.T) {
	testCases := []struct {
		thumbnailType SongThumbnailType
		str           string
	}{
		{SongThumbnailTypeJacket, "jacket"},
		{SongThumbnailTypeBanner, "banner"},
	}
	for _, tc := range testCases {
		if tc.thumbnailType.String() != tc.str {
			t.Errorf("Expected song thumbnail type string '%s', got '%s'", tc.str, tc.thumbnailType.String())
		}
	}
}

func TestRegion(t *testing.T) {
	testCases := []struct {
		region Region
		idx    int
		str    string
	}{
		{RegionJP, 0, "jp"},
		{RegionEN, 1, "en"},
		{RegionTW, 2, "tw"},
		{RegionCN, 3, "cn"},
		{RegionKR, 4, "kr"},
	}
	for _, tc := range testCases {
		if tc.region.Idx() != tc.idx {
			t.Errorf("Expected region index %d, got %d", tc.idx, tc.region.Idx())
		}
		if tc.region.String() != tc.str {
			t.Errorf("Expected region string '%s', got '%s'", tc.str, tc.region.String())
		}
	}
}
