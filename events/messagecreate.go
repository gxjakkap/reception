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

package events

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/gxjakkap/reception/followups"
	"github.com/gxjakkap/reception/utils"
)

func (c *EventContext) MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	// ignore dms, for now.
	if m.Member == nil {
		return
	}

	log.Printf("checkpoint: new msg event from %v (%v) in %v (%v): %v", m.Author.Username, m.Author.ID, utils.GetGuildNameFromState(s, m.GuildID), m.GuildID, m.Content)

	// if user have pending interaction
	pen, err := c.ps.UserInPendingList(m.Author.ID)

	if err != nil {
		log.Printf("error checking in user '%v' (%v) have pending interaction: %v", m.Author.Username, m.Author.ID, err)
		return
	}

	// user have pending interaction
	if pen {
		inter, err := c.ps.GetPendingByUser(m.Author.ID)

		log.Printf("checkpoint: check for pending interaction for %v (%v): %v", m.Author.Username, m.Author.ID, inter)

		if err != nil {
			log.Printf("error getting pending interaction for user '%v' (%v): %v", m.Author.Username, m.Author.ID, err)
			return
		}

		if fu, ok := followups.FollowUps[inter.Next]; ok {
			fu.Action(followups.NewFollowUpsCtx(c.gs, c.ps), s, m, inter.ID)
		}

		return
	}
}
