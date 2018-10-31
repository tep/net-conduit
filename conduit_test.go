package conduit

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"os"
	"os/exec"
	"testing"
)

// Testing diminsions:
//
// 		Transport:
// 			net.Conn
// 			os.File
// 			uintptr (raw FD)
//
//		Send:
//			net.Conn
//			os.File
//			uintptr
//
//		Receive:
//			net.Conn
//			os.File
//			uintptr

//----------------------------------------------------------------------------

func TestConn(t *testing.T) {
}

//----------------------------------------------------------------------------

func runTcpServer() error {
	cmd := exec.Command(os.Args[0], "-test.run=TestHelperProcess", "--")
	cmd.Env = []string{"ConduitTestHelperProcessAction=\"tcp-server\""}
}

// TestHelperProcess is not a real test; it's here to act as a child process
// for other tests.
func TestHelperProcess(t *testing.T) {
	var helper func(context.Context) error
	switch os.Getenv("ConduitTestHelperProcessAction") {
	case "tcp-server":
		helper = tcpServerHelper
	default:
		return
	}

	ctx := context.Background()

	if err := helper(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(1)
	}

	fmt.Fprintln(os.Stderr, "OK")
	os.Exit(0)
}

func tcpServerHelper(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ln, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		return err
	}
	defer ln.Close()

	fmt.Fprintf(os.Stderr, "ADDR: %s\n", ln.Addr())

	var shutdown bool

	for {
		err := func() error {
			conn, err := ln.Accept()
			if err != nil {
				return err
			}
			defer conn.Close()

			fmt.Fprintf(os.Stderr, "CONNECTION: %s\n", conn.RemoteAddr())

			scnr := bufio.NewScanner(conn)
			for scnr.Scan() {
				line := scnr.Text()
				fmt.Fprintf(os.Stderr, "RECEIVED: %s\n", line)

				if line == "CLOSE" {
					break
				}

				if line == "SHUTDOWN" {
					shutdown = true
					break
				}
			}
			return scnr.Err()
		}()

		if err != nil {
			return err
		}

		if shutdown {
			break
		}
	}

	return nil
}
