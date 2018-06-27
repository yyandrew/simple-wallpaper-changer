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

func detectDesktopEnv() string {
	var out bytes.Buffer
	// cmd := exec.Command("/bin/sh", "-c", "ls /usr/bin/")
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

func main() {
	var err error
	var cmd *exec.Cmd
	var out bytes.Buffer
	osName := runtime.GOOS
	imagesDir := flag.String("d", "", "Select you image folder")
	file := flag.String("f", "", "Give a image file")
	flag.Parse()
	if *file == "" {
		log.Fatal("you must give a image by -f options")
	}
	fmt.Println(*imagesDir)
	if osName == "linux" {
		// /usr/share/lubuntu/wallpapers/formula_1_ferrari_f2008-wallpaper-1366x768.jpg
		// /home/andrew/Pictures/desktop 1_001.png
		switch detectDesktopEnv() {
		case "LXDE":
			cmd = exec.Command("pcmanfm", "-w", *file)
		case "GNOME", "Unity":
			cmd = exec.Command("gsettings", "set", "org.gnome.desktop.background", "picture-uri", "file://"+*file)
		default:
			log.Fatal("Not support this desktop environment!")
		}
	} else if osName == "darwin" {
		cmd = exec.Command("/usr/bin/osascript", "-e", fmt.Sprintf(`tell application "Finder" to set desktop picture to POSIX file "%s"`, *file))
	}
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%q \n", out.String())
}
