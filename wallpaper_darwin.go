package main

import (
	"fmt"
	"os"
	"os/exec"
)

// changeWallpaper runs wallpaper change script
func changeWallpaper(displayNumber int, pictDir string, f os.FileInfo) error {
	script := `
	tell application "System Events"
		set desktopCount to count of desktops
		repeat with desktopNumber from 1 to desktopCount
			tell desktop desktopNumber
				if desktopNumber is %d or %d less than 1 then
					set picture to "%s/%s"
				end if
			end tell
		end repeat
	end tell
	`
	cmd := exec.Command("osascript", "-e",
		fmt.Sprintf(script, displayNumber, displayNumber, pictDir, f.Name()))
	return cmd.Run()
}
