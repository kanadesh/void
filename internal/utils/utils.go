package utils

import (
	"net/url"
	"path/filepath"
)

var AudioMap map[string]string = map[string]string{
	"audio/mpeg": "mp3",
	"audio/wav":  "wav",
}

func FixUrl(raw string, base string) string {
	if u, err := url.ParseRequestURI(raw); err != nil || u.Scheme == "" || u.Host == "" {
		link, _ := url.Parse(base)

		if raw[0:1] == "/" {
			link.Path = raw
		} else {
			link.Path = filepath.Join(link.Path, raw)
		}
		return link.String()
	}
	return raw
}
