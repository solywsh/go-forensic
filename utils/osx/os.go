package osx

import (
	"archive/tar"
	"fmt"
	"github.com/solywsh/go-forensic/utils/logger"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var (
	log = logger.NewLogger()
)

func TarDecompression(filePath, destinationDir string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("could not open tar file: %w", err)
	}
	tarReader := tar.NewReader(file)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Errorf("error reading tar header: %v", err)
			continue
		}
		// clean path to prevent absolute paths and path traversal attacks
		cleanName := filepath.Clean(header.Name)
		if strings.Contains(cleanName, "..") || filepath.IsAbs(cleanName) {
			log.Warnf("skipping suspicious path: %s", cleanName)
			continue
		}
		// replace illegal characters in path (especially for Windows)
		safePath := sanitizePath(cleanName)
		targetPath := filepath.Join(destinationDir, safePath)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(targetPath, os.FileMode(header.Mode)); err != nil {
				log.Errorf("error creating directory %s: %v", targetPath, err)
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
				log.Errorf("error creating parent dir for file %s: %v", targetPath, err)
				continue
			}
			outFile, err := os.Create(targetPath)
			if err != nil {
				log.Errorf("error creating file %s: %v", targetPath, err)
				continue
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				log.Errorf("error writing file %s: %v", targetPath, err)
			}
			_ = outFile.Close()
		case tar.TypeSymlink:
			linkPath := sanitizePath(filepath.Join(destinationDir, header.Linkname))
			if err := os.Symlink(linkPath, targetPath); err != nil {
				log.Errorf("error creating symlink %s: %v", targetPath, err)
			}
		case tar.TypeLink:
			linkPath := sanitizePath(filepath.Join(destinationDir, header.Linkname))
			if err := os.Link(linkPath, targetPath); err != nil {
				log.Errorf("error creating hard link %s: %v", targetPath, err)
			}
		default:
			log.Warnf("unknown type: %c in %s", header.Typeflag, header.Name)
		}
	}

	if err := file.Close(); err != nil {
		log.Errorf("error closing tar file: %v", err)
	}

	if err := os.Remove(filePath); err != nil {
		log.Errorf("error deleting tar file: %v", err)
	}

	return nil
}

// sanitizePath replaces illegal characters in file names for cross-platform compatibility.
func sanitizePath(p string) string {
	replacer := strings.NewReplacer(
		":", "_",
		"*", "_",
		"?", "_",
		"\"", "_",
		"<", "_",
		">", "_",
		"|", "_",
	)
	return replacer.Replace(p)
}

func IsTTY() bool {
	// some output needs to be run in the terminal of the tty, such as spinner
	return terminal.IsTerminal(int(os.Stdout.Fd()))
}
