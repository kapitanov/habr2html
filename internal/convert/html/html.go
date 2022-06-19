package html

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gabriel-vasile/mimetype"
	"github.com/rs/zerolog/log"

	"github.com/kapitanov/habr2html/internal/habr"
	"github.com/kapitanov/habr2html/internal/store"
)

func Do(ctx store.WriteArticleContext, article *habr.Article, title, description string) error {
	html := fmt.Sprintf(
		template,
		escape(title),
		escape(title),
		escape(article.Author.FullName),
		escape(article.ID),
		article.TextHTML)

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return err
	}

	doc.Find("*[xmlns]").Each(func(_ int, el *goquery.Selection) {
		val, _ := el.Attr("xmlns")
		if val == "http://www.w3.org/1999/xhtml" {
			el.RemoveAttr("xmlns")
		}
	})

	err = saveImages(ctx, doc)
	if err != nil {
		return err
	}

	html, err = doc.Html()
	if err != nil {
		return err
	}

	filename := fmt.Sprintf("%s.html", title)
	err = ctx.Write(filename, func(dst io.Writer) error {
		_, err := dst.Write([]byte(html))
		return err
	})
	if err != nil {
		return err
	}

	return nil
}

func escape(s string) string {
	var w bytes.Buffer
	xml.Escape(&w, []byte(s))
	return w.String()
}

const template = `<!DOCTYPE html>
<html lang="ru">
<head>
	<meta charset="UTF-8">
	<title>%s</title>
	<style type="text/css">
		body {
			font-family: Sans-serif;
		}
		code {
			font-style: monospace;
		}
</style>
</head>
<body>
	<header>
		<h1>%s</h1>
		<h2>%s</h2>
		<p>
			<a href="https://habr.com/ru/post/%s/">[ Ссылка на статью ]</a>
		</p>
	</header>
	<hr />
	<article>
		%s
	</article>
</body>
</html>
`

func saveImages(ctx store.WriteArticleContext, doc *goquery.Document) error {
	var err error = nil
	doc.Find("img").Each(func(_ int, img *goquery.Selection) {
		if err != nil {
			return
		}

		img.RemoveAttr("loading")
		img.RemoveAttr("srcset")

		src, _ := img.Attr("src")
		if !strings.HasPrefix(src, "http") {
			return
		}

		localFile, e := addImage(ctx, src)
		if e != nil {
			log.Warn().Str("src", src).Err(err).Msg("unable to load image")
			return
		}

		img.SetAttr("src", localFile)
	})

	return err
}

func addImage(ctx store.WriteArticleContext, sourceURL string) (string, error) {
	resp, err := http.Get(sourceURL)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("request \"GET %s\" failed with \"%s\"", sourceURL, resp.Status)
	}

	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	m := mimetype.Detect(bytes)

	h := sha1.New()
	_, err = h.Write(bytes)
	if err != nil {
		return "", err
	}

	filename := hex.EncodeToString(h.Sum(nil))
	filename = fmt.Sprintf("%s%s", filename, m.Extension())

	err = ctx.Write(filename, func(w io.Writer) error {
		_, err := w.Write(bytes)
		return err
	})
	if err != nil {
		return "", err
	}

	return filename, nil
}
