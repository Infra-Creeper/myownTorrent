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
		ipAddr, ipErr := GetLocalIP()
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
