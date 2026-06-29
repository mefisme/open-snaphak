package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var stdin = bufio.NewReader(os.Stdin)

// interactiveInstall is the double-click / no-args experience: auto-detect DOOM, confirm, install, then pause
// so the console window (spawned by a double-click) stays up long enough to read the result.
func interactiveInstall() {
	fmt.Println("SnapHak installer", version)
	fmt.Println()

	doom, err := resolveDoom("")
	if err != nil {
		fmt.Println("Could not auto-detect your DOOM 2016 install.")
		fmt.Print("Enter the path to your DOOM folder (the one with DOOMx64vk.exe), or leave blank to cancel: ")
		doom = readLine()
		if doom == "" || !hasDoomExe(doom) {
			fmt.Println("No valid DOOM folder provided -- cancelled.")
			pause()
			return
		}
	} else {
		fmt.Printf("Found DOOM: %s\n", doom)
		fmt.Print("Install SnapHak here? [Y/n] ")
		if !isYes(readLine()) {
			fmt.Println("Cancelled.")
			pause()
			return
		}
	}

	if err := cmdInstall(flags{doom: doom}); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
	}
	pause()
}

func readLine() string {
	s, _ := stdin.ReadString('\n')
	return strings.TrimSpace(s)
}

func isYes(s string) bool {
	s = strings.ToLower(strings.TrimSpace(s))
	return s == "" || s == "y" || s == "yes"
}

func pause() {
	fmt.Print("\nPress Enter to exit...")
	readLine()
}
