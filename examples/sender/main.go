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
	"bufio"
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
	fmt.Printf("Using socket %q for Conduit\n", name)

	conn, err := net.Dial("unix", name)
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
