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

	"github.com/fsnotify/fsnotify"
)

var on string
var sercpath = `c:\test\` //source path
var destpath = `c:\src\` //dest path
var onl = 86 //set this to what line you want to read the .XST

func main() {
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
						on, err = rrline(event.Name, onl)
						if err != nil {
							log.Fatal(err)
						}
						err = os.Rename(strings.TrimSuffix(event.Name, filepath.Ext(event.Name))+".tif", destpath+"so"+on+".tif")
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

	err = watcher.Add(sercpath)
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
