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
	"database/sql"
	"context"
	"github.com/DomesticMoth/manifold/src/id"
	"github.com/DomesticMoth/manifold/src/events"
	"golang.org/x/sync/errgroup"
)

type StrSlice []string

func (list StrSlice) Has(a string) bool {
    for _, b := range list {
        if b == a {
            return true
        }
    }
    return false
}

type UnitContextBuilder struct{
	dblock chan *sql.DB
	eventsChanSize int
	contexts []*UnitContext
}

func NewUnitContextBuilder(db *sql.DB, eventsChanSize int) UnitContextBuilder{
	dblock := make(chan *sql.DB, 1)
	dblock <- db
	return UnitContextBuilder{
		dblock,
		eventsChanSize,
		[]*UnitContext{},
	}
}

func (b *UnitContextBuilder) Build(Name string, BlockListInternal, BlockListExternal []id.Id) *UnitContext{
	context := NewUnitContext(
								Name,
								BlockListInternal,
								BlockListExternal,
								b.dblock,
								b.eventsChanSize,
	)
	for _, ctx := range b.contexts {
		context.Bind(ctx)
		ctx.Bind(&context)
	}
	ret := &context
	b.contexts = append(b.contexts, ret)
	return ret
}

type UnitContext struct{
	name string
	dblock chan *sql.DB
	db *sql.DB
	incomingRaw events.EventChan
	incoming events.EventChan
	outgoingRaw events.EventChan
	outgoing []events.EventChan
	outgoingNames []string
	userIdDoNotGet []id.Id
	userIdDoNotSend []id.Id
}

func NewUnitContext(name string, 
					userIdDoNotGet []id.Id, 
					userIdDoNotSend []id.Id, 
					dblock chan *sql.DB,
					EventsChanSize int) UnitContext{
	incomingRaw := make(events.EventChan, EventsChanSize)
	incoming := make(events.EventChan)
	outgoingRaw := make(events.EventChan)
	return UnitContext{
		name,
		dblock,
		nil,
		incomingRaw,
		incoming,
		outgoingRaw,
		[]events.EventChan{},
		[]string{},
		userIdDoNotGet,
		userIdDoNotSend,
	}
}

func (c *UnitContext) Bind(other *UnitContext) {
	name := other.GetName()
	if StrSlice(c.outgoingNames).Has(name) {return}
	c.outgoingNames = append(c.outgoingNames, name)
	c.outgoing = append(c.outgoing, other.incomingRaw)
}

func (c *UnitContext) filterInc() {
	for {
		event := <- c.incomingRaw
		if !StrSlice(event.Tags).Has(events.ALLTAG) {
			if !StrSlice(event.Tags).Has(c.name) { continue }
		}
		if event.Msgevent != nil {
			if id.IdSlice(c.userIdDoNotGet).Has(event.Msgevent.AuthorId) {
				event.Msgevent = nil
			}
		}
		if event.Msgevent == nil && event.Deletemsgevent == nil && event.Userevent == nil {
			continue
		}
		c.incoming <- event
	}
}

func (c *UnitContext) filterOut() {
	for {
		event := <- c.outgoingRaw
		if event.Msgevent != nil {
			if id.IdSlice(c.userIdDoNotSend).Has(event.Msgevent.AuthorId) {
				event.Msgevent = nil
			}
		}
		if event.Msgevent == nil && event.Deletemsgevent == nil && event.Userevent == nil {
			continue
		}
		for _, ev := range c.outgoing {
			ev <- event
		}
	}
}

func (c *UnitContext) Run() {
	go c.filterOut()
	c.filterInc()
	return
}

func (c *UnitContext) Sender() events.EventChan {
	return c.outgoingRaw
}

func (c *UnitContext) Receiver() events.EventChan {
	return c.incoming
}

func (c *UnitContext) GetDb() *sql.DB {
	db := <- c.dblock
	c.db = db
	return db
}

func (c *UnitContext) RetDb() {
	if c.db != nil {
		c.dblock <- c.db
	}
	c.db = nil
}

func (c *UnitContext) GetName() string {
	return c.name
}

func (c *UnitContext) Close() {
	c.RetDb()
}

type Unit interface{
	Init(context *UnitContext) error
	Run(group *errgroup.Group, ctx context.Context) error
	Stop() error
}
