package main

import (
	"flag"
	"github.com/lucas-clemente/quic-go/internal/utils"
	"io/ioutil"
)

func main() {
	verbose := flag.Bool("v", false, "verbose")
	s := flag.Bool("s", false, "Run server")

	streams := flag.Int("S", 1, "Number of streams")
	dataSize := flag.Int("size", 102400, "Data size on each stream")
	//c := flag.String("c", "", "Run client")
	//flag.Var(&bs, "bind", "bind to")
	file := flag.String("F", "", "Source data file")
	flag.Parse()

	logger := utils.DefaultLogger

	if *verbose {
		logger.SetLogLevel(utils.LogLevelDebug)
	} else {
		logger.SetLogLevel(utils.LogLevelInfo)
	}
	logger.SetLogTimeFormat("")

	if *s {
		data, err := ioutil.ReadFile(*file)
		if err != nil {
			panic(err)
		}
		logger.Infof("Starting server...")
		serverMain(data, *streams, *dataSize)
	} else {


		logger.Infof("Starting client...")
		clientMain()
	}

}
