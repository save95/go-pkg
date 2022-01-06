package storage

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPublicFrom(t *testing.T) {
	filenames := []string{
		"storage/public/abc/storage/def/1.png",
		"/storage/public/abc/storage/def/1.png",
	}
	path := "storage/public/abc/storage/def/1.png"
	dir := "storage/public/abc/storage/def"
	URL := "/storage/abc/storage/def/1.png"
	host := "https://wwww.domain.com/abc"
	for _, filename := range filenames {
		p := PublicFromFile(filename)

		assert.Equal(t, path, p.Path())
		assert.Equal(t, dir, p.Dir())
		assert.Equal(t, URL, p.URL())
		assert.Equal(t, "1.png", p.Filename())
		assert.True(t, strings.HasPrefix(p.URLWithHost(host), host))
	}

	p2 := PublicFromUrl(URL)
	assert.Equal(t, path, p2.Path())
	assert.Equal(t, URL, p2.URL())
	assert.Equal(t, dir, p2.Dir())
}
