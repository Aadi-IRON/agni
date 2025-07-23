package detectors

import (
	"fmt"
	"os"
	"os/exec"
)

// RunDeadCode checks if dead code is installed; if not, it installs it, then runs it
func RunDeadCode(path string) {
	fmt.Println("🔍 Scanning for dead code... Time to clean the skeletons from your closet 🧹")
	// Check if 'deadcode' is available in PATH
	_, err := exec.LookPath("deadcode")
	if err != nil {
		fmt.Println("⚠️  'deadcode' not found. Attempting to install it...")

		// Attempt to install it using `go install`
		installCmd := exec.Command("go", "install", "golang.org/x/tools/cmd/deadcode@latest")
		installCmd.Env = os.Environ() // inherit user's env
		output, installErr := installCmd.CombinedOutput()
		if installErr != nil {
			fmt.Println("❌ Failed to install 'deadcode':", string(output))
			return
		}
		fmt.Println("✅ 'deadcode' installed successfully.")
	}

	// Now run deadcode on the provided path
	cmd := exec.Command("deadcode", path)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("❌ Error running 'deadcode':", err)
	}
	if len(output) == 0 {
		fmt.Println("✅ No dead code found.")
	} else {
		fmt.Println("🧠 Deadcode report:")
		fmt.Println(string(output))
		fmt.Println("⚠️  Note: Some functions may be falsely flagged as unused.")
		fmt.Println(">>> Always cross-check before deletion to avoid accidental removal of valid code.")
	}
}
