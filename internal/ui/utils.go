package ui

import "os"

func removeFile(path string) error {
	return os.Remove(path)
}