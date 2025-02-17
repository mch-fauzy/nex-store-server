package handlers

import (
	"net/http"

	"github.com/nexmedis-be-technical-test/models/dto"
	"github.com/nexmedis-be-technical-test/utils/response"
)

func (h *Handler) ProductGetListByFilter(w http.ResponseWriter, r *http.Request) {
	var request dto.ProductGetByFilterRequest
	request.Page = r.URL.Query().Get("page")
	request.PageSize = r.URL.Query().Get("page_size")
	request.Name = r.URL.Query().Get("name")
	err := request.Validate()
	if err != nil {
		response.WithError(w, err)
		return
	}

	result, metadata, err := h.Service.ProductGetListByFilter(request)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithMetadata(w, http.StatusOK, result, metadata)
}
