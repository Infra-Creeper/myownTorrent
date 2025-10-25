package manageTFile

import (
	"encoding/json"
	"fmt"
	"os"
)

// joins the pieces to create a new file from the pieces and save it as `saveAs`, gets filename from metadata if saveAs is empty string i.e. ""
func JoinTorrentPieces(torrentfilename string, saveAs string) error {
	metadata, err := ScanTFile(torrentfilename)
	if err != nil {
		fmt.Println("ERROR in reading metadata")
		return err
	}
	fmt.Println("Joining", metadata.Name, "...")
	var data []byte
	for i := 0; i < metadata.Pieces; i++ {
		bindata, err := os.ReadFile(getBinPieceFileName(metadata.Name, i))
		if err != nil {
			fmt.Println("Error occured while joining piece indexed:", i)
			return err
		}
		data = append(data, bindata...)
	}
	var writeName string
	if saveAs != "" {
		writeName = saveAs
	} else {
		writeName = metadata.Name
	}
	errWrite := os.WriteFile(writeName, data, 0644)
	if errWrite != nil {
		fmt.Println("ERROR: Unable to write file")
		return err
	}
	fmt.Println("File Joined sucessfully as", writeName)
	return nil
}

func ScanTFile(filename string) (TFile, error) {
	var tfile TFile

	// Read file
	data, err := os.ReadFile(filename)
	if err != nil {
		return tfile, fmt.Errorf("error reading file: %v", err)
	}

	// Decode JSON into struct
	err = json.Unmarshal(data, &tfile)
	if err != nil {
		return tfile, fmt.Errorf("error decoding JSON: %v", err)
	}

	return tfile, nil

}
