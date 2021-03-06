package hn

import (
	"context"
	"fmt"

	"git.sr.ht/~adnano/go-gemini"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/lukakerr/hkn"
)

func HandlerIndex(baseUrl string, logger log.Logger) func(context.Context, gemini.ResponseWriter, *gemini.Request) {
	return func(ctx context.Context, w gemini.ResponseWriter, r *gemini.Request) {
		client := hkn.NewClient()
		ids, err := client.GetTopStories(25)
		if err != nil {
			level.Error(logger).Log("msg", "unable to get top items", "err", err)
			w.WriteHeader(gemini.StatusNotFound, "Not found")
			return
		}
		stories, err := client.GetItems(ids)
		if err != nil {
			level.Error(logger).Log("msg", "unable to get items", "err", err)
			w.WriteHeader(gemini.StatusNotFound, "Not found")
			return
		}
		sortByTime(stories)

		var text gemini.Text
		text = append(text, gemini.LineHeading1("Hacker News Mirror\n"))
		text = append(text, gemini.LineHeading2("Top 25 Stories\n"))

		for _, story := range stories {
			linkURL := fmt.Sprintf("gemini://%s/hn/item/%d", baseUrl, story.ID)
			linkName := fmt.Sprintf("%s [%d points | %d comments]", story.Title, story.Score, len(story.Kids))
			text = append(text, gemini.LineLink{
				URL:  linkURL,
				Name: linkName,
			})
			text = append(text, gemini.LineText(""))
		}

		text = append(text, gemini.LineHeading1("Navigation\n"))
		text = append(text, gemini.LineLink{
			URL:  fmt.Sprintf("gemini://%s", baseUrl),
			Name: "Home",
		})

		w.Write([]byte(text.String()))
	}
}
