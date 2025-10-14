package createTFile

import (
	"errors"
	"io"
	"os"
	"strconv"
)

type TFile struct {
	name   string
	length int
	pieces int
}

const pieceSize int = 64

func CreateTorrent(fileName string) error {
	fobj, ferr := os.Open(fileName)
	metadata := TFile{fileName, 0, 0}
	defer fobj.Close()
	if ferr != nil {
		return ferr
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
		createPiece(pieceBuf, metadata.pieces, fileName)
		metadata.length += n
		metadata.pieces++
	}
	return nil
}

func createPiece(data []byte, pid int, filename string) error {
	var binName string = filename + "_x_" + strconv.Itoa(pid) + ".bin"
	err := os.WriteFile(binName, data, 0644)
	if err != nil {
		return errors.New("ERROR: Unable to create piece")
	}

	return nil
}
