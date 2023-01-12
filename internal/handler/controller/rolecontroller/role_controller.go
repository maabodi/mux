package rolecontroller

import (
	"encoding/json"
	data "muxapp/internal/database"
	"muxapp/internal/model"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

var ResponseError = model.ResponseError
var ResponseJson = model.ResponseJson
var Roleset string

func IndexRole(w http.ResponseWriter, r *http.Request) {
	var res []model.Role

	if err := data.DB.Find(&res).Error; err != nil {
		ResponseError(w, http.StatusInternalServerError, map[string]interface{}{
			"code":    http.StatusInternalServerError,
			"message": "Error",
		})
		return
	}

	ResponseJson(w, http.StatusOK, map[string]interface{}{
		"code":    http.StatusOK,
		"message": "Success",
		"data":    res,
	})
}

var CreateRole = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	if Roleset != "operator" {
		ResponseError(w, http.StatusUnauthorized, map[string]interface{}{
			"code":    http.StatusUnauthorized,
			"message": "Error",
		})
		return
	}

	var res model.Role

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&res); err != nil {
		ResponseError(w, http.StatusBadRequest, map[string]interface{}{
			"code":    http.StatusBadRequest,
			"message": "Error",
		})
		return
	}

	defer r.Body.Close()
	if err := data.DB.Create(&res).Error; err != nil {
		ResponseError(w, http.StatusInternalServerError, map[string]interface{}{
			"code":    http.StatusInternalServerError,
			"message": "Error",
		})
	}

	ResponseJson(w, http.StatusOK, map[string]interface{}{
		"code":    http.StatusOK,
		"message": "Success",
	})
})

var ShowRole = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		ResponseError(w, http.StatusBadRequest, map[string]interface{}{
			"code":    http.StatusBadRequest,
			"message": "Error",
		})
		return
	}

	var res model.Role

	if err := data.DB.First(&res, id).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			ResponseError(w, http.StatusNotFound, map[string]interface{}{
				"code":    http.StatusNotFound,
				"message": "Error",
			})
			return
		default:
			ResponseError(w, http.StatusInternalServerError, map[string]interface{}{
				"code":    http.StatusInternalServerError,
				"message": "Error",
			})
			return
		}
	}

	ResponseJson(w, http.StatusOK, map[string]interface{}{
		"code":    http.StatusOK,
		"message": "Success",
		"data":    res,
	})
})

var UpdateRole = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	if Roleset != "operator" {
		ResponseError(w, http.StatusUnauthorized, map[string]interface{}{
			"code":    http.StatusUnauthorized,
			"message": "Error",
		})
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		ResponseError(w, http.StatusBadRequest, map[string]interface{}{
			"code":    http.StatusNotFound,
			"message": "Error",
		})
		return
	}

	var res model.Role

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&res); err != nil {
		ResponseError(w, http.StatusInternalServerError, map[string]interface{}{
			"code":    http.StatusInternalServerError,
			"message": "Error",
		})
		return
	}

	defer r.Body.Close()

	if data.DB.Where("role_id = ?", id).Updates(&res).RowsAffected == 0 {
		ResponseError(w, http.StatusNotFound, map[string]interface{}{
			"code":    http.StatusNotFound,
			"message": "Error",
		})
	}

	ResponseJson(w, http.StatusOK, map[string]interface{}{
		"code":    http.StatusOK,
		"message": "Success",
	})
})

var DeleteRole = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	if Roleset != "operator" {
		ResponseError(w, http.StatusUnauthorized, map[string]interface{}{
			"code":    http.StatusUnauthorized,
			"message": "Error",
		})
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		ResponseError(w, http.StatusBadRequest, map[string]interface{}{
			"code":    http.StatusBadRequest,
			"message": "Error",
		})
		return
	}

	var res model.Role

	if data.DB.Delete(&res, id).RowsAffected == 0 {
		ResponseError(w, http.StatusInternalServerError, map[string]interface{}{
			"code":    http.StatusInternalServerError,
			"message": "Error",
		})
		return
	}

	ResponseJson(w, http.StatusOK, map[string]interface{}{
		"code":    http.StatusOK,
		"message": "Success",
	})
})
