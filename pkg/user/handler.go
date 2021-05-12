package user

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/users-api/cmd/server"
	"net/http"
)

type IHandler interface {
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	Find(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)

	GetLocation(w http.ResponseWriter, r *http.Request)
}

const (
	ErrorCodeInvalidParams    string = "INVALID_PARAMS"
	ErrorMessageSizeInvalid   string = "The size is invalid"
	ErrorMessageOffsetInvalid string = "The offset is invalid"
)

type Handler struct {
	service IService
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var user *User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		server.BadRequest(w, r, ErrorCodeInvalidParams, err.Error())
		return
	}

	if err := user.Validate(); err != nil {
		server.BadRequest(w, r, ErrorCodeInvalidParams, err.Error())
		return
	}

	id, err := h.service.Create(user)
	if err != nil {
		server.InternalServerError(w, r, err)
		return
	}

	server.Created(w, r, fmt.Sprintf("user/%d", id))
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	userId, err := server.GetIntFromPath(r, "id")
	if err != nil {
		server.BadRequest(w, r, "")
		return
	}

	var user *User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		server.BadRequest(w, r, ErrorCodeInvalidParams, err.Error())
		return
	}

	if err := user.Validate(); err != nil {
		server.BadRequest(w, r, ErrorCodeInvalidParams, err.Error())
		return
	}

	if err := h.service.Update(userId, user); err != nil {
		handlerException(w, r, err)
		return
	}

	server.OK(w, r, nil)

}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	userId, err := server.GetIntFromPath(r, "id")
	if err != nil {
		server.BadRequest(w, r, "")
		return
	}

	if err := h.service.Delete(userId); err != nil {
		handlerException(w, r, err)
		return
	}

	server.OK(w, r, nil)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	userId, err := server.GetIntFromPath(r, "id")
	if err != nil {
		server.BadRequest(w, r, "")
		return
	}

	user, err := h.service.Get(userId)

	if err != nil {
		handlerException(w, r, err)
		return
	}

	server.OK(w, r, user)
}

func (h *Handler) Find(w http.ResponseWriter, r *http.Request) {
	name := server.GetStringParam(r, "name", "")
	if name == "" {
		server.BadRequest(w, r, ErrorCodeInvalidParams, "name is mandatory")
		return
	}

	size, err := server.GetIntParam(r, "size", 0)
	if err != nil {
		server.BadRequest(w, r, ErrorCodeInvalidParams, ErrorMessageSizeInvalid)
		return
	}

	offset, err := server.GetIntParam(r, "offset", 0)
	if err != nil {
		server.BadRequest(w, r, ErrorCodeInvalidParams, ErrorMessageOffsetInvalid)
		return
	}

	users, err := h.service.Find(name, size, offset)

	if err != nil {
		handlerException(w, r, err)
		return
	}

	server.OK(w, r, users)

}

func (h *Handler) GetLocation(w http.ResponseWriter, r *http.Request) {
	userId, err := server.GetIntFromPath(r, "id")
	if err != nil {
		server.BadRequest(w, r, "")
		return
	}

	location, err := h.service.GetLocation(userId)

	if err != nil {
		handlerException(w, r, err)
		return
	}

	server.OK(w, r, location)
}

func handlerException(w http.ResponseWriter, r *http.Request, err error) {
	switch err.(type) {
	case *NotFoundError:
		server.NotFound(w, r, err.Error())
	default:
		logrus.Error(err)
		server.InternalServerError(w, r, err)
	}
}

func RegisterRoutes(s *server.Server) {

	handler := newHandler()

	s.AddRoute("/v{version}/users", handler.Create, http.MethodPost)
	s.AddRoute("/v{version}/users/{id}", handler.Update, http.MethodPut)
	s.AddRoute("/v{version}/users/{id}", handler.Get, http.MethodGet)
	s.AddRoute("/v{version}/users", handler.Find, http.MethodGet)
	s.AddRoute("/v{version}/users/{id}", handler.Delete, http.MethodDelete)

	s.AddRoute("/v{version}/users/{id}/locations", handler.GetLocation, http.MethodGet)
}

func newHandler() IHandler {
	return &Handler{
		service: NewService(),
	}
}
