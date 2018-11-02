package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"net"
	"time"
)

type Packet struct {
	SequenceNum uint64
	Timestamp   int64
}

type Stats struct {
	rcvdTime int64
	packet   *Packet
}

func client(config *Config) {
	protocol := "tcp"
	if config.udp {
		protocol = "udp"
	}

	conn, err := net.Dial(protocol, config.address+":"+config.port)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ret := make(chan error)
	read := make(chan Stats, 1)

	// Stat reporter
	go func() {
		var tick uint64
		var numRcvd uint64
		var numSent uint64
		var count float64
		var meanDelay float64
		var dSquared float64

		ticker := time.NewTicker(time.Millisecond * 1000)

		fmt.Println("Tick UnixTimeMs Size #Received #Sent AvgDelayMs DelayStd")
		for {
			select {
			case stats := <-read:
				numRcvd++

				if numSent < stats.packet.SequenceNum {
					numSent = stats.packet.SequenceNum
				}

				count++
				delay := float64(stats.rcvdTime-stats.packet.Timestamp) / 1e6
				newMeanDelay := meanDelay + (delay-meanDelay)/count
				dSquared = dSquared + (delay-newMeanDelay)*(delay-meanDelay)
				meanDelay = newMeanDelay
			case <-ticker.C:
				tick++
				fmt.Printf("%v", tick)
				fmt.Printf(" %v", time.Now().UnixNano()/1e6)
				fmt.Printf(" %v", config.size)
				fmt.Printf(" %v", numRcvd)
				fmt.Printf(" %v", numSent)
				fmt.Printf(" %v", meanDelay)
				fmt.Printf(" %v", math.Sqrt(dSquared/(count-1)))
				fmt.Println("")

				count = 0
				meanDelay = 0
			}
		}
	}()

	// Writer
	go func() {
		buf := NewPaddedBuffer(config.size)

		// The sequence number represents how many we have tried to send
		seqNum := uint64(1)
		for {
			packet := Packet{SequenceNum: seqNum, Timestamp: time.Now().UnixNano()}
			seqNum++

			err := binary.Write(buf, binary.LittleEndian, &packet)
			if err != nil {
				ret <- err
				return
			}

			n, err := buf.WriteTo(conn)
			if err != nil {
				ret <- err
				return
			}
			if n != int64(config.size) {
				ret <- fmt.Errorf("Wrote %v bytes of %v", n, config.size)
				return
			}
		}
	}()

	// Reader
	go func() {
		buf := NewPaddedBuffer(config.size)

		for {
			n, err := buf.ReadFrom(conn)
			if err != nil {
				ret <- err
				return
			}
			if n != int64(config.size) {
				ret <- fmt.Errorf("Read %v bytes of %v", n, config.size)
				return
			}

			stats := Stats{rcvdTime: time.Now().UnixNano()}

			var packet Packet
			err = binary.Read(buf, binary.LittleEndian, &packet)
			if err != nil {
				ret <- err
				return
			}
			stats.packet = &packet

			read <- stats
		}
	}()

	log.Fatal(<-ret)
}
