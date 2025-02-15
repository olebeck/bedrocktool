package utils

import (
	"archive/zip"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func UnpackZip(r io.ReaderAt, size int64, unpackFolder string) {
	zr, _ := zip.NewReader(r, size)
	for _, srcFile := range zr.File {
		srcName := strings.ReplaceAll(srcFile.Name, "\\", "/")
		outPath := path.Join(unpackFolder, srcName)
		if srcFile.Mode().IsDir() {
			os.Mkdir(outPath, 0o755)
		} else {
			os.MkdirAll(path.Dir(outPath), 0o755)
			fr, _ := srcFile.Open()
			f, _ := os.Create(path.Join(unpackFolder, srcName))
			io.Copy(f, fr)
		}
	}
}

func ZipFolder(filename, folder string) error {
	f, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	zw := zip.NewWriter(f)
	err = filepath.WalkDir(folder, func(path string, d fs.DirEntry, err error) error {
		if !d.Type().IsDir() {
			rel := path[len(folder)+1:]
			zwf, _ := zw.Create(rel)
			data, err := os.ReadFile(path)
			if err != nil {
				fmt.Println(err)
			}
			zwf.Write(data)
		}
		return nil
	})
	zw.Close()
	f.Close()
	return err
}
