package fsutil

import (
	"io/fs"
	"log"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var _tree = []string{
	"test/1.md",
	"test/2.md",
	"test/3.md",
	"test/a/1.md",
	"test/a/a/1.md",
	"test/a/a/2.md",
	"test/a/b/1.md",
	"test/a/c/1.md",
	"test/b/a/b/1.md",
	"test/test/a/b/1.md",
}

func _makeFileTree() {
	for _, s := range _tree {
		dir := path.Dir(s)
		if len(dir) > 0 {
			_ = os.MkdirAll(dir, fs.ModePerm)
		}

		if _, err := os.Create(s); nil != err {
			log.Fatalln(err)
		}
	}
}

func _clear() {
	_ = os.RemoveAll("test")
	_ = os.RemoveAll("test2")
}

func TestExist(t *testing.T) {
	assert.True(t, Exist("filepath.go"))
	assert.False(t, Exist("filepath.go.no-exist"))
	assert.False(t, Exist("../maker"))
}

func TestPathExist(t *testing.T) {
	assert.True(t, PathExist("../maker"))
	assert.False(t, PathExist("../no-exist"))
	assert.False(t, PathExist("../maker/trans_no.go"))
}

func TestClear(t *testing.T) {
	_makeFileTree()

	err := Clear("test")
	assert.Nil(t, err)

	for _, s := range _tree {
		_, err := os.Stat(s)
		assert.True(t, os.IsNotExist(err))
	}
}

func TestCopy(t *testing.T) {
	_makeFileTree()

	total, err := Copy(_tree[0], "test2/copied/1")
	assert.Nil(t, err)
	info, err := os.Stat(_tree[0])
	assert.Nil(t, err)
	assert.Equal(t, total, info.Size())

	_, err = Copy(_tree[0], "test2/copied")
	assert.NotNil(t, err)

	_clear()
}

func TestCopyPath(t *testing.T) {
	_makeFileTree()

	err := CopyPath("test", "test2")
	assert.Nil(t, err)

	assert.True(t, PathExist("test2"))
	for _, s := range _tree {
		assert.True(t, Exist(strings.Replace(s, "test", "test2", 1)))
	}

	_clear()
}
