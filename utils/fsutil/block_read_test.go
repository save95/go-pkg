package fsutil

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pkg/errors"
)

var (
	// mac/linux 使用以下命令直接生成大文件
	// dd if=/dev/zero of=big-file bs=1024 count=100000
	_bigFile = "big-file"
)

func TestBlockRead(t *testing.T) {
	// 打开待写入的新文件
	newFilename := "block-read-file-test.txt"
	fd, err := os.OpenFile(newFilename, os.O_WRONLY|os.O_CREATE, 0666)
	assert.Nil(t, err)
	defer func() {
		_ = fd.Close()
	}()

	// 读文件写入新文件
	err = BlockRead(_bigFile, func(data []byte) error {
		if _, err := fd.Write(data); nil != err {
			return errors.Wrap(err, "write failed")
		}
		return nil
	})
	assert.Nil(t, err)

	// 比较两个文件大小
	ofi, err := os.Stat(_bigFile)
	assert.Nil(t, err)
	dfi, err := os.Stat(newFilename)
	assert.Nil(t, err)
	assert.Equal(t, ofi.Size(), dfi.Size())

	// 清理文件
	_ = os.Remove(newFilename)
}
