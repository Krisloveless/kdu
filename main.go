package main

import (
	"context"
	"flag"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"time"
)

// kdu: linux-command-like du
/*
	this program triggers a large amount of goroutines, with graceful shutdown
*/

func dirents(dir string) []fs.FileInfo {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "kdu error: %v\n", err)
		return nil
	}
	return files
}

func walkDir(dir string, filesizes *fileSizeChannel, wg *sync.WaitGroup, routineLock chan struct{}) {
	defer wg.Done()
	for _, entry := range dirents(dir) {
		if entry.IsDir() {
			newPath := filepath.Join(dir, entry.Name())
			wg.Add(1)
			routineLock <- struct{}{}
			go walkDir(newPath, filesizes, wg, routineLock)
			<-routineLock
		} else {
			if filesizes.isClosed {
				// interrupted
				return
			}
			filesizes.c <- entry.Size()
		}
	}
}

type fileSizeChannel struct {
	isClosed bool
	c        chan int64
}

func main() {
	start := time.Now()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()
	// to inform wg.Wait()
	waitCh := make(chan struct{})
	maxGoroutine := 10000
	routineLock := make(chan struct{}, maxGoroutine)
	flag.Parse()
	dir := flag.Args()
	if len(dir) == 0 {
		dir = []string{"."}
	}
	// keep track of current walkDir number
	var wg sync.WaitGroup
	// this waitgroup waits until goroutine prints and ends
	var wgPrint sync.WaitGroup
	fileSizes := fileSizeChannel{c: make(chan int64), isClosed: false}
	for _, value := range dir {
		wg.Add(1)
		routineLock <- struct{}{}
		go walkDir(value, &fileSizes, &wg, routineLock)
		<-routineLock
	}

	go func() {
		defer wgPrint.Done()
		wgPrint.Add(1)
		var total int64
		var fileNo int64
		for value := range fileSizes.c {
			fileNo++
			total += value
		}
		fmt.Printf("number of files: %v, size: %.2f GB, time elapsed: %v\n", fileNo, float64(total)/1e9, time.Since(start))
	}()

	go func() {
		wg.Wait()
		close(waitCh)
	}()

	select {
	case <-waitCh:
	case <-ctx.Done():
		fileSizes.isClosed = true
		stop()
	}
	close(fileSizes.c)
	wgPrint.Wait()
}
