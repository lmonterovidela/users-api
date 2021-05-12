package server

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Config struct {
	Port    int    `json:"port"`
	Version string `json:"version"`
}

type Server struct {
	httpSrv *http.Server
	router  *mux.Router
	cfg     *Config
	Version string
}

func New(c *Config) *Server {
	r := mux.NewRouter()
	return &Server{
		cfg: c,
		httpSrv: &http.Server{
			Handler:      r,
			Addr:         fmt.Sprintf(":%d", c.Port),
			WriteTimeout: 60 * time.Second,
			ReadTimeout:  15 * time.Second,
		},
		router:  r,
		Version: c.Version,
	}
}

func (s *Server) AddRoute(path string, h http.HandlerFunc, methods ...string) {

	r := handlers.RecoveryHandler(handlers.PrintRecoveryStack(true))(h)

	g := http.TimeoutHandler(r, time.Duration(10)*time.Second, "response timeout exceeded")

	s.router.Handle(path, g).Methods(methods...)
}

func (s *Server) ListenAndServe() {
	if err := s.httpSrv.ListenAndServe(); err != nil {
		panic(err)
	}

}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func Render(w http.ResponseWriter, r *http.Request, obj interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if obj != nil {
		js, err := json.Marshal(obj)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = w.Write(js)
		if err != nil {
			logrus.Errorf(err.Error())
		}
	}
}

func OK(w http.ResponseWriter, r *http.Request, obj interface{}) {
	Render(w, r, obj, http.StatusOK)
}

func Created(w http.ResponseWriter, r *http.Request, location string) {
	w.Header().Add("location", location)
	Render(w, r, nil, http.StatusCreated)
}

type errorResponse struct {
	Messages []string `json:"messages"`
	Code     string   `json:"code"`
}

func BadRequest(w http.ResponseWriter, r *http.Request, code string, messages ...string) {
	err := &errorResponse{
		Code:     code,
		Messages: messages,
	}
	Render(w, r, err, http.StatusBadRequest)
}

func NotFound(w http.ResponseWriter, r *http.Request, messages ...string) {
	err := &errorResponse{
		Code:     "NOT_FOUND",
		Messages: messages,
	}
	Render(w, r, err, http.StatusNotFound)
}

func InternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	Render(w, r, &errorResponse{
		Code:     "INTERNAL_SERVER_ERROR",
		Messages: []string{err.Error()},
	}, http.StatusInternalServerError)
}

func GetIntParam(r *http.Request, key string, defaultValue int) (int, error) {
	keys, ok := r.URL.Query()[key]

	if !ok || len(keys[0]) < 1 {
		return defaultValue, nil
	}

	value, err := strconv.ParseInt(keys[0], 10, 0)

	if err != nil {
		return defaultValue,
			fmt.Errorf("%s is not a valid int value for %s", keys[0], key)
	}

	return int(value), nil
}

func GetStringParam(r *http.Request, key string, defaultValue string) string {
	keys, ok := r.URL.Query()[key]

	if !ok || len(keys[0]) < 1 {
		return defaultValue
	}

	return keys[0]
}

func GetIntFromPath(r *http.Request, key string) (int, error) {
	str := mux.Vars(r)[key]
	val, err := strconv.ParseInt(str, 10, 64)

	if err != nil {
		return 0, fmt.Errorf("%s is not a valid float value for %s", str, key)
	}

	return int(val), nil
}
