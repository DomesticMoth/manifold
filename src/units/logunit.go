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
    log "github.com/sirupsen/logrus"
)


type LogUnit struct{
    context UnitContext
}

func (l *LogUnit) Init(context UnitContext) error {
    l.context = context
    return nil
}

func (l *LogUnit) Run() error {
    defer l.context.Close()
    for {
        event, err := l.context.GetEvent()
        if err != nil {
            return err
        }
        log.Info("Event: ", event.ToString())
    }
}

func (l *LogUnit) Stop() error { return nil }
