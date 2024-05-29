package camera

import (
	"log"
	"os"
)

var tmpdir string

func SetupTempDir() error {
	dir, err := os.MkdirTemp("", os.Getenv("STILL_IMG_DIR"))
	if err != nil {
		log.Println("Error creating temp dir:", err)
		return err
	}
	log.Println("Created temp dir:", dir)

	tmpdir = dir
	return nil
}
