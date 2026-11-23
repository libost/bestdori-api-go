package api

import "testing"

func TestIndexHelper(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected int
	}{
		{"JP", "jp", 0},
		{"EN", "en", 1},
		{"TW", "tw", 2},
		{"CN", "cn", 3},
		{"KR", "kr", 4},
		{"Powerful", "powerful", 0},
		{"Cool", "cool", 1},
		{"Happy", "happy", 2},
		{"Pure", "pure", 3},
		{"Special", "special", 4},
		{"Easy", "easy", 0},
		{"Normal", "normal", 1},
		{"Hard", "hard", 2},
		{"Expert", "expert", 3},
		{"Invalid", "invalid", -1},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := indexHelper(tc.input)
			if result != tc.expected {
				t.Errorf("Expected %d for input '%s', got %d", tc.expected, tc.input, result)
			}
		})
	}
}

func Test_regionHelper(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		idx  int
		want Region
	}{
		{"JP", 0, RegionJP},
		{"EN", 1, RegionEN},
		{"TW", 2, RegionTW},
		{"CN", 3, RegionCN},
		{"KR", 4, RegionKR},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := regionHelper(tt.idx)
			if got != tt.want {
				t.Errorf("regionHelper() = %v, want %v", got, tt.want)
			}
		})
	}
}
