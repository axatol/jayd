package server

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	fp "github.com/axatol/go-utils/functional"
	"github.com/axatol/jayd/pkg/downloader"
	"github.com/axatol/jayd/pkg/downloader/miniodriver"
	"github.com/axatol/jayd/pkg/youtube"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
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

	info, err := downloader.GetInfoJSON(r.Context(), id.VideoID)
	if err != nil {
		log.Error().Str("id", id.VideoID).Err(err).Msg("failed to request youtube metadata")
		responseErr(w, err_FailedRequest, http.StatusBadRequest)
		return
	}

	log.Debug().Str("id", info.VideoID).Msg("metadata request successful")
	responseOk(w, info)
}

func handler_QueueVideoDownload(w http.ResponseWriter, r *http.Request) {
	log := log.With().Logger()

	target := r.URL.Query().Get("target")
	if target == "" {
		responseErr(w, err_MissingTarget, http.StatusBadRequest)
		return
	}

	format := r.URL.Query().Get("format")
	log = log.With().Str("format", format).Logger()
	if format == "" {
		responseErr(w, err_InvalidFormat, http.StatusBadRequest)
		return
	}

	overwrite := r.URL.Query().Get("overwrite") == "true"
	log = log.With().Bool("overwrite", overwrite).Logger()

	unescaped, err := url.QueryUnescape(target)
	if err != nil {
		log.Error().Str("target", target).Err(err).Msg("failed to unescape target")
		responseErr(w, err_InvalidTarget, http.StatusBadRequest)
		return
	}

	metadata, err := youtube.ParseURL(unescaped)
	log = log.With().Str("video_id", metadata.VideoID).Logger()
	if err != nil {
		log.Error().Str("unescaped", unescaped).Err(err).Msg("failed to parse unescaped url")
		responseErr(w, err_InvalidTarget, http.StatusBadRequest)
		return
	}

	info, err := downloader.GetInfoJSON(r.Context(), metadata.VideoID)
	if err != nil {
		log.Error().Err(err).Msg("failed to fetch metadata")
		responseErr(w, err_FetchMetadata, http.StatusBadRequest)
		return
	}

	matching := fp.Filter(info.Formats, func(e downloader.Format, i int) bool { return e.FormatID == format })
	if len(matching) < 1 {
		log.Error().Err(err).Msg("invalid format")
		responseErr(w, err_InvalidFormat, http.StatusBadRequest)
		return
	}

	go func() {
		if err := downloader.Download(context.Background(), *info, format, overwrite); err != nil {
			log.Error().Err(err).Msg("failed to download")
		}
	}()

	log.Debug().Msg("queued download")
	responseOk[any](w, nil)
}

func handler_ListDownloadQueue(w http.ResponseWriter, r *http.Request) {
	target := r.URL.Query().Get("target")
	if target == "" {
		responseOk(w, downloader.Cache.Entries())
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

	format := r.URL.Query().Get("format")
	if format == "" {
		responseErr(w, err_MissingFormat, http.StatusBadRequest)
		return
	}

	id := downloader.CacheItemID(metadata.VideoID, format)
	items := downloader.Cache.Get(id)
	if items == nil {
		responseErr(w, err_NotFound, http.StatusNotFound)
		return
	}

	responseOk(w, downloader.Cache)
}

func handler_DeleteQueueItem(w http.ResponseWriter, r *http.Request) {
	target := r.URL.Query().Get("target")
	if target == "" {
		responseErr(w, err_MissingTarget, http.StatusBadRequest)
		return
	}

	format := r.URL.Query().Get("format")
	if format == "" {
		responseErr(w, err_MissingFormat, http.StatusBadRequest)
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

	id := downloader.CacheItemID(metadata.VideoID, format)
	downloader.Cache.Remove(id)
	responseOk[any](w, nil)
}

func handler_PresignedContent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	key := chi.URLParam(r, "key")
	if key == "" {
		responseErr(w, err_MissingKey, http.StatusNotFound)
		return
	}

	client, err := miniodriver.AssertClient(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to assert storage driver")
		responseErr(w, err_GenericError, http.StatusInternalServerError)
		return
	}

	url, err := client.GetPresignedURL(ctx, key)
	if err != nil {
		log.Error().Err(err).Str("key", key).Msg("failed to create presigned url")
		responseErr(w, err_GenericError, http.StatusInternalServerError)
		return
	}

	// TODO setup cors on backend
	w.Header().Add("location", url.String())
	w.WriteHeader(http.StatusSeeOther)
	w.Write([]byte(url.String()))
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
