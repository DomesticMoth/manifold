/*
    This file is part of manifold.
    Manifold is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.
    Manifold is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.
    You should have received a copy of the GNU General Public License
    along with manifold.  If not, see <https://www.gnu.org/licenses/>.
*/
package tg

import (
	"context"
	"golang.org/x/sync/errgroup"
	"github.com/DomesticMoth/manifold/src/id"
	"github.com/DomesticMoth/manifold/src/events"
	"github.com/DomesticMoth/manifold/src/units"
	tele "gopkg.in/telebot.v3"
	"time"
	log "github.com/sirupsen/logrus"
)

type TgUnitConfig struct {
	Token string
	ChatId int64
}

type TgUnit struct {
	config TgUnitConfig
	ucontext *units.UnitContext
	bot *tele.Bot
	chat *tele.Chat
	incoming chan tele.Context
}

func NewTgUnit(config TgUnitConfig) TgUnit{
	incoming := make(chan tele.Context)
	return TgUnit{ config, nil, nil, nil, incoming }
}

func (t *TgUnit) Init(ucontext *units.UnitContext) error {
	t.ucontext = ucontext
	pref := tele.Settings{
		Token:  t.config.Token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil { return err }
	t.bot = b

	chat, err:= t.bot.ChatByID(t.config.ChatId)
	if err != nil { return err }
	t.chat = chat

	t.bot.Handle(tele.OnText, func(c tele.Context) error {
		if c.Chat().ID == t.config.ChatId {
			t.incoming <- c
		}else{
			log.Debug(c.Chat().ID)
		}
		return nil
	})

	return nil
}

func GetUserName(user tele.User) string {
	ret := ""
	ret += user.FirstName
	if ret != "" { ret += " "}
	ret += user.LastName
	if ret == "" {
		ret += user.Username
	}
	ret = "[Tg] "+ret
	return ret
}

func (t *TgUnit) Run(group *errgroup.Group, ctx context.Context) error {
	go t.bot.Start()
	// c.Send(c.Text())
	for {
		select{
			case <- ctx.Done():
				log.Debug(t.ucontext.GetName(), " Stopping")
				return nil
			case inc := <- t.incoming:
				msgid, err := id.NewID()
				if err != nil { return err }
				authorid, err := id.NewID()
				if err != nil { return err }
				name := GetUserName(*inc.Sender())
				date := time.Now().Unix()
				msgevent := events.MsgEvent{
					MsgId: msgid,
					AuthorId: authorid,
					AuthorName: name,
					CreateTime: date,
					RedactTime: date,
					Text: inc.Text(),
					Images: []events.Image{},
					Messages: []events.MsgEvent{},
				}
				event := events.Event{
					Tags: []string{events.ALLTAG},
					Msgevent: &msgevent,
					Deletemsgevent: nil,
					Userevent: nil,
				}
				t.ucontext.Sender() <- event
			case event := <- t.ucontext.Receiver():
				if event.Msgevent == nil { continue }
				text := "< BY '"+event.Msgevent.AuthorName+"' >\n"+event.Msgevent.Text
				_, err := t.bot.Send(t.chat, text)
				if err != nil { log.Error(err) }
		}
	}
	return nil
}

func (t *TgUnit) Stop() error {
	t.bot.Stop()
	t.ucontext.Close()
	return nil
}
