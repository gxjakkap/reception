// File: utils/cattemp.go
// Project: Reception
// Author: Jakkaphat Chalermphanaphan <gunt@guntxjakka.me>
// Created: 2026-04-14 14:31+07
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

package utils

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

type CategoryTemplateChannel struct {
	Name string
	Type *discordgo.ChannelType
}

/**
	Template:
	"channelType:channelName;channelType:channelName;..."

	Valid channelType: text, voice, announcement, stage, forum

	Example: "text:general;voice:voice-chat;"
**/

func ParseCTemp(s string) []CategoryTemplateChannel {
	var channels []CategoryTemplateChannel

	segments := strings.Split(s, ";")
	for _, segment := range segments {
		segment = strings.TrimSpace(segment)
		if segment == "" {
			continue
		}

		parts := strings.SplitN(segment, ":", 2)
		if len(parts) != 2 {
			continue
		}

		cTypeStr := strings.TrimSpace(strings.ToLower(parts[0]))
		cName := strings.TrimSpace(parts[1])

		var cType discordgo.ChannelType

		switch cTypeStr {
		case "text":
			cType = discordgo.ChannelTypeGuildText
		case "voice":
			cType = discordgo.ChannelTypeGuildVoice
		case "announcement":
			cType = discordgo.ChannelTypeGuildNews
		case "stage":
			cType = discordgo.ChannelTypeGuildStageVoice
		case "forum":
			cType = discordgo.ChannelTypeGuildForum
		default:
			// Unknown type, skip or handle error. We will skip.
			continue
		}

		channels = append(channels, CategoryTemplateChannel{
			Name: cName,
			Type: &cType,
		})
	}

	return channels
}

func ConcatTemplate(prev string, new string) string {
	if prev == "" {
		return new
	}
	return prev + new + ";"
}
