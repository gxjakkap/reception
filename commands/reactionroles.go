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

package commands

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/gxjakkap/reception/models"
	"github.com/gxjakkap/reception/utils"
)

type RRInitPendingData struct {
	ChannelID string `json:"channel_id"`
}

var ReactionRoles = &discordgo.ApplicationCommand{
	Name:        "reactionroles",
	Description: "Set up reaction roles.",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionChannel,
			Name:        "channel",
			Description: "Channel to set up reaction roles in.",
			Required:    true,
			ChannelTypes: []discordgo.ChannelType{
				discordgo.ChannelTypeGuildText,
			},
		},
	},
}

func (sc *StoreCtx) ReactionRolesCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Member == nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "This command only works in server!",
			},
		})
		return
	}

	p, err := s.State.UserChannelPermissions(i.Member.User.ID, i.ChannelID)

	if err != nil {
		log.Printf("err while checking permission of %v (%v) in %v (%v) for rr_init: %v", i.Member.User.Username, i.Member.User.ID, utils.GetGuildNameFromState(s, i.GuildID), i.GuildID, err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Error!",
			},
		})
		return
	}

	if p&discordgo.PermissionManageGuild == 0 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Not Enough Permission! You need `PermissionManageGuild`.",
			},
		})
		return
	}

	targetCh := i.ApplicationCommandData().GetOption("channel").ChannelValue(s).ID

	datab, err := json.Marshal(&RRInitPendingData{
		ChannelID: targetCh,
	})

	if err != nil {
		log.Printf("err while marshalling data for %v (%v) in guild %v (%v) in rr_init: %v", i.Member.User.Username, i.Member.User.ID, utils.GetGuildNameFromState(s, i.GuildID), i.GuildID, err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Error!",
			},
		})
		return
	}

	pen := &models.Pending{
		GuildID:   i.GuildID,
		UserID:    i.Member.User.ID,
		Type:      "rr_init",
		Next:      "rr_set_embed_message",
		Data:      datab,
		CreatedAt: time.Now(),
		ExpiredAt: time.Now().Add(time.Minute * 5),
	}
	err, ec := sc.ps.Create(pen)

	if err != nil {
		switch ec {
		case -1:
			log.Printf("err while querying for current pending of %v (%v) in %v (%v) for rr_init: %v", i.Member.User.Username, i.Member.User.ID, utils.GetGuildNameFromState(s, i.GuildID), i.GuildID, err)
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Error!",
				},
			})
		case 1:
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "You already have pending interactions! Finish or cancel that interaction first or wait for it to expire.",
				},
			})
		case 2:
			log.Printf("err while creating pending row for %v (%v) in %v (%v) for rr_init: %v", i.Member.User.Username, i.Member.User.ID, utils.GetGuildNameFromState(s, i.GuildID), i.GuildID, err)
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Error!",
				},
			})

		}
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("You are creating reaction roles in <#%v>. Type in message to attach the reaction role message down below. (Expires in 5 minutes)", targetCh),
		},
	})
}
