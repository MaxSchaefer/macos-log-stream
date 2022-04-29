package main

import (
	"flag"
	"fmt"
	"github.com/MaxSchaefer/macos-log-stream/pkg/mls"
	"os"
	"os/signal"
	"syscall"
)

type Cmd struct {
	fs        *flag.FlagSet
	Predicate string
}

func NewCmd() *Cmd {
	cmd := &Cmd{
		fs: flag.NewFlagSet("cmd", flag.ExitOnError),
	}
	cmd.fs.StringVar(&cmd.Predicate, "predicate", "", `
take a look at '> man log' or visit
https://eclecticlight.co/2016/10/17/log-a-primer-on-predicates/
`)
	return cmd
}

func (c *Cmd) Exec(args []string) error {
	if err := c.fs.Parse(args); err != nil {
		return err
	}
	return c.run()
}

func (c *Cmd) run() error {
	logs := mls.NewLogs()

	logs.Predicate = c.Predicate

	signalChannel := make(chan os.Signal)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalChannel
		logs.StopGathering()
		os.Exit(0)
	}()

	if err := logs.StartGathering(); err != nil {
		panic(err)
	}

	for log := range logs.Channel {
		fmt.Println(log.EventMessage)
	}

	return nil
}
