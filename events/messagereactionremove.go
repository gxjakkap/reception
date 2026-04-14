// File: events/messagereactionremove.go
// Project: Reception
// Author: Jakkaphat Chalermphanaphan <gunt@guntxjakka.me>
// Created: 2026-04-14 19:24+07
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

package events

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/gxjakkap/reception/utils"
)

func (c *EventContext) MessageReactionRemove(s *discordgo.Session, m *discordgo.MessageReactionRemove) {
	// check cache first to avoid query
	if !c.gs.IsReactionRoleMessage(m.MessageID) {
		return
	}

	log.Printf("checkpoint: reaction remove event in %v (%v): msgid:%v emoji:%v", utils.GetGuildNameFromState(s, m.GuildID), m.GuildID, m.MessageID, m.MessageReaction.Emoji)

	// check if message is set up for reaction roles
	rrs, err := c.gs.GetReactionRolesFromMessage(m.MessageID)

	if err != nil {
		log.Printf("error getting reaction roles for message '%v': %v", m.MessageID, err)
		return
	}

	for _, rr := range rrs {
		if rr.Emoji == m.Emoji.APIName() {
			err := s.GuildMemberRoleRemove(m.GuildID, m.UserID, rr.RoleID)
			if err != nil {
				log.Printf("error removing role '%v' from user '%v' in guild '%v': %v", rr.RoleID, m.UserID, m.GuildID, err)
			}
			return
		}
	}
}
