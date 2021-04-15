package usershand

import (
	"encoding/json"
	"net/http"

	"golang-seed/apps/auth/pkg/messagesconst"
	"golang-seed/apps/auth/pkg/service/usersserv"
	"golang-seed/pkg/httperror"
	"golang-seed/pkg/messages"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type UsersHandler struct {
	usersService *usersserv.UsersService
}

func NewUsersHandler(usersService *usersserv.UsersService) *UsersHandler {
	return &UsersHandler{usersService: usersService}
}

func (u UsersHandler) Get(w http.ResponseWriter, r *http.Request) error {
	pathVars := mux.Vars(r)
	id := pathVars["id"]

	uid, err := uuid.Parse(id)
	if err != nil {
		return httperror.NewHTTPError(err, http.StatusBadRequest, messages.Get(messagesconst.GeneralErrorInvalidID))
	}

	user, err := u.usersService.Get(uid)
	if err != nil {
		return httperror.NewHTTPError(err, http.StatusInternalServerError, messages.Get(messagesconst.GeneralErrorGetting, messages.Get(messagesconst.UsersUser)))
	}

	body, err := json.Marshal(user)
	if err != nil {
		return httperror.NewHTTPError(err, http.StatusInternalServerError, messages.Get(messagesconst.GeneralErrorMarshal))
	}

	w.WriteHeader(http.StatusOK)
	w.Write(body)
	return nil
}

func (u UsersHandler) GetAll(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (u UsersHandler) GetAllPaged(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (u UsersHandler) Create(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (u UsersHandler) Update(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (u UsersHandler) Delete(w http.ResponseWriter, r *http.Request) error {
	return nil
}
