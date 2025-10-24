package main

import (
	"errors"
	"flag"
	"fmt"
	"myownTorrent/manageTFile"
	"net"
	"os"
)

func main() {
	//hostCmd := flag.NewFlagSet("host", flag.ExitOnError)
	fileCmd := flag.NewFlagSet("file", flag.ExitOnError)

	// Flags for file subcommand
	filename := fileCmd.String("f", "", "Input filename (required)")
	outfile := fileCmd.String("o", *filename, "Output file name")
	joinFlag := fileCmd.Bool("j", false, "Join torrent pieces")
	splitFlag := fileCmd.Bool("s", false, "Create torrent/split file")

	// Check if subcommand is provided
	if len(os.Args) < 2 {
		fmt.Println("Usage: program <host|file> [options]")
		os.Exit(1)
	}
	if os.Args[1] == "host" {
		ipAddr, ipErr := GetLocalIP()
		if ipErr != nil {
			panic(ipErr)
		}
		fmt.Println("Hosting at", ipAddr)
	} else if os.Args[1] == "file" {
		fileCmd.Parse(os.Args[2:])
		if *filename == "" {
			fmt.Println("Error: -f flag (filename) is required")
			fileCmd.PrintDefaults()
			os.Exit(1)
		}
		if *joinFlag && *splitFlag {
			fmt.Println("ERROR:Both join and split flags are passed")
			os.Exit(1)
		}
		if *joinFlag {
			err := manageTFile.JoinTorrentPieces(*filename, *outfile)
			if err != nil {
				panic(err)
			}
			fmt.Println("Files joined successfully as", *outfile)
		} else if *splitFlag {
			err := manageTFile.CreateTorrent(*filename)
			if err != nil {
				panic(err)
			}
			fmt.Println("Torrent file and pieces created sucessfully")
		} else {
			fmt.Println("No split/join flags passed")
		}

	} else {
		fmt.Println("Invalid Arguments passed")
	}
}

func GetLocalIP() (string, error) {
	// 1) Try interfaces
	ifaces, err := net.Interfaces()
	if err == nil {
		for _, iface := range ifaces {
			// skip down or loopback interfaces
			if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
				continue
			}
			addrs, err := iface.Addrs()
			if err != nil {
				continue
			}
			for _, addr := range addrs {
				var ip net.IP
				switch v := addr.(type) {
				case *net.IPNet:
					ip = v.IP
				case *net.IPAddr:
					ip = v.IP
				}
				if ip == nil {
					continue
				}
				// prefer IPv4
				ip4 := ip.To4()
				if ip4 == nil {
					continue
				}
				return ip4.String(), nil
			}
		}
	}

	// 2) Fallback: use UDP dial to determine outbound IP (no packets are sent).
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", errors.New("could not determine local IP address")
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	if localAddr.IP == nil {
		return "", errors.New("could not determine local IP address")
	}
	return localAddr.IP.String(), nil
}

func isFlagPassed(name string) bool {
	found := false
	flag.VisitAll(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}
