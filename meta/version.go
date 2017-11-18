package meta

import "fmt"

var (
	// Version holds the current version.
	Version = "dev"
	// BuildDate holds the build date.
	BuildDate = "I don't remember exactly"
)

// Display PRM version
func DisplayVersion() {
	fmt.Printf("Version: %s, %s\n", Version, BuildDate)
}
