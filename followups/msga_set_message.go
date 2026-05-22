// File: followups/msga_set_message.go
// Project: Reception
// Author: Jakkaphat Chalermphanaphan <gunt@guntxjakka.me>
// Created: 2026-05-22 18:46+07
// Copyright 2026 Jakkaphat Chalermphanaphan
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

package followups

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/gxjakkap/reception/commands"
	"github.com/gxjakkap/reception/utils"
)

var MsgaSetMessage = &FollowUp{
	Name:   "msga_set_message",
	Action: (*FollowUpsCtx).MsgaSetMessage,
}

func (c *FollowUpsCtx) MsgaSetMessage(s *discordgo.Session, m *discordgo.MessageCreate, id uint) {
	pen, err := c.ps.Fulfill(id)

	if err != nil {
		log.Printf("error while popping pending data for %v (%v) in %v (%v) in msga_set_message: %v", m.Author.Username, m.Author.ID, utils.GetGuildNameFromState(s, m.GuildID), m.GuildID, err)
		return
	}

	var data commands.MessageAsInitData

	err = json.Unmarshal(pen.Data, &data)
	if err != nil {
		log.Printf("error while unmarshaling data for %v (%v) in %v (%v) in msga_set_message: %v", m.Author.Username, m.Author.ID, utils.GetGuildNameFromState(s, m.GuildID), m.GuildID, err)
		return
	}

	if strings.ToLower(m.Content) == "_cancel" {
		s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
			Content:   "Cancelled.",
			Reference: m.MessageReference,
		})
		s.MessageReactionAdd(m.ChannelID, m.ID, "🫡")
		return
	}

	_, err = s.ChannelMessageSend(data.ChannelID, m.Content)
	if err != nil {
		log.Printf("error while sending message to channel %v for %v (%v) in %v (%v) in msga_set_message: %v", data.ChannelID, m.Author.Username, m.Author.ID, utils.GetGuildNameFromState(s, m.GuildID), m.GuildID, err)
		s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
			Content:   "Error sending message!",
			Reference: m.MessageReference,
		})
		return
	}

	s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Content:   "Message sent to <#" + data.ChannelID + "> successfully!",
		Reference: m.MessageReference,
	})
}
