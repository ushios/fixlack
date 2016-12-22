package fixlack

import (
	"encoding/hex"
	"errors"
	"io"
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
	unicodeRegexp = regexp.MustCompile("(u[0-9a-fA-F]{4})")
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
			// log.Println(file.Name(), "is directory")
			continue
		}

		switch file.Mode() {
		case os.ModeSymlink:
			// log.Println(file.Name(), "is symblic link")
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

		destFileName := unicodeRegexp.ReplaceAllStringFunc(
			file.Name(),
			func(s string) string {
				h := s[1:5]
				dec, err := hex.DecodeString(h)
				if err != nil {
					panic(err)
				}
				i := int(dec[0])*16*16 + int(dec[1])

				return string(i)
			},
		)

		destFullPath := filepath.Join(path, destFileName)
		dest, err := os.Create(destFullPath)
		if err != nil {
			return err
		}
		defer dest.Close()

		if _, err := io.Copy(dest, origin); err != nil {
			return err
		}

		log.Println("create: ", destFullPath)
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
