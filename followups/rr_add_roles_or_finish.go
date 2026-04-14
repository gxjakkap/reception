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
	"log"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/gxjakkap/reception/models"
	"github.com/gxjakkap/reception/utils"
)

var RRAddRolesOrFinish = &FollowUp{
	Name:   "rr_add_roles_or_finish",
	Action: (*FollowUpsCtx).AddRolesOrFinish,
}

func (c *FollowUpsCtx) AddRolesOrFinish(s *discordgo.Session, m *discordgo.MessageCreate, id uint) {
	pen, err := c.ps.Fulfill(id)

	if err != nil {
		log.Printf("error while popping pending data for %v (%v) in %v (%v) in rr_add_roles_or_finish: %v", m.Author.Username, m.Author.ID, utils.GetGuildNameFromState(s, m.GuildID), m.GuildID, err)
		return
	}

	var data RRSetEmbedMessageOrAddRolesOrFinishData

	err = json.Unmarshal(pen.Data, &data)

	if err != nil {
		log.Printf("error while unmarshaling data for %v (%v) in %v (%v) in rr_add_roles_or_finish: %v", m.Author.Username, m.Author.ID, utils.GetGuildNameFromState(s, m.GuildID), m.GuildID, err)
		return
	}

	if m.Content == "end" {
		s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
			Content:   "Reaction roles have been added successfully!",
			Reference: m.MessageReference,
		})
		return
	}

	var validRoles []*models.ReactionRoles
	var validEmoji []string
	var validRoleIDs []string

	lines := strings.Split(m.Content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 2 {
			s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
				Content:   "Invalid format! Please use `@Role :emoji:` format for each line.",
				Reference: m.MessageReference,
			})
			c.ps.Create(&models.Pending{
				GuildID:   m.GuildID,
				UserID:    m.Member.User.ID,
				Type:      "rr_add_roles_or_finish",
				Next:      "rr_add_roles_or_finish",
				Data:      pen.Data,
				CreatedAt: time.Now(),
				ExpiredAt: time.Now().Add(time.Minute * 5),
			})
			return
		}

		roleStr := parts[0]
		emojiStr := parts[1]

		roleID := utils.ExtractRoleID(roleStr)
		if roleID == "" {
			s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
				Content:   "Invalid role mention! Please mention a valid role.",
				Reference: m.MessageReference,
			})
			c.ps.Create(&models.Pending{
				GuildID:   m.GuildID,
				UserID:    m.Member.User.ID,
				Type:      "rr_add_roles_or_finish",
				Next:      "rr_add_roles_or_finish",
				Data:      pen.Data,
				CreatedAt: time.Now(),
				ExpiredAt: time.Now().Add(time.Minute * 5),
			})
			return
		}

		validRoles = append(validRoles, &models.ReactionRoles{
			GuildID:   m.GuildID,
			Emoji:     emojiStr,
			RoleID:    roleID,
			MessageID: data.MessageID,
		})
		validEmoji = append(validEmoji, emojiStr)
		validRoleIDs = append(validRoleIDs, roleID)
	}

	if len(validRoles) == 0 {
		return
	}

	for _, rr := range validRoles {
		c.gs.AddReactionRole(m.GuildID, rr)
		s.MessageReactionAdd(data.ChannelID, data.MessageID, rr.Emoji)
	}

	data.AddedEmoji = append(data.AddedEmoji, validEmoji...)
	data.AddedRoles = append(data.AddedRoles, validRoleIDs...)

	datab, err := json.Marshal(data)

	if err != nil {
		log.Printf("error while marshaling data for %v (%v) in %v (%v) in rr_add_roles_or_finish: %v", m.Author.Username, m.Author.ID, utils.GetGuildNameFromState(s, m.GuildID), m.GuildID, err)
		return
	}

	s.MessageReactionAdd(m.ChannelID, m.MessageReference.MessageID, "✅")

	c.ps.Create(&models.Pending{
		GuildID:   m.GuildID,
		UserID:    m.Member.User.ID,
		Type:      "rr_add_roles_or_finish",
		Next:      "rr_add_roles_or_finish",
		Data:      datab,
		CreatedAt: time.Now(),
		ExpiredAt: time.Now().Add(time.Minute * 5),
	})
}
