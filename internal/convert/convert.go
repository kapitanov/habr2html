package convert

import (
	"sort"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/rs/zerolog/log"

	"github.com/kapitanov/habr2html/internal/convert/html"
	"github.com/kapitanov/habr2html/internal/habr"
	"github.com/kapitanov/habr2html/internal/store"
)

func HTML(article *habr.Article, ctx store.StoreArticleContext) error {
	log.Debug().Str("id", article.ID).Msg("converting article")

	// Title
	title := article.TitleHTML

	// Category
	var category string
	if len(article.Hubs) > 0 {
		sort.Slice(article.Hubs, func(i, j int) bool {
			return hubLessComparator(article.Hubs[i], article.Hubs[j])
		})

		category = article.Hubs[0].TitleHTML
	} else if len(article.Tags) > 0 {
		sort.Slice(article.Tags, func(i, j int) bool {
			return strings.Compare(article.Tags[i].TitleHTML, article.Tags[j].TitleHTML) < 0
		})

		category = article.Tags[0].TitleHTML
	} else {
		category = "UNDEFINED"
	}

	// Description
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(article.LeadData.TextHTML))
	if err != nil {
		return err
	}

	description:= doc.Text()

	err = html.Do(ctx.For(category), article, title, description)
	if err != nil {
		return err
	}

	return nil
}

func hubLessComparator(x, y *habr.Hub) bool {
	if hubLessComparator_Type(x, y) {
		return true
	}

	if hubLessComparator_Alias(x, y) {
		return true
	}

	if hubLessComparator_ID(x, y) {
		return true
	}

	return false
}

func hubLessComparator_Type(x, y *habr.Hub) bool {
	if x.Type == y.Type {
		return false
	}

	if x.Type == habr.CollectiveHub {
		return true
	}

	return false
}

func hubLessComparator_Alias(x, y *habr.Hub) bool {
	if x.Alias == y.Alias {
		return false
	}

	if strings.Compare(x.Alias, y.Alias) < 0 {
		return true
	}

	return false
}

func hubLessComparator_ID(x, y *habr.Hub) bool {
	if x.ID == y.ID {
		return false
	}

	if strings.Compare(x.ID, y.ID) < 0 {
		return true
	}

	return false
}
