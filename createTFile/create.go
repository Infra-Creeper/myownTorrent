package createTFile

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
)

type TFile struct {
	name   string
	length int
	pieces int
}

const pieceSize int = 128

func CreateTorrent(fileName string) error {
	fobj, ferr := os.Open(fileName)
	metadata := TFile{fileName, 0, 0}
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
		createPiece(pieceBuf, metadata.pieces, fileLoc)
		metadata.length += n
		metadata.pieces++
	}
	fmt.Printf("%v", metadata)
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
