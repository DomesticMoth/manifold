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
	"errors"
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

type Puppet struct{
	Vk api.VK
	PeerId int
}

type PuppetConfig struct{
	Token string
	PeerId int
}

type IdComparisonInc struct{
	Vk int
	Local id.Id
}

type ComparisonOutg struct{
	Local id.Id
	Puppet int
}

type VkUnitConfig struct{
	Token string
	PeerId int
	UsersInc []IdComparisonInc
	UsersOutg []ComparisonOutg
	Puppet []PuppetConfig
}

type VkUnit struct{
	config VkUnitConfig
	ucontxt *units.UnitContext
	vk *api.VK
	group *api.GroupsGetByIDResponse
	lp **longpoll.LongPoll
	incoming chan vkevents.MessageNewObject
	usersIncoming map[int]id.Id
	puppets []Puppet
	comparisonOutg map[id.Id]*Puppet
}

func NewVkUnit(config VkUnitConfig) VkUnit{
	usersIncoming := make(map[int]id.Id)
	comparisonOutg := make(map[id.Id]*Puppet)
	puppets := []Puppet{}
	for _, v := range config.UsersInc {
		usersIncoming[v.Vk] = v.Local
	}
	return VkUnit{
		config,
		nil,
		nil,
		nil,
		nil,
		make(chan vkevents.MessageNewObject),
		usersIncoming,
		puppets,
		comparisonOutg,
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

	for _, p := range v.config.Puppet{
		puppet := Puppet{
			*api.NewVK(p.Token),
			p.PeerId,
		}
		v.puppets = append(v.puppets, puppet)
	}
	for _, u := range v.config.UsersOutg{
		v.comparisonOutg[u.Local] = &v.puppets[u.Puppet]
	}
	
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
	if len(users) < 1 {
		return "", errors.New("B0t")
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

func (v *VkUnit) GetIncId(vkid int) (id.Id, error) {
	if id, ok := v.usersIncoming[vkid]; ok {
	    return id, nil
	}
	return id.NewID()
}

func (v *VkUnit) GetApi(local id.Id) *api.VK {
	if cmp, ok := v.comparisonOutg[local]; ok {
	    return &cmp.Vk
	}
	return v.vk
}

func (v *VkUnit) GetPeerId(local id.Id) int {
	if cmp, ok := v.comparisonOutg[local]; ok {
	    return cmp.PeerId
	}
	return v.config.PeerId
}

func (v *VkUnit) GetHeader(msg events.MsgEvent) string{
	if _, ok := v.comparisonOutg[msg.AuthorId]; ok {
		return ""
	}
	return "< BY '"+msg.AuthorName+"' >\n"
}

func (v *VkUnit) Run(group *errgroup.Group, ctx context.Context) error {
	go (*v.lp).Run()
	b := params.NewMessagesSendBuilder()
	b.Message("Bridge is turned on")
	b.RandomID(0)
	b.PeerID(v.config.PeerId)
	_, err := v.vk.MessagesSend(b.Params)
	if err != nil { log.Error(err) }
	for {
		select {
			case <- ctx.Done():
				log.Debug(v.ucontxt.GetName(), " Stopping")
				b := params.NewMessagesSendBuilder()
				b.Message("Bridge is turned off")
				b.RandomID(0)
				b.PeerID(v.config.PeerId)
				v.vk.MessagesSend(b.Params)
				return nil
			case msg := <- v.incoming:
				msgid, err := id.NewID()
				if err != nil { return err }
				authorid, err := v.GetIncId(msg.Message.FromID)
				if err != nil { return err }
				name, err := v.getName(msg.Message.FromID)
				if err != nil { continue }
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
				text := v.GetHeader(*event.Msgevent)+event.Msgevent.Text
				b.Message(text)
				b.RandomID(0)
				b.PeerID(v.GetPeerId(event.Msgevent.AuthorId))
				_, err := v.GetApi(event.Msgevent.AuthorId).MessagesSend(b.Params)
				if err != nil { log.Error(err) }
		}
	}
}

func (v *VkUnit) Stop() error {
	(*v.lp).Shutdown()
	v.ucontxt.Close()
	return nil
}
