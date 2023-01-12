package userprofilecontroller

import (
	"encoding/json"
	data "muxapp/internal/database"
	"muxapp/internal/model"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

var ResponseJson = model.ResponseJson
var ResponseError = model.ResponseError

func IndexUserProfile(w http.ResponseWriter, r *http.Request) {
	res := []model.UserProfile{}

	if err := data.DB.Preload("Avatar").Find(&res).Error; err != nil {
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

func CreateUserProfile(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		ResponseError(w, http.StatusBadRequest, map[string]interface{}{
			"code":    http.StatusBadRequest,
			"message": "Error",
		})
		return
	}
	var res model.UserProfile

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&res); err != nil {
		ResponseError(w, http.StatusBadRequest, map[string]interface{}{
			"code":    http.StatusBadRequest,
			"message": "Error",
		})
		return
	}

	defer r.Body.Close()

	res.UserId = id
	if err := data.DB.Omit("Avatars").Create(&res).Error; err != nil {
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
}

func ShowUserProfile(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		ResponseError(w, http.StatusBadRequest, map[string]interface{}{
			"code":    http.StatusBadRequest,
			"message": "Error",
		})
		return
	}

	var res model.UserProfile
	if err := data.DB.Preload("Avatar").First(&res, id).Error; err != nil {
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
}

func UpdateUserProfile(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		ResponseError(w, http.StatusBadRequest, map[string]interface{}{
			"code":    http.StatusBadRequest,
			"message": "Error",
		})
		return
	}

	var user model.UserProfile

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		ResponseError(w, http.StatusInternalServerError, map[string]interface{}{
			"code":    http.StatusInternalServerError,
			"message": "Error",
		})
		return
	}

	defer r.Body.Close()

	if data.DB.Where("user_profile_id = ?", id).Updates(&user).RowsAffected == 0 {
		ResponseError(w, http.StatusNotFound, map[string]interface{}{
			"code":    http.StatusNotFound,
			"message": "Error",
		})
		return
	}

	ResponseJson(w, http.StatusCreated, map[string]interface{}{
		"code":    http.StatusCreated,
		"message": "Success",
	})
}

func DeleteUserProfile(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		ResponseError(w, http.StatusBadRequest, map[string]interface{}{
			"code":    http.StatusBadRequest,
			"message": "Error",
		})
		return
	}

	var user model.UserProfile
	if data.DB.Delete(&user, id).RowsAffected == 0 {
		ResponseError(w, http.StatusBadRequest, map[string]interface{}{
			"code":    http.StatusBadRequest,
			"message": "Error",
		})
		return
	}

	ResponseJson(w, http.StatusCreated, map[string]interface{}{
		"code":    http.StatusCreated,
		"message": "Success",
	})
}
