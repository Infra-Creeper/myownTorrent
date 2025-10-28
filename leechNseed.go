package main

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"myownTorrent/TorrentNet"
	"myownTorrent/manageTFile"
	"time"

	"github.com/purehyperbole/dht"
)

const (
	ttl time.Duration = time.Duration(5 * time.Minute)
)

func downloadAllPieces(torrentfilename string, node *dht.DHT) error {
	metadata, Scanerr := manageTFile.ScanTFile(torrentfilename)
	if Scanerr != nil {
		return Scanerr
	}
	downloadedPieces := make([]bool, metadata.Pieces)
	//var pieceError error
	var i int = 0
	var leastPopPiece int = i
	for isAllTrue(downloadedPieces) {
		if i > metadata.Pieces {
			i = 0
			downloadedPieces[leastPopPiece] = true
			hash := metadata.Hashes[leastPopPiece]
			downErr := downloadPiece(i, hash, metadata.Name, node)
			if downErr != nil {
				eMsg := fmt.Sprintf("ERROR DOWNLOADING PIECE %d \n %v", leastPopPiece, downErr)
				return errors.New(eMsg)
			}
			continue
		}
		seeds, err := TorrentNet.GetSeeds(node, metadata.Hashes[i], 5*time.Second)
		if err != nil {
			fmt.Println("Unable to get seeds info for piece", i)
			downloadedPieces[i] = true
		}
		if !downloadedPieces[i] && len(seeds) < leastPopPiece {
			leastPopPiece = i
		}
		i++
	}
	return nil
}

func downloadPiece(index int, hash string, filename string, node *dht.DHT) error {
	seeds, seedsErr := TorrentNet.GetSeeds(node, hash, 10*time.Second)
	if seedsErr != nil {
		return seedsErr
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
