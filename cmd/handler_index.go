package main

import (
	"context"
	"fmt"

	"git.sr.ht/~adnano/go-gemini"
	"github.com/go-kit/kit/log"
)

func HandlerIndex(baseUrl string, logger log.Logger) func(context.Context, gemini.ResponseWriter, *gemini.Request) {
	return func(ctx context.Context, w gemini.ResponseWriter, r *gemini.Request) {
		var text gemini.Text

		text = append(text, gemini.LineHeading1("Eberl.se\n"))
		text = append(text, gemini.LineText("This is my Gemini site. Very much work-in-progress.\n"))
		text = append(text, gemini.LineLink{
			URL:  "https://github.com/geberl/go-gemini-site",
			Name: "Source code is available at GitHub",
		})
		text = append(text, gemini.LineText(""))
		text = append(text, gemini.LineLink{
			URL:  "https://eberl.se",
			Name: "My website",
		})
		text = append(text, gemini.LineText(""))

		text = append(text, gemini.LineHeading1("Navigation\n"))
		text = append(text, gemini.LineLink{
			URL:  fmt.Sprintf("gemini://%s/gemlog/", baseUrl),
			Name: "Gemlog",
		})
		text = append(text, gemini.LineLink{
			URL:  fmt.Sprintf("gemini://%s/hn/", baseUrl),
			Name: "Hacker News Mirror",
		})

		w.Write([]byte(text.String()))
	}
}
