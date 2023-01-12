package avatarcontroller

import (
	"muxapp/internal/bucket"
	data "muxapp/internal/database"
	"muxapp/internal/model"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

var ResponseError = model.ResponseError
var ResponseJson = model.ResponseJson

func IndexAvatar(w http.ResponseWriter, r *http.Request) {
	var res []model.Avatar

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

func ChangeAvatar(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		ResponseError(w, http.StatusBadRequest, map[string]interface{}{
			"code":    http.StatusBadRequest,
			"message": "Error",
		})
		return
	}

	// To rename file
	currentTime := time.Now().Format("2006#01#02#15#04#05")
	code := strings.ReplaceAll(currentTime, "#", "")

	// Get file
	file, handler, err := r.FormFile("avatar")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Procced upload file to minio
	url, err := bucket.Open(file, handler, "images", "avatar", code)
	if err != nil {
		ResponseError(w, http.StatusBadRequest, map[string]interface{}{
			"code":    http.StatusBadRequest,
			"message": "Error 5",
		})
		return
	}
	res := model.Avatar{}
	res.UserId = id
	res.AvatarUrl = url

	if err := data.DB.Where("user_id = ? ", id).Updates(&res).Error; err != nil {
		ResponseError(w, http.StatusInternalServerError, map[string]interface{}{
			"code":    http.StatusInternalServerError,
			"message": err,
		})
		return
	}

	ResponseJson(w, http.StatusOK, map[string]interface{}{
		"code":    http.StatusOK,
		"message": "Success",
	})
	return
}
