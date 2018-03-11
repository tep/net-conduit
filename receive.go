package conduit

import (
	"errors"
	"net"
	"os"
	"unsafe"

	"golang.org/x/sys/unix"
)

// ReceiveFD receives and returns a single open file descriptor from the
// Conduit along with a nil error. If an error is returned it will be a
// conduit.Error with its type set according to the following conditions.
//
// 		ErrFailedTransfer: if the message cannot be recieved.
//
// 		ControlMessageError: if the control message cannot be parsed, more than
// 		one control message is sent or more than one file descriptor is
// 		transfered.
//
func (c *Conduit) ReceiveFD() (uintptr, error) {
	fd := int(c.file.Fd())
	buf := make([]byte, unix.CmsgSpace(int(unsafe.Sizeof(fd))))

	_, n, _, _, err := unix.Recvmsg(fd, nil, buf, 0)
	if err != nil {
		return 0, transferError(err)
	}

	cmsgs, err := unix.ParseSocketControlMessage(buf[:n])
	if err != nil {
		return 0, controlMessageError(err)
	}

	if len(cmsgs) != 1 {
		return 0, controlMessageError(errors.New("bad control message count"))
	}

	fds, err := unix.ParseUnixRights(&(cmsgs[0]))

	if len(fds) != 1 {
		return 0, controlMessageError(errors.New("bad fd count"))
	}

	return uintptr(fds[0]), nil
}

// ReceiveFile returns a *os.File associated with the open file descripted
// recevied through the Conduit. The provided name will be attached to the now
// File object. See ReceiveFD() for a discussion of possible error conditions.
func (c *Conduit) ReceiveFile() (*os.File, error) {
	fd, err := c.ReceiveFD()
	if err != nil {
		return nil, err
	}

	return os.NewFile(fd, ""), nil
}

// ReceiveConn returns a net.Conn associated with the open file descriptor
// received through the Conduit.
//
// In addition to the errors described for ReceiveFD, the following are also
// possible.  The act of receiving the Conn requires a clone of an underlying
// File object. If this fails, a conduit.Error of type ErrBadClone is returned.
// Prior to returning the Conn, the original File will be closed. If this close
// results in an error, a conduit.Error of type ErrFailedClosed is returned.
func (c *Conduit) ReceiveConn() (net.Conn, error) {
	f, err := c.ReceiveFile()
	if err != nil {
		return nil, err
	}

	conn, err := net.FileConn(f)
	if err != nil {
		return nil, cloneError(err)
	}

	if err := f.Close(); err != nil {
		return nil, closeError(err)
	}

	return conn, nil
}
