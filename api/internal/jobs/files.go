package jobs

func UploadPath(jobID, filename string) string {
	return "var/uploads/" + jobID + "-" + filename
}

func MaskPath(jobID string) string {
	return "var/masks/" + jobID + ".pgm"
}

func ResultPath(jobID, mode string) string {
	return "var/results/" + jobID + "-" + mode + ".ppm"
}
