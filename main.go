package main

import (
	"flag"
	"fmt"
	"myownTorrent/TorrentNet"
	"myownTorrent/manageTFile"
	"os"
)

func main() {
	//netCmd := flag.NewFlagSet("net", flag.ExitOnError)
	joinCmd := flag.NewFlagSet("join", flag.ExitOnError)
	splitCmd := flag.NewFlagSet("split", flag.ExitOnError)

	// Flags for file subcommand
	filename := splitCmd.String("f", "", "Input filename (required)")

	tfile := joinCmd.String("t", "", "Name of Torrent file (required)")
	outfile := joinCmd.String("o", "", "Output file name")

	// Check if subcommand is provided
	if len(os.Args) < 2 {
		fmt.Println("Usage: program <host|join|split> [options]")
		os.Exit(1)
	}
	if os.Args[1] == "net" {
		ipAddr, ipErr := TorrentNet.GetLocalIP()
		if ipErr != nil {
			panic(ipErr)
		}
		fmt.Println("Hosting at", ipAddr)
	} else if os.Args[1] == "join" {
		joinCmd.Parse(os.Args[2:])
		err := manageTFile.JoinTorrentPieces(*tfile, *outfile)
		if err != nil {
			println("LOG: tfile=", *tfile, "outfile=", *outfile)
			panic(err)
		}

	} else if os.Args[1] == "split" {
		splitCmd.Parse(os.Args[2:])
		err := manageTFile.CreateTorrent(*filename)
		if err != nil {
			println("LOG: filename=", *filename)
			panic(err)
		}
	} else {
		fmt.Println("Invalid Arguments passed")
	}
}
