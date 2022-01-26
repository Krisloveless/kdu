package kdu

import "io/fs"

func pushToFileChannel(lock chan struct{}, filesizes *fileSizeChannel, basePath string, entry fs.FileInfo) {
	filesizes.c <- entry.Size()
}
