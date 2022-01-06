# 本地文件存储 Storage

`storage` 包是基于项目的文件存储管理工具。默认有公开文件管理工具 `storage.Public()`，和内部文件管理工具 `storage.Disk(string)`，
以及系统临时文件管理工具 `storage.Temp()`。

公开文件和内部文件的存储，均相对于项目的根目录（运行时目录）。而，临时文件则存储于系统的临时目录，受操作系统影响。

## 基础接口 IStorage

```
	Dir() string      // 获得文件存储的目录
	Path() string     // 获得文件存储的路径
	Filename() string // 获得文件名
```


