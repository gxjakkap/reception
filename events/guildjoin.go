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
	"github.com/gxjakkap/reception/models"
)

func (c *EventContext) GuildJoin(s *discordgo.Session, g *discordgo.GuildCreate) {
	if g.Unavailable {
		log.Printf("joined unavailable guild %v (%v)", g.Guild.Name, g.Guild.ID)
		return
	}

	log.Printf("joined %v (%v) owner: %v, usrcnt: %v", g.Guild.Name, g.Guild.ID, g.Owner, g.MemberCount)
	guild := &models.Guilds{
		ID:        g.ID,
		GuildName: g.Name,
	}

	err := c.gs.Create(guild)

	if err != nil {
		log.Printf("error adding guild %v (%v) to db!: %v", g.Name, g.ID, err)
	}
}
