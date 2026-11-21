package api

import "testing"

func TestNewRegionTier(t *testing.T) {
	regionTier, err := NewRegionTier(RegionJP, 100)
	if err != nil {
		t.Fatalf("NewRegionTier() error = %v", err)
	}
	if regionTier.Region != RegionJP || regionTier.Tier != 100 {
		t.Fatalf("NewRegionTier() = %+v, want region jp and tier 100", regionTier)
	}
}

func TestNewRegionTierRejectsUnsupportedTier(t *testing.T) {
	if _, err := NewRegionTier(RegionTW, 10000); err == nil {
		t.Fatal("NewRegionTier() accepted a tier unsupported by the region")
	}
}

func TestRegionTiers(t *testing.T) {
	if got := RegionKR.Tiers(); len(got) != 1 || got[0] != 100 {
		t.Fatalf("RegionKR.Tiers() = %v, want [100]", got)
	}
	if got := Region("unknown").Tiers(); got != nil {
		t.Fatalf("unknown Region.Tiers() = %v, want nil", got)
	}
}
