package api

import (
	"context"
	"encoding/json/v2"
	"errors"
	"fmt"
	"net/http"
)

var (
	// ErrSongNotFound is returned when the song ID is not found in the song meta.
	ErrSongNotFound = errors.New("song not found")
	// ErrDifficultyNotFound is returned when the difficulty is not found in the song meta.
	ErrDifficultyNotFound = errors.New("difficulty not found")
	// ErrSkillLengthNotFound is returned when the skill length is not found in the song meta.
	ErrSkillLengthNotFound = errors.New("skill length not found")
)

type SongMeta map[string]map[string]map[string]SongMetaFactors

type SongMetaFactors [4]float64

// NoFeverOut gives you the value of a, which is the factor for no Fever and outside of skill.
func (m SongMetaFactors) NoFeverOut() float64 { return m[0] }

// NoFeverIn gives you the value of b, which is the factor for no Fever and inside of skill.
func (m SongMetaFactors) NoFeverIn() float64 { return m[1] }

// FeverOut gives you the value of c, which is the factor for Fever and outside of skill.
func (m SongMetaFactors) FeverOut() float64 { return m[2] }

// FeverIn gives you the value of d, which is the factor for Fever and inside of skill.
func (m SongMetaFactors) FeverIn() float64 { return m[3] }

// ScoreFactor calculates the score factor based on whether Fever is active, the skill boost(default: 1.0), and the accuracy(default: 1.1).
func (m SongMetaFactors) ScoreFactor(fever bool, skillBoost, accuracy float64) float64 {
	if fever {
		return accuracy * (m.FeverOut() + m.FeverIn()*(1+skillBoost))
	}
	return accuracy * (m.NoFeverOut() + m.NoFeverIn()*(1+skillBoost))
}

// SkillReliance calculates the skill reliance based on whether Fever is active. It returns a value between 0 and 1, where 0 means no reliance on skill and 1 means full reliance on skill.
func (m SongMetaFactors) SkillReliance(fever bool) float64 {
	if fever {
		return m.FeverIn() / (m.FeverOut() + m.FeverIn())
	}
	return m.NoFeverIn() / (m.NoFeverOut() + m.NoFeverIn())
}

// Efficiency calculates the efficiency based on whether Fever is active, the skill boost(default: 1.0), the accuracy(default: 1.1), and the length of the song in seconds.
//
// length is the length of the song in seconds, and the formula is: efficiency = score factor * 60 / (length + 30). The +30 is a constant that represents the time it takes to start and end the song, which is not affected by Fever or skill boost.
func (m SongMetaFactors) Efficiency(fever bool, skillBoost, accuracy, length float64) float64 {
	score := m.ScoreFactor(fever, skillBoost, accuracy)
	eff := score * 60 / (length + 30)
	return eff
}

func (a *Api) InitSongMeta(ctx context.Context) (*SongMeta, error) {
	resp, err := a.do(ctx, http.MethodGet, songsMetaURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var songMeta SongMeta
	if err = json.UnmarshalRead(resp.Body, &songMeta); err != nil {
		return nil, err
	}
	return &songMeta, nil
}

// GetSongMetaFactors retrieves the SongMetaFactors for a given song ID, difficulty, and skill length. It returns an error if the song ID, difficulty, or skill length is not found in the SongMeta.
//
// songId is the ID of the song,
//
// difficulty is the difficulty level
//
// skillLength is the length of the skill in seconds (pass "7" if you have no idea what it is).
func (m *SongMeta) GetSongMetaFactors(songId, skillLength string, difficulty Difficulty) (*SongMetaFactors, error) {
	if _, ok := (*m)[songId]; !ok {
		return nil, ErrSongNotFound
	}
	diffint := fmt.Sprintf("%d", difficulty.Idx())
	if _, ok := (*m)[songId][diffint]; !ok {
		return nil, ErrDifficultyNotFound
	}
	factors, ok := (*m)[songId][diffint][skillLength]
	if !ok {
		return nil, ErrSkillLengthNotFound
	}
	return &factors, nil
}

// GetSongMetaFactors retrieves the SongMetaFactors for a given difficulty, and skill length. It returns an error if the difficulty, or skill length is not found in the SongMeta.
//
// difficulty is the difficulty level
//
// skillLength is the length of the skill in seconds (pass "7" if you have no idea what it is).
func (s *Song) GetSongMetaFactors(ctx context.Context, difficulty Difficulty, skillLength string) (*SongMetaFactors, error) {
	songMeta, err := s.api.InitSongMeta(ctx)
	if err != nil {
		return nil, err
	}
	return songMeta.GetSongMetaFactors(s.id, skillLength, difficulty)
}

// GetSongMetaFactors retrieves the SongMetaFactors for a given difficulty, and skill length. It returns an error if the difficulty, or skill length is not found in the SongMeta.
//
// difficulty is the difficulty level
//
// skillLength is the length of the skill in seconds (pass "7" if you have no idea what it is).
func (s *Songs) GetSongMetaFactors(ctx context.Context, difficulty Difficulty, skillLength string) (*SongMetaFactors, error) {
	songMeta, err := s.api.InitSongMeta(ctx)
	if err != nil {
		return nil, err
	}
	return songMeta.GetSongMetaFactors(s.id, skillLength, difficulty)
}
