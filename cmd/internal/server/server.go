package server

import (
	"encoding/json"
	"fmt"
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
	s.httpSrv.ListenAndServe()
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func Render(w http.ResponseWriter, r *http.Request, obj interface{}, status int) {
	js, err := json.Marshal(obj)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
}

func GetStringFromPath(r *http.Request, key string, defaultValue string) string {
	str := mux.Vars(r)[key]

	if len(str) < 1 {
		return defaultValue
	}

	return str
}

func OK(w http.ResponseWriter, r *http.Request, obj interface{}) {
	Render(w, r, obj, http.StatusOK)
}

func OkNotContent(w http.ResponseWriter, r *http.Request) {
	Render(w, r, nil, http.StatusNoContent)
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

// GetBoolParam parses a bool parameter from query param
func GetBoolParamPointer(r *http.Request, key string) (*bool, error) {
	keys, ok := r.URL.Query()[key]

	if !ok || len(keys[0]) < 1 {
		return nil, nil
	}

	value, err := strconv.ParseBool(keys[0])

	if err != nil {
		return nil, fmt.Errorf("%s is not a valid bool value for %s", keys[0], key)
	}

	return &value, nil
}

func NotFound(w http.ResponseWriter, r *http.Request, messages ...string) {
	err := &errorResponse{
		Code:     "NOT_FOUND",
		Messages: messages,
	}
	Render(w, r, err, http.StatusNotFound)
}

func Forbidden(w http.ResponseWriter, r *http.Request, messages ...string) {
	err := &errorResponse{
		Code:     "FORBIDDEN",
		Messages: messages,
	}
	Render(w, r, err, http.StatusForbidden)
}

func InternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	Render(w, r, &errorResponse{
		Code:     "INTERNAL_SERVER_ERROR",
		Messages: []string{err.Error()},
	}, http.StatusInternalServerError)
}

// GetIntParam parses a int parameter from url
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

// GetUintParam parses a uint (unsigned integer type) parameter from url
func GetUintParam(r *http.Request, key string, defaultValue uint) (uint, error) {
	keys, ok := r.URL.Query()[key]

	if !ok || len(keys[0]) < 1 {
		return defaultValue, nil
	}

	value, err := strconv.ParseUint(keys[0], 10, 0)

	if err != nil {
		return defaultValue,
			fmt.Errorf("%s is not a valid uint value for %s", keys[0], key)
	}

	return uint(value), nil
}

// GetInt32Param parses a int parameter from url
func GetInt32Param(r *http.Request, key string, defaultValue int32) (int32, error) {
	keys, ok := r.URL.Query()[key]

	if !ok || len(keys[0]) < 1 {
		return defaultValue, nil
	}

	value, err := strconv.ParseInt(keys[0], 10, 32)

	if err != nil {
		return defaultValue,
			fmt.Errorf("%s is not a valid int value for %s", keys[0], key)
	}

	return int32(value), nil
}

// GetUint32Param parses a uint32 (unsigned integer type) parameter from url
func GetUint32Param(r *http.Request, key string, defaultValue uint32) (uint32, error) {
	keys, ok := r.URL.Query()[key]

	if !ok || len(keys[0]) < 1 {
		return defaultValue, nil
	}

	value, err := strconv.ParseUint(keys[0], 10, 32)

	if err != nil {
		return defaultValue,
			fmt.Errorf("%s is not a valid uint32 value for %s", keys[0], key)
	}

	return uint32(value), nil
}

// GetInt64Param parses a int parameter from url
func GetInt64Param(r *http.Request, key string, defaultValue int64) (int64, error) {
	keys, ok := r.URL.Query()[key]

	if !ok || len(keys[0]) < 1 {
		return defaultValue, nil
	}

	value, err := strconv.ParseInt(keys[0], 10, 64)

	if err != nil {
		return defaultValue,
			fmt.Errorf("%s is not a valid int value for %s", keys[0], key)
	}

	return value, nil
}

// GetUint64Param parses a uint64 (unsigned integer type) parameter from url
func GetUint64Param(r *http.Request, key string, defaultValue uint64) (uint64, error) {
	keys, ok := r.URL.Query()[key]

	if !ok || len(keys[0]) < 1 {
		return defaultValue, nil
	}

	value, err := strconv.ParseUint(keys[0], 10, 64)

	if err != nil {
		return defaultValue,
			fmt.Errorf("%s is not a valid uint64 value for %s", keys[0], key)
	}

	return uint64(value), nil
}

// GetStringParam returns a string param from query string
func GetStringParam(r *http.Request, key string, defaultValue string) string {
	keys, ok := r.URL.Query()[key]

	if !ok || len(keys[0]) < 1 {
		return defaultValue
	}

	return keys[0]
}

// GetBoolParam parses a bool parameter from url
func GetBoolParam(r *http.Request, key string, defaultValue bool) (bool, error) {
	keys, ok := r.URL.Query()[key]

	if !ok || len(keys[0]) < 1 {
		return defaultValue, nil
	}

	value, err := strconv.ParseBool(keys[0])

	if err != nil {
		return defaultValue,
			fmt.Errorf("%s is not a valid bool value for %s", keys[0], key)
	}

	return value, nil
}

// GetFloat32FromPath : gets a float32 value from a named path part
func GetFloat32FromPath(r *http.Request, key string) (float32, error) {
	str := mux.Vars(r)[key]
	val, err := strconv.ParseFloat(str, 32)

	if err != nil {
		return 0, fmt.Errorf("%s is not a valid float value for %s", str, key)
	}

	return float32(val), nil
}

// GetFloat64FromPath : gets a float64 value from a named path part
func GetFloat64FromPath(r *http.Request, key string) (float64, error) {
	str := mux.Vars(r)[key]
	val, err := strconv.ParseFloat(str, 64)

	if err != nil {
		return 0, fmt.Errorf("%s is not a valid float value for %s", str, key)
	}

	return val, nil
}

// GetFloat64PtrParam returns a *float64 or nil from query string by key
func GetFloat64PtrParam(r *http.Request, key string) (*float64, error) {
	keys, ok := r.URL.Query()[key]

	if !ok || len(keys) == 0 {
		return nil, nil
	}

	value, err := strconv.ParseFloat(keys[0], 64)

	if err != nil {
		return nil,
			fmt.Errorf("%s is not a valid float value for %s", keys[0], key)
	}

	return &value, nil
}

// ValidateParams checks if all params have a unique value and are in
// valid params list
func ValidateParams(r *http.Request, validParams ...string) error {
	for k, v := range r.URL.Query() {

		if len(v) > 1 {
			return fmt.Errorf("multiple values for %s param", k)
		}

		valid := false
		for _, v := range validParams {

			if v == k {
				valid = true
			}
		}

		if !valid {
			return fmt.Errorf("unknown param: %s", k)
		}

	}

	return nil
}
