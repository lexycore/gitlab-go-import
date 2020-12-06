package version

import "fmt"

const major = 0
const minor = 0
const patch = 1
const relInfo = "master"

const (
	// Description provides tidy naming service
	Description = "go-import meta fixer for gitlab"
	// Usage describes what aims it will be used for
	Usage = "Allows importing subprojects and subdirectories of projects"
)

// Version returns version of this service
func Version() string {
	return fmt.Sprintf("v%d.%d.%d#%s", major, minor, patch, relInfo)
}
