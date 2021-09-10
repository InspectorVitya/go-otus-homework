package main

import (
	"errors"
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"io"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	file, err := os.Stat(fromPath)
	if err != nil {
		return fmt.Errorf("cannot read stat file %q: %w", fromPath, err)
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
		return fmt.Errorf("cannot open file %q: %w", fromPath, err)
	}
	defer src.Close()

	dst, err := os.Create(toPath)
	if err != nil {
		return  fmt.Errorf("create file %q: %w", toPath, err)
	}
	defer dst.Close()

	_, err = src.Seek(offset, io.SeekStart)
	if err != nil {
		return fmt.Errorf("cannot execute seek: %w", err)
	}


	bar := pb.Full.Start64(limit)
	barReader := bar.NewProxyReader(src)
	_, err = io.CopyN(dst, barReader, limit)
	if err != nil && !errors.Is(err, io.EOF) {
		return err
	}
	bar.Finish()
	return nil
}
