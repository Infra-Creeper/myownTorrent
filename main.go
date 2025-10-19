package main

import (
	"errors"
	"fmt"
	"myownTorrent/createTFile"
	"net"
)

func main() {
	fname := "shorttext.txt"
	ipAddr, ipErr := GetLocalIP()
	if ipErr != nil {
		fmt.Println("ERROR GETTING IP", ipErr)
	}
	fmt.Println("IP Address of the system: ", ipAddr)
	err := createTFile.JoinTorrentPieces(fname+".TRRNTjson", "shortout.txt")
	//err := createTFile.CreateTorrent(fname)
	if err != nil {
		fmt.Println("ERROR:", err)
	} else {
		fmt.Println("Files created succesfully")
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
