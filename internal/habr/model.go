package habr

import "time"

type ArticleListResponse struct {
	PagesCount  int                    `json:"pagesCount"`
	ArticleIds  []string               `json:"articleIds"`
	ArticleRefs map[string]*ArticleRef `json:"articleRefs"`
}

type ArticleRef struct {
	ID    string `json:"id"`
	Title string `json:"titleHtml"`
}

type Article struct {
	ID            string    `json:"id"`
	TimePublished time.Time `json:"timePublished"`
	TitleHTML     string    `json:"titleHtml"`
	Author        *Author   `json:"author"`
	LeadData      *LeadData `json:"leadData"`
	Hubs          []*Hub    `json:"hubs"`
	TextHTML      string    `json:"textHtml"`
	Tags          []*Tag    `json:"tags"`
	Metadata      *Metadata `json:"metadata"`
}

type Author struct {
	ID         string `json:"id"`
	Alias      string `json:"alias"`
	FullName   string `json:"fullname"`
	Speciality string `json:"speciality"`
}

type HubType string

const (
	CollectiveHub HubType = "collective"
)

type Hub struct {
	ID         string  `json:"id"`
	Alias      string  `json:"alias"`
	Type       HubType `json:"type"`
	TitleHTML  string  `json:"titleHtml"`
	IsProfiled bool    `json:"isProfiled"`
}

type Tag struct {
	TitleHTML string `json:"titleHtml"`
}

type Metadata struct {
	ShareImageUrl string `json:"shareImageUrl"`
}

type LeadData struct {
	TextHTML string `json:"textHtml"`
}
