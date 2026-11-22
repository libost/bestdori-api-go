package api

import (
	"context"
	"encoding/json/v2"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Skills struct {
	api               *Api
	SimpleDescription []string  `json:"simpleDescription"`
	Description       []string  `json:"description"`
	Duration          []float64 `json:"duration"`
	ActivationEffect  struct {
		ActivateEffectTypes struct {
			Score SkillActivationEffectTypesValue `json:"score"`
			Judge SkillActivationEffectTypesValue `json:"judge"`
		}
	} `json:"activationEffect"`
	OnceEffect OnceEffect `json:"onceEffect"`
}

type SkillActivationEffectTypesValue struct {
	ActivateEffectValue     []int  `json:"activateEffectValue"`
	ActivateEffectValueType string `json:"activateEffectValueType"`
	ActivateCondition       string `json:"activateCondition"`
}

type OnceEffect struct {
	OnceEffectType      string `json:"onceEffectType"`
	OnceEffectValueType string `json:"onceEffectValueType"`
	OnceEffectValue     []int  `json:"onceEffectValue"`
}

type SkillsMap map[string]Skills

// GetSkill retrieves the Skills struct for the Cards struct. It returns an error if the skill ID is not found in the SkillsMap.
func (c *Cards) GetSkill(ctx context.Context) (*Skills, error) {
	return getSkill(ctx, c.api, strconv.Itoa(c.SkillID))
}

// GetSkill retrieves the Skills struct for the Card struct. It returns an error if the skill ID is not found in the SkillsMap.
func (c *Card) GetSkill(ctx context.Context) (*Skills, error) {
	return getSkill(ctx, c.api, strconv.Itoa(c.SkillID))
}

func getSkill(ctx context.Context, api *Api, skillId string) (*Skills, error) {
	resp, err := api.do(ctx, http.MethodGet, skillsURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var skills SkillsMap
	if err = json.UnmarshalRead(resp.Body, &skills); err != nil {
		return nil, err
	}
	skill, ok := skills[skillId]
	if !ok {
		return nil, fmt.Errorf("skill %q not found", skillId)
	}
	skill.api = api
	return &skill, nil
}

// IsAvailable checks if the skill is available in the given region. It returns true if the skill has a non-empty description for the specified region, indicating that it is available.
//
// It is recommended to use this method before calling GetDescription or GetSimpleDescription to ensure that the skill is available in the desired region.
func (s *Skills) IsAvailable(region Region) bool {
	return s.GetDescription(region) != ""
}

// GetDescription returns the description of the skill for the given region. If the skill is not available in the specified region, it returns an empty string.
//
// Be aware that the string returned may contain placeholders (e.g., "{0}") that need to be replaced with actual values, such as the skill duration. Use GetDescriptionFull to get the full description with placeholders replaced.
func (s *Skills) GetDescription(region Region) string {
	return s.Description[region.Idx()]
}

// GetSimpleDescription returns the simple description of the skill for the given region. If the skill is not available in the specified region, it returns an empty string.
func (s *Skills) GetSimpleDescription(region Region) string {
	return s.SimpleDescription[region.Idx()]
}

// GetDuration returns the duration values of the skill for each level. The returned slice contains the duration values in seconds, corresponding to each level of the skill.
func (s *Skills) GetDuration() []float64 {
	return s.Duration
}

// GetDescriptionFull returns the full description of the skill, replacing the {0} placeholder with the duration value for the given level.
//
// The level parameter should be between 1 and the length of the Duration slice. If the level is out of range, it will return the description without replacing the placeholder.
func (s *Skills) GetDescriptionFull(region Region, level int) string {
	desc := s.GetDescription(region)
	durations := s.GetDuration()
	duration := strings.TrimSuffix(fmt.Sprintf("%.1f", durations[level-1]), ".0")
	onceEffect := s.GetOnceEffect()
	if onceEffect.OnceEffectValue != nil && len(onceEffect.OnceEffectValue) >= level {
		desc = strings.ReplaceAll(desc, "{0}", strconv.Itoa(onceEffect.OnceEffectValue[level-1]))
		desc = strings.ReplaceAll(desc, "{1}", duration)
	} else {
		desc = strings.ReplaceAll(desc, "{0}", duration)
	}
	return desc
}

// GetOnceEffect returns the OnceEffect struct for the skill, which contains information about the one-time effect of the skill, including its type, value type, and values.
//
// Be aware that the OnceEffect may not be applicable to all skills, and its values may vary depending on the skill's configuration.
func (s *Skills) GetOnceEffect() OnceEffect {
	return s.OnceEffect
}
