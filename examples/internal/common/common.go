package common

import (
	"os"
	"path/filepath"
	"runtime"
)

func SocketName() string {
	if runtime.GOOS == "linux" {
		return "@net-conduit-example"
	}
	return filepath.Join(os.TempDir(), "net-conduit-example")
}
