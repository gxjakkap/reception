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

package followups

import (
	"encoding/json"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/gxjakkap/reception/commands"
	"github.com/gxjakkap/reception/models"
	"github.com/gxjakkap/reception/utils"
)

var RRSetEmbedMessage = &FollowUp{
	Name:   "rr_set_embed_message",
	Action: (*FollowUpsCtx).SetEmbedMessage,
}

type RRSetEmbedMessageOrAddRolesOrFinishData struct {
	ChannelID  string
	MessageID  string
	AddedEmoji []string
	AddedRoles []string
}

func (c *FollowUpsCtx) SetEmbedMessage(s *discordgo.Session, m *discordgo.MessageCreate, id uint) {
	pen, err := c.ps.Fulfill(id)

	if err != nil {
		log.Printf("error while popping pending data for %v (%v) in %v (%v) in rr_set_embed_message: %v", m.Author.Username, m.Author.ID, utils.GetGuildNameFromState(s, m.GuildID), m.GuildID, err)
		return
	}

	var data commands.RRInitPendingData

	err = json.Unmarshal(pen.Data, &data)

	if err != nil {
		log.Printf("error while unmarshaling data for %v (%v) in %v (%v) in rr_set_embed_message: %v", m.Author.Username, m.Author.ID, utils.GetGuildNameFromState(s, m.GuildID), m.GuildID, err)
		return
	}

	// send embed message to target channel
	t, err := s.ChannelMessageSend(data.ChannelID, m.Content)

	if err != nil {
		log.Printf("error while sending target message for %v (%v) in %v (%v) in rr_set_embed_message: %v", m.Author.Username, m.Author.ID, utils.GetGuildNameFromState(s, m.GuildID), m.GuildID, err)
		return
	}

	rep := &discordgo.MessageSend{
		Content:   "I have sent the message you provided to <%v>.\n\nNext, please input each role and corresponding discord emoji by mentioning a role followed by the emoji. e.g.:\n'@role' :book:\nWhen finished, type `end`.",
		Reference: m.MessageReference,
	}

	s.ChannelMessageSendComplex(m.ChannelID, rep)

	npd := &RRSetEmbedMessageOrAddRolesOrFinishData{
		ChannelID:  data.ChannelID,
		MessageID:  t.MessageReference.MessageID,
		AddedEmoji: []string{},
		AddedRoles: []string{},
	}

	npdb, err := json.Marshal(npd)

	if err != nil {
		log.Printf("error while marshaling data for %v (%v) in %v (%v) in rr_set_embed_message: %v", m.Author.Username, m.Author.ID, utils.GetGuildNameFromState(s, m.GuildID), m.GuildID, err)
		s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
			Content:   "Error occured! Please try running `/reactionroles` again.",
			Reference: m.MessageReference,
		})
		return
	}

	nPen := &models.Pending{
		GuildID:   m.GuildID,
		UserID:    m.Member.User.ID,
		Type:      "rr_set_embed_message",
		Next:      "rr_add_roles_or_finish",
		Data:      npdb,
		CreatedAt: time.Now(),
		ExpiredAt: time.Now().Add(time.Minute * 5),
	}

	err, ec := c.ps.Create(nPen)

	if err != nil {
		switch ec {
		case -1:
			log.Printf("err while querying for current pending of %v (%v) in %v (%v) for rr_set_embed_message: %v", m.Author.Username, m.Author.ID, utils.GetGuildNameFromState(s, m.GuildID), m.GuildID, err)
		case 1:
			log.Printf("err user of %v (%v) in %v (%v) already have pending interaction for rr_set_embed_message: %v", m.Author.Username, m.Author.ID, utils.GetGuildNameFromState(s, m.GuildID), m.GuildID, err)
		case 2:
			log.Printf("err while creating pending row for %v (%v) in %v (%v) for rr_set_embed_message: %v", m.Author.Username, m.Author.ID, utils.GetGuildNameFromState(s, m.GuildID), m.GuildID, err)
		}

		s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
			Content:   "Error occured! Please try running `/reactionroles` again.",
			Reference: m.MessageReference,
		})

		return
	}
}
