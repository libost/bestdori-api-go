package api

import (
	"context"
	"encoding/json/v2"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

var (
	ErrBandNotFound   = errors.New("band not found")
	ErrNotValidBandId = errors.New("not a valid band id")
)

type Bands struct {
	BandName []string `json:"bandName"`
}

type BandsMap map[string]Bands

// GetBandName returns the band name for the given region. It returns an error if the band ID is not found or if the region is invalid.
func (c *Character) GetBandName(ctx context.Context, region Region) (string, error) {
	return getBandName(ctx, c.api, c.BandID, region)
}

// GetBandName returns the band name for the given region. It returns an error if the band ID is not found or if the region is invalid.
func (c *Characters) GetBandName(ctx context.Context, region Region) (string, error) {
	return getBandName(ctx, c.api, c.BandID, region)
}

// GetBandName returns the band name for the given region. It returns an error if the band ID is not found or if the region is invalid.
func (s *Songs) GetBandName(ctx context.Context, region Region) (string, error) {
	return getBandName(ctx, s.api, s.BandID, region)
}

// GetBandName returns the band name for the given region. It returns an error if the band ID is not found or if the region is invalid.
func (s *Song) GetBandName(ctx context.Context, region Region) (string, error) {
	return getBandName(ctx, s.api, s.BandID, region)
}

func getBandName(ctx context.Context, api *Api, bandID int, region Region) (string, error) {
	resp, err := api.do(ctx, http.MethodGet, bandsURL)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()
	var bandsMap BandsMap
	if err = json.UnmarshalRead(resp.Body, &bandsMap); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}
	if band, exists := bandsMap[strconv.Itoa(bandID)]; exists {
		return band.BandName[indexHelper(region.String())], nil
	}
	return "", ErrBandNotFound
}
