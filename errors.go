package conduit

import (
	"errors"
)

// ErrType differentiates disparate Conduit errors.
type ErrType int

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

// Error encapsulates a conduit related error providing a Type method to
// discern the type of error.
type Error struct {
	error
	errType ErrType
}

// Type returns the conduit related error type indicated by the returned
// ErrType value.
func (e *Error) Type() ErrType {
	return e.errType
}

func noFDError() error {
	return &Error{errors.New("cannot access underlying file descriptor"), ErrNoFD}
}

func transferError(err error) error {
	return &Error{err, ErrFailedTransfer}
}

func closeError(err error) error {
	return &Error{err, ErrFailedClose}
}

func cloneError(err error) error {
	return &Error{err, ErrBadClone}
}

func controlMessageError(err error) error {
	return &Error{err, ErrBadCtrlMesg}
}
