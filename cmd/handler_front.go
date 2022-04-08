package main

import (
	"context"
	"fmt"

	"git.sr.ht/~adnano/go-gemini"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/lukakerr/hkn"
)

func frontHandler(baseUrl string, logger log.Logger) func(context.Context, gemini.ResponseWriter, *gemini.Request) {
	return func(ctx context.Context, w gemini.ResponseWriter, r *gemini.Request) {
		client := hkn.NewClient()
		ids, err := client.GetTopStories(25)
		if err != nil {
			level.Error(logger).Log("msg", "unable to get top items", "err", err)
		}
		stories, err := client.GetItems(ids)
		if err != nil {
			level.Error(logger).Log("msg", "unable to get items", "err", err)
		}
		sortByTime(stories)

		var text gemini.Text
		text = append(text, gemini.LineHeading1("Hacker News\n"))
		text = append(text, gemini.LineHeading2("Top 25 Stories\n"))

		for _, story := range stories {
			text = append(text, gemini.LineLink{
				URL:  fmt.Sprintf("gemini://%s/item/%d", baseUrl, story.ID),
				Name: story.Title,
			})
			text = append(text, gemini.LineText(fmt.Sprintf("%d score | %d comments\n", story.Score, len(story.Kids))))
		}

		text = append(text, gemini.LineHeading1("Navigation\n"))
		text = append(text, gemini.LineLink{
			URL:  fmt.Sprintf("gemini://%s/about", baseUrl),
			Name: "About",
		})

		w.Write([]byte(text.String()))
	}
}
