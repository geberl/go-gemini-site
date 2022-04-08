package main

import (
	"context"
	"fmt"

	"git.sr.ht/~adnano/go-gemini"
	"github.com/go-kit/kit/log"
)

func aboutHandler(baseUrl string, logger log.Logger) func(context.Context, gemini.ResponseWriter, *gemini.Request) {
	return func(ctx context.Context, w gemini.ResponseWriter, r *gemini.Request) {
		var text gemini.Text

		text = append(text, gemini.LineHeading1("About\n"))
		text = append(text, gemini.LineText("This is a Gemini mirror of Hacker News. Very much work-in-progress.\n"))
		text = append(text, gemini.LineLink{
			URL:  "https://github.com/geberl/go-gemini-hn",
			Name: "Source code is available at GitHub",
		})
		text = append(text, gemini.LineText(""))
		text = append(text, gemini.LineLink{
			URL:  "https://eberl.se",
			Name: "My main website is eberl.se",
		})
		text = append(text, gemini.LineText(""))
		text = append(text, gemini.LineHeading2("Navigation\n"))
		text = append(text, gemini.LineLink{
			URL:  fmt.Sprintf("gemini://%s/", baseUrl),
			Name: "Home",
		})

		w.Write([]byte(text.String()))
	}
}