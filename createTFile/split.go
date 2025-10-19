package createTFile

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
)

type TFile struct {
	Name   string   `json:"name"`
	Length int      `json:"length"`
	Pieces int      `json:"pieces"`
	Hashes []string `json:"hashes"`
}

const pieceSize int = 16000

// creates the torrent file of the given filename
func CreateTorrent(fileName string) error {
	fobj, ferr := os.Open(fileName)
	var metadata TFile
	metadata.Name = fileName
	defer fobj.Close()
	if ferr != nil {
		return ferr
	}
	dir, fileLoc := getFolderString(fileName)

	dirErr := os.Mkdir(dir, 0755)
	if dirErr != nil {
		return dirErr
	}

	pieceBuf := make([]byte, pieceSize)
	for {
		n, err := fobj.Read(pieceBuf)
		if err == io.EOF {
			break // end of file
		}
		if err != nil {
			return err
		}
		createPiece(pieceBuf[:n], metadata.Pieces, fileLoc)
		hashBytes := sha1.Sum(pieceBuf[:n])
		var hashStr string = hex.EncodeToString(hashBytes[:])

		metadata.Length += n
		metadata.Pieces++
		metadata.Hashes = append(metadata.Hashes, hashStr)
	}
	//fmt.Printf("%v\n", metadata)
	metaErr := createMeta(metadata)
	if metaErr != nil {
		return metaErr
	}
	return nil
}

// creates a piece from the data and file name
func createPiece(data []byte, pid int, filename string) error {
	var binName string = filename + "_x_" + strconv.Itoa(pid) + ".bin"
	err := os.WriteFile(binName, data, 0644)
	if err != nil {
		return errors.New("ERROR: Unable to create piece")
	}

	return nil
}

// creates JSON meta data of the struct
func createMeta(tfile TFile) error {
	// Marshal struct to JSON (pretty print)
	data, err := json.MarshalIndent(tfile, "", "  ")
	if err != nil {
		return fmt.Errorf("error encoding struct: %v", err)
	}
	var filename string = getTorrentFileName(tfile.Name)

	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return fmt.Errorf("error writing file: %v", err)
	}

	return nil
}
