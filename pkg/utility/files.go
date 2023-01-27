package utility

import (
	"os"
)

func MakeDirIfNotExists(dirpath string) error {
	if _, err := os.Stat(dirpath); os.IsNotExist(err) {
		err := os.Mkdir(dirpath, os.ModeDir|os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}
