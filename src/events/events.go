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
package events

import (
	Id "github.com/DomesticMoth/manifold/src/id"
)

type EventChan chan Event

type Image struct{
	Url string
}

type Event struct {
	Tags []string
	Msgevent *MsgEvent
	Deletemsgevent *DeleteMsgEvent
	Userevent *UserEvent
}

type MsgEvent struct {
	MsgId Id.Id
	AuthorId Id.Id
	AuthorName string
	CreateTime int64
	RedactTime int64
	Text string
	Images []Image
	Messages []MsgEvent
}

type DeleteMsgEvent struct {
	MsgId Id.Id
}

type UserEvent struct {
	ExecutorId Id.Id
	ExecutorName string
	TargetId Id.Id
	TargetName string
}
