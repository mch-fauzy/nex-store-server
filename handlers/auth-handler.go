package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/nexmedis-be-technical-test/models"
	"github.com/nexmedis-be-technical-test/models/dto"
	"github.com/nexmedis-be-technical-test/utils/response"
)

func (h *Handler) AuthRegisterUser(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		response.WithError(w, err)
		return
	}

	var request dto.AuthRegisterRequest
	err = json.Unmarshal(body, &request)
	if err != nil {
		response.WithError(w, err)
		return
	}

	request.RoleId = models.MasterUserRoleId.User
	err = request.Validate()
	if err != nil {
		response.WithError(w, err)
		return
	}

	msg, err := h.Service.AuthRegister(request)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithMessage(w, http.StatusCreated, msg)
}

func (h *Handler) AuthRegisterAdmin(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		response.WithError(w, err)
		return
	}

	var request dto.AuthRegisterRequest
	err = json.Unmarshal(body, &request)
	if err != nil {
		response.WithError(w, err)
		return
	}

	request.RoleId = models.MasterUserRoleId.Admin
	err = request.Validate()
	if err != nil {
		response.WithError(w, err)
		return
	}

	msg, err := h.Service.AuthRegister(request)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithMessage(w, http.StatusCreated, msg)
}

func (h *Handler) AuthLogin(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		response.WithError(w, err)
		return
	}

	var request dto.AuthLoginRequest
	err = json.Unmarshal(body, &request)
	if err != nil {
		response.WithError(w, err)
		return
	}

	err = request.Validate()
	if err != nil {
		response.WithError(w, err)
		return
	}

	result, err := h.Service.AuthLogin(request)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithData(w, http.StatusOK, result)
}
