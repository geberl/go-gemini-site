package main

import (
	"context"
	"fmt"
	"strings"

	"git.sr.ht/~adnano/go-gemini"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/lukakerr/hkn"
)

func userHandler(baseUrl string, logger log.Logger) func(context.Context, gemini.ResponseWriter, *gemini.Request) {
	return func(ctx context.Context, w gemini.ResponseWriter, r *gemini.Request) {
		username := strings.TrimPrefix(r.URL.Path, "/user/")

		client := hkn.NewClient()
		user, err := client.GetUser(username)
		if err != nil {
			level.Error(logger).Log("msg", "unable to get user", "err", err)
		}

		submittedItems, err := client.GetItems(user.Submitted)
		if err != nil {
			level.Error(logger).Log("msg", "unable to get items", "err", err)
		}
		sortByTime(submittedItems)

		var text gemini.Text

		text = append(text, gemini.LineHeading1("User Profile\n"))

		text = append(text, gemini.LineHeading2("Links\n"))
		text = append(text, gemini.LineLink{URL: fmt.Sprintf("https://news.ycombinator.com/user?id=%s", user.ID)})

		text = append(text, gemini.LineHeading2("Metadata\n"))
		text = append(text, gemini.LineText(fmt.Sprintf("Name: %s", user.ID)))
		text = append(text, gemini.LineText(fmt.Sprintf("Created: %s", timestamp(int(user.Created)))))
		text = append(text, gemini.LineText(fmt.Sprintf("Karma: %d", user.Karma)))
		text = append(text, gemini.LineText(fmt.Sprintf("Delay: %d", user.Delay)))
		text = append(text, gemini.LineText(fmt.Sprintf("Submitted: %d", len(user.Submitted))))
		text = append(text, gemini.LineText(""))

		text = append(text, gemini.LineHeading2("About\n"))
		if len(user.About) > 0 {
			aboutLines := plaintext(user.About, logger)
			for _, line := range aboutLines {
				text = append(text, gemini.LineQuote(line))
			}
			text = append(text, gemini.LineText(""))
		} else {
			text = append(text, gemini.LineQuote("n/a"))
		}

		text = append(text, gemini.LineHeading2("Submitted\n"))

		text = append(text, gemini.LineHeading3("Stories\n"))
		for _, submittedItem := range submittedItems {
			if submittedItem.Type == "story" {
				text = append(text, gemini.LineLink{
					URL:  fmt.Sprintf("gemini://%s/item/%d", baseUrl, submittedItem.ID),
					Name: submittedItem.Title,
				})
				text = append(text, gemini.LineText(fmt.Sprintf(
					"%d score | %d comments | %s\n",
					submittedItem.Score,
					len(submittedItem.Kids),
					timestamp(int(submittedItem.Time)),
				)))
			}
		}

		text = append(text, gemini.LineHeading3("Comments\n"))
		for _, submittedItem := range submittedItems {
			if submittedItem.Type == "comment" {
				text = append(text, gemini.LineLink{
					URL:  fmt.Sprintf("gemini://%s/item/%d", baseUrl, submittedItem.ID),
					Name: submittedItem.Title,
				})
				text = append(text, gemini.LineText(fmt.Sprintf(
					"%d score | %d comments | %s\n",
					submittedItem.Score,
					len(submittedItem.Kids),
					timestamp(int(submittedItem.Time)),
				)))
			}
		}

		text = append(text, gemini.LineHeading2("Navigation\n"))
		text = append(text, gemini.LineLink{
			URL:  fmt.Sprintf("gemini://%s/", baseUrl),
			Name: "Home",
		})
		w.Write([]byte(text.String()))
	}
}
