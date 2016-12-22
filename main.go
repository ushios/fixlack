package main

import "github.com/ushios/fixlack/lib/fixlack"

func main() {
	if err := fixlack.Fixlack("/Users/ushio/Downloads"); err != nil {
		panic(err)
	}
}
