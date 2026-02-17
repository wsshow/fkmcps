package update

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// UnzipCallback defines the type for the progress callback function used in Unzip.
type UnzipCallback func(processed int, total int, fileName string, isDir bool)

// Unzip extracts a zip archive to a specified destination directory with an optional progress callback.
func Unzip(source, destination string, callback UnzipCallback) error {
	reader, err := zip.OpenReader(source)
	if err != nil {
		return err
	}
	defer reader.Close()

	destDir := filepath.Clean(destination)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return err
	}

	totalFiles := len(reader.File)

	for i, f := range reader.File {
		// Trigger progress callback
		if callback != nil {
			callback(i+1, totalFiles, f.Name, f.FileInfo().IsDir())
		}

		// Extract individual file logic to facilitate resource management with defer
		err := extractFile(f, destDir)
		if err != nil {
			return fmt.Errorf("failed to extract file %s: %w", f.Name, err)
		}
	}
	return nil
}

// extractFile handles the extraction logic for a single file
func extractFile(f *zip.File, destDir string) error {
	fpath := filepath.Join(destDir, f.Name)

	if !strings.HasPrefix(fpath, destDir+string(os.PathSeparator)) {
		return fmt.Errorf("illegal path (Zip Slip): %s", fpath)
	}

	// If it's a directory, create it directly
	if f.FileInfo().IsDir() {
		return os.MkdirAll(fpath, os.ModePerm)
	}

	// Ensure parent directory exists (prevent issues if files appear before directories in the zip)
	if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
		return err
	}

	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	mode := f.Mode()
	if mode&0200 == 0 {
		mode |= 0200
	}

	if f.Mode()&os.ModeSymlink != 0 {
		return nil
	}

	outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, rc)
	return err
}
