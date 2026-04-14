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

package followups

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/gxjakkap/reception/commands"
	"github.com/gxjakkap/reception/models"
	"github.com/gxjakkap/reception/utils"
)

var CCTAddChannelsOrFinish = &FollowUp{
	Name:   "cct_add_channels_or_finish",
	Action: (*FollowUpsCtx).CCTAddChannels,
}

func (c *FollowUpsCtx) CCTAddChannels(s *discordgo.Session, m *discordgo.MessageCreate, id uint) {
	pen, err := c.ps.Fulfill(id)

	if err != nil {
		log.Printf("error while popping pending data for %v (%v) in %v (%v) in cct_add_channels_or_finish: %v", m.Author.Username, m.Author.ID, utils.GetGuildNameFromState(s, m.GuildID), m.GuildID, err)
		return
	}

	var data commands.CCTInitData

	err = json.Unmarshal(pen.Data, &data)

	if err != nil {
		log.Printf("error while unmarshaling data for %v (%v) in %v (%v) in cct_add_channels_or_finish: %v", m.Author.Username, m.Author.ID, utils.GetGuildNameFromState(s, m.GuildID), m.GuildID, err)
		return
	}

	if strings.ToLower(m.Content) == "end" {
		err := c.gs.SetGuildCategoryTemplate(m.GuildID, data.Name, data.Template)

		if err != nil {
			log.Printf("error while saving category template for %v (%v) in %v (%v) in cct_add_channels_or_finish: %v", m.Author.Username, m.Author.ID, utils.GetGuildNameFromState(s, m.GuildID), m.GuildID, err)
			s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
				Content:   "Error adding template!",
				Reference: m.MessageReference,
			})
			return
		}

		s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
			Content:   "Template added successfully!",
			Reference: m.MessageReference,
		})
		s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
		return
	}

	parts := strings.SplitN(m.Content, " ", 2)
	if len(parts) != 2 {
		s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
			Content:   "Invalid format! Please type `channel_type channel_name` or `end` to finish.",
			Reference: m.MessageReference,
		})

		// Renew pending so user can try again
		c.ps.Create(&models.Pending{
			GuildID:   m.GuildID,
			UserID:    m.Author.ID,
			Type:      "cct_add_channels_or_finish",
			Next:      "cct_add_channels_or_finish",
			Data:      pen.Data,
			CreatedAt: time.Now(),
			ExpiredAt: time.Now().Add(time.Minute * 5),
		})
		s.MessageReactionAdd(m.ChannelID, m.ID, "❌")
		return
	}

	cType := strings.ToLower(strings.TrimSpace(parts[0]))
	cName := strings.TrimSpace(parts[1])

	validTypes := map[string]bool{
		"text": true, "voice": true, "announcement": true, "stage": true, "forum": true, "media": true,
	}

	if !validTypes[cType] {
		s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
			Content:   "Invalid channel type! Valid types are: text, voice, announcement, stage, forum, media.",
			Reference: m.MessageReference,
		})
		c.ps.Create(&models.Pending{
			GuildID:   m.GuildID,
			UserID:    m.Author.ID,
			Type:      "cct_add_channels_or_finish",
			Next:      "cct_add_channels_or_finish",
			Data:      pen.Data,
			CreatedAt: time.Now(),
			ExpiredAt: time.Now().Add(time.Minute * 5),
		})
		s.MessageReactionAdd(m.ChannelID, m.ID, "❌")
		return
	}

	if data.Template != "" {
		data.Template += ";"
	}
	data.Template += fmt.Sprintf("%s:%s", cType, cName)

	datab, err := json.Marshal(data)

	if err != nil {
		log.Printf("error while marshaling data for %v (%v) in %v (%v) in cct_add_channels_or_finish: %v", m.Author.Username, m.Author.ID, utils.GetGuildNameFromState(s, m.GuildID), m.GuildID, err)
		return
	}

	s.MessageReactionAdd(m.ChannelID, m.ID, "✅")

	c.ps.Create(&models.Pending{
		GuildID:   m.GuildID,
		UserID:    m.Author.ID,
		Type:      "cct_add_channels_or_finish",
		Next:      "cct_add_channels_or_finish",
		Data:      datab,
		CreatedAt: time.Now(),
		ExpiredAt: time.Now().Add(time.Minute * 5),
	})
}
