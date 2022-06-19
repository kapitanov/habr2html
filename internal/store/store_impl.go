package store

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
)

type store struct {
	rootDir string
}

func (s *store) StoreArticle(articleId string, fn func(ctx StoreArticleContext) error) error {
	return s.withRootInfo(func(info *rootInfoJSON, infoFilePath string) error {
		if _, ok := info.Articles[articleId]; ok {
			log.Debug().Str("id", articleId).Msg("article is already exported")
			return nil
		}

		log.Debug().Str("id", articleId).Msg("exporting article")
		err := fn(storeArticleContext{articleId: articleId, rootDir: s.rootDir})
		if err != nil {
			return err
		}

		info.Articles[articleId] = &articleInfoJSON{}

		err = info.Save(infoFilePath)
		if err != nil {
			return err
		}

		return nil
	})
}

func (s *store) withRootInfo(fn func(info *rootInfoJSON, infoFilePath string) error) error {
	infoFilePath := filepath.Join(s.rootDir, RootInfoFileName)
	info, err := loadRootInfoJSON(infoFilePath)
	if err != nil {
		return err
	}

	err = fn(info, infoFilePath)
	if err != nil {
		return err
	}

	return nil
}

type storeArticleContext struct {
	articleId string
	rootDir   string
}

func (ctx storeArticleContext) For(dirName string) WriteArticleContext {
	dirName = normalize(dirName)
	dirName = filepath.Join(ctx.rootDir, dirName)

	return writeArticleContext{
		articleId: ctx.articleId,
		dirName:   dirName,
	}
}

type writeArticleContext struct {
	articleId string
	dirName   string
}

func (ctx writeArticleContext) Write(fileName string, fn func(w io.Writer) error) error {
	if len(fileName) > 96 {
		ext := filepath.Ext(fileName)
		fileName = fileName[:96-len(ext)-1] + ext
	}

	fileName = normalize(fileName)
	fileName = filepath.Join(ctx.dirName, fileName)

	err := os.MkdirAll(ctx.dirName, 0)
	if err != nil {
		log.Error().Str("id", ctx.articleId).Str("path", fileName).Err(err).Msg("unable to write article")
		return err
	}

	f, err := os.Create(fileName)
	if err != nil {
		log.Error().Str("id", ctx.articleId).Str("path", fileName).Err(err).Msg("unable to write article")
		return err
	}
	defer f.Close()

	err = fn(f)
	if err != nil {
		return err
	}

	log.Trace().Str("id", ctx.articleId).Str("path", fileName).Err(err).Msg("new file written")
	return nil
}

func normalize(str string) string {
	str = strings.ReplaceAll(str, "/", "-")
	str = strings.ReplaceAll(str, ":", "-")
	return str
}
