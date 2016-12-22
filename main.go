package main

import (
	"os"
	"path/filepath"

	"github.com/ushios/fixlack/lib/fixlack"
)

func main() {
	home := os.Getenv("HOME")
	if err := fixlack.Fixlack(filepath.Join(home, "Downloads")); err != nil {
		panic(err)
	}
}
