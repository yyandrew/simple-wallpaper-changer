package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"strings"
)

// DetectDesktopEnv returns which Desktop environment you are using.It should return one of:
// Unity
// LXDE
// GNOME
// KDE
// MATE
// XFCE
func DetectDesktopEnv() string {
	var out bytes.Buffer
	cmd := exec.Command("/bin/sh", "-c", "ls /usr/bin | grep session")
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err.Error())
		log.Fatal("Unable support window platform")
	}
	output := out.String()
	if strings.Contains(output, "lxsession") {
		return "LXDE"
	} else if strings.Contains(output, "gnome-session") {
		return "GNOME"
	} else {
		return ""
	}
}

// ChangeWallPaper returns an error
//
// It will change you wall paper
func ChangeWallPaper(file string) error {
	var cmd *exec.Cmd
	osName := runtime.GOOS
	if osName == "linux" {
		switch DetectDesktopEnv() {
		case "LXDE":
			cmd = exec.Command("pcmanfm", "-w", file)
		case "GNOME", "Unity":
			cmd = exec.Command("gsettings", "set", "org.gnome.desktop.background", "picture-uri", "file://"+file)
		default:
			log.Fatal("Not support this desktop environment!")
		}
	} else if osName == "darwin" {
		cmd = exec.Command("/usr/bin/osascript", "-e", fmt.Sprintf(`tell application "Finder" to set desktop picture to POSIX file "%s"`, file))
	}
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func main() {
	imagesDir := flag.String("d", "", "Select you image folder")
	file := flag.String("f", "", "Give a image file")
	flag.Parse()
	if *file == "" {
		log.Fatal("You must give a image by -f options")
	}
	fmt.Println(*imagesDir)
	if err := ChangeWallPaper(*file); err != nil {
		log.Fatal(err.Error())
	}
}
