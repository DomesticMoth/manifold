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
package main

import (
	"fmt"
	"github.com/DomesticMoth/confer"
	"github.com/creasty/defaults"
	"github.com/DomesticMoth/manifold/src/id"
	"github.com/DomesticMoth/manifold/src/events"
	"github.com/DomesticMoth/manifold/src/units"
	"github.com/DomesticMoth/manifold/src/units/vk"
	"github.com/DomesticMoth/manifold/src/units/tg"
)

const DEFAULT_GLOBAL_PATH string = "/etc/manifold/config.toml"

type UnitConfig struct {
	Name string
	BlockListInternal []id.Id
	BlockListExternal []id.Id
	Log *units.LogUnitConfig
	Ping *units.PingUnitConfig
	Vk *vk.VkUnitConfig
	Tg *tg.TgUnitConfig
}

type Config struct {
	LogLevel string			`default:"Info"`
	Db string				`default:"/etc/manifold/manifold.db"`
	EventsChanSize int		`default:"100"`
	BlockList []id.Id
	Unit []UnitConfig
}

func LoadConfig() (Config, error) {
	var conf Config
	if err := defaults.Set(&conf); err != nil {
		return conf, err
	}
	err := confer.LoadConfig([]string{DEFAULT_GLOBAL_PATH}, &conf)
	for _, unit := range conf.Unit {
		for _, ID := range conf.BlockList {
			unit.BlockListInternal = append(unit.BlockListInternal, ID)
			unit.BlockListExternal = append(unit.BlockListExternal, ID)
		}
	}
	for a, unitA := range conf.Unit {
		if unitA.Name == events.ALLTAG {
			err = fmt.Errorf("You cannot use '%s' as the unit name.", events.ALLTAG)
			return conf, err
		}
		unitSpecified := 0
		if unitA.Log != nil { unitSpecified += 1}
		if unitA.Ping != nil { unitSpecified += 1}
		if unitA.Vk != nil { unitSpecified += 1}
		if unitA.Tg != nil { unitSpecified += 1}
		if unitSpecified < 1 {
			err = fmt.Errorf("Unit '%s' type must be specified.", unitA.Name)
			return conf, err
		}
		if unitSpecified > 1 {
			err = fmt.Errorf("Unit '%s' should have only one type.", unitA.Name)
			return conf, err
		}
		for b, unitB := range conf.Unit {
			if a == b { continue }
			if unitA.Name == unitB.Name {
				err = fmt.Errorf("More than one unit has the same name '%s'. Unit names must be unique.", unitA.Name)
				return conf, err
			}
		}
	}
	return conf, err
}
