package api

import (
	"context"
	"encoding/json/v2"
	"fmt"
	"io"
	"maps"
	"net/http"
	"strconv"
	"time"
)

type Songs struct {
	id   string
	api  *Api
	etag string
	// `normal` for og, `anime` for covers and `tie-up` for extras
	Tag         string                         `json:"tag"`
	BandID      int                            `json:"bandId"`
	JacketImage []string                       `json:"jacketImage"`
	MusicTitle  []string                       `json:"musicTitle"`
	PublishedAt []Timestamp                    `json:"publishedAt"`
	ClosedAt    []Timestamp                    `json:"closedAt"`
	Difficulty  map[string]SongsDifficultyData `json:"difficulty"`
}

type SongsDifficultyData struct {
	PlayLevel   int         `json:"playLevel"`
	PublishedAt []Timestamp `json:"publishedAt"`
}

type SongsMap map[string]Songs

type Song struct {
	id           string
	api          *Api
	BgmID        string                    `json:"bgmId"`
	BgmFile      string                    `json:"bgmFile"`
	Tag          string                    `json:"tag"`
	BandID       int                       `json:"bandId"`
	Achievements []SongAchievementsData    `json:"achievements"`
	JacketImage  []string                  `json:"jacketImage"`
	Seq          int                       `json:"seq"`
	MusicTitle   []string                  `json:"musicTitle"`
	Ruby         []string                  `json:"ruby"`
	Phonetic     []string                  `json:"phonetic"`
	Lyricist     []string                  `json:"lyricist"`
	Composer     []string                  `json:"composer"`
	Arranger     []string                  `json:"arranger"`
	HowToGet     []string                  `json:"howToGet"`
	PublishedAt  []Timestamp               `json:"publishedAt"`
	ClosedAt     []Timestamp               `json:"closedAt"`
	Description  []string                  `json:"description"`
	Difficulty   map[string]SongDifficulty `json:"difficulty"`
	Length       Timeduration              `json:"length"`
	Notes        map[string]int            `json:"notes"`
	Bpm          map[string]SongBpm        `json:"bpm"`
}

type SongAchievementsData struct {
	MusicID         int    `json:"musicId"`
	AchievementType string `json:"achievementType"`
	RewardType      string `json:"rewardType"`
	Quantity        int    `json:"quantity"`
}

type SongDifficulty struct {
	PlayLevel         int                                        `json:"playLevel"`
	MultiLiveScoreMap map[string]SongDifficultyMultiLiveScoreMap `json:"multiLiveScoreMap"`
	NotesQuantity     int                                        `json:"notesQuantity"`
	ScoreC            int                                        `json:"scoreC"`
	ScoreB            int                                        `json:"scoreB"`
	ScoreA            int                                        `json:"scoreA"`
	ScoreS            int                                        `json:"scoreS"`
	ScoreSS           int                                        `json:"scoreSS"`
	PublishedAt       []Timestamp                                `json:"publishedAt"`
}

type SongDifficultyMultiLiveScoreMap struct {
	MusicID                 int    `json:"musicId"`
	MusicDifficulty         string `json:"musicDifficulty"`
	MusicLiveDifficultyID   int    `json:"musicLiveDifficultyId"`
	ScoreS                  int    `json:"scoreS"`
	ScoreA                  int    `json:"scoreA"`
	ScoreB                  int    `json:"scoreB"`
	ScoreC                  int    `json:"scoreC"`
	MultiLiveDifficultyType string `json:"multiLiveDifficultyType"`
	ScoreSS                 int    `json:"scoreSS"`
	ScoreSSS                int    `json:"scoreSSS"`
}

type SongBpm []struct {
	Bpm   int     `json:"bpm"`
	Start float64 `json:"start"`
	End   float64 `json:"end"`
}

// InitSongs initializes the SongsMap by fetching the song data from the API. It returns a pointer to the SongsMap and an error if any occurs during the API request or JSON unmarshalling.
func (a *Api) InitSongs(ctx context.Context) (*SongsMap, error) {
	resp, err := a.do(ctx, http.MethodGet, songsURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var songs SongsMap
	if err = json.UnmarshalRead(resp.Body, &songs); err != nil {
		return nil, err
	}
	for id, song := range songs {
		song.id = id
		song.api = a
		song.etag = resp.Header.Get("ETag")
		songs[id] = song
	}
	return &songs, nil
}

// UpdateSongs updates the SongsMap with the latest data from the API. It uses the ETag header to check if the data has changed since the last update. If the data has not changed, it does nothing. If the data has changed, it updates the SongsMap with the new data and initializes the api field for each song in the map.
func (s *SongsMap) UpdateSongs(ctx context.Context) error {
	var etag string
	var api *Api
	for _, song := range *s {
		etag = song.etag
		if etag == "" {
			continue
		}
		api = song.api
		if api == nil {
			continue
		}
		break
	}
	resp, err := api.update(ctx, http.MethodGet, songsURL, etag)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotModified {
		return nil
	}
	var songs SongsMap
	if err = json.UnmarshalRead(resp.Body, &songs); err != nil {
		return err
	}
	for id, song := range songs {
		song.id = id
		song.api = api
		song.etag = resp.Header.Get("ETag")
		songs[id] = song
	}
	maps.Copy(*s, songs)
	return nil
}

func getSongs(ctx context.Context, api *Api, songId string) (*Songs, error) {
	resp, err := api.do(ctx, http.MethodGet, songsURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var songs SongsMap
	if err = json.UnmarshalRead(resp.Body, &songs); err != nil {
		return nil, err
	}
	song, ok := songs[songId]
	if !ok {
		return nil, fmt.Errorf("song %q not found", songId)
	}
	song.api = api
	song.id = songId
	return &song, nil
}

// GetSongs retrieves the Songs struct for the RecentSongs struct. It returns an error if the song ID is not found in the SongsMap.
//
// Be aware that the Songs struct that is returned is a summary of the song data, and does not include all the details that can be found in the Song struct.
func (r *RecentSongs) GetSongs(ctx context.Context) (*Songs, error) {
	return getSongs(ctx, r.api, r.id)
}

// GetSongs retrieves the Song struct for the SongsMap struct. It returns an error if the song ID is not found in the SongsMap.
func (r *RecentSongs) GetSong(ctx context.Context) (*Song, error) {
	return r.api.GetSong(ctx, r.id)
}

// GetSongs retrieves the Song struct for the given song ID. It returns an error if the song ID is not found in the SongsMap.
//
// Be aware that the Songs struct that is returned is a summary of the song data, and does not include all the details that can be found in the Song struct.
func (s *SongsMap) GetSong(songId string) (*Songs, error) {
	song, ok := (*s)[songId]
	if !ok {
		return nil, fmt.Errorf("song %q not found", songId)
	}
	song.api = (*s)[songId].api
	song.id = songId
	return &song, nil
}

// GetSongs retrieves the Song struct for the given song ID. It returns an error if the song ID is not found in the SongsMap.
func (a *Api) GetSong(ctx context.Context, songId string) (*Song, error) {
	url := "https://bestdori.com/api/songs/" + songId + ".json"
	resp, err := a.do(ctx, http.MethodGet, url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var song Song
	if err = json.UnmarshalRead(resp.Body, &song); err != nil {
		return nil, err
	}
	song.id = songId
	song.api = a
	return &song, nil
}

// GetSongs retrieves the Song struct for the Songs struct. It returns an error if any error occurs.
func (s *Songs) GetSong(ctx context.Context) (*Song, error) {
	return s.api.GetSong(ctx, s.id)
}

// GetType returns the type of the song, which can be "normal" for original songs, "anime" for cover songs, or "tie-up" for extra songs.
func (s *Songs) GetType() string {
	return s.Tag
}

// GetType returns the type of the song, which can be "normal" for original songs, "anime" for cover songs, or "tie-up" for extra songs.
func (s *Song) GetType() string {
	return s.Tag
}

// IsAvailable checks if the song is available in the given region. It returns true if the song has been published in the specified region.
//
// Be aware that even if a song is available, it may not be playable due to being closed in that region. Use IsPlayable to check if the song is both available and playable.
//
// This function is useful for determining if a song can be accessed in a specific region, but it does not guarantee that the song can be played. For checking playability, refer to the IsPlayable method.
func (s *Songs) IsAvailable(region Region) bool {
	if s.PublishedAt[region.Idx()].Time.IsZero() {
		return false
	}
	return true
}

// IsAvailable checks if the song is available in the given region. It returns true if the song has been published in the specified region.
//
// Be aware that even if a song is available, it may not be playable due to being closed in that region. Use IsPlayable to check if the song is both available and playable.
//
// This function is useful for determining if a song can be accessed in a specific region, but it does not guarantee that the song can be played. For checking playability, refer to the IsPlayable method.
func (s *Song) IsAvailable(region Region) bool {
	if s.PublishedAt[region.Idx()].Time.IsZero() {
		return false
	}
	return true
}

// IsPlayable checks if the song is playable in the given region. It returns true if the song is available and has not been closed in that region.
func (s *Songs) IsPlayable(region Region) bool {
	if !s.IsAvailable(region) {
		return false
	}
	now := time.Now()
	if !s.ClosedAt[region.Idx()].Time.IsZero() && s.ClosedAt[region.Idx()].Time.Before(now) {
		return false
	}
	return true
}

// IsPlayable checks if the song is playable in the given region. It returns true if the song is available and has not been closed in that region.
func (s *Song) IsPlayable(region Region) bool {
	if !s.IsAvailable(region) {
		return false
	}
	now := time.Now()
	if !s.ClosedAt[region.Idx()].Time.IsZero() && s.ClosedAt[region.Idx()].Time.Before(now) {
		return false
	}
	return true
}

// GetTitle returns the title of the song for the given region.
func (s *Songs) GetTitle(region Region) string {
	return s.MusicTitle[region.Idx()]
}

// GetTitle returns the title of the song for the given region.
func (s *Song) GetTitle(region Region) string {
	return s.MusicTitle[region.Idx()]
}

// GetPublishTime returns the publish time of the song for the given region.
func (s *Songs) GetPublishTime(region Region) time.Time {
	return s.PublishedAt[region.Idx()].Time
}

// GetPublishTime returns the publish time of the song for the given region.
func (s *Song) GetPublishTime(region Region) time.Time {
	return s.PublishedAt[region.Idx()].Time
}

// GetCloseTime returns the close time of the song for the given region.
//
// It is not recommended to show the CloseTime to the user if the time is bigger than common sense.
func (s *Songs) GetCloseTime(region Region) time.Time {
	return s.ClosedAt[region.Idx()].Time
}

// GetCloseTime returns the close time of the song for the given region.
//
// It is not recommended to show the CloseTime to the user if the time is bigger than common sense.
func (s *Song) GetCloseTime(region Region) time.Time {
	return s.ClosedAt[region.Idx()].Time
}

// GetCoverImageURL returns the URL of the song cover images for the specified region.
func (s *Songs) GetCoverImageURL(region Region, thumbnailType SongThumbnailType) string {
	return getJacket(s.id, s.JacketImage[0], region.String(), thumbnailType)
}

// GetCoverImageURL returns the URL of the song cover images for the specified region.
func (s *Song) GetCoverImageURL(region Region, thumbnailType SongThumbnailType) string {
	return getJacket(s.id, s.JacketImage[0], region.String(), thumbnailType)
}

// GetCoverImageFile retrieves the song cover image file for the specified region and thumbnail type. It returns the image data as a byte slice and an error if any occurs during the API request or reading the response body.
//
// The format of the song cover image is PNG.
func (s *Songs) GetCoverImageFile(ctx context.Context, region Region, thumbnailType SongThumbnailType) ([]byte, error) {
	url := s.GetCoverImageURL(region, thumbnailType)
	if url == "" {
		return nil, fmt.Errorf("invalid cover image URL for song %q in region %q", s.id, region.String())
	}
	resp, err := s.api.do(ctx, http.MethodGet, url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

// GetCoverImageFile retrieves the song cover image file for the specified region and thumbnail type. It returns the image data as a byte slice and an error if any occurs during the API request or reading the response body.
//
// The format of the song cover image is PNG.
func (s *Song) GetCoverImageFile(ctx context.Context, region Region, thumbnailType SongThumbnailType) ([]byte, error) {
	url := s.GetCoverImageURL(region, thumbnailType)
	if url == "" {
		return nil, fmt.Errorf("invalid cover image URL for song %q in region %q", s.id, region.String())
	}
	resp, err := s.api.do(ctx, http.MethodGet, url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func getJacket(songId, musicJacket, region string, thumbnailType SongThumbnailType) string {
	songIdInt, err := strconv.Atoi(songId)
	if err != nil {
		return ""
	}
	musicJacketPack := ceilTo10(songIdInt)
	jacketURL := fmt.Sprintf("https://bestdori.com/assets/%s/musicjacket/musicjacket%d_rip/assets-star-forassetbundle-startapp-musicjacket-musicjacket%d-%s-jacket.png", region, musicJacketPack, musicJacketPack, musicJacket)
	thumbURL := fmt.Sprintf("https://bestdori.com/assets/%s/musicjacket/musicjacket%d_rip/assets-star-forassetbundle-startapp-musicjacket-musicjacket%d-%s-thumb.png", region, musicJacketPack, musicJacketPack, musicJacket)
	if thumbnailType == SongThumbnailTypeJacket {
		return jacketURL
	}
	if thumbnailType == SongThumbnailTypeBanner {
		return thumbURL
	}
	return ""
}

func ceilTo10(n int) int {
	return (n + 9) / 10 * 10
}

// HasSpecialDifficulty checks if the song has a special difficulty (difficulty 4) available in the given region.
func (s *Songs) HasSpecialDifficulty(region Region) bool {
	if len(s.Difficulty) < 5 {
		return false
	}
	if s.Difficulty["4"].PublishedAt[region.Idx()].Time.IsZero() {
		return false
	}
	return true
}

// HasSpecialDifficulty checks if the song has a special difficulty (difficulty 4) available in the given region.
func (s *Song) HasSpecialDifficulty(region Region) bool {
	if len(s.Difficulty) < 5 {
		return false
	}
	if s.Difficulty["4"].PublishedAt[region.Idx()].Time.IsZero() {
		return false
	}
	return true
}

// GetDifficulties retrieves the play levels of the available difficulties for the song in the given region. It returns a slice of integers representing the play levels of the difficulties.
func (s *Songs) GetDifficulties(region Region) []int {
	var difficulties []int
	idx := []string{"0", "1", "2", "3"}
	if s.HasSpecialDifficulty(region) {
		idx = append(idx, "4")
	}
	for _, i := range idx {
		difficulty, ok := s.Difficulty[i]
		if !ok {
			return nil
		}
		difficulties = append(difficulties, difficulty.PlayLevel)
	}
	return difficulties
}

// GetDifficulties retrieves the play levels of the available difficulties for the song in the given region. It returns a slice of integers representing the play levels of the difficulties.
func (s *Song) GetDifficulties(region Region) []int {
	var difficulties []int
	idx := []string{"0", "1", "2", "3"}
	if s.HasSpecialDifficulty(region) {
		idx = append(idx, "4")
	}
	for _, i := range idx {
		difficulty, ok := s.Difficulty[i]
		if !ok {
			return nil
		}
		difficulties = append(difficulties, difficulty.PlayLevel)
	}
	return difficulties
}

// GetDuration returns the duration of the song as a time.Duration value.
func (s *Song) GetDuration() time.Duration {
	return s.Length.Duration
}

// GetSongURL returns the URL of the song's audio file for the specified region. The URL points to the MP3 file of the song.
func (s *Song) GetSongURL(region Region) string {
	return fmt.Sprintf("https://bestdori.com/assets/%s/sound/bgm%s_rip/bgm%s.mp3", region.String(), s.id, s.id)
}

// GetSongURL returns the URL of the song's audio file for the specified region. The URL points to the MP3 file of the song.
func (s *Songs) GetSongURL(region Region) string {
	return fmt.Sprintf("https://bestdori.com/assets/%s/sound/bgm%s_rip/bgm%s.mp3", region.String(), s.id, s.id)
}

// GetSongFile retrieves the song's audio file for the specified region. It returns the audio data as a byte slice and an error if any occurs during the API request or reading the response body.
//
// The format of the song's audio file is MP3.
func (s *Song) GetSongFile(ctx context.Context, region Region) ([]byte, error) {
	url := s.GetSongURL(region)
	return getSongFile(ctx, s.api, url)
}

// GetSongFile retrieves the song's audio file for the specified region. It returns the audio data as a byte slice and an error if any occurs during the API request or reading the response body.
//
// The format of the song's audio file is MP3.
func (s *Songs) GetSongFile(ctx context.Context, region Region) ([]byte, error) {
	url := s.GetSongURL(region)
	return getSongFile(ctx, s.api, url)
}

func getSongFile(ctx context.Context, api *Api, url string) ([]byte, error) {
	resp, err := api.do(ctx, http.MethodGet, url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

// GetNotesCounts retrieves the counts of notes for each difficulty level of the song in the given region. It returns a slice of integers representing the counts of notes for each difficulty level.
func (s *Song) GetNotesCounts(region Region) []int {
	idx := []string{"0", "1", "2", "3"}
	if s.HasSpecialDifficulty(region) {
		idx = append(idx, "4")
	}
	counts := make([]int, 0, len(idx))
	for _, i := range idx {
		counts = append(counts, s.Notes[i])
	}
	return counts
}

// GetLyricist returns the lyricist's name localized for the given region.
func (s *Song) GetLyricist(region Region) string {
	return s.Lyricist[region.Idx()]
}

// GetComposer returns the composer's name localized for the given region.
func (s *Song) GetComposer(region Region) string {
	return s.Composer[region.Idx()]
}

// GetArranger returns the arranger's name localized for the given region.
func (s *Song) GetArranger(region Region) string {
	return s.Arranger[region.Idx()]
}

// GetHowToGet returns the information on how to obtain the song localized for the given region.
func (s *Song) GetHowToGet(region Region) string {
	return s.HowToGet[region.Idx()]
}

// GetDescription returns the description of the song localized for the given region.
func (s *Song) GetDescription(region Region) string {
	if len(s.Description) <= region.Idx() {
		return ""
	}
	if s.Description[region.Idx()] == "" {
		return ""
	}
	return s.Description[region.Idx()]
}

func (s *Song) GetBPM() int {
	for _, bpm := range s.Bpm {
		if bpm[0].Bpm != 0 {
			return bpm[0].Bpm
		}
	}
	return 0
}

type SongChart []SongChartSingle

type SongChartSingle struct {
	// Type can be "BPM"(gives the BPM), "System"(gives the Data, which is the song's file), "Long"(gives the Connections, which is the long note's connections), "Single"(gives nothing), "Slide"(gives the Connections, which is the slide note's connections)
	Type string `json:"type"`
	// Bpm is the song's BPM, only present if Type is "BPM"
	Bpm int `json:"bpm"`
	// Data is the song's file or fever sound effect, only present if Type is "System"
	Data string `json:"data"`
	// Lane is the lane of the note, only present if Type is "Single"
	Lane int `json:"lane"`
	// Connections is how the long note or slide note is connected, only present if Type is "Long" or "Slide"
	Connections []SongChartConnection `json:"connections"`
	// Beat is the beat of the note, only present if Type is "Single"
	Beat float64 `json:"beat"`
	// Flick is whether the note is a flick note, only present if Type is "Single"
	Flick bool `json:"flick"`
	// Skill is whether the note is a skill note, only present if Type is "Single"
	Skill bool `json:"skill"`
}

type SongChartConnection struct {
	// Lane is the lane of the note
	Lane int `json:"lane"`
	// Beat is the beat of the note
	Beat float64 `json:"beat"`
	// Flick is whether the note is a flick note, it is only present if true, and is omitted if false
	Flick bool `json:"flick"`
}

// GetChart retrieves the chart data for the song at the specified difficulty level. It returns a SongChart, which is a slice of SongChartSingle, and an error if any occurs during the API request or JSON unmarshalling.
//
// Always check if the song has special difficulty before calling this function with DifficultySpecial, as it may not be available for all songs.
func (s *Song) GetChart(ctx context.Context, difficulty Difficulty) (SongChart, error) {
	return getChart(ctx, s.api, s.id, difficulty)
}

// GetChart retrieves the chart data for the songs at the specified difficulty level. It returns a SongChart, which is a slice of SongChartSingle, and an error if any occurs during the API request or JSON unmarshalling.
//
// Always check if the songs have special difficulty before calling this function with DifficultySpecial, as it may not be available for all songs.
func (s *Songs) GetChart(ctx context.Context, difficulty Difficulty) (SongChart, error) {
	return getChart(ctx, s.api, s.id, difficulty)
}

func getChart(ctx context.Context, api *Api, songId string, difficulty Difficulty) (SongChart, error) {
	url := fmt.Sprintf("https://bestdori.com/api/charts/%s/%s.json", songId, difficulty.String())
	resp, err := api.do(ctx, http.MethodGet, url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var chart SongChart
	if err = json.UnmarshalRead(resp.Body, &chart); err != nil {
		return nil, err
	}
	return chart, nil
}
