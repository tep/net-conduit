package main

import (
	"fmt"
	"log"
	"net"

	"toolman.org/net/conduit"
)

func run() error {
	ln, err := net.Listen("unixpacket", "@give-get-test")
	if err != nil {
		return err
	}
	defer ln.Close()

	fmt.Println("Awaiting transfer...")

	conn, err := ln.Accept()
	if err != nil {
		return err
	}
	defer conn.Close()

	c, err := conduit.FromConn(conn)
	if err != nil {
		return err
	}

	f, err := c.ReceiveFile()
	if err != nil {
		return err
	}

	fmt.Printf("Got new file: %#v (fd:%d)\n", f, f.Fd())

	msg := "Hello from the other side"

	fmt.Printf("Sending: %q\n", msg)

	fmt.Fprintln(f, msg)
	fmt.Println("That is all.")

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
