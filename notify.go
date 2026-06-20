package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/fsnotify/fsnotify"
)

// According to the docs, directories should be watched instead of
// particular files.
func watch(ctx context.Context, dir, logSource string, signal chan any) error {
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

	for {
		select {
		case <-ctx.Done():
			return nil

		case evt := <-watcher.Events:
			// TODO: A write event could mean a write in progress, so we'd want a delayed way of knowing
			if evt.Has(fsnotify.Write) && evt.Name == logSource {
				log.Printf("%s [+] has new write\n", logSource)
				signal <- struct{}{}
			}

		case err := <-watcher.Errors:
			return err

		}

	}

}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Kill)
	defer cancel()
	if len(os.Args) < 3 {
		log.Fatal("notify <path/to/watch> <target-file>")
	}

	event := make(chan any)
	dir := os.Args[1]
	logSource := os.Args[2]


	go func() {
		defer close(event)
		if err := watch(ctx, dir, logSource, event); err != nil {
			log.Println(err)
		}
	}()

	i := 0
	for evt := range event {
		i++
		fmt.Println("number of mods made to file: ", i , evt)
	}
}
