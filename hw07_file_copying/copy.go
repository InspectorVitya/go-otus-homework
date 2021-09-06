package main

import (
	"io"
	"os"
	"time"

	"github.com/cheggaaa/pb/v3"
	"github.com/pkg/errors"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

// test
func Copy(fromPath, toPath string, offset, limit int64) error {
	file, err := os.Stat(fromPath)
	if err != nil {
		return err
	}
	// fileInfo validation
	if !file.Mode().IsRegular() {
		return ErrUnsupportedFile
	}
	// offset validation
	if offset > file.Size() {
		return ErrOffsetExceedsFileSize
	}
	// limit adjusting
	if limit == 0 || limit > file.Size()-offset {
		limit = file.Size() - offset
	}

	src, err := os.OpenFile(fromPath, os.O_RDONLY, os.ModeDir)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(toPath)
	if err != nil {
		return errors.Wrapf(err, "cannot create file for path: %s", toPath)
	}
	defer dst.Close()

	_, err = src.Seek(offset, io.SeekStart)
	if err != nil {
		return errors.Wrapf(err, "cannot execute seek")
	}

	var chunk int64 = 1024
	bar := pb.Full.Start64(limit)
	bar.Set(pb.Bytes, true)
	for {
		if chunk > limit {
			chunk = limit
		}

		if limit == 0 {
			break
		}
		written, err := io.CopyN(dst, src, chunk)

		if errors.Is(err, io.EOF) {
			bar.Add64(written)
			break
		}
		if err != nil {
			return errors.Wrapf(err, "cannot execute io.CopyN")
		}
		bar.Add64(written)
		limit -= written
		time.Sleep(time.Millisecond * 100)
	}
	bar.Finish()
	return nil
}
