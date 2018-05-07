package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

var on string

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
						fmt.Println(event.Name)
						time.Sleep(500 * time.Millisecond)
						on, err := rrline(event.Name, 86)
						if err != nil {
							log.Fatal(err)
						}
						fmt.Println("so" + on)
						err = os.Remove(event.Name)
						if err != nil {
							log.Fatal(err)
						}

					}
				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(`c:\test\`)
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
