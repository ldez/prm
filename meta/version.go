package meta

import (
	"fmt"
	"runtime"
)

var (
	version = "dev"
	commit  = "I don't remember exactly"
	date    = "I don't remember exactly"
)

// GetVersion returns the current version.
func GetVersion() string {
	return version
}

// DisplayVersion DisplayVersion PRM version
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
