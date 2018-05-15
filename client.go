package main

import (
	"github.com/lucas-clemente/quic-go"
	"crypto/tls"
	"sync"
	"github.com/lucas-clemente/quic-go/internal/utils"
	"io"
	"fmt"
)

var logger = utils.DefaultLogger

func clientMain(addr string) error {
	session, err := quic.DialAddr(addr, &tls.Config{InsecureSkipVerify: true}, nil)
	if err != nil {
		return err
	}
	defer session.Close(nil)

	var wg sync.WaitGroup
	//var total int64 = 0

	var nrRead int64 = 0

	for {
		stream, err := session.AcceptStream()
		if err != nil {
			logger.Errorf("--------Accept Error: " + err.Error())
			break
		}
		wg.Add(1)
		go func() {
			buf := make([]byte, 100)
			for {
				nr, er := stream.Read(buf)
				nrRead += int64(nr)
				if er != nil {
					if er != io.EOF {
						err = er
					}
					break
				}
			}

			logger.Infof("Total received %d", nrRead)
			stream.Close()
			nr, er := stream.Read(buf)
			nrRead += int64(nr)
			fmt.Printf("Error %s, nr=%d\n", er.Error(), nr)
			wg.Done()
		}()
	}

	wg.Wait()

	session.Close(nil)
	fmt.Println(nrRead)

	return nil
}
