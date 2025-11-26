
[![GoDoc](https://godoc.org/toolman.org/net/conduit?status.svg)](https://godoc.org/toolman.org/net/conduit) [![Go Report Card](https://goreportcard.com/badge/toolman.org/net/conduit)](https://goreportcard.com/report/toolman.org/net/conduit) [![Build Status](https://travis-ci.org/tep/net-conduit.svg?branch=master)](https://travis-ci.org/tep/net-conduit)

# conduit
`import "toolman.org/net/conduit"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)
* [Subdirectories](#pkg-subdirectories)

## <a name="pkg-overview">Overview</a>

## Install

``` sh
  go get toolman.org/net/conduit
```


## <a name="pkg-index">Index</a>
* [type Conduit](#Conduit)
  * [func FromConn(conn net.Conn) (*Conduit, error)](#FromConn)
  * [func FromFile(f *os.File) *Conduit](#FromFile)
  * [func New(fd uintptr, name string) *Conduit](#New)
  * [func (c *Conduit) Close() error](#Conduit.Close)
  * [func (c *Conduit) ReceiveConn() (net.Conn, error)](#Conduit.ReceiveConn)
  * [func (c *Conduit) ReceiveFD() (uintptr, error)](#Conduit.ReceiveFD)
  * [func (c *Conduit) ReceiveFile() (*os.File, error)](#Conduit.ReceiveFile)
  * [func (c *Conduit) TransferConn(conn net.Conn) error](#Conduit.TransferConn)
  * [func (c *Conduit) TransferFD(fd uintptr) error](#Conduit.TransferFD)
  * [func (c *Conduit) TransferFile(f *os.File) error](#Conduit.TransferFile)
* [type ErrType](#ErrType)
* [type Error](#Error)
  * [func (e *Error) Type() ErrType](#Error.Type)


#### <a name="pkg-files">Package files</a>
[conduit.go](/src/toolman.org/net/conduit/conduit.go) [errors.go](/src/toolman.org/net/conduit/errors.go) [receive.go](/src/toolman.org/net/conduit/receive.go) [send.go](/src/toolman.org/net/conduit/send.go) 






## <a name="Conduit">type</a> [Conduit](/src/target/conduit.go?s=476:537#L5)
``` go
type Conduit struct {
    // contains filtered or unexported fields
}
```
A Conduit is a mechanism for transfering open file descriptors between
cooperating processes. Transfers can take place over an os.File or net.Conn
but ultimately the transport descriptor must manifest as a socket capable of
carrying out-of-band control messages.







### <a name="FromConn">func</a> [FromConn](/src/target/conduit.go?s=2020:2066#L45)
``` go
func FromConn(conn net.Conn) (*Conduit, error)
```
FromConn creates a new Conduit from the provided net.Conn. The underlying
type for the provided Conn must be one having a method with the signature
"File() (*os.File, error)".  If not, a conduit.Error of type ErrNoFD will
be returned.

A cloned os.File object is created by FromConn. If this cloning fails,
a conduit.Error of type ErrBadClone} will be returned. Since a clone is
being created, you should be sure to call Close to avoid leaking a file
descriptor.


### <a name="FromFile">func</a> [FromFile](/src/target/conduit.go?s=1465:1499#L32)
``` go
func FromFile(f *os.File) *Conduit
```
FromFile creates a new Conduit from the provided os.File.


### <a name="New">func</a> [New](/src/target/conduit.go?s=1311:1353#L27)
``` go
func New(fd uintptr, name string) *Conduit
```
New creates a new Conduit. The provided file descriptor is the transport
over which other open FDs will be transferred and thus must be capable of
carrying out-of-band control messages. Note that this restriction is not
enforced here but will instead cause later transfer actions to fail. The
given name is as passed to os.NewFile.





### <a name="Conduit.Close">func</a> (\*Conduit) [Close](/src/target/conduit.go?s=870:901#L15)
``` go
func (c *Conduit) Close() error
```
Close is provided to allow the caller to close any cloned os.File objects
that may have been created while constructing a Conduit. If none were
created, then calling Close has no effect and will return nil. Therefore,
it's a good idea to always call Close when you're done with a Conduit.
Close implements io.Closer




### <a name="Conduit.ReceiveConn">func</a> (\*Conduit) [ReceiveConn](/src/target/receive.go?s=2039:2088#L59)
``` go
func (c *Conduit) ReceiveConn() (net.Conn, error)
```
ReceiveConn returns a net.Conn associated with the open file descriptor
received through the Conduit.

In addition to the errors described for ReceiveFD, the following are also
possible.  The act of receiving the Conn requires a clone of an underlying
File object. If this fails, a conduit.Error of type ErrBadClone is returned.
Prior to returning the Conn, the original File will be closed. If this close
results in an error, a conduit.Error of type ErrFailedClosed is returned.




### <a name="Conduit.ReceiveFD">func</a> (\*Conduit) [ReceiveFD](/src/target/receive.go?s=538:584#L12)
``` go
func (c *Conduit) ReceiveFD() (uintptr, error)
```
ReceiveFD receives and returns a single open file descriptor from the
Conduit along with a nil error. If an error is returned it will be a
conduit.Error with its type set according to the following conditions.


	ErrFailedTransfer: if the message cannot be recieved.
	
	ControlMessageError: if the control message cannot be parsed, more than
	one control message is sent or more than one file descriptor is
	transfered.




### <a name="Conduit.ReceiveFile">func</a> (\*Conduit) [ReceiveFile](/src/target/receive.go?s=1384:1433#L42)
``` go
func (c *Conduit) ReceiveFile() (*os.File, error)
```
ReceiveFile returns a *os.File associated with the open file descripted
recevied through the Conduit. The provided name will be attached to the now
File object. See ReceiveFD() for a discussion of possible error conditions.




### <a name="Conduit.TransferConn">func</a> (\*Conduit) [TransferConn](/src/target/send.go?s=1495:1546#L35)
``` go
func (c *Conduit) TransferConn(conn net.Conn) error
```
TransferConn send the open file descriptor associated with conn through the
Conduit.  If succesfully transfered, conn will be closed and may no longer
be used by the caller.  Nil is returned on success.

If conn's underlying type provides no way to discern its file descriptor,
a conduit.Error of type ErrNoFD is returned. As part of the transfer, an
os.File object is cloned from conn. If this fails, a conduit.Error of type
ErrBadClone is returned. Note that both conn and its clone are closed upon
a successful transfer.




### <a name="Conduit.TransferFD">func</a> (\*Conduit) [TransferFD](/src/target/send.go?s=325:371#L5)
``` go
func (c *Conduit) TransferFD(fd uintptr) error
```
TransferFD sends the open file descriptor fd through the Conduit. If
successfully transfered, fd will be closed and may no longer be used
by the caller.  On success, nil is returned.

If an error is returned, it will be of type conduit.Error.




### <a name="Conduit.TransferFile">func</a> (\*Conduit) [TransferFile](/src/target/send.go?s=862:910#L22)
``` go
func (c *Conduit) TransferFile(f *os.File) error
```
TransferFile send the open file descriptor associated with f through the
Conduit.  If succesfully transfered, f will be closed and may no longer be
used by the caller.  On success, nil is returned.

If an error is returned, it will be of type conduit.Error.




## <a name="ErrType">type</a> [ErrType](/src/target/errors.go?s=91:107#L1)
``` go
type ErrType int
```
ErrType differentiates disparate Conduit errors.


``` go
const (
    // ErrUnknown is an unknown error type; there are no errors of this type
    // (i.e. if you get one of these it's a bug)
    ErrUnknown ErrType = iota

    // ErrNoFD is returned when a Conduit method is unable to extrapolate an
    // underlying file descriptor from one of its arguments.
    ErrNoFD

    // ErrFailedTransfer is returned when a file descriptor transfer fails.
    ErrFailedTransfer

    // ErrFailedClose is returned when a Close method fails.
    ErrFailedClose

    // ErrBadClone is returned on a failed attempt to clone an os.File object.
    ErrBadClone

    // ErrBadCtrlMesg is returned for low level errors while constructing,
    // sending or receiving the out-of-band control message used to transfer
    // a file descriptor.
    ErrBadCtrlMesg
)
```

## <a name="Error">type</a> [Error](/src/target/errors.go?s=956:1001#L26)
``` go
type Error struct {
    // contains filtered or unexported fields
}
```
Error encapsulates a conduit related error providing a Type method to
discern the type of error.

### <a name="Error.Type">func</a> (\*Error) [Type](/src/target/errors.go?s=1094:1124#L33)
``` go
func (e *Error) Type() ErrType
```
Type returns the conduit related error type indicated by the returned
ErrType value.
