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

type IdComparisonInc struct{
	Tg int64
	Local id.Id
}

type Puppet = tele.Bot

type ComparisonOutg struct{
	Local id.Id
	Puppet int
}

type PuppetConfig struct {
	Token string
}

type TgUnitConfig struct {
	Token string
	ChatId int64
	UsersInc []IdComparisonInc
	Puppet []PuppetConfig
	UsersOutg []ComparisonOutg
}

type TgUnit struct {
	config TgUnitConfig
	ucontext *units.UnitContext
	bot *tele.Bot
	chat *tele.Chat
	incoming chan tele.Context
	puppets []Puppet
	comparisonOutg map[id.Id]*Puppet
	usersIncoming map[int64]id.Id
}

func NewTgUnit(config TgUnitConfig) TgUnit{
	incoming := make(chan tele.Context)
	comparisonOutg := make(map[id.Id]*Puppet)
	usersIncoming := make(map[int64]id.Id)
	return TgUnit{ config, nil, nil, nil, incoming, []Puppet{}, comparisonOutg, usersIncoming }
}

func (t *TgUnit) Init(ucontext *units.UnitContext) error {
	t.ucontext = ucontext

	for _, pair := range t.config.UsersInc {
		t.usersIncoming[pair.Tg]=pair.Local
	}
	
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

	for _, pc := range t.config.Puppet{
		pref := tele.Settings{
			Token:  pc.Token,
			Poller: &tele.LongPoller{Timeout: 10 * time.Second},
		}
		puppet, err := tele.NewBot(pref)
		if err != nil { return err }
		t.puppets = append(t.puppets, *puppet)
	}

	for _, cmp := range t.config.UsersOutg {
		t.comparisonOutg[cmp.Local] = &t.puppets[cmp.Puppet]
	}

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

func (t *TgUnit) Bot(local id.Id) *tele.Bot {
	if bot, ok := t.comparisonOutg[local]; ok {
	    return bot
	}
	return t.bot
}

func (t *TgUnit) Header(msg events.MsgEvent) string {
	if _, ok := t.comparisonOutg[msg.AuthorId]; ok {
	    return ""
	}
	return "< BY '"+msg.AuthorName+"' >\n"
}

func (t *TgUnit) GetIncId(tgid int64) (id.Id, error) {
	if id, ok := t.usersIncoming[tgid]; ok {
	    return id, nil
	}
	return id.NewID()
}

func (t *TgUnit) Run(group *errgroup.Group, ctx context.Context) error {
	go t.bot.Start()
	_, err := t.bot.Send(t.chat, "Bridge is turned on")
	if err != nil { log.Error(err) }
	for {
		select{
			case <- ctx.Done():
				log.Debug(t.ucontext.GetName(), " Stopping")
				t.bot.Send(t.chat, "Bridge is turned off")
				return nil
			case inc := <- t.incoming:
				msgid, err := id.NewID()
				if err != nil { return err }
				authorid, err := t.GetIncId(inc.Sender().ID)
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
				text := t.Header(*event.Msgevent)+event.Msgevent.Text
				_, err := t.Bot(event.Msgevent.AuthorId).Send(t.chat, text)
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
