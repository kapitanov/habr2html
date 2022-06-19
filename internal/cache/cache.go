package cache

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

type C interface {
	For(articleID string) W
	Close()
}

type W interface {
	Write(name string, content []byte) (string, error)
}

func New(rootDir string, keepCacheFiles bool) (C, error) {
	return &cache{
		rootDir:        rootDir,
		keepCacheFiles: keepCacheFiles,
	}, nil
}

type cache struct {
	rootDir        string
	keepCacheFiles bool
	dirs           []string
}

func (c *cache) For(articleID string) W {
	return &cacheBranch{
		cache: c,
		dir:   filepath.Join(c.rootDir, articleID),
	}
}

func (c *cache) Close() {
	if c.keepCacheFiles {
		return
	}

	for _, dir := range c.dirs {
		err := os.RemoveAll(dir)
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			log.Error().Str("dir", dir).Err(err).Msg("unable to remove directory")
		}
	}
}

type cacheBranch struct {
	cache *cache
	dir   string
}

func (b *cacheBranch) Write(name string, content []byte) (string, error) {
	err := os.MkdirAll(b.dir, 0)
	if err != nil {
		log.Error().Str("dir", b.dir).Err(err).Msg("unable to create directory")
		return "", err
	}

	path := filepath.Join(b.dir, name)
	path, err = filepath.Abs(path)
	if err != nil {
		log.Error().Str("path", path).Err(err).Msg("unable to write file")
		return "", err
	}

	err = ioutil.WriteFile(path, content, 0666)
	if err != nil {
		log.Error().Str("path", path).Err(err).Msg("unable to write file")
		return "", err
	}

	return path, nil
}
