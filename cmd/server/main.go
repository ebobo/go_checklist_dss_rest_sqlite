package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ebobo/dss_checklist/pkg/model"
	"github.com/ebobo/dss_checklist/pkg/server"
	sqlitestore "github.com/ebobo/dss_checklist/pkg/store/sqlite"
	"github.com/jessevdk/go-flags"
)

var opt struct {
	HTTPAddr   string `short:"h" long:"http-addr" default:":9099" description:"http listen address" required:"yes"`
	SqliteFile string `long:"sqlite-file" env:"SQLITE_FILE" default:"items.db" description:"sqlite file"`
}

func main() {
	_, err := flags.ParseArgs(&opt, os.Args)
	if err != nil {
		log.Fatalf("error parsing flags: %v", err)
	}

	db, created, err := sqlitestore.New(opt.SqliteFile)
	if err != nil {
		log.Fatalf("error connect to sqlite: %v", err)
	}

	//some test data
	if created {
		db.AddItem(model.Item{ID: "01", Name: "Door-A", Position: 0, Tag: "door-a", Status: false})
		db.AddItem(model.Item{ID: "02", Name: "Door-B", Position: 1, Tag: "door-b", Status: false})
		db.AddItem(model.Item{ID: "03", Name: "Door-C", Position: 2, Tag: "door-c", Status: false})
		db.AddItem(model.Item{ID: "04", Name: "Door-D", Position: 3, Tag: "door-d", Status: false})
	} else {
		log.Println("db already exists")
	}

	server := server.New(server.Config{
		HTTPListenAddr: opt.HTTPAddr,
		DB:             db,
	})

	e := server.Start()
	if e != nil {
		log.Fatalf("error starting server: %v", e)
	}

	// Block forever
	// Capture Ctrl-C
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	server.Shutdown()
}
