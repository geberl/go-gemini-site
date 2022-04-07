package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"git.sr.ht/~adnano/go-gemini"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/lukakerr/hkn"
)

func itemHandler(baseUrl string, logger log.Logger) func(context.Context, gemini.ResponseWriter, *gemini.Request) {
	return func(ctx context.Context, w gemini.ResponseWriter, r *gemini.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/item/")
		cleanId, err := strconv.Atoi(id)
		if err != nil {
			level.Error(logger).Log("msg", "unable to convert item id to integer", "err", err)
		}

		client := hkn.NewClient()
		item, err := client.GetItem(cleanId)
		if err != nil {
			level.Error(logger).Log("msg", "unable to get item", "err", err)
		}

		firstLevelComments, err := client.GetItems(item.Kids)
		if err != nil {
			level.Error(logger).Log("msg", "unable to get items", "err", err)
		}
		sortByTime(firstLevelComments)

		var text gemini.Text

		if item.Type == "story" {
			text = append(text, gemini.LineHeading1(item.Title))
		} else {
			text = append(text, gemini.LineHeading1(fmt.Sprintf("%s %d", strings.Title(item.Type), item.ID)))
		}

		text = append(text, gemini.LineHeading2("Links"))
		if item.URL != "" {
			text = append(text, gemini.LineLink{URL: item.URL})
		}
		text = append(text, gemini.LineLink{URL: fmt.Sprintf("https://news.ycombinator.com/item?id=%d", item.ID)})

		text = append(text, gemini.LineHeading2("Metadata"))
		text = append(text, gemini.LineText(fmt.Sprintf("By: %s", item.By)))
		text = append(text, gemini.LineText(fmt.Sprintf("Comments: %d", item.Descendants)))
		text = append(text, gemini.LineText(fmt.Sprintf("Score: %d", item.Score)))
		text = append(text, gemini.LineText(fmt.Sprintf("Time: %s", timestamp(int(item.Time)))))

		if len(item.Text) > 0 {
			text = append(text, gemini.LineHeading2("Text"))
			storyLines := plaintext(item.Text, logger)
			for _, line := range storyLines {
				text = append(text, gemini.LineQuote(line))
			}
		}

		text = append(text, gemini.LineHeading2("Comments"))
		for _, firstLevelComment := range firstLevelComments {
			commentLines := plaintext(firstLevelComment.Text, logger)
			if commentLines[0] != "n/a" {
				text = append(text, gemini.LineHeading3(fmt.Sprintf("%s by %s", timestamp(int(firstLevelComment.Time)), firstLevelComment.By)))

				for _, line := range commentLines {
					text = append(text, gemini.LineQuote(line))
				}

				if len(firstLevelComment.Kids) > 0 {
					text = append(text, gemini.LineLink{
						URL:  fmt.Sprintf("gemini://%s/item/%d", baseUrl, firstLevelComment.ID),
						Name: fmt.Sprintf("%d Responses", len(firstLevelComment.Kids)),
					})
				}
			}
		}

		text = append(text, gemini.LineHeading2("Navigation"))
		if item.Parent > 0 {
			text = append(text, gemini.LineLink{
				URL:  fmt.Sprintf("gemini://%s/item/%d", baseUrl, item.Parent),
				Name: "Parent",
			})
		}
		text = append(text, gemini.LineLink{
			URL:  fmt.Sprintf("gemini://%s/", baseUrl),
			Name: "Home",
		})

		w.Write([]byte(text.String()))
	}
}
