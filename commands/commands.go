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
	"github.com/bwmarrin/discordgo"
	"github.com/gxjakkap/reception/store"
)

type StoreCtx struct {
	gs *store.GuildsStore
	ps *store.PendingStore
}

func NewStoreCtx(gs *store.GuildsStore, ps *store.PendingStore) *StoreCtx {
	return &StoreCtx{
		gs: gs,
		ps: ps,
	}
}

func (c *StoreCtx) GetHandlers() map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		Ping.Name:          c.PingCommandHandler,
		ReactionRoles.Name: c.ReactionRolesCommandHandler,
		Welcome.Name:       c.WelcomeCommandHandler,
		TempC.Name:         c.TempCCommandHandler,
		CreateCTemp.Name:   c.CreateCTempCommandHandler,
		ListCTemp.Name:     c.ListCTempCommandHandler,
	}
}

var Infos = []*discordgo.ApplicationCommand{
	Ping,
	ReactionRoles,
	Welcome,
	TempC,
	CreateCTemp,
	ListCTemp,
}
