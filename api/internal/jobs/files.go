package jobs

import (
	"io"
	"os"
	"path/filepath"
)

func UploadPath(jobID, filename string) string {
	return "var/uploads/" + jobID + "-" + filename
}

func MaskPath(jobID string) string {
	return "var/masks/" + jobID + ".pgm"
}

func ResultPath(jobID, mode string) string {
	return "var/results/" + jobID + "-" + mode + ".ppm"
}

func SaveFile(path string, source io.Reader) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, source)
	return err
}
