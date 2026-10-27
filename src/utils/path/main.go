package path

import (
	. "github.com/pixel2175/errs"
	"os"
	"path/filepath"
)

func CopyFile(src, dst string) (err error) {
	defer Catch(&err)

	info := TryV(os.Stat(src))
	data := TryV(os.ReadFile(src))
	TryE(os.MkdirAll(filepath.Dir(dst), 0755))
	TryE(os.WriteFile(dst, data, info.Mode().Perm()))
	TryE(os.Chmod(dst, info.Mode().Perm()))

	return
}

func Getwd() (dir string) {
	var err error
	defer Catch(&err, func(*error) { dir = filepath.Clean(".") })
	dir = TryV(os.Getwd())
	return
}
