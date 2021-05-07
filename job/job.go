package job

// IJob job 约定
type IJob interface {
	Run() error
}
