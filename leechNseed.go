package main

import (
	"encoding/hex"
	"myownTorrent/manageTFile"
	"time"

	"github.com/purehyperbole/dht"
)

const ttl time.Duration = time.Duration(5 * time.Minute)

func downloadAllPieces(torrentfilename string, node *dht.DHT) error {
	metadata, Scanerr := manageTFile.ScanTFile(torrentfilename)
	if Scanerr != nil {
		return Scanerr
	}
	var pieceError error
	for _, hash := range metadata.Hashes {
		hashbyte, decodeErr := hex.DecodeString(hash)
		if decodeErr != nil {
			return decodeErr
		}

	}
	return nil
}

func downloadPiece(pieceNumber int) error {

}
