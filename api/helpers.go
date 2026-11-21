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

func tzHelper(idx int) string {
	switch idx {
	case 0:
		return "Asia/Tokyo"
	case 1:
		return "America/Los_Angeles"
	case 2:
		return "Asia/Taipei"
	case 3:
		return "Asia/Shanghai"
	case 4:
		return "Asia/Seoul"
	default:
		return ""
	}
}
