package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

func server(config *Config) {
	ret := make(chan error)

	go func() {
		log.Println("Listening on TCP port ", config.port)
		ln, err := net.Listen("tcp", ":"+config.port)
		if err != nil {
			ret <- err
			return
		}
		defer ln.Close()

		log.Println("Waiting for clients to connect")
		for {
			conn, err := ln.Accept()
			if err != nil {
				ret <- err
				return
			}

			log.Println("Client " + conn.RemoteAddr().String() + " connected")
			go func() {
				io.Copy(conn, conn)
				conn.Close()
				log.Println("Client " + conn.RemoteAddr().String() + " disconnected")
			}()
		}
	}()

	go func() {
		log.Println("Listening on UDP port ", config.port)
		pc, err := net.ListenPacket("udp", ":"+config.port)
		if err != nil {
			ret <- err
			return
		}
		defer pc.Close()

		buf := make([]byte, config.size)
		for {
			n, addr, err := pc.ReadFrom(buf)
			if err != nil {
				ret <- err
				return
			}
			if n != config.size {
				ret <- fmt.Errorf("Read only %v of %v", n, 1024)
				return
			}

			n, err = pc.WriteTo(buf, addr)
			if err != nil {
				ret <- err
				return
			}
			if n != config.size {
				ret <- fmt.Errorf("Wrote only %v of %v", n, 1024)
				return
			}
		}
	}()

	log.Fatal(<-ret)
}
