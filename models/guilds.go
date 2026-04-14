// Copyright 2026 Jakkaphat Chalermphanaphan

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     https://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Guilds struct {
	ID            string          `gorm:"primaryKey;type:text"`
	GuildName     string          `gorm:"not null"`
	JoinedAt      time.Time       `gorm:"autoCreateTime"`
	DeletedAt     gorm.DeletedAt  `gorm:"index"`
	Prefix        string          `gorm:"default:'&&'"`
	Settings      datatypes.JSON  `gorm:"type:jsonb;default:'{}'"`
	ReactionRoles []ReactionRoles `gorm:"foreignKey:GuildID;constraint:OnDelete:CASCADE"`
	History       []History       `gorm:"foreignKey:GuildID;constraint:OnDelete:CASCADE"`
}

type GuildSettings struct {
	Features         map[string]bool    `json:"features"`
	Welcome          WelcomeConfig      `json:"welcome"`
	Tickets          TicketConfig       `json:"tickets"`
	CategoryTemplate []CategoryTemplate `json:"category_template"`
}

type WelcomeConfig struct {
	Enabled          bool   `json:"enabled"`
	ChannelID        string `json:"channel_id"`
	Message          string `json:"message"`
	IncludeImage     bool   `json:"include_image"`
	CustomBackground string `json:"custom_background"`
	TextColor        string `json:"text_color"`
}

type TicketConfig struct {
	CategoryID        string `json:"category_id"`
	ChannelNamePrefix string `json:"channel_name_prefix"`
}

type CategoryTemplate struct {
	Name     string `json:"name"`
	Template string `json:"template"`
}
