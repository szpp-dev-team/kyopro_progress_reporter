package util

import "os"

func Exists(name string) bool {
	_, err := os.Stat(name)
	return err == nil
}
