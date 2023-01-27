package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ebobo/dss_checklist/pkg/server"
	"github.com/jessevdk/go-flags"
)

var opt struct {
	HTTPAddr string `short:"h" long:"http-addr" default:":9099" description:"http listen address" required:"yes"`
}

func main() {
	_, err := flags.ParseArgs(&opt, os.Args)
	if err != nil {
		log.Fatalf("error parsing flags: %v", err)
	}

	server := server.New(server.Config{
		HTTPListenAddr: opt.HTTPAddr,
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
