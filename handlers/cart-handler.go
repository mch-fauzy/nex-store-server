package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/nexmedis-be-technical-test/models/dto"
	"github.com/nexmedis-be-technical-test/utils/constant"
	"github.com/nexmedis-be-technical-test/utils/response"
)

func (h *Handler) UserCartAddItem(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		response.WithError(w, err)
		return
	}

	var request dto.UserCartAddItemRequest
	err = json.Unmarshal(body, &request)
	if err != nil {
		response.WithError(w, err)
		return
	}

	request.UserId = r.Header.Get(constant.UserIdHeader)
	request.Email = r.Header.Get(constant.EmailHeader)
	err = request.Validate()
	if err != nil {
		response.WithError(w, err)
		return
	}

	_, err = h.Service.UserCartAddItem(request)
	if err != nil {
		response.WithError(w, err)
		return
	}

	result, err := h.Service.UserCartGetList(dto.UserCartGetListRequest{UserId: r.Header.Get(constant.UserIdHeader)})
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithData(w, http.StatusOK, result)
}
