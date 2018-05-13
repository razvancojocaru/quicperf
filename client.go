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

func clientMain() error {
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
			//nrRead, err := io.Copy(ioutil.Discard, stream)
			//
			////_, err = stream.Write([]byte(message))
			//if err != nil {
			//	logger.Errorf("--------Error: " + err.Error())
			//}

			logger.Infof("Total received %d", nrRead)
			//total += nrRead
			defer stream.Close()
			defer wg.Done()
		}()
	}

	wg.Wait()

	fmt.Println(nrRead)

	return nil
}
