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

package conduit

import (
	"net"
	"os"

	"golang.org/x/sys/unix"
)

// TransferFD sends the open file descriptor fd through the Conduit. If
// successfully transferred, fd will be closed and may no longer be used
// by the caller.  On success, nil is returned.
//
// If an error is returned, it will be of type conduit.Error.
func (c *Conduit) TransferFD(fd uintptr) error {
	if err := unix.Sendmsg(int(c.file.Fd()), nil, unix.UnixRights(int(fd)), nil, 0); err != nil {
		return transferError(err)
	}

	if err := unix.Close(int(fd)); err != nil {
		return closeError(err)
	}

	return nil
}

// TransferFile send the open file descriptor associated with f through the
// Conduit.  If successfully transferred, f will be closed and may no longer be
// used by the caller.  On success, nil is returned.
//
// If an error is returned, it will be of type conduit.Error.
func (c *Conduit) TransferFile(f *os.File) error {
	return c.TransferFD(f.Fd())
}

// TransferConn send the open file descriptor associated with conn through the
// Conduit.  If successfully transferred, conn will be closed and may no longer
// be used by the caller.  Nil is returned on success.
//
// If conn's underlying type provides no way to discern its file descriptor,
// a conduit.Error of type ErrNoFD is returned. As part of the transfer, an
// os.File object is cloned from conn. If this fails, a conduit.Error of type
// ErrBadClone is returned. Note that both conn and its clone are closed upon
// a successful transfer.
func (c *Conduit) TransferConn(conn net.Conn) error {
	cf, ok := conn.(filer)
	if !ok {
		return noFDError()
	}

	f, err := cf.File()
	if err != nil {
		return cloneError(err)
	}

	defer conn.Close()

	return c.TransferFile(f)
}
