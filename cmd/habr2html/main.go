package main

import (
	"errors"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/kapitanov/habr2html/internal/convert"
	"github.com/kapitanov/habr2html/internal/habr"
	pkgstore "github.com/kapitanov/habr2html/internal/store"
)

var rootCommand = &cobra.Command{
	Use:          "habr2html",
	SilenceUsage: true,
}

func init() {
	userID := rootCommand.Flags().StringP("user-id", "u", "", "Habr user ID")
	outputDir := rootCommand.Flags().StringP("output", "o", "./out", "Path to output directory")

	rootCommand.PreRun = func(cmd *cobra.Command, args []string) {
		log.Logger = log.Output(zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
			w.Out = os.Stderr
			w.TimeFormat = "15:04:05"
		}))
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	rootCommand.RunE = func(cmd *cobra.Command, args []string) error {
		if *userID == "" {
			return errors.New("missing user ID")
		}

		habrClient := habr.New()

		store, err := pkgstore.New(*outputDir)
		if err != nil {
			return err
		}

		err = habr.ListFavorites(habrClient, *userID, func(c habr.ArticleContext) error {
			err := store.StoreArticle(c.ID(), func(ctx pkgstore.StoreArticleContext) error {
				article, err := c.Get()
				if err != nil {
					return err
				}

				if article == nil {
					return nil
				}

				err = convert.HTML(article, ctx)
				if err != nil {
					return err
				}

				log.Info().Str("id", c.ID()).Msg("article has been exported")
				return nil
			})

			if err != nil {
				return err
			}

			return nil
		})
		if err != nil {
			return err
		}

		return nil
	}
}

func main() {
	err := rootCommand.Execute()
	if err != nil {
		os.Exit(-1)
	}
}
