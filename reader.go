package main

import (
	// "errors"
	// "io"
	// "os"
	"log"
)

// func main() {
// 	f, err := os.OpenFile("./main.go", os.O_RDWR, 600)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
//
// 	info, _ := f.Stat()
// 	fsize := info.Size()
// 	pos, err := f.Seek(fsize, 0)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
//
// 	log.Println("EOF->", pos)
// 	n, err := f.WriteAt([]byte("\n\nwhat about confidenitiality? You have ...?\n"), pos)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
//
// 	if err := f.Sync(); err != nil {
// 		log.Fatal(err)
// 	}
// 	new_pos := pos + (int64(n))
// 	log.Println("newPos: ", new_pos)
//
// 	log.Println("bytes written: ", n)
// 	buff := make([]byte, n)
//
// 	// when we write to a file i assum seek automatically moves the
// 	// pointer to the last byte? so where does read start from if we're truly
// 	// reading from the OG position of the file before writing?
// 	lastWritten, err := f.ReadAt(buff, pos)
//
// 	if err != nil && !errors.Is(err, io.EOF) {
// 		log.Fatal("read error: ", err, lastWritten, string(buff))
// 	}
//
// 	log.Println("lastWritten:: ", lastWritten)
// 	log.Println("content: ", string(buff))
// }

func (c *Cho) readLastLog(originalFileSize int64) (int64, []byte, error) {
	fileInfo, err := c.source.Stat()
	if err != nil {
		return 0, nil, err
	}

	newFileSize := fileInfo.Size()

	bytesWritten := newFileSize - originalFileSize

	// currPosition, err := c.source.Seek(newFileSize, 0)
	// if err != nil {
	// 	return 0, nil, err
	// }

	buffer := make([]byte, bytesWritten)
	n, err := c.source.ReadAt(buffer, originalFileSize)
	if err != nil {
		return 0, nil, err
	}

	log.Println("number of bytesRead: ", n)

	return newFileSize, buffer, nil
}
