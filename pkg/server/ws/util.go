package ws

import (
	"net/http"
	"net/url"
	"strings"
	"unicode/utf8"

	"github.com/axatol/jayd/pkg/config"
	"github.com/rs/zerolog/log"
)

// source: https://github.com/gorilla/websocket/blob/master/util.go#L176-L198
// equalASCIIFold returns true if s is equal to t with ASCII case folding as
// defined in RFC 4790.
func equalASCIIFold(s, t string) bool {
	for s != "" && t != "" {
		sr, size := utf8.DecodeRuneInString(s)
		s = s[size:]
		tr, size := utf8.DecodeRuneInString(t)
		t = t[size:]
		if sr == tr {
			continue
		}
		if 'A' <= sr && sr <= 'Z' {
			sr = sr + 'a' - 'A'
		}
		if 'A' <= tr && tr <= 'Z' {
			tr = tr + 'a' - 'A'
		}
		if sr != tr {
			return false
		}
	}
	return s == t
}

func CheckOrigin(r *http.Request) bool {
	origin := r.Header["Origin"]
	if len(origin) == 0 {
		return true
	}

	u, err := url.Parse(origin[0])
	if err != nil {
		return false
	}

	allowedList := strings.Split(config.ServerCORSList, ",")
	if len(allowedList) == 0 {
		return true
	}

	for _, allowed := range allowedList {
		allowedURL, err := url.Parse(allowed)
		if err != nil {
			log.Error().Err(err).Str("allowed", allowed).Msg("failed to parse url from allow list")
			return false
		}

		if equalASCIIFold(u.Host, allowedURL.Host) {
			return true
		}
	}

	return false
}
