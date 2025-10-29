package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"myownTorrent/TorrentNet"
	"myownTorrent/manageTFile"
	"os"
	"strings"
	"time"

	"github.com/purehyperbole/dht"
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
	bootstrap := flag.String("b", "", "Bootstrap node address (e.g., 127.0.0.1:9000)")
	port := flag.String("p", "9000", "port to start DHT node at")

	// Check if subcommand is provided
	if len(os.Args) < 2 {
		fmt.Println("Usage: program <host|join|split> [options]")
		os.Exit(1)
	}
	if os.Args[1] == "net" {
		if ipErr != nil {
			panic(ipErr)
		}
		//fmt.Println("Hosting at", &ipAddr)
		go TorrentNet.StartServer(*ipAddr, "8080", ".", nil)
		storage := TorrentNet.NewCustomStorage()
		cfg := &dht.Config{
			ListenAddress: *ipAddr + ":" + *port,
			Listeners:     4,
			Timeout:       time.Minute / 2,
			Storage:       storage,
		}

		if *bootstrap != "" {
			cfg.BootstrapAddresses = []string{*bootstrap}
		}
		reader := bufio.NewReader(os.Stdin)
		node, err := dht.New(cfg)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Node Started listenting at", ipAddr)

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
				allerr := downloadAllPieces(tfilename, node, *ipAddr)
				if allerr != nil {
					log.Println("ERROR DOWNLOADING", allerr)
					continue
				}
				log.Println("All files Downloaded Successfully")
			}
			if strings.HasPrefix(text, "post ") {
				tfilename := strings.TrimPrefix(text, "post ")
				posterr := PostTorrentFile(node, tfilename, *ipAddr)
				if posterr != nil {
					log.Println("ERROR POSTING", posterr)
					continue
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
