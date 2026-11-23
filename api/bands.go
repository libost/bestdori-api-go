package api

import (
	"context"
	"encoding/json/v2"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

var (
	ErrBandNotFound   = errors.New("band not found")
	ErrNotValidBandId = errors.New("not a valid band id")
	ErrNoBandLogo     = errors.New("no band logo available")
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

// GetBandLogoURL returns the URL of the band logo for the given character.
//
// If the character's band ID is invalid, an empty string is returned.
func (c *Character) GetBandLogoURL() string {
	if c.BandID < 1 || c.BandID > len(bandFrames) {
		return ""
	}
	return bandFrames[c.BandID-1]
}

// GetBandLogoURL returns the URL of the band logo for the given character.
//
// If the character's band ID is invalid, an empty string is returned.
func (c *Characters) GetBandLogoURL() string {
	if c.BandID < 1 || c.BandID > len(bandFrames) {
		return ""
	}
	return bandFrames[c.BandID-1]
}

// GetBandLogoURL returns the URL of the band logo for the given song.
//
// If the song's band ID is invalid, an empty string is returned.
func (s *Songs) GetBandLogoURL() string {
	if s.BandID < 1 || s.BandID > len(bandFrames) {
		return ""
	}
	return bandFrames[s.BandID-1]
}

// GetBandLogoURL returns the URL of the band logo for the given song.
//
// If the song's band ID is invalid, an empty string is returned.
func (s *Song) GetBandLogoURL() string {
	if s.BandID < 1 || s.BandID > len(bandFrames) {
		return ""
	}
	return bandFrames[s.BandID-1]
}

// GetBandLogoFile retrieves the band logo image file for the given character.
//
// The format of the band logo image is SVG.
//
// If the character's band ID is invalid, an ErrNoBandLogo error is returned.
func (c *Character) GetBandLogoFile(ctx context.Context) ([]byte, error) {
	url := c.GetBandLogoURL()
	if url == "" {
		return nil, ErrNoBandLogo
	}
	resp, err := c.api.do(ctx, http.MethodGet, url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

// GetBandLogoFile retrieves the band logo image file for the given character.
//
// The format of the band logo image is SVG.
//
// If the character's band ID is invalid, an ErrNoBandLogo error is returned.
func (c *Characters) GetBandLogoFile(ctx context.Context) ([]byte, error) {
	url := c.GetBandLogoURL()
	if url == "" {
		return nil, ErrNoBandLogo
	}
	resp, err := c.api.do(ctx, http.MethodGet, url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

// GetBandLogoFile retrieves the band logo image file for the given song.
//
// The format of the band logo image is SVG.
//
// If the song's band ID is invalid, an ErrNoBandLogo error is returned.
func (s *Songs) GetBandLogoFile(ctx context.Context) ([]byte, error) {
	url := s.GetBandLogoURL()
	if url == "" {
		return nil, ErrNoBandLogo
	}
	resp, err := s.api.do(ctx, http.MethodGet, url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

// GetBandLogoFile retrieves the band logo image file for the given song.
//
// The format of the band logo image is SVG.
//
// If the song's band ID is invalid, an ErrNoBandLogo error is returned.
func (s *Song) GetBandLogoFile(ctx context.Context) ([]byte, error) {
	url := s.GetBandLogoURL()
	if url == "" {
		return nil, ErrNoBandLogo
	}
	resp, err := s.api.do(ctx, http.MethodGet, s.GetBandLogoURL())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}
