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
	logger             zerolog.Logger
	wsUpgrader         websocket.Upgrader
	connections        map[string][]*WebsocketConnection
	connectionMutex    sync.Mutex
	websocketIdCounter uint64
}

type WebsocketConnection struct {
	id     uint64
	conn   *websocket.Conn
	logger zerolog.Logger
}

func (handler *WebsocketHandler) newWebsocketConnection(conn *websocket.Conn, logger zerolog.Logger) *WebsocketConnection {
	id := handler.websocketIdCounter
	handler.websocketIdCounter++
	return &WebsocketConnection{
		id:     id,
		conn:   conn,
		logger: logger.With().Uint64("id", id).Str("component", "websocketConnection").Logger(),
	}
}

func NewWebsocketHandler(logger zerolog.Logger) *WebsocketHandler {
	connections := make(map[string][]*WebsocketConnection)
	connections[WebsocketCommandGetLogs] = make([]*WebsocketConnection, 0)
	connections[WebsocketCommandGetStatus] = make([]*WebsocketConnection, 0)
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
	wrapperConn := handler.newWebsocketConnection(wsConn, handler.logger)

	defer func() {
		handler.closeConnection(wrapperConn)
		wsConn.Close()
	}()
	for {
		ctx := context.Background()
		req := WebsocketRequest{}
		err := wsConn.ReadJSON(&req)
		if err != nil {
			handler.logger.Err(err).Msgf("cannot read from websocket")
			return err
		}
		handler.logger.Debug().Msgf("received command %q from websocket", req.Command)
		switch req.Command {
		case WebsocketCommandGetLogs:
			server.runWebsocketStepLogs(ctx, wrapperConn)
			handler.addConnection(WebsocketCommandGetLogs, wrapperConn)

		case WebsocketCommandGetStatus:
			server.runWebsocketStepStatus(ctx, wrapperConn)
			handler.addConnection(WebsocketCommandGetStatus, wrapperConn)
		}
	}
}

func (handler *WebsocketHandler) closeConnection(targetConn *WebsocketConnection) {
	handler.connectionMutex.Lock()
	defer handler.connectionMutex.Unlock()
	for cmd, conns := range handler.connections {
		for index, conn := range conns {
			if conn.id == targetConn.id {
				handler.connections[cmd] = append(conns[:index], conns[index+1:]...)
				break
			}
		}
	}
}

func (handler *WebsocketHandler) addConnection(command string, conn *WebsocketConnection) error {
	handler.connectionMutex.Lock()
	defer handler.connectionMutex.Unlock()
	handler.connections[command] = append(handler.connections[command], conn)
	return nil
}

func (conn *WebsocketConnection) sendMessage(ctx context.Context, command string, data interface{}) {
	conn.logger.Debug().Msgf("send message %s to %s", command, conn.conn.RemoteAddr())
	err := conn.conn.WriteJSON(&data)
	if err != nil {
		conn.logger.Err(err).Msgf("error sending command %s", command)
	}
}

func (handler *WebsocketHandler) sendMessage(ctx context.Context, command string, data interface{}, conn *WebsocketConnection) error {
	handler.connectionMutex.Lock()
	defer handler.connectionMutex.Unlock()
	if conn != nil {
		conn.sendMessage(ctx, command, data)
	} else {
		for _, conn := range handler.connections[command] {
			conn.sendMessage(ctx, command, data)
		}
	}
	return nil
}

func (handler *WebsocketHandler) SendStatus(ctx context.Context, status ProcessStatus, conn *WebsocketConnection) {
	resp := &WebsocketStatusResponse{
		Type: "status",
		Status: HTTPStatus{
			Code: int(status),
			Text: status.Text(),
		},
	}
	handler.sendMessage(ctx, WebsocketCommandGetStatus, resp, conn)
}

func (handler *WebsocketHandler) SendLogs(ctx context.Context, logs []string, newest uint64, conn *WebsocketConnection) {
	resp := &WebsocketLogsResponse{
		Type:   "logs",
		Newest: newest,
		Lines:  logs,
	}
	handler.sendMessage(ctx, WebsocketCommandGetLogs, &resp, conn)
}
