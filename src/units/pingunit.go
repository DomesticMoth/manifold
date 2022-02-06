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
package units

import (
	"time"
	"context"
	"golang.org/x/sync/errgroup"
    "github.com/DomesticMoth/manifold/src/id"
    log "github.com/sirupsen/logrus"
    "github.com/DomesticMoth/manifold/src/events"
)


type PingUnitConfig struct {}

type PingUnit struct{
    context *UnitContext
}

func NewPingUnit(config PingUnitConfig) PingUnit {
	return PingUnit{}
}

func (l *PingUnit) Init(context *UnitContext) error {
    l.context = context
    return nil
}

func (p *PingUnit) Run(group *errgroup.Group, ctx context.Context) error {
    defer p.context.Close()
    for {
    	select {
    		case <- ctx.Done():
    			log.Debug(p.context.GetName(), " Stopping")
    			return nil
    		case event := <- p.context.Receiver():
		        if event.Msgevent != nil {
		        	if event.Msgevent.Text == "ping" {
		        		msgid, err := id.NewID()
						authorid, err := id.NewID()
						if err != nil { return err }
						msgevent := events.MsgEvent{
							MsgId: msgid,
							AuthorId: authorid,
							AuthorName: "[BOT] "+p.context.GetName(),
							CreateTime: time.Now().Unix(),
							RedactTime: time.Now().Unix(),
							Text: "pong",
							Images: []events.Image{},
							Messages: []events.MsgEvent{*event.Msgevent},
						}
						event := events.Event{
							Tags: []string{events.ALLTAG},
							Msgevent: &msgevent,
							Deletemsgevent: nil,
							Userevent: nil,
						}
						p.context.Sender() <- event
		        	} else if event.Msgevent.Text == "pong" {
		        		msgid, err := id.NewID()
						authorid, err := id.NewID()
						if err != nil { return err }
						msgevent := events.MsgEvent{
							MsgId: msgid,
							AuthorId: authorid,
							AuthorName: "[BOT] "+p.context.GetName(),
							CreateTime: time.Now().Unix(),
							RedactTime: time.Now().Unix(),
							Text: "ping",
							Images: []events.Image{},
							Messages: []events.MsgEvent{*event.Msgevent},
						}
						event := events.Event{
							Tags: []string{events.ALLTAG},
							Msgevent: &msgevent,
							Deletemsgevent: nil,
							Userevent: nil,
						}
						p.context.Sender() <- event
		        	}
		        }
    	}
    }
}

func (l *PingUnit) Stop() error { return nil }
