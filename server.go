package midimonster

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	config     *Config
	controller *Controller
}

type HTTPError struct {
	Error string
}

type HTTPReponse struct {
	Body interface{}
	Code int
}

type HTTPConfigWrite struct {
	Content string
}

func NewServer(config *Config, controller *Controller) *Server {
	return &Server{
		config:     config,
		controller: controller,
	}
}

func (server *Server) Start() error {
	router := mux.NewRouter()
	router.HandleFunc("/api/reload", func(w http.ResponseWriter, r *http.Request) {
		response := server.handleConfigReload(w, r)
		w.WriteHeader(response.Code)
		encoder := json.NewEncoder(w)
		encoder.Encode(&response.Body)
	})
	router.HandleFunc("/api/write", func(w http.ResponseWriter, r *http.Request) {
		response := server.handleConfigWrite(w, r)
		w.WriteHeader(response.Code)
		encoder := json.NewEncoder(w)
		encoder.Encode(&response.Body)
	})
	router.HandleFunc("/api/config", func(w http.ResponseWriter, r *http.Request) {
		response := server.handleConfigGet(w, r)
		w.WriteHeader(response.Code)
		encoder := json.NewEncoder(w)
		encoder.Encode(&response.Body)
	})
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(indexHTML))
	})
	router.HandleFunc("/main.css", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mainCSS))
	})
	router.HandleFunc("/main.js", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mainJS))
	})
	handler := http.NewServeMux()
	handler.Handle("/", router)
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", server.config.BindAddr, server.config.Port), handler)
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
