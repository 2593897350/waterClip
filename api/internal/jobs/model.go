package jobs

type Job struct {
	ID         string
	Kind       string
	SourcePath string
	MaskPath   string
	ResultPath string
	Status     string
	Mode       string
	Error      string
}
