// API client for habr.com

package habr

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/rs/zerolog/log"
)

type Client interface {
	ListFavorites(userId string, page int) (*ArticleListResponse, error)
	GetArticle(articleId string) (*Article, error)
}

func New() Client {
	return &client{}
}

type client struct{}

func (c client) ListFavorites(userId string, page int) (*ArticleListResponse, error) {
	u := fmt.Sprintf("/articles/?user=%s&user_bookmarks=true&fl=ru&hl=ru&page=%d", url.QueryEscape(userId), page)

	var result ArticleListResponse
	err := c.getJSON(u, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c client) GetArticle(articleId string) (*Article, error) {
	u := fmt.Sprintf("/articles/%s", url.PathEscape(articleId))

	var result Article
	err := c.getJSON(u, &result)
	if err != nil {
		if e, ok := err.(httpError); ok && e.Status == 403 {
			log.Error().Str("article", articleId).Msg("article is not available")
			return nil, nil
		}

		return nil, err
	}

	return &result, nil
}

type httpError struct {
	Status  int
	Message string
}

func (e httpError) Error() string {
	return e.Message
}

func (c client) getJSON(u string, v interface{}) error {
	u = fmt.Sprintf("https://habr.com/kek/v2%s", u)

	resp, err := http.Get(u)
	if err != nil {
		log.Error().Str("url", u).Err(err).Msg("unable to execute http request")
		return err
	}

	if resp.StatusCode != 200 {
		log.Error().Str("url", u).Int("status", resp.StatusCode).Msg("http request executed")

		return httpError{
			Status:  resp.StatusCode,
			Message: fmt.Sprintf("request \"GET %s\" returned \" %s\"", u, resp.Status),
		}
	}

	defer resp.Body.Close()

	log.Debug().Str("url", u).Int("status", resp.StatusCode).Msg("http request executed")
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error().Str("url", u).Int("status", resp.StatusCode).Err(err).Msg("unable to read http response")
		return err
	}

	err = json.Unmarshal(bytes, v)
	if err != nil {
		log.Error().Str("url", u).Int("status", resp.StatusCode).Err(err).Msg("unable to deserialize http response")
		return err
	}

	return nil
}
