package storage

const storageRoot = "storage"

type IStorage interface {
	Dir() string      // 获得文件存储的目录
	Path() string     // 获得文件存储的路径
	Filename() string // 获得文件名
}

type IPrivateStorage interface {
	IStorage

	AppendDir(dirs ...string) IPrivateStorage // 追加存储目录
	SetName(name string) IPrivateStorage      // 设置文件名
}

type IPublicStorage interface {
	IStorage

	AppendDir(dirs ...string) IPublicStorage // 追加存储目录
	SetName(name string) IPublicStorage      // 设置文件名

	URL() string                    // 获得文件的访问链接（不含 host）
	URLWithHost(host string) string // 获得文件的访问链接（含 host）
}
