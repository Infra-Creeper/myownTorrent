package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"myownTorrent/TorrentNet"
	"myownTorrent/manageTFile"
	"os"
	"strconv"
	"strings"
)

func main() {
	defaultIpAddr, ipErr := TorrentNet.GetLocalIP()
	manageTFile.LocalAddr = defaultIpAddr

	netCmd := flag.NewFlagSet("net", flag.ExitOnError)
	joinCmd := flag.NewFlagSet("join", flag.ExitOnError)
	splitCmd := flag.NewFlagSet("split", flag.ExitOnError)

	// Flags for file subcommand
	filename := splitCmd.String("f", "", "Input filename (required)")

	tfile := joinCmd.String("t", "", "Name of Torrent file (required)")
	outfile := joinCmd.String("o", "", "Output file name")

	//flags for bet subcommands
	ipAddr := netCmd.String("ip", defaultIpAddr, "IP Address of this machine(optional)")
	bootstrap := netCmd.String("b", "http://10.0.0.1:8080", "Bootstrap node address (e.g., 127.0.0.1:9000)")
	//port := netCmd.String("p", "9000", "port to start DHT node at")

	// Check if subcommand is provided
	if len(os.Args) < 2 {
		fmt.Println("Usage: program <net|join|split> [options]")
		os.Exit(1)
	}
	if os.Args[1] == "net" {
		if ipErr != nil {
			panic(ipErr)
		}
		//fmt.Println("Hosting at", &ipAddr)
		go TorrentNet.StartServer(*ipAddr, "8080", ".", nil)

		// go func() {
		// 	ticker := time.NewTicker(10 * time.Second)
		// 	defer ticker.Stop()
		// 	for range ticker.C {
		// 		storage.PrintAll()
		// 	}
		// }()
		var URL string = *bootstrap
		downloadIP := ""
		reader := bufio.NewReader(os.Stdin)
		for {
			fmt.Print("myownTorrent> ")
			text, _ := reader.ReadString('\n')
			text = strings.TrimSpace(text)

			if text == "q" || text == "quit" {
				log.Println("Exiting...")
				break
			}
			if strings.HasPrefix(text, "get ") {
				tfilename := strings.TrimPrefix(text, "get ")
				allerr := downloadAllPieces(tfilename, URL, *ipAddr)
				if allerr != nil {
					log.Println("ERROR DOWNLOADING", allerr)
					continue
				}
				log.Println("All files Downloaded Successfully")
			}
			if strings.HasPrefix(text, "post ") {
				tfilename := strings.TrimPrefix(text, "post ")
				posterr := PostTorrentFile(URL, tfilename, *ipAddr)
				if posterr != nil {
					log.Println("ERROR POSTING", posterr)
					continue
				}

			}
			if strings.HasPrefix(text, "download ") {
				fileoptions := strings.Split(text, " ")
				filename := fileoptions[1]
				index, err := strconv.Atoi(fileoptions[2])
				if err != nil {
					fmt.Println(err)
					continue
				}
				donwloadLoc := manageTFile.GetBinPieceFileName(filename, index)
				Requesterr := TorrentNet.RequestFile(downloadIP, donwloadLoc, donwloadLoc)
				if Requesterr != nil {
					fmt.Println(Requesterr)
				} else {
					fmt.Println("Piece Downloaded sucessfully")
				}
			}
			if strings.HasPrefix(text, "server ") {
				downloadIP = strings.TrimPrefix(text, "server ")
				fmt.Println("server IP set to", downloadIP)
			}
			if strings.HasPrefix(text, "join ") {
				options := strings.Fields(strings.TrimPrefix(text, "join "))
				joinCmd.Parse(options)
				err := manageTFile.JoinTorrentPieces(*tfile, *outfile)
				if err != nil {
					println("LOG: tfile=", *tfile, "outfile=", *outfile)
					println(err)
				}
			}
			if strings.HasPrefix(text, "split ") {
				options := strings.Fields(strings.TrimPrefix(text, "split "))
				splitCmd.Parse(options)
				err := manageTFile.CreateTorrent(*filename)
				if err != nil {
					println("LOG: filename=", *filename)
					println(err)
				}
			}
		}

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
