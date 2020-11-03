package system

import (
	"os"
	"path/filepath"
)

func GetPackagePath() (string, error) {
	return filepath.Abs(filepath.Dir(os.Args[0]))
}
