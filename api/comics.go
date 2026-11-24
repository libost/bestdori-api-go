package api

import (
	"context"
	"encoding/json/v2"
	"errors"
	"fmt"
	"io"
	"maps"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	ErrComicNotFound   = errors.New("comic not found")
	ErrNotValidComicId = errors.New("not a valid comic id")
)

type ComicsMap map[string]Comics

type Comics struct {
	id              string
	api             *Api
	etag            string
	AssetBundleName string      `json:"assetBundleName"`
	Title           []string    `json:"title"`
	SubTitle        []string    `json:"subTitle"`
	PublicStartAt   []Timestamp `json:"publicStartAt"`
	CharacterId     []int       `json:"characterId"`
}

// InitComics initializes the ComicsMap by fetching the comics data from the API.
//
// It returns a pointer to the ComicsMap and an error if any occurs during the API request or JSON unmarshalling.
func (a *Api) InitComics(ctx context.Context) (*ComicsMap, error) {
	resp, err := a.do(ctx, http.MethodGet, comicsURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var comicsMap ComicsMap
	if err = json.UnmarshalRead(resp.Body, &comicsMap); err != nil {
		return nil, err
	}
	for comicID, comic := range comicsMap {
		comic.id = comicID
		comic.api = a
		comic.etag = resp.Header.Get("ETag")
		comicsMap[comicID] = comic
	}
	return &comicsMap, nil
}

// UpdateComics updates the ComicsMap with the latest data from the API.
//
// It uses the ETag header to check if the data has changed since the last update. If the data has not changed, it does nothing. If the data has changed, it updates the ComicsMap with the new data and initializes the api field for each comic in the map.
func (c *ComicsMap) UpdateComics(ctx context.Context) error {
	var etag string
	var api *Api
	for _, comic := range *c {
		etag = comic.etag
		if etag == "" {
			continue
		}
		api = comic.api
		if api == nil {
			continue
		}
		break
	}
	resp, err := api.update(ctx, http.MethodGet, comicsURL, etag)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotModified {
		return nil
	}
	var comicsMap ComicsMap
	if err = json.UnmarshalRead(resp.Body, &comicsMap); err != nil {
		return err
	}
	for comicID, comic := range comicsMap {
		comic.id = comicID
		comic.api = api
		comic.etag = resp.Header.Get("ETag")
		comicsMap[comicID] = comic
	}
	maps.Copy(*c, comicsMap)
	return nil
}

// GetComic retrieves a Comics struct from the ComicsMap based on the provided comic ID.
//
// If the comic ID is not found in the ComicsMap, it returns an error.
func (c *ComicsMap) GetComic(comicId string) (*Comics, error) {
	if _, err := strconv.Atoi(comicId); err != nil {
		return nil, ErrNotValidComicId
	}
	comic, exists := (*c)[comicId]
	if !exists {
		return nil, ErrComicNotFound
	}
	return &comic, nil
}

type ComicTitleMap struct {
	Id    string
	Title []string
}

// GetComicTitleMap returns a slice of ComicTitleMap containing the comic IDs and their corresponding titles from the ComicsMap.
func (c *ComicsMap) GetComicTitleMap() []ComicTitleMap {
	titles := make([]ComicTitleMap, 0, len(*c))
	for id, comic := range *c {
		titles = append(titles, ComicTitleMap{
			Id:    id,
			Title: comic.Title,
		})
	}
	return titles
}

// IsAvailable checks if the comic is available in the given region.
//
// It is recommended to use this method before calling any other methods that require the comic to be available in the specified region.
func (c *Comics) IsAvailable(region Region) bool {
	if c.PublicStartAt[region.Idx()].Time.IsZero() {
		return false
	}
	return true
}

// GetTitle returns the title of the comic for the given region.
func (c *Comics) GetTitle(region Region) string {
	return c.Title[region.Idx()]
}

// GetSubTitle returns the subtitle of the comic for the given region.
func (c *Comics) GetSubTitle(region Region) string {
	return c.SubTitle[region.Idx()]
}

// GetPublishTime returns the publish time of the comic for the given region.
//
// Be aware that if the time is 1970-01-01 00:00:00.001 UTC, it means that the comic is published with the game together. In this case, it is recommended to not show the PublishTime to the user, as it is not a real publish time.
func (c *Comics) GetPublishTime(region Region) time.Time {
	return c.PublicStartAt[region.Idx()].Time
}

// GetComicURL returns the URL of the comic image for the given region.
func (c *Comics) GetComicURL(region Region) string {
	parts := strings.Split(c.AssetBundleName, "_")
	if parts[1] == "fourframe" {
		return fmt.Sprintf("https://bestdori.com/assets/%s/comic/comic_fourframe/%s_rip/%s.png", region.String(), c.AssetBundleName, c.AssetBundleName)
	}
	return fmt.Sprintf("https://bestdori.com/assets/%s/comic/comic_singleframe/%s_rip/%s.png", region.String(), c.AssetBundleName, c.AssetBundleName)
}

// GetComicFile retrieves the comic image file for the given region.
//
// The format of the comic image is PNG.
//
// It fetches the comic image from the URL returned by GetComicURL and returns the image data as a byte slice. If any error occurs during the API request or reading the response body, it returns an error.
func (c *Comics) GetComicFile(ctx context.Context, region Region) ([]byte, error) {
	url := c.GetComicURL(region)
	resp, err := c.api.do(ctx, http.MethodGet, url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}
