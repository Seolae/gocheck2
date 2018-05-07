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
	Filetype   string
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
						if err = os.Rename(strings.TrimSuffix(event.Name, filepath.Ext(event.Name))+conf.Paths.Filetype, conf.Paths.Destpath+"so"+on+conf.Paths.Filetype); err != nil {
							log.Fatal(err)
						}
						if err = os.Remove(event.Name); err != nil {
							log.Fatal(err)
						}
						fmt.Println("Moved " + strings.TrimSuffix(event.Name, filepath.Ext(event.Name)) + conf.Paths.Filetype + " to " + conf.Paths.Destpath + "so" + on + conf.Paths.Filetype)

					}
				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	if err = watcher.Add(conf.Paths.Sourcepath); err != nil {
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
