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
			w.WriteHeader(gemini.StatusBadRequest, "Bad request")
			return
		}

		client := hkn.NewClient()
		item, err := client.GetItem(cleanId)
		if err != nil {
			level.Error(logger).Log("msg", "unable to get item", "err", err)
			w.WriteHeader(gemini.StatusNotFound, "Not found")
			return
		}

		firstLevelComments, err := client.GetItems(item.Kids)
		if err != nil {
			level.Error(logger).Log("msg", "unable to get items", "err", err)
			w.WriteHeader(gemini.StatusNotFound, "Not found")
			return
		}
		sortByTime(firstLevelComments)

		var text gemini.Text

		if item.Type == "story" {
			text = append(text, gemini.LineHeading1(fmt.Sprintf("%s\n", item.Title)))
		} else {
			text = append(text, gemini.LineHeading1(fmt.Sprintf("%s\n", strings.Title(item.Type))))
		}

		text = append(text, gemini.LineHeading2("Links\n"))
		if item.URL != "" {
			text = append(text, gemini.LineLink{URL: item.URL})
			text = append(text, gemini.LineText(""))
		}
		text = append(text, gemini.LineLink{URL: fmt.Sprintf("https://news.ycombinator.com/item?id=%d", item.ID)})
		text = append(text, gemini.LineText(""))

		text = append(text, gemini.LineHeading2("Metadata\n"))
		text = append(text, gemini.LineLink{
			URL:  fmt.Sprintf("gemini://%s/user/%s", baseUrl, item.By),
			Name: fmt.Sprintf("By: %s", item.By),
		})
		text = append(text, gemini.LineText(fmt.Sprintf("Id: %d", item.ID)))
		text = append(text, gemini.LineText(fmt.Sprintf("Comments: %d", item.Descendants)))
		text = append(text, gemini.LineText(fmt.Sprintf("Score: %d", item.Score)))
		text = append(text, gemini.LineText(fmt.Sprintf("Type: %s", strings.Title(item.Type))))
		text = append(text, gemini.LineText(fmt.Sprintf("Created: %s", timestamp(int(item.Time)))))
		text = append(text, gemini.LineText(""))

		if len(item.Text) > 0 {
			text = append(text, gemini.LineHeading2("Text\n"))
			storyLines := plaintext(item.Text, logger)
			for _, line := range storyLines {
				text = append(text, gemini.LineQuote(line))
			}
			text = append(text, gemini.LineText(""))
		}

		text = append(text, gemini.LineHeading2("Comments\n"))
		for _, firstLevelComment := range firstLevelComments {
			commentLines := plaintext(firstLevelComment.Text, logger)
			if commentLines[0] != "n/a" {
				text = append(text, gemini.LineHeading3(fmt.Sprintf("%s\n", timestamp(int(firstLevelComment.Time)))))

				text = append(text, gemini.LineLink{
					URL:  fmt.Sprintf("gemini://%s/user/%s", baseUrl, firstLevelComment.By),
					Name: fmt.Sprintf("By: %s", item.By),
				})

				if len(firstLevelComment.Kids) > 0 {
					text = append(text, gemini.LineLink{
						URL:  fmt.Sprintf("gemini://%s/item/%d", baseUrl, firstLevelComment.ID),
						Name: fmt.Sprintf("Responses: %d", len(firstLevelComment.Kids)),
					})
				}

				text = append(text, gemini.LineText(""))
				for _, line := range commentLines {
					text = append(text, gemini.LineQuote(line))
				}
			}
			text = append(text, gemini.LineText(""))
		}

		text = append(text, gemini.LineHeading1("Navigation\n"))
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
		text = append(text, gemini.LineLink{
			URL:  fmt.Sprintf("gemini://%s/about", baseUrl),
			Name: "About",
		})

		w.Write([]byte(text.String()))
	}
}
