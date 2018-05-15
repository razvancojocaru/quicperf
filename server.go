package main

import (
	"crypto/tls"
	"crypto/rsa"
	"crypto/rand"
	"crypto/x509"
	"math/big"
	"encoding/pem"
	"github.com/lucas-clemente/quic-go"
	"sync"
	"fmt"
	"time"
)

var numStreams = 1
var sendData []byte
var maxSec int
const BUFSIZE = 5000

// Setup a bare-bones TLS config for the server
func generateTLSConfig() *tls.Config {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}
	return &tls.Config{Certificates: []tls.Certificate{tlsCert}}
}

func writeToStream(s quic.Stream) (size int, err error){
	size = 0
	intervalSize := 0
	start := 0
	end := BUFSIZE
	elapsedSeconds := 0

	var mBytes float64

	timeStart := time.Now()
	intervalStart := timeStart
	for {
		var nr int
		if end < len(sendData) {
			nr, err = s.Write(sendData[start:end])
		} else {
			nr, err = s.Write(sendData[start:])
			size += nr
			if err != nil {
				break
			}
			end = end - len(sendData)
			nr, err = s.Write(sendData[:end])
		}
		size += nr
		if err != nil {
			break
		}
		start = end
		end = end + BUFSIZE

		now := time.Now()
		elapsed := now.Sub(intervalStart)
		if elapsed.Seconds() >= 1 {
			sec := intervalStart.Sub(timeStart).Seconds()
			mBytes = float64(size-intervalSize) / 1000000

			fmt.Printf("[ID %d] %.1f - %.1f sec %.2f MBytes %.2f Mbits/sec\n",
				s.StreamID(), sec, sec+1, mBytes, mBytes*8)
			intervalStart = now
			intervalSize = size

			elapsedSeconds++
			if maxSec == elapsedSeconds {
				totalMBytes := float64(size) / 1000000
				fmt.Printf("[ID %d] %.1f - %.1f sec %.2f MBytes %.2f Mbits/sec\n",
					s.StreamID(), 0.0, float64(maxSec), totalMBytes, totalMBytes*8 / float64(maxSec))
				break
			}
		}
	}
	return
}

func handleSession(session quic.Session) {
	defer session.Close(nil)
	var wg sync.WaitGroup
	wg.Add(numStreams)
	for i := 0 ; i < numStreams; i++ {
		go func() {
			stream, err := session.OpenStreamSync()
			if err != nil {
				return
			}

			//numWritten, err := stream.Write(sendData)
			_, err = writeToStream(stream)

			//fmt.Printf("Sent %d on stream %d\n", numWritten, (i+1)*2)
			if err != nil {
				logger.Errorf(err.Error())
			}

			stream.Context().Done()
			stream.Close()
			wg.Done()
		} ()
	}
	wg.Wait()
}

func serverMain(addr string, data []byte, streams int, ticks int) error {
	sendData = data
	numStreams = streams
	maxSec = ticks
	listener, err := quic.ListenAddr(addr, generateTLSConfig(), nil)
	if err != nil {
		return err
	}

	for {
		session, err := listener.Accept()
		if err != nil {
			return err
		}

		go handleSession(session)

	}
	return err
}
