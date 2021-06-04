package fsutil

import (
	"os"

	"github.com/pkg/errors"
)

// BlockRead 分块读取大文件，每次读取 4M 内容
func BlockRead(fileName string, handle func([]byte) error) error {
	f, err := os.Open(fileName)
	if err != nil {
		return errors.Errorf("can't opened this file: %s", fileName)
	}
	defer func() {
		_ = f.Close()
	}()

	s := make([]byte, 4096)
	for {
		switch nr, err := f.Read(s[:]); true {
		case nr < 0:
			return errors.Errorf("cat: error reading: %s", err.Error())
			//os.Exit(1)
		case nr == 0: // EOF
			return nil
		case nr > 0:
			if err := handle(s[0:nr]); nil != err {
				return errors.Wrap(err, "handle failed")
			}
		}
	}
}
