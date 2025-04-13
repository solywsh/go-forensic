package osx

import (
	"archive/tar"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"os"
	"path/filepath"
)

//func TarDecompression(filePath, destinationDir string) error {
//	cmd := exec.Command("tar", "-xf", filePath, "-C", destinationDir)
//	cmd.Run()
//	if err := os.Remove(filePath); err != nil {
//		fmt.Printf("Failed to delete local tar file: %s\n", err)
//	}
//	return nil
//}

func TarDecompression(filePath, destinationDir string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("could not open tar file: %w", err)
	}
	defer file.Close()
	// create tar reader
	tarReader := tar.NewReader(file)
	for {
		// read the next file header
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("could not read tar header: %w", err)
		}
		// generate the target path
		targetPath := filepath.Join(destinationDir, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			// create the directory if it doesn't exist
			if err := os.MkdirAll(targetPath, os.FileMode(header.Mode)); err != nil {
				return fmt.Errorf("could not create directory: %w", err)
			}
		case tar.TypeReg:
			// create the file
			if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
				return fmt.Errorf("could not create directory for file: %w", err)
			}
			outFile, err := os.Create(targetPath)
			if err != nil {
				return fmt.Errorf("could not create file: %w", err)
			}
			// copy the file content
			if _, err := io.Copy(outFile, tarReader); err != nil {
				return fmt.Errorf("could not copy file content: %w", err)
			}
			outFile.Close()
		case tar.TypeSymlink:
			// handle symlink
			linkPath := filepath.Join(destinationDir, header.Linkname)
			if err := os.Symlink(linkPath, targetPath); err != nil {
				return fmt.Errorf("could not create symlink: %w", err)
			}
		case tar.TypeLink:
			// handle hard link
			linkPath := filepath.Join(destinationDir, header.Linkname)
			if err := os.Link(linkPath, targetPath); err != nil {
				return fmt.Errorf("could not create hard link: %w", err)
			}
		default:
			fmt.Printf("Encountered unknown type: %c in %s\n", header.Typeflag, header.Name)
		}
	}
	// remove the tar file after extraction
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("failed to delete local tar file: %s\n", err)
	}
	return nil
}

func IsTTY() bool {
	// some output needs to be run in the terminal of the tty, such as spinner
	return terminal.IsTerminal(int(os.Stdout.Fd()))
}
