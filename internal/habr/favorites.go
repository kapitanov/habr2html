package habr

type ArticleContext interface {
	ID() string
	Get() (*Article, error)
}

func ListFavorites(
	client Client,
	userId string,
	process func(ctx ArticleContext) error) error {

	page := 0

	for {
		page++
		list, err := client.ListFavorites(userId, page)
		if err != nil {
			return err
		}

		if len(list.ArticleIds) == 0 {
			break
		}

		for _, ref := range list.ArticleRefs {
			err = process(articleContext{client: client, ref: ref})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

type articleContext struct {
	client Client
	ref    *ArticleRef
}

func (ctx articleContext) ID() string {
	return ctx.ref.ID
}

func (ctx articleContext) Get() (*Article, error) {
	article, err := ctx.client.GetArticle(ctx.ref.ID)
	if err != nil {
		return nil, err
	}

	return article, nil
}
