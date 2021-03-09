package dupver

import (
	"fmt"
)

// Print the global preferences structure
func (prefs Preferences) Print() {
	fmt.Printf("Editor: %s\n", prefs.Editor)
	fmt.Printf("Diff tool: %s\n", prefs.DiffTool)
	fmt.Printf("Default repository path: %s\n", prefs.DefaultRepo)
}

// Print the current global preferences as JSON
func (prefs Preferences) PrintJson() {
	PrintJson(prefs)
}