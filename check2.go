package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	_ "github.com/qodrorid/godaemon"

	"github.com/BurntSushi/toml"
	"github.com/fsnotify/fsnotify"
)

var on string

//Config loads the settings from settings.conf
type Config struct {
	Paths paths
}

type paths struct {
	Sourcepath string
	Destpath   string
	Onl        int
}

func main() {
	var conf Config
	if _, err := toml.DecodeFile("settings.conf", &conf); err != nil {
		log.Println(err)
	}
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Create == fsnotify.Create {
					if strings.Contains(event.Name, ".XST") {
						fmt.Println("Found new .XST waiting 3 seconds")
						time.Sleep(3 * time.Second)
						on, err = rrline(event.Name, conf.Paths.Onl)
						if err != nil {
							log.Fatal(err)
						}
						err = os.Rename(strings.TrimSuffix(event.Name, filepath.Ext(event.Name))+".tif", conf.Paths.Destpath+"so"+on+".tif")
						if err != nil {
							log.Fatal(err)
						}
						err = os.Remove(event.Name)
						if err != nil {
							log.Fatal(err)
						}
						fmt.Println("Moved " + strings.TrimSuffix(event.Name, filepath.Ext(event.Name)) + ".tff" + " to " + "c:\\src\\" + "so" + on + ".tif")

					}
				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(conf.Paths.Sourcepath)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}

func rrline(fn string, n int) (string, error) {
	f, err := os.Open(fn)
	if err != nil {
		return "", err
	}
	defer f.Close()
	bf := bufio.NewReader(f)
	var line string
	for lnum := 0; lnum < n; lnum++ {
		line, err = bf.ReadString('\n')
		if err != nil {
			return "", err
		}
	}
	reg, err := regexp.Compile("[^0-9]")
	if err != nil {
		log.Fatal(err)
	}
	var on = reg.ReplaceAllString(line, "")
	return on, nil
}
