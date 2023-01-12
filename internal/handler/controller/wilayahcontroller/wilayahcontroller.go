package wilayahcontroller

import (
	data "muxapp/internal/database"
	"muxapp/internal/model"
	"net/http"
)

var ResponseError = model.ResponseError
var ResponseJson = model.ResponseJson

func Index(w http.ResponseWriter, r *http.Request) {
	res := []model.WProvinsi{}

	if err := data.DB.Preload("Kota.Kecamatan.Kelurahan").Find(&res).Error; err != nil {
		ResponseError(w, http.StatusInternalServerError, map[string]interface{}{
			"code":    http.StatusInternalServerError,
			"message": "Error",
		})
	}

	ResponseJson(w, http.StatusOK, map[string]interface{}{
		"code":    http.StatusOK,
		"message": "Success",
		"data":    res,
	})
}
