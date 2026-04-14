// File: commands/tempc.go
// Project: Reception
// Author: Jakkaphat Chalermphanaphan <gunt@guntxjakka.me>
// Created: 2026-04-14 14:26+07
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
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/gxjakkap/reception/utils"
)

var TempC = &discordgo.ApplicationCommand{
	Name:        "tempc",
	Description: "Create new category with predefined channels",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "name",
			Description: "Name of the category.",
			Required:    true,
		},
		{
			Type:        discordgo.ApplicationCommandOptionInteger,
			Name:        "idx",
			Description: "Index of the category template to use. (Default to 0 if not specified)",
		},
	},
}

func (sc *StoreCtx) TempCCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
		log.Printf("err while checking permission of %v (%v) in %v (%v) for tempc: %v", i.User.Username, i.User.ID, utils.GetGuildNameFromState(s, i.GuildID), i.GuildID, err)
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

	idx := 0
	if opt, ok := optionMap["idx"]; ok {
		idx = int(opt.IntValue())
	}

	catTemp, tempName, err := sc.gs.GetGuildCategoryTemplate(i.GuildID, idx)

	if err != nil {
		log.Printf("err while getting guild category template for %v (%v) in guild %v (%v) in tempc_init: %v", i.User.Username, i.User.ID, utils.GetGuildNameFromState(s, i.GuildID), i.GuildID, err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Error!",
			},
		})
		return
	}

	name := optionMap["name"].StringValue()

	cat, err := s.GuildChannelCreate(i.GuildID, name, discordgo.ChannelTypeGuildCategory)

	if err != nil {
		log.Printf("err while creating category for %v (%v) in guild %v (%v) in tempc_init: %v", i.User.Username, i.User.ID, utils.GetGuildNameFromState(s, i.GuildID), i.GuildID, err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Error!",
			},
		})
		return
	}

	for _, ch := range catTemp {
		cData := &discordgo.GuildChannelCreateData{
			Name:     ch.Name,
			Type:     *ch.Type,
			ParentID: cat.ID,
		}
		_, err := s.GuildChannelCreateComplex(i.GuildID, *cData)
		if err != nil {
			log.Printf("err while creating channel for %v (%v) in guild %v (%v) in tempc_init: %v", i.User.Username, i.User.ID, utils.GetGuildNameFromState(s, i.GuildID), i.GuildID, err)
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Error creating channel!",
				},
			})
			return
		}
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Created category %v (%v) with %v channels.", name, tempName, len(catTemp)),
		},
	})
}
