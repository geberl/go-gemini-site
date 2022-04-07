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
		text = append(text, gemini.LineHeading1("HN - Top 25 Stories"))

		for _, story := range stories {
			text = append(text, gemini.LineLink{
				URL:  fmt.Sprintf("gemini://%s/item/%d", baseUrl, story.ID),
				Name: story.Title})
			text = append(text, gemini.LineText(fmt.Sprintf("%d score | %d comments\n", story.Score, len(story.Kids))))
		}

		w.Write([]byte(text.String()))
	}
}
