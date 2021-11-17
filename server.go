package midimonster

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
)

type ServerHandlerFunc func(server *Server, w http.ResponseWriter, r *http.Request) *HTTPReponse

type Server struct {
	config     *Config
	controller *Controller
	logger     zerolog.Logger
}

type HTTPError struct {
	Error string
}

type HTTPReponse struct {
	Body interface{}
	Code int
}

type HTTPLogs struct {
	Logs   []string
	Newest uint64
}

type HTTPConfigWrite struct {
	Content string
}

type HTTPStatus struct {
	Code int
	Text string
}

func NewServer(config *Config, controller *Controller, logger zerolog.Logger) *Server {
	server := &Server{
		config:     config,
		controller: controller,
		logger:     logger.With().Str("component", "server").Logger(),
	}
	return server
}

func (server *Server) handleResponse(w http.ResponseWriter, r *http.Request, response *HTTPReponse) {
	w.WriteHeader(response.Code)
	encoder := json.NewEncoder(w)
	encoder.Encode(&response.Body)
}

func (server *Server) Start() error {
	router := mux.NewRouter()
	router.HandleFunc("/api/reload", func(w http.ResponseWriter, r *http.Request) {
		response := server.handleConfigReload(w, r)
		server.handleResponse(w, r, response)
	})
	router.HandleFunc("/api/write", func(w http.ResponseWriter, r *http.Request) {
		response := server.handleConfigWrite(w, r)
		server.handleResponse(w, r, response)
	})
	router.HandleFunc("/api/config", func(w http.ResponseWriter, r *http.Request) {
		response := server.handleConfigGet(w, r)
		server.handleResponse(w, r, response)
	})
	router.HandleFunc("/api/status", func(w http.ResponseWriter, r *http.Request) {
		response := server.handleStatus(w, r)
		server.handleResponse(w, r, response)
	})
	router.HandleFunc("/api/logs", func(w http.ResponseWriter, r *http.Request) {
		response := server.handleLogs(w, r)
		server.handleResponse(w, r, response)
	})
	muxer := http.NewServeMux()
	muxer.Handle("/api/", router)
	muxer.Handle("/", http.RedirectHandler("/web/", http.StatusPermanentRedirect))
	if server.config.Development {
		muxer.Handle("/web/", http.StripPrefix("/web", http.FileServer(http.Dir("web"))))
	} else {
		muxer.Handle("/web/", http.FileServer(http.FS(webContent)))
	}
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", server.config.BindAddr, server.config.Port), muxer)
	if err != nil {
		return err
	}

	return nil
}

func (server *Server) handleConfigReload(w http.ResponseWriter, r *http.Request) *HTTPReponse {
	err := server.controller.Midimonster.Restart(r.Context())
	var resp HTTPReponse
	if err != nil {
		resp = HTTPReponse{
			Body: &HTTPError{
				Error: err.Error(),
			},
			Code: http.StatusInternalServerError,
		}
	} else {
		resp = HTTPReponse{
			Body: struct{}{},
			Code: http.StatusOK,
		}
	}
	return &resp
}

func (server *Server) handleConfigWrite(w http.ResponseWriter, r *http.Request) *HTTPReponse {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	content := HTTPConfigWrite{}
	err := decoder.Decode(&content)
	if err != nil {
		return &HTTPReponse{
			Body: HTTPError{
				Error: err.Error(),
			},
			Code: http.StatusBadRequest,
		}
	}
	err = server.controller.Midimonster.ReplaceConfig(r.Context(), content.Content)
	if err != nil {
		return &HTTPReponse{
			Body: HTTPError{
				Error: err.Error(),
			},
			Code: http.StatusInternalServerError,
		}
	}
	return &HTTPReponse{
		Body: struct{}{},
		Code: http.StatusOK,
	}
}

func (server *Server) handleConfigGet(w http.ResponseWriter, r *http.Request) *HTTPReponse {
	defer r.Body.Close()
	content := HTTPConfigWrite{
		Content: server.controller.Midimonster.CurrentConfig,
	}
	return &HTTPReponse{
		Body: &content,
		Code: http.StatusOK,
	}
}

func (server *Server) handleStatus(w http.ResponseWriter, r *http.Request) *HTTPReponse {
	defer r.Body.Close()
	status, err := server.controller.Midimonster.ProcessController.Status(r.Context())
	if err != nil {
		return &HTTPReponse{
			Body: HTTPError{
				Error: err.Error(),
			},
			Code: http.StatusInternalServerError,
		}
	}
	return &HTTPReponse{
		Body: &HTTPStatus{
			Code: int(status),
			Text: status.Text(),
		},
		Code: http.StatusOK,
	}
}

func (server *Server) handleLogs(w http.ResponseWriter, r *http.Request) *HTTPReponse {
	defer r.Body.Close()
	oldestString, ok := r.URL.Query()["oldest"]
	var oldest uint64
	var err error
	if ok {
		if len(oldestString) > 0 {
			oldest, err = strconv.ParseUint(oldestString[0], 10, 64)
			if err != nil {
				return &HTTPReponse{
					Body: &HTTPError{
						Error: err.Error(),
					},
					Code: http.StatusBadRequest,
				}
			}
		}
	}
	logs, newest, err := server.controller.Midimonster.ProcessController.Logs(r.Context(), oldest)
	if err != nil {
		return &HTTPReponse{
			Body: HTTPError{
				Error: err.Error(),
			},
			Code: http.StatusInternalServerError,
		}
	}
	return &HTTPReponse{
		Body: &HTTPLogs{
			Logs:   logs,
			Newest: newest,
		},
		Code: http.StatusOK,
	}
}
