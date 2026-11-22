package api

import "strings"

func indexHelper(str string) int {
	str = strings.ToLower(str)
	switch str {
	case "jp", "powerful", "easy":
		return 0
	case "en", "cool", "normal":
		return 1
	case "tw", "happy", "hard":
		return 2
	case "cn", "pure", "expert":
		return 3
	case "kr", "special":
		return 4
	default:
		return -1
	}
}

func regionHelper(idx int) Region {
	switch idx {
	case 0:
		return RegionJP
	case 1:
		return RegionEN
	case 2:
		return RegionTW
	case 3:
		return RegionCN
	case 4:
		return RegionKR
	default:
		return ""
	}
}
