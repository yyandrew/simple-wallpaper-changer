package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"runtime"
	"strings"
	"time"
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

// SetInterval update wallpaper intervally
func SetInterval(worker func(), delay time.Duration) chan bool {
	quit := make(chan bool)
	go func() {
		for {
			worker()
			select {
			case <-time.After(delay):
			case <-quit:
				return
			}
		}
	}()

	return quit
}

// SetTimeout Delay worker after delay seconds
func SetTimeout(worker func(), delay time.Duration) {
	time.AfterFunc(delay, worker)
}

func main() {
	// quit := make(chan bool)
	imagesDir := flag.String("d", "", "Select you image folder")
	file := flag.String("f", "", "Give a image file")
	interval := flag.Int("i", 1, "Set the interval time(Minute)")
	flag.Parse()
	if *imagesDir != "" {
		images, err := ioutil.ReadDir(*imagesDir)
		if err != nil {
			log.Fatal(err.Error())
		}
		for {
			for _, image := range images {
				filename := *imagesDir + "/" + image.Name()
				if strings.Contains(filename, ".jpg") {
					fmt.Println(filename)
					ChangeWallPaper(filename)
					time.Sleep(time.Minute * time.Duration(*interval))
				}
			}
		}
	} else if *file != "" {
		if err := ChangeWallPaper(*file); err != nil {
			log.Fatal(err.Error())
		}
	} else {
		fmt.Println("You shoud give a image path with -f or images dir with -d.")
	}
}
