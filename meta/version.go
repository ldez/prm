package meta

import "fmt"

var (
	// Version holds the current version.
	Version = "dev"
	// BuildDate holds the build date.
	BuildDate = "I don't remember exactly"
)

func DisplayVersion() {
	fmt.Printf("Version: %s, %s", Version, BuildDate)
}
