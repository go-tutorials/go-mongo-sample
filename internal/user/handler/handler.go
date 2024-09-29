package handler

import (
	"context"
	"fmt"
	"net/http"
	"reflect"

	"github.com/core-go/core"
	"github.com/core-go/search"

	"go-service/internal/user/model"
	"go-service/internal/user/service"
)

func NewUserHandler(service service.UserService, logError func(context.Context, string, ...map[string]interface{}), validate func(context.Context, interface{}) ([]core.ErrorMessage, error), action *core.ActionConfig) *UserHandler {
	userType := reflect.TypeOf(model.User{})
	parameters := search.CreateParameters(reflect.TypeOf(model.UserFilter{}), userType)
	attributes := core.CreateAttributes(userType, logError, action)
	return &UserHandler{service: service, Validate: validate, Attributes: attributes, Parameters: parameters}
}

type UserHandler struct {
	service  service.UserService
	Validate func(context.Context, interface{}) ([]core.ErrorMessage, error)
	*core.Attributes
	*search.Parameters
}

func (h *UserHandler) All(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.All(r.Context())
	if err != nil {
		h.Error(r.Context(), fmt.Sprintf("Error: %s", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	core.JSON(w, http.StatusOK, users)
}
func (h *UserHandler) Load(w http.ResponseWriter, r *http.Request) {
	id, err := core.GetRequiredString(w, r)
	if err == nil {
		user, err := h.service.Load(r.Context(), id)
		if err != nil {
			h.Error(r.Context(), fmt.Sprintf("Error to get user '%s': %s", id, err.Error()))
			http.Error(w, core.InternalServerError, http.StatusInternalServerError)
			return
		}
		core.JSON(w, core.IsFound(user), user)
	}
}
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var user model.User
	er1 := core.Decode(w, r, &user)
	if er1 == nil {
		errors, er2 := h.Validate(r.Context(), &user)
		if !core.HasError(w, r, errors, er2, h.Error, user, h.Log, h.Resource, h.Action.Create) {
			res, er3 := h.service.Create(r.Context(), &user)
			if er3 != nil {
				h.Error(r.Context(), er3.Error(), core.MakeMap(user))
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			if res > 0 {
				core.JSON(w, http.StatusCreated, user)
			} else {
				core.JSON(w, http.StatusConflict, res)
			}
		}
	}
}
func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	var user model.User
	er1 := core.DecodeAndCheckId(w, r, &user, h.Keys, h.Indexes)
	if er1 == nil {
		errors, er2 := h.Validate(r.Context(), &user)
		if !core.HasError(w, r, errors, er2, h.Error, user, h.Log, h.Resource, h.Action.Update) {
			res, er3 := h.service.Update(r.Context(), &user)
			if er3 != nil {
				h.Error(r.Context(), er3.Error(), core.MakeMap(user))
				http.Error(w, er3.Error(), http.StatusInternalServerError)
				return
			}
			if res > 0 {
				core.JSON(w, http.StatusOK, user)
			} else if res == 0 {
				core.JSON(w, http.StatusNotFound, res)
			} else {
				core.JSON(w, http.StatusConflict, res)
			}
		}
	}
}
func (h *UserHandler) Patch(w http.ResponseWriter, r *http.Request) {
	var user model.User
	r, jsonUser, er1 := core.BuildMapAndCheckId(w, r, &user, h.Keys, h.Indexes)
	if er1 == nil {
		errors, er2 := h.Validate(r.Context(), &user)
		if !core.HasError(w, r, errors, er2, h.Error, jsonUser, h.Log, h.Resource, h.Action.Patch) {
			res, er3 := h.service.Patch(r.Context(), jsonUser)
			if er3 != nil {
				h.Error(r.Context(), er3.Error(), core.MakeMap(jsonUser))
				http.Error(w, er3.Error(), http.StatusInternalServerError)
				return
			}
			if res > 0 {
				core.JSON(w, http.StatusOK, jsonUser)
			} else if res == 0 {
				core.JSON(w, http.StatusNotFound, res)
			} else {
				core.JSON(w, http.StatusConflict, res)
			}
		}
	}
}
func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := core.GetRequiredString(w, r)
	if err == nil {
		res, err := h.service.Delete(r.Context(), id)
		if err != nil {
			h.Error(r.Context(), fmt.Sprintf("Error to delete user '%s': %s", id, err.Error()))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if res > 0 {
			core.JSON(w, http.StatusOK, res)
		} else {
			core.JSON(w, http.StatusNotFound, res)
		}
	}
}
func (h *UserHandler) Search(w http.ResponseWriter, r *http.Request) {
	filter := model.UserFilter{Filter: &search.Filter{}}
	err := search.Decode(r, &filter, h.ParamIndex, h.FilterIndex)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	offset := search.GetOffset(filter.Limit, filter.Page)
	users, total, err := h.service.Search(r.Context(), &filter, filter.Limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	core.JSON(w, http.StatusOK, &search.Result{List: &users, Total: total})
}
