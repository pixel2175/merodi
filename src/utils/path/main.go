package path

import (
	. "merodi/src/utils/errs"
	"os"
	"path/filepath"
)

func CopyFile(src, dst string) (err error) {
	defer Handle(&err)

	info := CheckV(os.Stat(src))
	data := CheckV(os.ReadFile(src))
	CheckE(os.MkdirAll(filepath.Dir(dst), 0755))
	CheckE(os.WriteFile(dst, data, info.Mode().Perm()))
	CheckE(os.Chmod(dst, info.Mode().Perm()))

	return
}

func Getwd() (dir string) {
	var err error
	defer Handle(&err, func(*error) { dir = filepath.Clean(".") })
	dir = CheckV(os.Getwd())
	return
}
