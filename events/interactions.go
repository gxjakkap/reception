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
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/gxjakkap/reception/utils"
)

func (c *EventContext) Interactions(s *discordgo.Session, i *discordgo.InteractionCreate, hs map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)) {
	cn := i.ApplicationCommandData().Name
	h, ok := hs[cn]
	if !ok {
		return
	}

	var lp string
	un := "Unknown"
	if i.Member != nil && i.Member.User != nil {
		un = i.Member.User.Username
	} else if i.User != nil {
		un = i.User.Username
	}

	if i.GuildID == "" {
		lp = fmt.Sprintf("usr: %s", un)
	} else {
		gn := utils.GetGuildNameFromState(s, i.GuildID)
		lp = fmt.Sprintf("src: %s (%s) usr: %s", gn, i.GuildID, un)
	}

	log.Printf("%s int: %s", lp, cn)

	h(s, i)
}
