// File: commands/messageas.go
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

var MessageAs = &discordgo.ApplicationCommand{
	Name:        "messageas",
	Description: "Send message as bot",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionChannel,
			Name:        "channel",
			Description: "Channel to send the message to.",
			Required:    true,
		},
	},
}

type MessageAsInitData struct {
	ChannelID string
	Message   string
}

func (sc *StoreCtx) MessageAsCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
		log.Printf("err while checking permission of %v (%v) in %v (%v) for msga_init: %v", i.Member.User.Username, i.Member.User.ID, utils.GetGuildNameFromState(s, i.GuildID), i.GuildID, err)
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

	channel := optionMap["channel"].ChannelValue(s)

	data := &MessageAsInitData{
		ChannelID: channel.ID,
	}

	datab, err := json.Marshal(data)

	if err != nil {
		log.Printf("err while marshaling data for %v (%v) in guild %v (%v) in msga_init: %v", i.Member.User.Username, i.Member.User.ID, utils.GetGuildNameFromState(s, i.GuildID), i.GuildID, err)
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
		Type:      "msga_init",
		Next:      "msga_set_message",
		Data:      datab,
		CreatedAt: time.Now(),
		ExpiredAt: time.Now().Add(time.Minute * 5),
	}
	err, ec := sc.ps.Create(pen)

	if err != nil {
		switch ec {
		case -1:
			log.Printf("err while querying for current pending of %v (%v) in %v (%v) for msga_init: %v", i.Member.User.Username, i.Member.User.ID, utils.GetGuildNameFromState(s, i.GuildID), i.GuildID, err)
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
			Content: fmt.Sprintf("You are sending message to <#%v>. Type your message in chat or type `_cancel` to cancel.\n\nThis interaction will expire in 5 minutes.", channel.ID),
		},
	})
}
