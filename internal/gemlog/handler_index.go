package gemlog

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"mime"
	"net/url"
	"path"
	"sort"
	"strings"

	"git.sr.ht/~adnano/go-gemini"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/log/level"
)

func FileServer(baseUrl string, logger log.Logger, pattern string, fsys fs.FS) gemini.Handler {
	return fileServer{baseUrl, logger, pattern, fsys}
}

type fileServer struct {
	baseUrl string
	logger  log.Logger
	pattern string
	fs.FS
}

func (fsys fileServer) ServeGemini(ctx context.Context, w gemini.ResponseWriter, r *gemini.Request) {
	const indexPage = "/index.gmi"

	url := path.Clean(r.URL.Path)

	if strings.HasSuffix(url, indexPage) {
		w.WriteHeader(gemini.StatusPermanentRedirect, strings.TrimSuffix(url, "index.gmi"))
		return
	}

	name := url
	if name == fmt.Sprintf("/%s", fsys.pattern) {
		name = "."
	} else {
		name = strings.TrimPrefix(name, fmt.Sprintf("/%s/", fsys.pattern))
	}

	level.Info(fsys.logger).Log("msg", "serve gemlog", "url", url, "name", name)

	f, err := fsys.Open(name)
	if err != nil {
		level.Error(fsys.logger).Log("msg", "unable to open file", "url", url, "err", err)
		w.WriteHeader(toGeminiError(err))
		return
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		level.Error(fsys.logger).Log("msg", "unable to get file info", "url", url, "err", err)
		w.WriteHeader(toGeminiError(err))
		return
	}

	if len(r.URL.Path) != 0 {
		if stat.IsDir() {
			target := url
			if target != "/" {
				target += "/"
			}
			if len(r.URL.Path) != len(target) || r.URL.Path != target {
				w.WriteHeader(gemini.StatusPermanentRedirect, target)
				return
			}
		} else if r.URL.Path[len(r.URL.Path)-1] == '/' {
			w.WriteHeader(gemini.StatusPermanentRedirect, url)
			return
		}
	}

	if stat.IsDir() {
		dirList(w, f, fsys.baseUrl, fsys.logger)
		return
	}

	ext := path.Ext(name)
	mimetype := mime.TypeByExtension(ext)
	w.SetMediaType(mimetype)
	io.Copy(w, f)
}

func dirList(w gemini.ResponseWriter, f fs.File, baseUrl string, logger log.Logger) {
	var entries []fs.DirEntry
	var err error
	d, ok := f.(fs.ReadDirFile)
	if ok {
		entries, err = d.ReadDir(-1)
	}
	if !ok || err != nil {
		level.Error(logger).Log("msg", "unable to read directory", "err", err)
		w.WriteHeader(toGeminiError(err))
		return
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() > entries[j].Name()
	})

	var text gemini.Text
	text = append(text, gemini.LineHeading1("Günther Eberl's Gemlog\n"))
	text = append(text, gemini.LineHeading2("My Posts\n"))

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filename := entry.Name()

		components := strings.Split(filename, "_")
		if len(components) != 2 {
			continue
		}

		linkdate := components[0]

		linkname := components[1]
		linkname = strings.Replace(linkname, "-", " ", -1)
		linkname = strings.TrimSuffix(linkname, ".gmi")
		linkname = strings.Title(linkname)

		text = append(text, gemini.LineLink{
			Name: fmt.Sprintf("%s - %s", linkdate, linkname),
			URL:  (&url.URL{Path: filename}).EscapedPath(),
		})
	}
	text = append(text, gemini.LineText(""))

	text = append(text, gemini.LineHeading2("Other Gemlogs I Enjoy\n"))
	text = append(text, gemini.LineText("tba"))
	text = append(text, gemini.LineText(""))

	text = append(text, gemini.LineHeading1("Navigation\n"))
	text = append(text, gemini.LineLink{
		URL:  fmt.Sprintf("gemini://%s", baseUrl),
		Name: "Home",
	})

	w.Write([]byte(text.String()))
}

func toGeminiError(err error) (status gemini.Status, meta string) {
	if errors.Is(err, fs.ErrNotExist) {
		return gemini.StatusNotFound, "Not found"
	}
	if errors.Is(err, fs.ErrPermission) {
		return gemini.StatusNotFound, "Forbidden"
	}
	return gemini.StatusTemporaryFailure, "Internal server error"
}
