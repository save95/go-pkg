package ziputil

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// CompressPath 压缩一个指定目录
func CompressPath(dst, src string) error {
	// 创建准备写入的文件
	fw, err := os.Create(dst)
	defer func() {
		_ = fw.Close()
	}()
	if err != nil {
		return err
	}

	// 通过 fw 来创建 zip.Write
	zipW := zip.NewWriter(fw)
	defer func() {
		_ = zipW.Close()
	}()

	return filepath.Walk(src, func(path string, fi os.FileInfo, errBack error) (err error) {
		if errBack != nil {
			return errBack
		}

		// 通过文件信息，创建 zip 的文件信息
		fh, err := zip.FileInfoHeader(fi)
		if err != nil {
			return err
		}

		// 替换文件名中的相对路径
		fh.Name = strings.TrimPrefix(strings.TrimPrefix(path, src), string(filepath.Separator))

		// 这步开始没有加，会发现解压的时候说它不是个目录
		if fi.IsDir() {
			fh.Name += "/"
		}

		// 写入文件信息，并返回一个 Write 结构
		w, err := zipW.CreateHeader(fh)
		if err != nil {
			return
		}

		// 检测，如果不是标准文件就只写入头信息，不写入文件数据到 w
		// 如目录，也没有数据需要写
		if !fh.Mode().IsRegular() {
			return nil
		}

		// 打开要压缩的文件
		fr, err := os.Open(path)
		defer func() {
			_ = fr.Close()
		}()
		if err != nil {
			return err
		}

		// 将打开的文件 Copy 到 w
		if _, err := io.Copy(w, fr); err != nil {
			return err
		}

		// 输出压缩的内容
		//fmt.Printf("成功压缩文件： %s, 共写入了 %d 个字符的数据\n", path, n)

		return nil
	})
}
