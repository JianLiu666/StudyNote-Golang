package main

import (
	"encoding/json"
	"net/http"
	"websocketbenchmark/model"
	"websocketbenchmark/util"

	"github.com/sirupsen/logrus"
	"nhooyr.io/websocket"
)

// echoServer is the WebSocket echo server implementation.
// It ensures the client speaks the echo subprotocol and
// only allows one message every 100ms with a 10 message burst.
type echoServer struct {
}

func (s echoServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		Subprotocols: []string{"echo"},
	})
	if err != nil {
		logrus.Error(err)
		return
	}
	defer c.Close(websocket.StatusInternalError, "the sky is falling")

	if c.Subprotocol() != "echo" {
		c.Close(websocket.StatusPolicyViolation, "client must speak the echo subprotocol")
		return
	}

	for {
		mt, message, err := c.Read(r.Context())
		if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
			logrus.Info("client closed connection.")
			return
		}
		if err != nil {
			logrus.Errorf("failed to read: %v", err)
			return
		}

		var recv model.Payload
		err = json.Unmarshal(message, &recv)
		if err != nil {
			logrus.Error("json unmarshal:", err)
			return
		}

		b, err := util.GetEvent(recv.Count)
		if err != nil {
			logrus.Errorf("failed to generate binary data: %v", err)
			return
		}

		c.Write(r.Context(), mt, b)
		if err != nil {
			logrus.Errorf("failed to write: %v", err)
			return
		}
	}
}
