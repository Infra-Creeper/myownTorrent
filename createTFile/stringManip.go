package createTFile

import (
	"fmt"
	"path/filepath"
)

func getFolderString(fname string) (string, string) {
	dirname := fmt.Sprintf("TRRNT[%s]", fname)

	return dirname, filepath.Join(dirname, fname)
}
