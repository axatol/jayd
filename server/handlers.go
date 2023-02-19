package server

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/axatol/jayd/downloader"
	"github.com/axatol/jayd/youtube"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

var (
	yt youtube.Client
)

func handler_GetVideoMetadata(w http.ResponseWriter, r *http.Request) {
	target := r.URL.Query().Get("target")
	if target == "" {
		responseErr(w, err_MissingTarget, http.StatusBadRequest)
		return
	}

	unescaped, err := url.QueryUnescape(target)
	if err != nil {
		log.Error().Str("target", target).Err(err).Msg("failed to unescape target")
		responseErr(w, err_InvalidTarget, http.StatusBadRequest)
		return
	}

	id, err := youtube.ParseURL(unescaped)
	if err != nil {
		log.Error().Str("unescaped", unescaped).Err(err).Msg("failed to parse unescaped url")
		responseErr(w, err_InvalidTarget, http.StatusBadRequest)
		return
	}

	metadata, err := yt.Video(r.Context(), id.VideoID)
	if err != nil {
		log.Error().Str("id", id.VideoID).Err(err).Msg("failed to request youtube metadata")
		responseErr(w, err_FailedRequest, http.StatusBadRequest)
		return
	}

	log.Debug().Str("id", metadata.ID).Msg("metadata request successful")
	responseOk(w, metadata)
}

func handler_QueueVideoDownload(w http.ResponseWriter, r *http.Request) {
	target := r.URL.Query().Get("target")
	if target == "" {
		responseErr(w, err_MissingTarget, http.StatusBadRequest)
		return
	}

	format := (downloader.Format)(r.URL.Query().Get("format"))
	if !format.Valid() {
		responseErr(w, err_InvalidFormat, http.StatusBadRequest)
		return
	}

	unescaped, err := url.QueryUnescape(target)
	if err != nil {
		log.Error().Str("target", target).Err(err).Msg("failed to unescape target")
		responseErr(w, err_InvalidTarget, http.StatusBadRequest)
		return
	}

	metadata, err := youtube.ParseURL(unescaped)
	if err != nil {
		log.Error().Str("unescaped", unescaped).Err(err).Msg("failed to parse unescaped url")
		responseErr(w, err_InvalidTarget, http.StatusBadRequest)
		return
	}

	go downloader.Download(metadata.VideoID, string(format))

	log.Debug().Str("id", metadata.VideoID).Msg("queued download")
	responseOk[any](w, nil)
}

func handler_ListDownloadQueue(w http.ResponseWriter, r *http.Request) {
	responseOk(w, downloader.Jobs.Entries())
}

func handler_StaticContent(root string) http.HandlerFunc {
	dir := http.Dir(root)

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := chi.RouteContext(r.Context())
		prefix := strings.TrimSuffix(ctx.RoutePattern(), "/*")
		fs := http.StripPrefix(prefix, http.FileServer(dir))
		fs.ServeHTTP(w, r)
	}
}
