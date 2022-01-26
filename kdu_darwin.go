package kdu

import (
	"fmt"
	"golang.org/x/sys/unix"
	"io/fs"
	"os"
	"path/filepath"
)

func pushToFileChannel(lock chan struct{}, filesizes *fileSizeChannel, basePath string, entry fs.FileInfo) {
	// additional lock for fd in unix
	lock <- struct{}{}
	filePath := filepath.Join(basePath, entry.Name())
	statRes := unix.Stat_t{}
	err := unix.Stat(filePath, &statRes)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unix.Stat error: %v\n", err)
	}
	filesizes.c <- int64(statRes.Blocks * 512)
	<-lock
}
