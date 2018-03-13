// Copyright © 2018 Timothy E. Peoples <eng@toolman.org>
//
// This program is free software; you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// as published by the Free Software Foundation; either version 2
// of the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"toolman.org/net/conduit"
	"toolman.org/net/conduit/examples/internal/common"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	name := common.SocketName()
	if name[0] != '@' {
		if err := os.Remove(name); err == nil {
			fmt.Printf("Removed orphaned unix domain socket: %s", name)
		}
	}

	ln, err := net.Listen("unix", name)
	if err != nil {
		return err
	}
	defer ln.Close()

	fmt.Printf("Awaiting transfer via %q... ", name)

	conn, err := ln.Accept()
	if err != nil {
		return err
	}
	defer conn.Close()

	fmt.Printf("\nGot connection: %q\n", conn.RemoteAddr())

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
