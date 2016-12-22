package fixlack

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

var (
	ErrPathIsNotDirectory = errors.New("ErrPathIsNotDirectory")
	unicodeRegexp         *regexp.Regexp
)

func init() {
	unicodeRegexp = regexp.MustCompile("(u([0-9a-fA-F]{4}))")
}

func Fixlack(path string) error {
	_, err := directory(path)
	if err != nil {
		return err
	}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			log.Println(file.Name(), "is directory")
			continue
		}

		switch file.Mode() {
		case os.ModeSymlink:
			log.Println(file.Name(), "is symblic link")
			continue
		}

		if err := fixlack(path, file); err != nil {
			return err
		}
	}

	return nil
}

func fixlack(path string, file os.FileInfo) error {
	b := []byte(file.Name())
	if unicodeRegexp.Match(b) {
		fullpath := filepath.Join(path, file.Name())
		origin, err := os.Open(fullpath)
		if err != nil {
			return err
		}
		defer origin.Close()

		destName := unicodeRegexp.ReplaceAllString(file.Name(), `\u$2`)
		fmt.Println("destName", destName)
		fmt.Println("Go\u8a00\u8a9e\u306b\u3088\u308bWeb\u30a2\u30d5\u309a\u30ea\u30b1\u30fc\u30b7\u30e7\u30f3\u958b\u767a.pdf")

		log.Println("path: ", fullpath)
	}

	return nil
}

func directory(path string) (os.FileInfo, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if !info.IsDir() {
		return nil, ErrPathIsNotDirectory
	}

	return info, nil
}
