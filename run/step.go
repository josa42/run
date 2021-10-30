package run

type Step interface {
	Run(tasks Tasks) (chan struct{}, CancelFunc)
	SetDir(dir string)
}
