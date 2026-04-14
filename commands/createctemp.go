// File: commands/createctemp.go
// Project: Reception
// Author: Jakkaphat Chalermphanaphan <gunt@guntxjakka.me>
// Created: 2026-04-14 15:50+07
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

var CreateCTemp = &discordgo.ApplicationCommand{
	Name:        "createctemp",
	Description: "Create new category template",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "name",
			Description: "Name of the category template.",
			Required:    true,
		},
	},
}

type CCTInitData struct {
	Name     string
	Template string
}

func (sc *StoreCtx) CreateCTempCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Member == nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "This command only works in server!",
			},
		})
		return
	}

	p, err := s.State.UserChannelPermissions(i.User.ID, i.ChannelID)

	if err != nil {
		log.Printf("err while checking permission of %v (%v) in %v (%v) for cct_init: %v", i.User.Username, i.User.ID, utils.GetGuildNameFromState(s, i.GuildID), i.GuildID, err)
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

	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	name := optionMap["name"].StringValue()

	data := &CCTInitData{
		Name:     name,
		Template: "",
	}

	datab, err := json.Marshal(data)

	if err != nil {
		log.Printf("err while marshaling data for %v (%v) in guild %v (%v) in cct_init: %v", i.User.Username, i.User.ID, utils.GetGuildNameFromState(s, i.GuildID), i.GuildID, err)
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
		Type:      "cct_init",
		Next:      "cct_add_channels_or_finish",
		Data:      datab,
		CreatedAt: time.Now(),
		ExpiredAt: time.Now().Add(time.Minute * 5),
	}
	err, ec := sc.ps.Create(pen)

	if err != nil {
		switch ec {
		case -1:
			log.Printf("err while querying for current pending of %v (%v) in %v (%v) for cct_init: %v", i.Member.User.Username, i.Member.User.ID, utils.GetGuildNameFromState(s, i.GuildID), i.GuildID, err)
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
			log.Printf("err while creating pending row for %v (%v) in %v (%v) for cct_init: %v", i.Member.User.Username, i.Member.User.ID, utils.GetGuildNameFromState(s, i.GuildID), i.GuildID, err)
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
			Content: fmt.Sprintf("You are creating a new category template named %v. Type `channel_type channel_name` in chat one by one.\n\nValid channel type:\n- text\n- voice\n- announcement\n- stage\n- forum\n- media\n\nWhen finished, type `end`.\n\nThis interaction will expire in 5 minutes.", i.ApplicationCommandData().Options[0].StringValue()),
		},
	})
}
