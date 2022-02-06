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
package vk

import (
	"context"
	"github.com/DomesticMoth/manifold/src/id"
	"github.com/DomesticMoth/manifold/src/events"
	"github.com/DomesticMoth/manifold/src/units"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	vkevents "github.com/SevereCloud/vksdk/v2/events"
	"github.com/SevereCloud/vksdk/v2/longpoll-bot"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

type VkUnitConfig struct{
	Token string
	PeerId int
}

type VkUnit struct{
	config VkUnitConfig
	ucontxt *units.UnitContext
	vk *api.VK
	group *api.GroupsGetByIDResponse
	lp **longpoll.LongPoll
	incoming chan vkevents.MessageNewObject
}

func NewVkUnit(config VkUnitConfig) VkUnit{
	return VkUnit{
		config,
		nil,
		nil,
		nil,
		nil,
		make(chan vkevents.MessageNewObject),
	}
}

func (v *VkUnit) initDb() error {
	_ = v.ucontxt.GetDb()
	defer v.ucontxt.RetDb()
	// TODO create table "vk msg id to local id"
	// TODO create table "vk user id to local id"
	return nil
}

func (v *VkUnit) Init(ucontxt *units.UnitContext) error {
	v.ucontxt = ucontxt
	err := v.initDb()
	if err != nil { return err }
	v.vk = api.NewVK(v.config.Token)
	group, err := v.vk.GroupsGetByID(nil)
	if err != nil { return err }
	v.group = &group
	lp, err := longpoll.NewLongPoll(v.vk, group[0].ID)
	if err != nil { return err }
	v.lp = &lp
	lp.MessageNew(func(_ context.Context, obj vkevents.MessageNewObject) {
		if obj.Message.PeerID == v.config.PeerId {
			v.incoming <- obj
		}else{
			log.Debug("Msg from ", obj.Message.PeerID)
		}
	})
	return nil
}

func (v *VkUnit) getName(ID int) (string, error) {
	users, err := v.vk.UsersGet(api.Params{
		"user_ids": ID,
	})
	if err != nil {
		return "", err
	}
	ret := ""
	ret += users[0].FirstName
	if ret != "" { ret += " "}
	ret += users[0].LastName
	if ret == "" {
		ret += users[0].Nickname
	}
	ret = "[VK] "+ret
	return ret, nil
}

func (v *VkUnit) Run(group *errgroup.Group, ctx context.Context) error {
	go (*v.lp).Run()
	for {
		select {
			case <- ctx.Done():
				log.Debug(v.ucontxt.GetName(), " Stopping")
				return nil
			case msg := <- v.incoming:
				msgid, err := id.NewID()
				if err != nil { return err }
				authorid, err := id.NewID()
				if err != nil { return err }
				name, err := v.getName(msg.Message.FromID)
				if err != nil { return err }
				msgevent := events.MsgEvent{
					MsgId: msgid,
					AuthorId: authorid,
					AuthorName: name,
					CreateTime: int64(msg.Message.Date),
					RedactTime: int64(msg.Message.Date),
					Text: msg.Message.Text,
					Images: []events.Image{},
					Messages: []events.MsgEvent{},
				}
				event := events.Event{
					Tags: []string{events.ALLTAG},
					Msgevent: &msgevent,
					Deletemsgevent: nil,
					Userevent: nil,
				}
				v.ucontxt.Sender() <- event
			case event := <- v.ucontxt.Receiver():
				if event.Msgevent == nil { continue }
				b := params.NewMessagesSendBuilder()
				text := "< BY '"+event.Msgevent.AuthorName+"' >\n"+event.Msgevent.Text
				b.Message(text)
				b.RandomID(0)
				b.PeerID(v.config.PeerId)
				_, err := v.vk.MessagesSend(b.Params)
				if err != nil { return err }
		}
	}
	return nil
}

func (v *VkUnit) Stop() error {
	v.ucontxt.Close()
	return nil
}
