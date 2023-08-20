package server

import (
	"net/http"
	"time"

	"github.com/axatol/jayd/pkg/downloader"
	"github.com/axatol/jayd/pkg/server/ws"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

var (
	wsWriteDeadline = time.Second * 5
	wsUpgrader      = websocket.Upgrader{CheckOrigin: ws.CheckOrigin}
	wsCloseCodes    = []int{websocket.CloseNormalClosure, websocket.CloseGoingAway}
)

func handler_QueueEvents(w http.ResponseWriter, r *http.Request) {
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		panic(err)
	}

	defer conn.Close()
	conn.EnableWriteCompression(true)

	done := make(chan struct{})
	go func() {
		for {
			_, _, err := conn.ReadMessage()
			if websocket.IsCloseError(err, wsCloseCodes...) {
				break
			}

			if err != nil {
				log.Error().Err(err).Msg("error reading message from client")
				break
			}
		}

		close(done)
	}()

	requestID := r.Context().Value(middleware.RequestIDKey).(string)
	downloader.Cache.Subscribe(requestID, func(event downloader.CacheEvent) { conn.WriteJSON(event) })
	defer downloader.Cache.Unsubscribe(requestID)

	select {
	case <-done:
	case <-r.Context().Done():
	}

	if err := conn.WriteControl(websocket.CloseMessage, []byte{}, time.Now().Add(wsWriteDeadline)); err != nil && err != websocket.ErrCloseSent {
		log.Error().Err(err).Msg("failure writing ws close message")
	}
}
