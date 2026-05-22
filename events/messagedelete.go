// File: events/messagedelete.go
// Project: Reception
// Author: Jakkaphat Chalermphanaphan <gunt@guntxjakka.me>
// Created: 2026-05-22 16:51+07
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

func (c *EventContext) MessageDelete(s *discordgo.Session, m *discordgo.MessageDelete) {
	// ignore dms, for now.
	if m.Member == nil {
		return
	}

	log.Printf("checkpoint: new msg delete event id %v in %v (%v): %v", m.Message.ID, utils.GetGuildNameFromState(s, m.GuildID), m.GuildID, m.Content)

}
