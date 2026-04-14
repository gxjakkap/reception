// File: commands/welcome.go
// Project: Reception
// Author: Jakkaphat Chalermphanaphan <gunt@guntxjakka.me>
// Created: 2026-02-05 01:05+07
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
	"context"
	"fmt"
	"log"
	"net/http"
	"path"

	"github.com/bwmarrin/discordgo"
	"github.com/gxjakkap/reception/utils"
)

var Welcome = &discordgo.ApplicationCommand{
	Name:        "welcome",
	Description: "Set up welcome message.",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name:        "enable",
			Description: "Toggle welcome message functionality.",
			Type:        discordgo.ApplicationCommandOptionSubCommand,
		},
		{
			Name:        "settings",
			Description: "Query or change settings for welcome message.",
			Type:        discordgo.ApplicationCommandOptionSubCommandGroup,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "channel",
					Description: "Set the channel for welcome messages.",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "chan",
							Description: "The channel.",
							Type:        discordgo.ApplicationCommandOptionChannel,
							Required:    true,
						},
					},
				},
				{
					Name:        "message",
					Description: "Set the welcome message text. Use {user} for mention and {server} for guild name.",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "msg",
							Description: "The message text.",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    true,
						},
					},
				},
				{
					Name:        "include_image",
					Description: "Toggle whether to include a generated image in the welcome message.",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "inc",
							Description: "Include image?",
							Type:        discordgo.ApplicationCommandOptionBoolean,
							Required:    true,
						},
					},
				},
				{
					Name:        "textcolor",
					Description: "Get or set HEX value for text in the welcome message image.",
					Type:        discordgo.ApplicationCommandOptionSubCommandGroup,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "get",
							Description: "Get current HEX value.",
							Type:        discordgo.ApplicationCommandOptionSubCommand,
						},
						{
							Name:        "set",
							Description: "Set new HEX value.",
							Type:        discordgo.ApplicationCommandOptionSubCommand,
							Options: []*discordgo.ApplicationCommandOption{
								{
									Name:        "value",
									Description: "HEX value for text in the welcome message image.",
									Type:        discordgo.ApplicationCommandOptionString,
									Required:    true,
								},
							},
						},
					},
				},
				{
					Name:        "background",
					Description: "Get or set custom background image for the welcome message image.",
					Type:        discordgo.ApplicationCommandOptionSubCommandGroup,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "get",
							Description: "Get current background image.",
							Type:        discordgo.ApplicationCommandOptionSubCommand,
						},
						{
							Name:        "set",
							Description: "Set new background image.",
							Type:        discordgo.ApplicationCommandOptionSubCommandGroup,
							Options: []*discordgo.ApplicationCommandOption{
								{
									Name:        "image",
									Type:        discordgo.ApplicationCommandOptionSubCommand,
									Description: "Set new background image using attachment.",
									Options: []*discordgo.ApplicationCommandOption{
										{
											Name:     "img",
											Type:     discordgo.ApplicationCommandOptionAttachment,
											Required: true,
										},
									},
								},
								{
									Name:        "imagelink",
									Type:        discordgo.ApplicationCommandOptionSubCommand,
									Description: "Set new background image using URL.",
									Options: []*discordgo.ApplicationCommandOption{
										{
											Name:        "link",
											Description: "URL of the image.",
											Type:        discordgo.ApplicationCommandOptionString,
											Required:    true,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	},
}

func (sc *StoreCtx) WelcomeCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
		log.Printf("err while checking permission of %v (%v) in %v (%v) for wc_settings: %v", i.User.Username, i.User.ID, utils.GetGuildNameFromState(s, i.GuildID), i.GuildID, err)
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
	switch options[0].Name {
	case "enable":
		cur, err := sc.gs.ToggleWelcomeMessage(i.GuildID)
		if err != nil {
			log.Printf("err while toggling WelcomeSettings.Enable of in %v (%v) by %v (%v) for wc_settings: %v", utils.GetGuildNameFromState(s, i.GuildID), i.GuildID, i.User.Username, i.User.ID, err)
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Error!",
				},
			})
		}

		var nsm string

		if cur {
			nsm = "enabled"
		} else {
			nsm = "disabled"
		}

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("WelcomeSettings.Enabled has been set to `%v`", nsm),
			},
		})

		return
	case "settings":
		for _, opt := range options[0].Options {
			switch opt.Name {
			case "channel":
				ch := opt.Options[0].ChannelValue(s)
				err := sc.gs.SetWelcomeChannel(i.GuildID, ch.ID)
				if err != nil {
					log.Printf("err while setting WelcomeSettings.ChannelID of in %v (%v) by %v (%v) for wc_settings: %v", utils.GetGuildNameFromState(s, i.GuildID), i.GuildID, i.User.Username, i.User.ID, err)
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "Error!",
						},
					})
					return
				}
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: fmt.Sprintf("Welcome message channel has been set to <#%v>", ch.ID),
					},
				})
				return
			case "message":
				msg := opt.Options[0].StringValue()
				err := sc.gs.SetWelcomeMessage(i.GuildID, msg)
				if err != nil {
					log.Printf("err while setting WelcomeSettings.Message of in %v (%v) by %v (%v) for wc_settings: %v", utils.GetGuildNameFromState(s, i.GuildID), i.GuildID, i.User.Username, i.User.ID, err)
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "Error!",
						},
					})
					return
				}
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: fmt.Sprintf("Welcome message has been set to:\n\n%v", msg),
					},
				})
				return
			case "include_image":
				inc := opt.Options[0].BoolValue()
				err := sc.gs.SetWelcomeIncludeImage(i.GuildID, inc)
				if err != nil {
					log.Printf("err while setting WelcomeSettings.IncludeImage of in %v (%v) by %v (%v) for wc_settings: %v", utils.GetGuildNameFromState(s, i.GuildID), i.GuildID, i.User.Username, i.User.ID, err)
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "Error!",
						},
					})
					return
				}
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: fmt.Sprintf("Welcome message image inclusion has been set to `%v`", inc),
					},
				})
				return
			case "textcolor":
				for _, sOpt := range opt.Options {
					switch sOpt.Name {
					case "get":
						tc, err := sc.gs.GetWelcomeTextColor(i.GuildID)
						if err != nil {
							log.Printf("err while getting WelcomeSettings.TextColor of in %v (%v) by %v (%v) for wc_settings: %v", utils.GetGuildNameFromState(s, i.GuildID), i.GuildID, i.User.Username, i.User.ID, err)
							s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
								Type: discordgo.InteractionResponseChannelMessageWithSource,
								Data: &discordgo.InteractionResponseData{
									Content: "Error!",
								},
							})
							return
						}
						s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
							Type: discordgo.InteractionResponseChannelMessageWithSource,
							Data: &discordgo.InteractionResponseData{
								Content: fmt.Sprintf("Current Welcome Image Text Color of `%v` is `%v`", utils.GetGuildNameFromState(s, i.GuildID), tc),
							},
						})
						return
					case "set":
						color := sOpt.GetOption("value").StringValue()

						if color == "" {
							return
						}

						err := sc.gs.SetWelcomeTextColor(i.GuildID, color)

						if err != nil {
							if err.Error() == "invalid hex color string" {
								s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
									Type: discordgo.InteractionResponseChannelMessageWithSource,
									Data: &discordgo.InteractionResponseData{
										Content: "Invalid color option! You must use valid HEX color string.",
									},
								})
							} else {
								log.Printf("err while setting WelcomeSettings.TextColor of in %v (%v) by %v (%v) for wc_settings: %v", utils.GetGuildNameFromState(s, i.GuildID), i.GuildID, i.User.Username, i.User.ID, err)
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
								Content: fmt.Sprintf("Welcome Image Text Color of `%v` has been set to `%v`", utils.GetGuildNameFromState(s, i.GuildID), color),
							},
						})
						return
					}
				}
			case "background":
				for _, subOpt := range opt.Options {
					switch subOpt.Name {
					case "get":
						img, err := sc.gs.GetCurrentWelcomeImageBackground(i.GuildID)

						if err != nil {
							log.Printf("err while getting WelcomeSettings.CustomImage of %v (%v) by %v (%v) for wc_settings: %v", utils.GetGuildNameFromState(s, i.GuildID), i.GuildID, i.User.Username, i.User.ID, err)
							s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
								Type: discordgo.InteractionResponseChannelMessageWithSource,
								Data: &discordgo.InteractionResponseData{
									Content: "Error!",
								},
							})
							return
						}

						if img == "" {
							s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
								Type: discordgo.InteractionResponseChannelMessageWithSource,
								Data: &discordgo.InteractionResponseData{
									Content: "There's no custom background image set.",
								},
							})
							return
						}

						st, err := utils.NewS3Client()

						if err != nil {
							log.Printf("err while getting custom image from s3 for %v (%v) by %v (%v) for wc_settings: %v", utils.GetGuildNameFromState(s, i.GuildID), i.GuildID, i.User.Username, i.User.ID, err)
							s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
								Type: discordgo.InteractionResponseChannelMessageWithSource,
								Data: &discordgo.InteractionResponseData{
									Content: "Error getting image from storage.",
								},
							})
							return
						}

						s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
							Type: discordgo.InteractionResponseChannelMessageWithSource,
							Data: &discordgo.InteractionResponseData{
								Embeds: []*discordgo.MessageEmbed{
									{
										Title: fmt.Sprintf("Current custom image for welcome message of `%v`", utils.GetGuildNameFromState(s, i.GuildID)),
										Image: &discordgo.MessageEmbedImage{
											URL: st.GetFileUrl(img),
										},
									},
								},
							},
						})
					case "set":
						st, err := utils.NewS3Client()
						if err != nil {
							log.Printf("err while setting up s3 client for %v (%v) by %v (%v) for wc_settings background: %v", utils.GetGuildNameFromState(s, i.GuildID), i.GuildID, i.User.Username, i.User.ID, err)
							s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
								Type: discordgo.InteractionResponseChannelMessageWithSource,
								Data: &discordgo.InteractionResponseData{
									Content: "Storage service unavailable.",
								},
							})
							return
						}

						for _, subSubOpt := range subOpt.Options {
							switch subSubOpt.Name {
							case "image":
								attID := subSubOpt.GetOption("img").Value.(string)
								att := i.ApplicationCommandData().Resolved.Attachments[attID]

								resp, err := http.Get(att.URL)
								if err != nil {
									s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
										Type: discordgo.InteractionResponseChannelMessageWithSource,
										Data: &discordgo.InteractionResponseData{
											Content: "Failed to download the attachment.",
										},
									})
									return
								}
								defer resp.Body.Close()

								key := fmt.Sprintf("welcome_bg/%s%s", i.GuildID, path.Ext(att.Filename))
								_, err = st.UploadFile(context.Background(), key, resp.Body, resp.Header.Get("Content-Type"))
								if err != nil {
									log.Printf("err while uploading background image for %v (%v): %v", i.GuildID, i.User.ID, err)
									s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
										Type: discordgo.InteractionResponseChannelMessageWithSource,
										Data: &discordgo.InteractionResponseData{
											Content: "Error uploading image to storage.",
										},
									})
									return
								}

								sc.gs.SetWelcomeCustomBackground(i.GuildID, key)
								s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
									Type: discordgo.InteractionResponseChannelMessageWithSource,
									Data: &discordgo.InteractionResponseData{
										Content: "Welcome background image updated successfully!",
									},
								})

							case "imagelink":
								link := subSubOpt.GetOption("link").StringValue()
								resp, err := http.Get(link)
								if err != nil {
									s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
										Type: discordgo.InteractionResponseChannelMessageWithSource,
										Data: &discordgo.InteractionResponseData{
											Content: "Failed to access the image link.",
										},
									})
									return
								}
								defer resp.Body.Close()

								key := fmt.Sprintf("welcome_bg/%s%s", i.GuildID, path.Ext(link))
								if path.Ext(link) == "" {
									key += ".png" // default extension if missing
								}

								_, err = st.UploadFile(context.Background(), key, resp.Body, resp.Header.Get("Content-Type"))
								if err != nil {
									log.Printf("err while uploading background link for %v (%v): %v", i.GuildID, i.User.ID, err)
									s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
										Type: discordgo.InteractionResponseChannelMessageWithSource,
										Data: &discordgo.InteractionResponseData{
											Content: "Error uploading image to storage.",
										},
									})
									return
								}

								sc.gs.SetWelcomeCustomBackground(i.GuildID, key)
								s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
									Type: discordgo.InteractionResponseChannelMessageWithSource,
									Data: &discordgo.InteractionResponseData{
										Content: "Welcome background image updated successfully!",
									},
								})
							}
						}
					}
				}
			}
		}
	}
}
