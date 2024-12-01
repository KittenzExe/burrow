package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
)

const (
	blue  = "\033[34m"
	reset = "\033[0m"
)

func main() {
	all := flag.Bool("all", false, "Watch all Go files in the directory")
	flag.Parse()

	if len(flag.Args()) < 1 {
		log.Fatal("Please provide the file to watch")
	}

	fileToWatch := flag.Args()[0]
	absPath, err := filepath.Abs(fileToWatch)
	if err != nil {
		log.Fatalf("Failed to get absolute path: %v", err)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	var mu sync.Mutex
	var cmd *exec.Cmd
	var lastEventTime time.Time

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					mu.Lock()
					if time.Since(lastEventTime) > 100*time.Millisecond {
						fmt.Printf("%sFile modified: %s%s\n", blue, event.Name, reset)
						fmt.Println()
						if cmd != nil && cmd.Process != nil {
							cmd.Process.Kill()
						}
						cmd = runFile(absPath)
						lastEventTime = time.Now()
					}
					mu.Unlock()
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("Error:", err)
			}
		}
	}()

	if *all {
		dir := filepath.Dir(absPath)
		err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && strings.HasSuffix(info.Name(), ".go") {
				err = watcher.Add(path)
				if err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err = watcher.Add(absPath)
		if err != nil {
			log.Fatal(err)
		}
	}

	cmd = runFile(absPath)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		if cmd != nil && cmd.Process != nil {
			cmd.Process.Kill()
		}
		close(done)
	}()

	<-done
}

func runFile(filePath string) *exec.Cmd {
	cmd := exec.Command("go", "run", filePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		log.Printf("Failed to start file: %v", err)
		return cmd
	}

	err = cmd.Wait()
	if err != nil {
		log.Printf("Failed to run file: %v", err)
	}

	return cmd
}
