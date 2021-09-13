package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func setup(t *testing.T, offset, limit int64) (string, string) {
	tmpDir, err := ioutil.TempDir("", "copy")
	if err != nil {
		t.Fatal("can't create temp dir: ", err)
	}
	outPath := filepath.Join(tmpDir, "out_offset"+strconv.Itoa(int(offset))+"_limit"+strconv.Itoa(int(limit))+".txt")
	return outPath, tmpDir
}

func TestCopy_Valid(t *testing.T) {
	tests := []struct {
		name     string
		fromPath string
		toPath   string
		offset   int64
		limit    int64
		err      error
	}{
		{
			name:     "test case when offset=0 and limit=0",
			fromPath: "testdata/input.txt",
			offset:   0,
			limit:    0,
			toPath:   "out_offset0_limit0.txt",
		},
		{
			name:     "test case when offset=0 and limit=10",
			fromPath: "testdata/input.txt",
			offset:   0,
			limit:    10,
			toPath:   "out_offset0_limit10.txt",
		},
		{
			name:     "test case when offset=0 and limit=1000",
			fromPath: "testdata/input.txt",
			offset:   0,
			limit:    1000,
			toPath:   "out_offset0_limit1000.txt",
		},
		{
			name:     "test case when offset=0 and limit=10000",
			fromPath: "testdata/input.txt",
			offset:   0,
			limit:    10000,
			toPath:   "out_offset0_limit10000.txt",
		},
		{
			name:     "test case when offset=100 and limit=1000",
			fromPath: "testdata/input.txt",
			offset:   100,
			limit:    1000,
			toPath:   "out_offset100_limit1000.txt",
		},
		{
			name:     "test case when offset=6000 and limit=1000",
			fromPath: "testdata/input.txt",
			offset:   6000,
			limit:    1000,
			toPath:   "out_offset6000_limit1000.txt",
		},
		{
			name:     "test case when limit more than file size",
			fromPath: "testdata/input.txt",
			offset:   0,
			limit:    10000000000000,
			toPath:   "out_offset0_limit0.txt",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			outPath, tmpDir := setup(t, tt.offset, tt.limit)
			defer os.RemoveAll(tmpDir)
			err := Copy(tt.fromPath, outPath, tt.offset, tt.limit)
			require.Nil(t, err)

			expected, err := ioutil.ReadFile("testdata/" + tt.toPath)
			require.NoError(t, err)

			actual, err := ioutil.ReadFile(outPath)
			require.NoError(t, err)

			require.Equal(t, expected, actual)
		})
	}
}

func TestCopy_InValid(t *testing.T) {
	t.Run("test case when returns error offset exceeds file size", func(t *testing.T) {
		err := Copy("testdata/input.txt", "out.txt", 8000, 6000)
		require.EqualError(t, err, ErrOffsetExceedsFileSize.Error())
	})
	t.Run("src file is not regular", func(t *testing.T) {
		err := Copy("/dev/urandom", "out.txt", 0, 0)
		require.EqualError(t, err, ErrUnsupportedFile.Error())
	})
	t.Run("src file not exists err", func(t *testing.T) {
		err := Copy("non_existent_file.txt", "out.txt", 0, 0)
		require.Error(t, err)
	})
	t.Run("offset more than file size", func(t *testing.T) {
		err := Copy("testdata/input.txt", "out.txt", 10000000000000, 0)
		require.Error(t, err)
	})
}
