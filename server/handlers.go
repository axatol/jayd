package server

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	ds "github.com/axatol/go-utils/datastructures"
	"github.com/axatol/jayd/downloader"
	"github.com/axatol/jayd/youtube"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/websocket"
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
	target := r.URL.Query().Get("target")
	if target == "" {
		responseErr(w, err_MissingTarget, http.StatusBadRequest)
		return
	}

	format := r.URL.Query().Get("format")
	if format == "" {
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

	info, err := downloader.GetInfoJSON(r.Context(), metadata.VideoID)
	if err != nil {
		log.Error().Err(err).Str("target", info.VideoID).Msg("failed to fetch metadata")
		responseErr(w, err_FetchMetadata, http.StatusBadRequest)
		return
	}

	exists := false
	for _, cached := range info.Formats {
		if cached.FormatID == format {
			exists = true
			break
		}
	}

	if format != downloader.FormatDefaultAudio && format != downloader.FormatDefaultVideo && !exists {
		log.Error().Err(err).Str("target", info.VideoID).Str("format", format).Msg("invalid format")
		responseErr(w, err_InvalidFormat, http.StatusBadRequest)
		return
	}

	go func() {
		if err := downloader.Download(*info); err != nil {
			log.Error().Err(err).Str("target", info.VideoID).Msg("failed to download")
		}
	}()

	log.Debug().Str("target", info.VideoID).Msg("queued download")
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

	items := downloader.Cache.Get(metadata.VideoID + r.URL.Query().Get("format"))
	if items == nil {
		responseErr(w, err_NotFound, http.StatusNotFound)
		return
	}

	responseOk(w, items)
}

func handler_DeleteDownloadQueueItem(w http.ResponseWriter, r *http.Request) {
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

	downloader.Cache.Remove(metadata.VideoID + format)
	responseOk[any](w, nil)
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

func handler_QueueEvents(w http.ResponseWriter, r *http.Request) {
	writeDeadline := time.Second * 5
	upgrader := websocket.Upgrader{}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	defer conn.Close()
	conn.EnableWriteCompression(true)

	requestID := r.Context().Value(middleware.RequestIDKey)
	// events := make(chan ds.AsyncMapEvent[downloader.CacheItem], 1)
	events := make(chan ds.AsyncMapEvent[downloader.InfoJSON], 1)
	downloader.CacheEvents.Subscribe((requestID).(string), events)

	done := make(chan struct{})
	go func() {
		for {
			_, _, err := conn.ReadMessage()
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				break
			}

			if err != nil {
				log.Error().Err(err).Msg("error reading message from client")
				break
			}
		}

		close(done)
	}()

	for loop := true; loop; {
		select {
		case <-done:
			loop = false
		case event := <-events:
			conn.WriteJSON(event)
		}
	}

	if err := conn.WriteControl(websocket.CloseMessage, []byte{}, time.Now().Add(writeDeadline)); err != nil && err != websocket.ErrCloseSent {
		log.Error().Err(err).Msg("failure writing ws close message")
	}
}
