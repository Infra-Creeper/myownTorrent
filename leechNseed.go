package main

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"myownTorrent/TorrentNet"
	"myownTorrent/manageTFile"
	"time"
	//"github.com/purehyperbole/dht"
)

const (
	ttl      time.Duration = time.Duration(5 * time.Minute)
	maxPeers int           = 6
)

func downloadAllPieces(torrentfilename string, URL string, ipAddr string) error {
	metadata, Scanerr := manageTFile.ScanTFile(torrentfilename)
	if Scanerr != nil {
		return Scanerr
	}
	for i, hash := range metadata.Hashes {
		err := downloadPiece(i, hash, metadata.Name, URL)
		if err != nil {
			fmt.Printf("ERROR DOWNLOADING PIECES %d\n", i)
		}
	}
	return nil
}

func downloadPiece(index int, hash string, filename string, URL string) error {
	seeds, seedsErr := TorrentNet.GetValues(URL, hash)
	if seedsErr != nil {
		return seedsErr
	}
	if len(seeds) == 0 {
		fmt.Println("ERROR:No seeds found")
		return errors.New("no seeds found")
	}
	seedIP := seeds[rand.IntN(len(seeds))]
	pieceLoc := manageTFile.GetBinPieceFileName(filename, index)
	err := TorrentNet.RequestFile(seedIP, pieceLoc, "")
	return err
}

func isAllTrue(arr []bool) bool {
	for _, v := range arr {
		if !v {
			return false
		}
	}
	return true
}

func PostTorrentFile(URL string, filename string, ipAddr string) error {
	metadata, metaerr := manageTFile.ScanTFile(filename)
	if metaerr != nil {
		return metaerr
	}
	for i, hash := range metadata.Hashes {
		Posterr := TorrentNet.PostKeyValue(URL, hash, ipAddr+":8080")
		if Posterr != nil {
			return fmt.Errorf("ERROR POSTING %d HASH\n %v", i, Posterr)
		}
		fmt.Printf("Posted piece=%d hash=0x%s\n", i, hash)
	}
	return nil
}
