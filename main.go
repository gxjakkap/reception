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

package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/gxjakkap/reception/commands"
	"github.com/gxjakkap/reception/db"
	"github.com/gxjakkap/reception/events"
	"github.com/gxjakkap/reception/store"
	"github.com/gxjakkap/reception/utils"
	"github.com/joho/godotenv"
)

var (
	removeCommand     bool
	enableHealthCheck bool
	token             string
	s                 *discordgo.Session
)

func init() {
	var err error

	err = godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	removeCommand = utils.StringToBool(utils.GetEnv("REM_CMD", "true"))
	enableHealthCheck = utils.StringToBool(utils.GetEnv("HEALTHCHECK", "false"))
	token = utils.GetEnv("TOKEN", "")

	s, err = discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
}

func main() {
	// init db
	d := db.New()
	gs := store.NewGuildsStore(d)
	ps := store.NewPendingStore(d)

	e := events.NewEventContext(gs, ps)

	s.AddHandler(e.Ready)

	err := s.Open()

	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	log.Println("Registering commands: ")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands.Infos))
	for i, v := range commands.Infos {
		log.Printf("Registering '%s'", v.Name)
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, "", v)
		if err != nil {
			log.Panicf("Error registering '%v': %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	// add event handler

	s.AddHandler(e.GuildJoin)
	s.AddHandler(e.GuildDelete)

	c := commands.NewStoreCtx(gs, ps)
	ch := c.GetHandlers()

	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		e.Interactions(s, i, ch)
	})

	if enableHealthCheck {
		mux := http.NewServeMux()
		mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("pong!"))
			} else {
				w.WriteHeader(http.StatusMethodNotAllowed)
				w.Write([]byte(""))
			}
		})

		sv := http.Server{
			Addr:    ":8000",
			Handler: mux,
		}

		sv.ListenAndServe()

		go func() {
			log.Println("Starting Healthcheck server on :8000")
			if err := sv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("Healthcheck server error: %v", err)
			}
		}()
	}

	defer s.Close()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	log.Println("Press Ctrl+C to exit")
	<-stop

	// clean up
	if removeCommand {
		log.Println("Removing commands:")

		for _, v := range registeredCommands {
			log.Printf("Removing '%s'", v.Name)
			err := s.ApplicationCommandDelete(s.State.User.ID, "", v.ID)
			if err != nil {
				log.Panicf("Cannot remove '%v': %v", v.Name, err)
			}
		}
	}
}
