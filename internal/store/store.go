package store

import (
	"errors"
	"io"
	"os"

	"github.com/rs/zerolog/log"
)

type S interface {
	StoreArticle(articleId string, fn func(ctx StoreArticleContext) error) error
}

type StoreArticleContext interface {
	For(dirName string) WriteArticleContext
}

type WriteArticleContext interface {
	Write(fileName string, fn func(w io.Writer) error) error
}

func New(rootDir string) (S, error) {
	err := os.MkdirAll(rootDir, 0)
	if err != nil && !errors.Is(err, os.ErrExist) {
		log.Error().Str("dir", rootDir).Err(err).Msg("unable to create output directory")
		return nil, err
	}

	log.Debug().Str("dir", rootDir).Msg("output directory")
	return &store{rootDir: rootDir}, nil
}
