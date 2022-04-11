package hn

import (
	"sort"
	"strings"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/lukakerr/hkn"
	"jaytaylor.com/html2text"
)

func plaintext(html string, logger log.Logger) []string {
	var lines []string

	if html == "" {
		return append(lines, "n/a")
	}

	paragraphs1 := strings.Split(html, "\n")
	for _, para1 := range paragraphs1 {
		text, err := html2text.FromString(para1, html2text.Options{TextOnly: true})
		if err != nil {
			level.Error(logger).Log("msg", "unable to convert html into text", "err", err)
		}
		paragraphs2 := strings.Split(text, "\n")
		for _, para2 := range paragraphs2 {
			lines = append(lines, para2)
		}

	}
	return lines
}

func timestamp(unix int) string {
	return time.Unix(int64(unix), 0).Format("2006-01-02 15:04:05")
}

func sortByTime(items []hkn.Item) {
	sort.Slice(items, func(i, j int) bool {
		return items[i].Time < items[j].Time
	})
}
