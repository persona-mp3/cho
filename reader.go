package main

import (
	"errors"
	"io"
)

// readLastLog reads the latest appended log by evaluating the number of bytes written
// It returns the newFileSize, the contents written and an error if [Stat] call failed, or
// [ReadAt] failed
func (c *Cho) readLastLog(originalFileSize int64) (int64, []byte, error) {
	fileInfo, err := c.source.Stat()
	if err != nil {
		return 0, nil, err
	}

	newFileSize := fileInfo.Size()

	bytesWritten := newFileSize - originalFileSize

	buffer := make([]byte, bytesWritten)
	n, err := c.source.ReadAt(buffer, originalFileSize)
	_ = n
	if err != nil && !errors.Is(err, io.EOF){
		return 0, nil, err
	}

	return newFileSize, buffer, nil
}
