package main

import (
	"flag"
	"github.com/lucas-clemente/quic-go/internal/utils"
	"io/ioutil"
	"fmt"
)

const DEFAULT_ADDR = "0.0.0.0:4443"

func usage() {
	fmt.Println("Usage: quicperf [-s [host] -F seed_file |-c host] [options]")
}

func main() {
	verbose := flag.Bool("v", false, "verbose")
	s := flag.Bool("s", false, "Run server")
	c := flag.String("c", "", "Run client")
	streams := flag.Int("S", 1, "Number of streams")
	addr := flag.String("bind", DEFAULT_ADDR, "bind to")
	file := flag.String("F", "", "Source data file")
	ticks := flag.Int("t", 10, "Time in seconds to transmit for (default 10 secs)")
	flag.Parse()

	logger := utils.DefaultLogger

	if *verbose {
		logger.SetLogLevel(utils.LogLevelDebug)
	} else {
		logger.SetLogLevel(utils.LogLevelInfo)
	}
	logger.SetLogTimeFormat("")

	if *s {
		if *file == "" {
			usage()
			return
		}
		data, err := ioutil.ReadFile(*file)
		if err != nil {
			panic(err)
		}
		logger.Infof("Starting server...")
		serverMain(*addr, data, *streams, *ticks)
	} else {
		if *c == "" {
			usage()
			return
		}
		logger.Infof("Starting client...")
		clientMain(*c)
	}
}
