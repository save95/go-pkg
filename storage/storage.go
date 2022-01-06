package storage

import (
	"os"
	"path"
)

type storage struct {
	root []string

	dirs []string
	name string
}

func newStorage(dir string) *storage {
	return &storage{
		root: []string{storageRoot, dir},
		dirs: make([]string, 0),
	}
}

func newTempStorage() *storage {
	return &storage{
		root: []string{os.TempDir(), "go-pkg"},
		dirs: make([]string, 0),
	}
}

func (p *storage) AppendDir(dirs ...string) IPrivateStorage {
	p.dirs = append(p.dirs, dirs...)

	return p
}

func (p *storage) SetName(name string) IPrivateStorage {
	if len(name) > 0 {
		p.name = name
	}

	return p
}

func (p *storage) Dir() string {
	paths := p.root[:]
	paths = append(paths, p.dirs...)

	return path.Join(paths...)
}

func (p *storage) Filename() string {
	return p.name
}

func (p *storage) Path() string {
	paths := p.root[:]
	paths = append(paths, p.dirs...)
	paths = append(paths, p.name)

	return path.Join(paths...)
}
