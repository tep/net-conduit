package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"

	"toolman.org/net/conduit"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	conn, err := net.Dial("unixpacket", "@give-get-test")
	if err != nil {
		return err
	}

	c, err := conduit.FromConn(conn)
	if err != nil {
		return err
	}

	r, w, err := os.Pipe()
	if err != nil {
		return err
	}

	fmt.Printf("Created new pipe; sending writer (fd:%d) across Conduit\n", w.Fd())

	if err := c.TransferFile(w); err != nil {
		ce := err.(*conduit.Error)

		return fmt.Errorf("%#v -- %v", ce, ce)
	}

	fmt.Println("Reading from pipe")

	scn := bufio.NewScanner(r)

	for scn.Scan() {
		fmt.Printf("Got: %q\n", scn.Text())
	}

	if err := scn.Err(); err != nil {
		return err
	}

	return nil
}
