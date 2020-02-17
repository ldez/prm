package meta

import (
	"fmt"
	"runtime"
)

var (
	version = "devel"
	commit  = "-"
	date    = "-"
)

// GetVersion returns the current version.
func GetVersion() string {
	return version
}

// DisplayVersion Display version information.
func DisplayVersion() {
	fmt.Printf(`prm:
 version     : %s
 commit      : %s
 build date  : %s
 go version  : %s
 go compiler : %s
 platform    : %s/%s
`, version, commit, date, runtime.Version(), runtime.Compiler, runtime.GOOS, runtime.GOARCH)
}
