package users

import (
	"encoding/json"
	"fmt"
	"github.com/users-api/cmd/internal/server"
	"net/http"
)

type IUserHandler interface {
	InsertUser(w http.ResponseWriter, r *http.Request)
	UpdateUser(w http.ResponseWriter, r *http.Request)
	GetUser(w http.ResponseWriter, r *http.Request)
	FindUsers(w http.ResponseWriter, r *http.Request)
	DeleteUser(w http.ResponseWriter, r *http.Request)
}

type UserHandler struct {
}

func (handler *UserHandler) InsertUser(w http.ResponseWriter, r *http.Request) {
	var user *models.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		server.BadRequest(w, r, server.ErrorCodeInvalidParams, err.Error())
		return
	}

	insertId, err := handler.userService.InsertUser(user)
	if err == nil {
		id := make(map[string]interface{})
		id["id"] = insertId
		server.OK(w, r, id)
	} else {
		server.InternalServerError(w, r, err)
	}
}

func (handler *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userId := server.GetStringFromPath(r, "user_id", "")
	if userId == "" {
		err := fmt.Errorf("user id is mandatory")
		server.BadRequest(w, r, server.ErrorCodeInvalidParams, err.Error())
		return
	}

	var user *models.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		server.BadRequest(w, r, server.ErrorCodeInvalidParams, err.Error())
		return
	}

	if err := handler.userService.UpdateUser(userId, user); err != nil {
		server.InternalServerError(w, r, err)
	} else {
		ok := make(map[string]interface{})
		ok["ok"] = true
		server.OK(w, r, ok)
	}
}

func (handler *UserHandler) DeliverUser(w http.ResponseWriter, r *http.Request) {
	userId := server.GetStringFromPath(r, "user_id", "")
	if userId == "" {
		err := fmt.Errorf("user id is mandatory")
		server.BadRequest(w, r, server.ErrorCodeInvalidParams, err.Error())
		return
	}

	if err := handler.userService.DeliverUser(userId); err != nil {
		server.InternalServerError(w, r, err)
	} else {
		ok := make(map[string]interface{})
		ok["ok"] = true
		server.OK(w, r, ok)
	}
}

func (handler *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	userId := server.GetStringFromPath(r, "user_id", "")
	if userId == "" {
		err := fmt.Errorf("user id is mandatory")
		server.BadRequest(w, r, server.ErrorCodeInvalidParams, err.Error())
		return
	}

	user, err := handler.userService.GetUser(userId)

	if err != nil {
		server.InternalServerError(w, r, err)
		return
	}

	if user == nil {
		server.NotFound(w, r, fmt.Sprintf("Not found user with id: %v", userId))
		return
	}

	server.OK(w, r, user)
}

func (handler *UserHandler) FindUsers(w http.ResponseWriter, r *http.Request) {
	routeId := server.GetStringParam(r, "routeId", "")

	active, err := server.GetBoolParamPointer(r, "active")
	if err != nil {
		server.BadRequest(w, r, server.ErrorCodeInvalidParams, err.Error())
		return
	}

	size, err := server.GetIntParam(r, "size", 0)
	if err != nil {
		server.BadRequest(w, r, server.ErrorCodeInvalidParams, ErrorMessageSizeInvalid)
		return
	}

	offset, err := server.GetIntParam(r, "offset", 0)
	if err != nil {
		server.BadRequest(w, r, server.ErrorCodeInvalidParams, ErrorMessageOffsetInvalid)
		return
	}

	users, err := handler.userService.FindUsers(routeId, active, size, offset)

	if err == nil {
		server.OK(w, r, users)
	} else {
		server.InternalServerError(w, r, err)
	}
}

func RegisterRoutes(s server.Server) IUserHandler {
	UserHandler{} &UserHandler{
		userService: userService,
	}

	s.AddRoute("/v1/orders", orderHandler.InsertOrder, http.MethodPost)
	s.AddRoute("/v1/orders/{order_id}", orderHandler.UpdateOrder, http.MethodPut)
	s.AddRoute("/v1/orders/{order_id}/delivery", orderHandler.DeliverOrder, http.MethodPut)
	s.AddRoute("/v1/orders/{order_id}", orderHandler.GetOrder, http.MethodGet)
	s.AddRoute("/v1/orders", orderHandler.FindOrders, http.MethodGet)

	return &UserHandler{
		userService: userService,
	}
}
