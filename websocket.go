package midimonster

import (
	"context"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
)

const (
	WebsocketCommandGetLogs   = "getLogs"
	WebsocketCommandGetStatus = "getStatus"
)

type WebsocketRequest struct {
	Command string `json:"command"`
}

type WebsocketLogsResponse struct {
	Type   string   `json:"type"`
	Newest uint64   `json:"newest"`
	Lines  []string `json:"lines"`
}

type WebsocketStatusResponse struct {
	Type   string     `json:"type"`
	Status HTTPStatus `json:"status"`
}

type WebsocketHandler struct {
	logger          zerolog.Logger
	wsUpgrader      websocket.Upgrader
	connections     map[string][]*websocket.Conn
	connectionMutex sync.Mutex
}

func NewWebsocketHandler(logger zerolog.Logger) *WebsocketHandler {
	connections := make(map[string][]*websocket.Conn)
	connections[WebsocketCommandGetLogs] = make([]*websocket.Conn, 0)
	connections[WebsocketCommandGetStatus] = make([]*websocket.Conn, 0)
	return &WebsocketHandler{
		logger:      logger.With().Str("component", "websocket").Logger(),
		wsUpgrader:  websocket.Upgrader{},
		connections: connections,
	}
}

func (handler *WebsocketHandler) Connect(server *Server, w http.ResponseWriter, r *http.Request) error {
	handler.logger.Debug().Msgf("got web socket connect")
	wsConn, err := handler.wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		handler.logger.Err(err).Msgf("cannot upgrade connection to websocket")
		return err
	}

	defer func() {
		handler.closeConnection(wsConn)
		wsConn.Close()
	}()
	for {
		req := WebsocketRequest{}
		err := wsConn.ReadJSON(&req)
		if err != nil {
			handler.logger.Err(err).Msgf("cannot read from websocket")
			return err
		}
		handler.logger.Debug().Msgf("received command %q from websocket", req.Command)
		switch req.Command {
		case WebsocketCommandGetLogs:
			server.runWebsocketStepLogs(context.TODO(), wsConn)
			handler.addConnection(WebsocketCommandGetLogs, wsConn)

		case WebsocketCommandGetStatus:
			handler.addConnection(WebsocketCommandGetStatus, wsConn)
		}
	}
}

func (handler *WebsocketHandler) closeConnection(targetConn *websocket.Conn) {
	handler.connectionMutex.Lock()
	defer handler.connectionMutex.Unlock()
	for cmd, conns := range handler.connections {
		for index, conn := range conns {
			if conn == targetConn {
				handler.connections[cmd] = append(conns[:index], conns[index+1:]...)
				break
			}
		}
	}
}

func (handler *WebsocketHandler) addConnection(command string, conn *websocket.Conn) error {
	handler.connectionMutex.Lock()
	defer handler.connectionMutex.Unlock()
	handler.connections[command] = append(handler.connections[command], conn)
	return nil
}

func (handler *WebsocketHandler) sendMessageToConnection(ctx context.Context, conn *websocket.Conn, command string, data interface{}) {
	handler.logger.Debug().Msgf("send message %s to %s", command, conn.RemoteAddr())
	err := conn.WriteJSON(&data)
	if err != nil {
		handler.logger.Err(err).Msgf("error sending command %s", command)
	}
}

func (handler *WebsocketHandler) sendMessage(ctx context.Context, command string, data interface{}) error {
	handler.connectionMutex.Lock()
	defer handler.connectionMutex.Unlock()
	for _, conn := range handler.connections[command] {
		handler.sendMessageToConnection(ctx, conn, command, data)
	}
	return nil
}

func (handler *WebsocketHandler) SendStatus(ctx context.Context, status ProcessStatus) {
	resp := &WebsocketStatusResponse{
		Type: "status",
		Status: HTTPStatus{
			Code: int(status),
			Text: status.Text(),
		},
	}
	handler.sendMessage(ctx, WebsocketCommandGetStatus, resp)
}

func (handler *WebsocketHandler) SendLogs(ctx context.Context, logs []string, newest uint64, conn *websocket.Conn) {
	resp := &WebsocketLogsResponse{
		Type:   "logs",
		Newest: newest,
		Lines:  logs,
	}
	if conn == nil {
		handler.sendMessage(ctx, WebsocketCommandGetLogs, &resp)
	} else {
		handler.sendMessageToConnection(ctx, conn, WebsocketCommandGetLogs, &resp)
	}
}
