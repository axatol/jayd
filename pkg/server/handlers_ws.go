package server

import (
	"net/http"
	"time"

	ds "github.com/axatol/go-utils/datastructures"
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

	requestID := r.Context().Value(middleware.RequestIDKey)
	events := make(chan ds.AsyncMapEvent[downloader.InfoJSON], 1)
	downloader.CacheEvents.Subscribe((requestID).(string), events)

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

	for loop := true; loop; {
		select {
		case <-done:
			loop = false
		case event := <-events:
			conn.WriteJSON(event)
		}
	}

	if err := conn.WriteControl(websocket.CloseMessage, []byte{}, time.Now().Add(wsWriteDeadline)); err != nil && err != websocket.ErrCloseSent {
		log.Error().Err(err).Msg("failure writing ws close message")
	}
}
