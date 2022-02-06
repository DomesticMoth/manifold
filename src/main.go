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
	"os"
	"os/signal"
	"syscall"
	"strings"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
	"github.com/DomesticMoth/manifold/src/units"
	"github.com/DomesticMoth/manifold/src/units/vk"
	"github.com/DomesticMoth/manifold/src/units/tg"
	"golang.org/x/sync/errgroup"
	"context"
)

func ConfigLogger(config Config){
	switch strings.ToLower(config.LogLevel) {
		case "trace":
			log.SetLevel(log.TraceLevel)
		case "debug":
			log.SetLevel(log.DebugLevel)
		case "info":
			log.SetLevel(log.InfoLevel)
		case "warn":
			log.SetLevel(log.WarnLevel)
		case "error":
			log.SetLevel(log.ErrorLevel)
		case "fatal":
			log.SetLevel(log.FatalLevel)
		case "panic":
			log.SetLevel(log.PanicLevel)
		default:
			log.WithFields(log.Fields{
			  "level": config.LogLevel,
			}).Fatal("Unknown logging level in config")
	}
}

func main() {
	log.Info("Starting")
	config, err := LoadConfig()
	if err != nil { log.Fatal(err) }
	ConfigLogger(config)
	log.Debug("Logger configured")
	db, err := sql.Open("sqlite3", config.Db)
	if err != nil { log.Fatal(err) }
	defer db.Close()
	log.Debug("Connected to database")
	unitsList := []units.Unit{}
	contextsList := []*units.UnitContext{}
	unitContextBuilder := units.NewUnitContextBuilder(db, config.EventsChanSize)
	for _, uc := range config.Unit {
		unitContext := unitContextBuilder.Build(
													uc.Name, 
													uc.BlockListInternal, 
													uc.BlockListExternal,
		)
		var unit units.Unit
		if uc.Log != nil{
			logunit := units.NewLogUnit(*uc.Log)
			unit = &logunit
		}else if uc.Ping != nil{
			pingunit := units.NewPingUnit(*uc.Ping)
			unit = &pingunit
		}else if uc.Vk != nil{
			vkuint := vk.NewVkUnit(*uc.Vk)
			unit = &vkuint
		}else if uc.Tg != nil{
			tguint := tg.NewTgUnit(*uc.Tg)
			unit = &tguint
		}else{
			db.Close()
			log.Fatal("Unknown unit type")
		}
		err := unit.Init(unitContext)
		if err != nil { log.Fatal(err) }
		unitsList = append(unitsList, unit)
		contextsList = append(contextsList, unitContext)
		log.Debug("Created unit '", uc.Name, "'")
	}
	for _, c := range contextsList {
		conext := c
		go conext.Run()
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	group, ctx := errgroup.WithContext(ctx)
	group.Go(func() error{
		sigc := make(chan os.Signal, 1)
		signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
		<- sigc
		log.Info("Stopping")
		cancel()
		return nil
	})
	for _, u := range unitsList {
		unit := u
		group.Go(func() error{
			return unit.Run(group, ctx)
		})
	}
	err = group.Wait()
	for _, u := range unitsList {
		unit := u
		unit.Stop()
	}
	if err != nil {
		db.Close()
		log.Fatal(err)
	}
}
