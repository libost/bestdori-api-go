package api

import (
	"context"
	"math"
	"testing"
	"time"
)

func TestGetSongById(t *testing.T) {
	client := NewApi("", 10*time.Second)
	// Initialize songs first
	songs, err := client.InitSongs(context.Background())
	if err != nil {
		t.Fatalf("Failed to initialize songs: %v", err)
	}
	if songs == nil {
		t.Fatal("Expected songs, got nil")
	}
	if err = songs.UpdateSongs(context.Background()); err != nil {
		t.Fatalf("Failed to update songs: %v", err)
	}
	sos, err := songs.GetSong("306")
	if err != nil {
		t.Fatalf("Failed to get song by ID: %v", err)
	}
	if sos == nil {
		t.Fatal("Expected song, got nil")
	}
	if !sos.IsAvailable(RegionCN) {
		t.Fatal("Expected song to be available in CN region")
	}
	if !sos.IsPlayable(RegionCN) {
		t.Fatal("Expected song to be playable in CN region")
	}
	if bandlogo, err := sos.GetBandLogoFile(context.Background()); err == nil || bandlogo != nil {
		t.Fatalf("Expected no band logo file for song, got bandlogo: %v, error: %v", bandlogo, err)
	}
	if sos.GetBandLogoURL() != "" {
		t.Fatalf("Expected no band logo URL for song, got '%s'", sos.GetBandLogoURL())
	}
	if bandname, err := sos.GetBandName(context.Background(), RegionCN); bandname != "RAISE A SUILEN×友希那" || err != nil {
		t.Fatalf("Expected band name 'RAISE A SUILEN×友希那' or error is nil, got '%s', error: %v", bandname, err)
	}
	if sos.GetCloseTime(RegionCN).UnixMilli() != 4133908800000 {
		t.Fatalf("Expected close time 4133908800000 for song, got %d", sos.GetCloseTime(RegionCN).UnixMilli())
	}
	if bannerfile, err := sos.GetCoverImageFile(context.Background(), RegionCN, SongThumbnailTypeJacket); err != nil || bannerfile == nil {
		t.Fatalf("Expected cover image file for song, got bannerfile: %v, error: %v", bannerfile, err)
	}
	if sos.GetCoverImageURL(RegionCN, SongThumbnailTypeJacket) != "https://bestdori.com/assets/cn/musicjacket/musicjacket310_rip/assets-star-forassetbundle-startapp-musicjacket-musicjacket310-306_savior_of_song-jacket.png" {
		t.Fatalf("Expected cover image URL 'https://bestdori.com/assets/cn/musicjacket/musicjacket310_rip/assets-star-forassetbundle-startapp-musicjacket-musicjacket310-306_savior_of_song-jacket.png', got '%s'", sos.GetCoverImageURL(RegionCN, SongThumbnailTypeJacket))
	}
	if lowresfile, err := sos.GetCoverImageFile(context.Background(), RegionCN, SongThumbnailTypeBanner); err != nil || lowresfile == nil {
		t.Fatalf("Expected cover image file for song, got lowresfile: %v, error: %v", lowresfile, err)
	}
	if sos.GetCoverImageURL(RegionCN, SongThumbnailTypeBanner) != "https://bestdori.com/assets/cn/musicjacket/musicjacket310_rip/assets-star-forassetbundle-startapp-musicjacket-musicjacket310-306_savior_of_song-thumb.png" {
		t.Fatalf("Expected cover image URL 'https://bestdori.com/assets/cn/musicjacket/musicjacket310_rip/assets-star-forassetbundle-startapp-musicjacket-musicjacket310-306_savior_of_song-thumb.png', got '%s'", sos.GetCoverImageURL(RegionCN, SongThumbnailTypeBanner))
	}
	if !sos.HasSpecialDifficulty(RegionCN) {
		t.Fatal("Expected song to have special difficulty in CN region")
	}
	if len(sos.GetDifficulties(RegionCN)) != 5 {
		t.Fatalf("Expected 5 difficulties for song, got %d", len(sos.GetDifficulties(RegionCN)))
	}
	var expectedDifficulties = []int{9, 14, 20, 28, 27}
	for i, difficulty := range sos.GetDifficulties(RegionCN) {
		if difficulty != expectedDifficulties[i] {
			t.Fatalf("Expected difficulty %d for song, got %d", expectedDifficulties[i], difficulty)
		}
	}
	if sos.GetPublishTime(RegionCN).UnixMilli() != 1622178000000 {
		t.Fatalf("Expected publish time 1622178000000 for song, got %d", sos.GetPublishTime(RegionCN).UnixMilli())
	}
	if songFile, err := sos.GetSongFile(context.Background(), RegionCN); err != nil || songFile == nil {
		t.Fatalf("Expected song file for song, got songFile: %v, error: %v", songFile, err)
	}
	if sos.GetSongURL(RegionCN) != "https://bestdori.com/assets/cn/sound/bgm306_rip/bgm306.mp3" {
		t.Fatalf("Expected song URL 'https://bestdori.com/assets/cn/sound/bgm306_rip/bgm306.mp3', got '%s'", sos.GetSongURL(RegionCN))
	}
	if sos.GetTitle(RegionCN) != "SAVIOR OF SONG" {
		t.Fatalf("Expected song title 'SAVIOR OF SONG', got '%s'", sos.GetTitle(RegionCN))
	}
	if sos.GetType() != "anime" {
		t.Fatalf("Expected song type 'anime', got '%s'", sos.GetType())
	}
	if chart, err := sos.GetChart(context.Background(), DifficultyExpert); err != nil || chart == nil {
		t.Fatalf("Expected chart for song, got chart: %v, error: %v", chart, err)
	}
	m1, err := sos.GetSongMetaFactors(context.Background(), DifficultyExpert, "7")
	if err != nil {
		t.Fatalf("Failed to get song meta factors: %v", err)
	}
	if m1 == nil {
		t.Fatal("Expected song meta factors, got nil")
	}
	meta, err := client.InitSongMeta(context.Background())
	if err != nil {
		t.Fatalf("Failed to initialize song meta: %v", err)
	}
	if meta == nil {
		t.Fatal("Expected song meta, got nil")
	}
	m, err := meta.GetSongMetaFactors("306", "7", DifficultyExpert)
	if err != nil {
		t.Fatalf("Failed to get song meta factors: %v", err)
	}
	if m == nil {
		t.Fatal("Expected song meta factors, got nil")
	}
	if m.FeverIn() != 2.9675 {
		t.Fatalf("Expected fever in 2.9675, got %f", m.FeverIn())
	}
	if m.FeverOut() != 1.9689 {
		t.Fatalf("Expected fever out 1.9689, got %f", m.FeverOut())
	}
	if m.NoFeverIn() != 2.2602 {
		t.Fatalf("Expected no fever in 2.2602, got %f", m.NoFeverIn())
	}
	if m.NoFeverOut() != 1.6794 {
		t.Fatalf("Expected no fever out 1.6794, got %f", m.NoFeverOut())
	}
	if math.Round(m.ScoreFactor(true, 1.0, 1.1)*100) != 869 {
		t.Fatalf("Expected score factor 869, got %f", m.ScoreFactor(true, 1.0, 1.1))
	}
	if math.Round(m.SkillReliance(true)*100) != 60 {
		t.Fatalf("Expected skill reliance 60, got %f", m.SkillReliance(true))
	}
	sosd, err := sos.GetSong(context.Background())
	if err != nil {
		t.Fatalf("Failed to get song details: %v", err)
	}
	if sosd == nil {
		t.Fatal("Expected song details, got nil")
	}
	if !sosd.IsAvailable(RegionCN) {
		t.Fatal("Expected song details to be available in CN region")
	}
	if !sosd.IsPlayable(RegionCN) {
		t.Fatal("Expected song details to be playable in CN region")
	}
	duration := sosd.GetDuration().Seconds()
	if duration != 90.048 {
		t.Fatalf("Expected duration 90.048, got %f", duration)
	}
	if math.Round(m.Efficiency(true, 1.0, 1.1, duration)*100) != 435 {
		t.Fatalf("Expected efficiency 435, got %f", m.Efficiency(true, 1.0, 1.1, duration))
	}
	if sosd.GetArranger(RegionCN) != "菊田大介（Elements Garden）" {
		t.Fatalf("Expected arranger '菊田大介（Elements Garden）', got '%s'", sosd.GetArranger(RegionCN))
	}
	if bandlogo, err := sosd.GetBandLogoFile(context.Background()); err == nil || bandlogo != nil {
		t.Fatalf("Expected no band logo file for song details, got bandlogo: %v, error: %v", bandlogo, err)
	}
	if sosd.GetBandLogoURL() != "" {
		t.Fatalf("Expected no band logo URL for song details, got '%s'", sosd.GetBandLogoURL())
	}
	if bandname, err := sosd.GetBandName(context.Background(), RegionCN); bandname != "RAISE A SUILEN×友希那" || err != nil {
		t.Fatalf("Expected band name 'RAISE A SUILEN×友希那' or error is nil, got '%s', error: %v", bandname, err)
	}
	if sosd.GetCloseTime(RegionCN).UnixMilli() != 4133908800000 {
		t.Fatalf("Expected close time 4133908800000 for song details, got %d", sosd.GetCloseTime(RegionCN).UnixMilli())
	}
	if sosd.GetCoverImageURL(RegionCN, SongThumbnailTypeJacket) != "https://bestdori.com/assets/cn/musicjacket/musicjacket310_rip/assets-star-forassetbundle-startapp-musicjacket-musicjacket310-306_savior_of_song-jacket.png" {
		t.Fatalf("Expected cover image URL 'https://bestdori.com/assets/cn/musicjacket/musicjacket310_rip/assets-star-forassetbundle-startapp-musicjacket-musicjacket310-306_savior_of_song-jacket.png', got '%s'", sosd.GetCoverImageURL(RegionCN, SongThumbnailTypeJacket))
	}
	if bannerfile, err := sosd.GetCoverImageFile(context.Background(), RegionCN, SongThumbnailTypeBanner); err != nil || bannerfile == nil {
		t.Fatalf("Expected cover image file for song details, got bannerfile: %v, error: %v", bannerfile, err)
	}
	if sosd.GetCoverImageURL(RegionCN, SongThumbnailTypeBanner) != "https://bestdori.com/assets/cn/musicjacket/musicjacket310_rip/assets-star-forassetbundle-startapp-musicjacket-musicjacket310-306_savior_of_song-thumb.png" {
		t.Fatalf("Expected cover image URL 'https://bestdori.com/assets/cn/musicjacket/musicjacket310_rip/assets-star-forassetbundle-startapp-musicjacket-musicjacket310-306_savior_of_song-thumb.png', got '%s'", sosd.GetCoverImageURL(RegionCN, SongThumbnailTypeBanner))
	}
	if lowresfile, err := sosd.GetCoverImageFile(context.Background(), RegionCN, SongThumbnailTypeBanner); err != nil || lowresfile == nil {
		t.Fatalf("Expected cover image file for song details, got lowresfile: %v, error: %v", lowresfile, err)
	}
	if sosd.GetPublishTime(RegionCN).UnixMilli() != 1622178000000 {
		t.Fatalf("Expected publish time 1622178000000 for song details, got %d", sosd.GetPublishTime(RegionCN).UnixMilli())
	}
	if songFile, err := sosd.GetSongFile(context.Background(), RegionCN); err != nil || songFile == nil {
		t.Fatalf("Expected song file for song details, got songFile: %v, error: %v", songFile, err)
	}
	if sosd.GetSongURL(RegionCN) != "https://bestdori.com/assets/cn/sound/bgm306_rip/bgm306.mp3" {
		t.Fatalf("Expected song URL 'https://bestdori.com/assets/cn/sound/bgm306_rip/bgm306.mp3', got '%s'", sosd.GetSongURL(RegionCN))
	}
	if sosd.GetTitle(RegionCN) != "SAVIOR OF SONG" {
		t.Fatalf("Expected song title 'SAVIOR OF SONG', got '%s'", sosd.GetTitle(RegionCN))
	}
	if sosd.GetDescription(RegionCN) != "" {
		t.Fatalf("Expected song description to be empty, got '%s'", sosd.GetDescription(RegionCN))
	}
	if !sosd.HasSpecialDifficulty(RegionCN) {
		t.Fatal("Expected song details to have special difficulty in CN region")
	}
	for i, difficulty := range sosd.GetDifficulties(RegionCN) {
		if difficulty != expectedDifficulties[i] {
			t.Fatalf("Expected difficulty %d for song details, got %d", expectedDifficulties[i], difficulty)
		}
	}
	if sosd.GetDuration().Seconds() != 90.048 {
		t.Fatalf("Expected duration 90.048 for song details, got %f", sosd.GetDuration().Seconds())
	}
	if sosd.GetHowToGet(RegionCN) != "在CiRCLE的乐曲商店中交换" {
		t.Fatalf("Expected how to get '在CiRCLE的乐曲商店中交换', got '%s'", sosd.GetHowToGet(RegionCN))
	}
	if sosd.GetLyricist(RegionCN) != "ナノ" {
		t.Fatalf("Expected lyricist 'ナノ', got '%s'", sosd.GetLyricist(RegionCN))
	}
	if len(sosd.GetNotesCounts(RegionCN)) != 5 {
		t.Fatalf("Expected 5 notes counts for song details, got %d", len(sosd.GetNotesCounts(RegionCN)))
	}
	var expectedNotesCounts = []int{133, 246, 401, 772, 727}
	for i, notesCount := range sosd.GetNotesCounts(RegionCN) {
		if notesCount != expectedNotesCounts[i] {
			t.Fatalf("Expected notes count %d for song details, got %d", expectedNotesCounts[i], notesCount)
		}
	}
	if sosd.GetPublishTime(RegionCN).UnixMilli() != 1622178000000 {
		t.Fatalf("Expected publish time 1622178000000 for song details, got %d", sosd.GetPublishTime(RegionCN).UnixMilli())
	}
	if songFile, err := sosd.GetSongFile(context.Background(), RegionCN); err != nil || songFile == nil {
		t.Fatalf("Expected song file for song details, got songFile: %v, error: %v", songFile, err)
	}
	if sosd.GetSongURL(RegionCN) != "https://bestdori.com/assets/cn/sound/bgm306_rip/bgm306.mp3" {
		t.Fatalf("Expected song URL 'https://bestdori.com/assets/cn/sound/bgm306_rip/bgm306.mp3', got '%s'", sosd.GetSongURL(RegionCN))
	}
	if m, err := sosd.GetSongMetaFactors(context.Background(), DifficultyExpert, "7"); err != nil || m == nil {
		t.Fatalf("Expected song meta factors for song details, got metaFactors: %v, error: %v", m, err)
	}
	if sosd.GetTitle(RegionCN) != "SAVIOR OF SONG" {
		t.Fatalf("Expected song title 'SAVIOR OF SONG', got '%s'", sosd.GetTitle(RegionCN))
	}
	if sosd.GetType() != "anime" {
		t.Fatalf("Expected song type 'anime', got '%s'", sosd.GetType())
	}
	if sosd.GetComposer(RegionCN) != "MY FIRST STORY" {
		t.Fatalf("Expected composer 'MY FIRST STORY', got '%s'", sosd.GetComposer(RegionCN))
	}
	if chart, err := sosd.GetChart(context.Background(), DifficultyExpert); err != nil || chart == nil {
		t.Fatalf("Expected chart for song details, got chart: %v, error: %v", chart, err)
	}
	sia, err := client.GetSong(context.Background(), "245")
	if err != nil {
		t.Fatalf("Failed to get song by ID: %v", err)
	}
	if sia == nil {
		t.Fatal("Expected song, got nil")
	}
	if sia.GetDescription(RegionCN) != "「Noble・Rose -歌至顶点-」活动追加歌曲" {
		t.Fatalf("Expected song description '「Noble・Rose -歌至顶点-」活动追加歌曲', got '%s'", sia.GetDescription(RegionCN))
	}
}
