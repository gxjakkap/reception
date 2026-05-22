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

import "github.com/bwmarrin/discordgo"

type FollowUp struct {
	Name   string
	Action func(ctx *FollowUpsCtx, s *discordgo.Session, m *discordgo.MessageCreate, id uint)
}

var FollowUps = map[string]*FollowUp{
	RRAddRolesOrFinish.Name:     RRAddRolesOrFinish,
	RRSetEmbedMessage.Name:      RRSetEmbedMessage,
	CCTAddChannelsOrFinish.Name: CCTAddChannelsOrFinish,
	MsgaSetMessage.Name:         MsgaSetMessage,
}
