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
	"github.com/DomesticMoth/manifold/src/id"
	"github.com/DomesticMoth/manifold/src/events"
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

type UnitContext struct{
	name string
	dblock chan sql.DB
	db *sql.DB
	ev events.EventChan
	userIdDoNotGet []id.Id
	userIdDoNotSend []id.Id
}

func (c *UnitContext) GetEvent() (events.Event, error) {
	for {
		// TODO Check context
		// TODO Check chan valid
		event := <- c.ev
		if !StrSlice(event.Tags).Has("<<ALL>>") {
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
		return event, nil
	}
}

func (c *UnitContext) SendEvent(event events.Event) error {
	// TODO Check context
	// TODO Check chan valid
	if event.Msgevent != nil {
		if id.IdSlice(c.userIdDoNotSend).Has(event.Msgevent.AuthorId) {
			event.Msgevent = nil
		}
	}
	if event.Msgevent == nil && event.Deletemsgevent == nil && event.Userevent == nil {
		return nil
	}
	c.ev <- event
	return nil
}

func (c *UnitContext) GetDb() (*sql.DB, error) {
	// TODO Check context
	// TODO Check chan valid
	db := <- c.dblock
	c.db = &db
	return &db, nil
}

func (c *UnitContext) RetDb() error {
	if c.db != nil {
		c.dblock <- *c.db
	}
	return nil
}

func (c *UnitContext) GetName() string {
	return c.name
}

func (c *UnitContext) Close() {
	c.RetDb()
}

type Unit interface{
	Init(context UnitContext) error
	Run() error
	Stop() error
}
