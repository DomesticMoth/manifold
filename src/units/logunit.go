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
	"context"
    log "github.com/sirupsen/logrus"
    "golang.org/x/sync/errgroup"
)


type LogUnitConfig struct{}

type LogUnit struct{
    context *UnitContext
    LogLevel string
}

func NewLogUnit(config LogUnitConfig) LogUnit {
	return LogUnit{}
}

func (l *LogUnit) Init(context *UnitContext) error {
    l.context = context
    return nil
}

func (l *LogUnit) Run(group *errgroup.Group, ctx context.Context) error {
    defer l.context.Close()
    for {
    	select {
    		case <- ctx.Done():
    			log.Debug(l.context.GetName(), " Stopping")
    			return nil
    		case event := <- l.context.Receiver():
		        s, err := event.ToString()
		        if err != nil { return err }
		        log.Info("Event: ", s)
    	}
    }
}

func (l *LogUnit) Stop() error { return nil }
