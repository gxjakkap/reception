// File: commands/listctemp.go
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
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/gxjakkap/reception/utils"
)

var ListCTemp = &discordgo.ApplicationCommand{
	Name:        "listctemp",
	Description: "List category template for this server",
}

func (sc *StoreCtx) ListCTempCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
		log.Printf("err while checking permission of %v (%v) in %v (%v) for listctemp: %v", i.Member.User.Username, i.Member.User.ID, utils.GetGuildNameFromState(s, i.GuildID), i.GuildID, err)
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

	temps, err := sc.gs.ListGuildCategoryTemplates(i.GuildID)
	if err != nil {
		log.Printf("err while getting guild category template list for %v (%v) in guild %v (%v) in listctemp: %v", i.Member.User.Username, i.Member.User.ID, utils.GetGuildNameFromState(s, i.GuildID), i.GuildID, err)
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: utils.StringPtr("Error: Template not found at that index."),
		})
		return
	}

	var content string
	content += fmt.Sprintf("### Category Templates for **%v**\n\n", utils.GetGuildNameFromState(s, i.GuildID))
	for idx, temp := range temps {
		content += fmt.Sprintf("**%v. %v**\n", idx+1, temp.Name)
		ch := utils.ParseCTemp(temp.Template)
		for _, c := range ch {
			content += fmt.Sprintf("> • `[%s]` %s\n", utils.ChannelTypeToString(c.Type), c.Name)
		}
		content += "\n"
	}

	if len(temps) == 0 {
		content = fmt.Sprintf("No category templates found for **%v**.", utils.GetGuildNameFromState(s, i.GuildID))
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
		},
	})
}
