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
	"bytes"
	"context"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/gxjakkap/reception/enum"
	"github.com/gxjakkap/reception/utils"
)

func (c *EventContext) GuildMemberAdd(s *discordgo.Session, g *discordgo.GuildMemberAdd) {
	gsett, err := c.gs.GetSettings(g.GuildID)

	if err != nil {
		log.Printf("error getting guild settings for %v (%v): %v", utils.GetGuildNameFromState(s, g.GuildID), g.GuildID, err)
		return
	}

	if gsett.Features[enum.GuildSetting[enum.WelcomeMessage]] {
		wconf := gsett.Welcome
		if wconf.ChannelID == "" {
			return
		}

		msg := wconf.Message
		if msg == "" {
			msg = "Welcome {user} to {server}!"
		}

		guildName := utils.GetGuildNameFromState(s, g.GuildID)
		msg = strings.ReplaceAll(msg, "{user}", g.Member.User.Mention())
		msg = strings.ReplaceAll(msg, "{server}", guildName)

		ms := &discordgo.MessageSend{
			Content: msg,
		}

		if wconf.IncludeImage {
			st, err := utils.NewS3Client()
			var bgImg image.Image
			if err == nil && wconf.CustomBackground != "" {
				body, err := st.DownloadFile(context.Background(), wconf.CustomBackground)
				if err == nil {
					defer body.Close()
					bgImg, _, _ = image.Decode(body)
				} else {
					log.Printf("error downloading custom background from s3 for %v: %v", g.GuildID, err)
				}
			}

			imgBytes, err := utils.GenerateWelcomeImage(
				g.Member.User.Username,
				guildName,
				g.Member.User.AvatarURL("256"),
				bgImg,
				wconf.TextColor,
			)

			if err != nil {
				log.Printf("error generating welcome image for %v in %v: %v", g.Member.User.Username, g.GuildID, err)
			} else {
				ms.Files = []*discordgo.File{
					{
						Name:        "welcome.png",
						ContentType: "image/png",
						Reader:      bytes.NewReader(imgBytes),
					},
				}
			}
		}

		_, err = s.ChannelMessageSendComplex(wconf.ChannelID, ms)
		if err != nil {
			log.Printf("error sending welcome message in %v (%v): %v", guildName, g.GuildID, err)
		}
	}
}
