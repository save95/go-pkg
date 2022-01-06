package storage

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type public struct {
	root []string

	dirs []string
	name string
}

func newPublic() *public {
	return &public{
		root: []string{storageRoot, "public"},
		dirs: make([]string, 0),
	}
}

func Public() IPublicStorage {
	return newPublic()
}

func PublicFromFile(filename string) IPublicStorage {
	return newPublic().withFile(filename)
}

func PublicFromUrl(fileURL string) IPublicStorage {
	return newPublic().withURL(fileURL)
}

func (p *public) withFile(filename string) *public {
	sep := strings.Join(p.root, string(os.PathSeparator))
	seps := strings.Split(filename, sep)
	paths := strings.Split(strings.TrimLeft(seps[(len(seps)-1)], string(os.PathSeparator)), string(os.PathSeparator))

	p.dirs = paths[:len(paths)-1]
	p.name = paths[len(paths)-1]

	return p
}

func (p *public) withURL(fileURL string) *public {
	urls := strings.SplitN(fileURL, fmt.Sprintf("%s/", storageRoot), 2)
	if len(urls) != 2 {
		return p
	}

	paths := strings.Split(urls[1], "/")
	p.dirs = paths[:len(paths)-1]
	p.name = paths[len(paths)-1]

	return p
}

func (p *public) AppendDir(dirs ...string) IPublicStorage {
	p.dirs = append(p.dirs, dirs...)

	return p
}

func (p *public) SetName(name string) IPublicStorage {
	if len(name) > 0 {
		p.name = name
	}

	return p
}

func (p *public) Dir() string {
	paths := p.root[:]
	paths = append(paths, p.dirs...)

	return path.Join(paths...)
}

func (p *public) Filename() string {
	return p.name
}

func (p *public) Path() string {
	paths := p.root[:]
	paths = append(paths, p.dirs...)
	paths = append(paths, p.name)

	return path.Join(paths...)
}

func (p *public) URL() string {
	paths := p.root[2:]
	paths = append(paths, p.dirs...)
	paths = append(paths, p.name)

	return fmt.Sprintf("/%s/%s", storageRoot, filepath.ToSlash(path.Join(paths...)))
}

func (p *public) URLWithHost(host string) string {
	if !strings.HasPrefix(host, "http") {
		return p.URL()
	}

	return fmt.Sprintf("%s%s", host, p.URL())
}
