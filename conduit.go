// Package conduit provides the Conduit type used for transferring open file
// descriptors between cooperating processes.

package conduit // import "toolman.org/net/conduit"

import (
	"net"
	"os"
)

// A Conduit is a mechanism for transferring open file descriptors between
// cooperating processes. Transfers can take place over an os.File or net.Conn
// but ultimately the transport descriptor must manifest as a socket capable of
// carrying out-of-band control messages.
type Conduit struct {
	file   *os.File
	closer func() error
}

// Close is provided to allow the caller to close any cloned os.File objects
// that may have been created while constructing a Conduit. If none were
// created, then calling Close has no effect and will return nil. Therefore,
// it's a good idea to always call Close when you're done with a Conduit.
// Close implements io.Closer
func (c *Conduit) Close() error {
	if c.closer != nil {
		return c.closer()
	}
	return nil
}

// New creates a new Conduit. The provided file descriptor is the transport
// over which other open FDs will be transferred and thus must be capable of
// carrying out-of-band control messages. Note that this restriction is not
// enforced here but will instead cause later transfer actions to fail. The
// given name is as passed to os.NewFile.
func New(fd uintptr, name string) *Conduit {
	return &Conduit{file: os.NewFile(fd, name)}
}

// FromFile creates a new Conduit from the provided os.File.
func FromFile(f *os.File) *Conduit {
	return &Conduit{file: f}
}

// FromConn creates a new Conduit from the provided net.Conn. The underlying
// type for the provided Conn must be one having a method with the signature
// "File() (*os.File, error)".  If not, a conduit.Error of type ErrNoFD will
// be returned.
//
// A cloned os.File object is created by FromConn. If this cloning fails,
// a conduit.Error of type ErrBadClone} will be returned. Since a clone is
// being created, you should be sure to call Close to avoid leaking a file
// descriptor.
func FromConn(conn net.Conn) (*Conduit, error) {
	cf, ok := conn.(filer)
	if !ok {
		return nil, noFDError()
	}

	f, err := cf.File()
	if err != nil {
		return nil, cloneError(err)
	}

	return &Conduit{f, f.Close}, nil
}

type filer interface {
	File() (*os.File, error)
}
