package createTFile

import (
	"fmt"
	"path/filepath"
	"strconv"
)

// returns (dirname, incomplete bin file name) from the given file name
func getFolderString(fname string) (string, string) {
	dirname := fmt.Sprintf("TRRNT[%s]", fname)

	return dirname, filepath.Join(dirname, fname)
}

// returns the filename of the torrent file from the file
func getTorrentFileName(fname string) string {
	return fmt.Sprintf("%s.TRRNTjson", fname)
}

// returns the filename of the 'index'th bin file from the actual file name
func getBinPiece(filename string, index int) string {
	dir, _ := getFolderString(filename)
	var out string = filepath.Join(dir, filename+"_x_"+strconv.Itoa(index)+".bin")
	return out
}
