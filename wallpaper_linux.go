package main

import (
	"fmt"
	"os"
	"os/exec"
)

// changeWallpaper runs wallpaper change script
func changeWallpaper(displayNumber int, pictDir string, f os.FileInfo) error {
	cmd := exec.Command("gsettings", "set", "org.gnome.desktop.background", "picture-uri",
		fmt.Sprintf(`"file://%s/%s"`, pictDir, f.Name()))
	return cmd.Run()
}
