// Copyright 2026 Jakkaphat Chalermphanap
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package store

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/gxjakkap/reception/enum"
	"github.com/gxjakkap/reception/models"
	"github.com/gxjakkap/reception/utils"
	"gorm.io/gorm"
)

type GuildsStore struct {
	db *gorm.DB
}

func NewGuildsStore(db *gorm.DB) *GuildsStore {
	return &GuildsStore{
		db: db,
	}
}

func (s *GuildsStore) Create(g *models.Guilds) error {
	if err := s.db.Create(g).Error; err != nil {
		return err
	}
	return nil
}

func (s *GuildsStore) Delete(gid string) error {
	if err := s.db.Where("guild_id = ?", gid).Delete(&models.Guilds{}).Error; err != nil {
		return err
	}
	return nil
}

func (s *GuildsStore) AddReactionRole(gid string, rr *models.ReactionRoles) error {
	if err := s.db.Model(&models.Guilds{ID: gid}).Association("ReactionRoles").Append(rr); err != nil {
		return err
	}
	return nil
}

func (s *GuildsStore) GetSettings(gid string) (models.GuildSettings, error) {
	var g models.Guilds
	err := s.db.Where("guild_id = ?", gid).First(&g).Error
	if err != nil {
		return models.GuildSettings{}, err
	}

	var sett models.GuildSettings
	err = json.Unmarshal(g.Settings, &sett)

	if err != nil {
		return models.GuildSettings{}, err
	}

	return sett, nil
}

func (s *GuildsStore) GetWelcomeConfig(gid string) (models.WelcomeConfig, error) {
	gs, err := s.GetSettings(gid)

	if err != nil {
		return models.WelcomeConfig{}, err
	}

	return gs.Welcome, nil
}

func (s *GuildsStore) ToggleWelcomeMessage(gid string) (bool, error) {
	var g models.Guilds
	if err := s.db.Where("guild_id = ?", gid).First(&g).Error; err != nil {
		return false, err
	}

	var sett models.GuildSettings
	if err := json.Unmarshal(g.Settings, &sett); err != nil {
		sett = models.GuildSettings{Features: map[string]bool{}}
	}

	if sett.Features == nil {
		sett.Features = make(map[string]bool)
	}

	key := enum.GuildSetting[enum.WelcomeMessage]

	nv := !sett.Features[key]
	sett.Features[key] = nv

	sett.Welcome.Enabled = nv

	gsb, err := json.Marshal(sett)
	if err != nil {
		return false, err
	}

	if err := s.db.Model(&models.Guilds{ID: gid}).Update("settings", gsb).Error; err != nil {
		return false, err
	}

	return nv, nil
}

func (s *GuildsStore) GetWelcomeTextColor(gid string) (string, error) {
	gs, err := s.GetSettings(gid)

	if err != nil {
		return "", err
	}

	return gs.Welcome.TextColor, nil
}

func (s *GuildsStore) SetWelcomeTextColor(gid string, color string) error {
	pat := `/^#?([0-9a-f]{3}|[0-9a-f]{6})$/i`
	re := regexp.MustCompile(pat)

	if !re.MatchString(color) {
		return fmt.Errorf("invalid hex color string")
	}

	gs, err := s.GetSettings(gid)
	if err != nil {
		return err
	}

	gs.Welcome.TextColor = color

	gsb, err := json.Marshal(gs)
	if err != nil {
		return err
	}

	return s.db.Model(&models.Guilds{ID: gid}).Update("settings", gsb).Error
}

func (s *GuildsStore) GetCurrentWelcomeImageBackground(gid string) (string, error) {
	gs, err := s.GetWelcomeConfig(gid)

	if err != nil {
		return "", err
	}

	return gs.CustomBackground, nil
}

func (s *GuildsStore) GetGuildCategoryTemplate(gid string, idx int) ([]utils.CategoryTemplateChannel, string, error) {
	var g models.Guilds
	err := s.db.Where("guild_id = ?", gid).First(&g).Error
	if err != nil {
		return nil, "", err
	}

	var sett models.GuildSettings
	err = json.Unmarshal(g.Settings, &sett)

	if err != nil {
		return nil, "", err
	}

	catTemp := utils.ParseCTemp(sett.CategoryTemplate[idx].Template)

	return catTemp, sett.CategoryTemplate[idx].Name, nil
}

func (s *GuildsStore) SetGuildCategoryTemplate(gid string, name string, catTemp string) error {
	var g models.Guilds
	err := s.db.Where("guild_id = ?", gid).First(&g).Error
	if err != nil {
		return err
	}

	var sett models.GuildSettings
	err = json.Unmarshal(g.Settings, &sett)

	if err != nil {
		return err
	}

	sett.CategoryTemplate = append(sett.CategoryTemplate, models.CategoryTemplate{
		Name:     name,
		Template: catTemp,
	})

	gsb, err := json.Marshal(sett)
	if err != nil {
		return err
	}

	return s.db.Model(&models.Guilds{ID: gid}).Update("settings", gsb).Error
}

func (s *GuildsStore) ListGuildCategoryTemplates(gid string) ([]models.CategoryTemplate, error) {
	var g models.Guilds
	err := s.db.Where("guild_id = ?", gid).First(&g).Error
	if err != nil {
		return nil, err
	}

	var sett models.GuildSettings
	err = json.Unmarshal(g.Settings, &sett)
	if err != nil {
		return nil, err
	}

	return sett.CategoryTemplate, nil
}
func (s *GuildsStore) SetWelcomeChannel(gid string, cid string) error {
	gs, err := s.GetSettings(gid)
	if err != nil {
		return err
	}

	gs.Welcome.ChannelID = cid

	gsb, err := json.Marshal(gs)
	if err != nil {
		return err
	}

	return s.db.Model(&models.Guilds{ID: gid}).Update("settings", gsb).Error
}

func (s *GuildsStore) SetWelcomeMessage(gid string, msg string) error {
	gs, err := s.GetSettings(gid)
	if err != nil {
		return err
	}

	gs.Welcome.Message = msg

	gsb, err := json.Marshal(gs)
	if err != nil {
		return err
	}

	return s.db.Model(&models.Guilds{ID: gid}).Update("settings", gsb).Error
}

func (s *GuildsStore) SetWelcomeIncludeImage(gid string, inc bool) error {
	gs, err := s.GetSettings(gid)
	if err != nil {
		return err
	}

	gs.Welcome.IncludeImage = inc

	gsb, err := json.Marshal(gs)
	if err != nil {
		return err
	}

	return s.db.Model(&models.Guilds{ID: gid}).Update("settings", gsb).Error
}

func (s *GuildsStore) SetWelcomeCustomBackground(gid string, url string) error {
	gs, err := s.GetSettings(gid)
	if err != nil {
		return err
	}

	gs.Welcome.CustomBackground = url

	gsb, err := json.Marshal(gs)
	if err != nil {
		return err
	}

	return s.db.Model(&models.Guilds{ID: gid}).Update("settings", gsb).Error
}
