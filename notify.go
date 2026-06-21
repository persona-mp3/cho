package main

import (
	"context"
	"fmt"
	"log"
	"os"
	// "os/signal"

	"github.com/fsnotify/fsnotify"
)

// According to the docs, directories should be watched instead of particular files.
// If the dir isn't a directory, a watcher error occured, or the  logSource doesn't exist, an error is returned,
func watch(ctx context.Context, dir, logSource string, signal chan struct{}) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	info, err := os.Stat(dir)
	if err != nil {
		return err
	}

	if !info.IsDir() {
		return fmt.Errorf("%s is not a directory. Provide directory to watch", dir)
	}

	if err := watcher.Add(dir); err != nil {
		return err
	}

	defer func() {
		if err := watcher.Close(); err != nil {
			log.Println("tried closing watcher after already been closed", err)
		}
	}()

	log.Println("watcher started....")
	log.Println("target_dir: ", dir)
	log.Println("logSource: ", logSource)
	for {
		select {
		case <-ctx.Done():
			return nil

		case evt := <-watcher.Events:
			if evt.Has(fsnotify.Write) {
				if evt.Name == logSource {
					log.Printf("[+] %s\n", logSource)
					signal <- struct{}{}
				}
			}

		case err := <-watcher.Errors:
			return err

		}

	}

}
