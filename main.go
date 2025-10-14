package main

import (
	"fmt"
	"myownTorrent/createTFile"
)

func main() {
	fname := "shorttext.txt"
	err := createTFile.CreateTorrent(fname)
	if err != nil {
		fmt.Println("ERROR", err)
	} else {
		fmt.Println("Files created succesfully")
	}
}
