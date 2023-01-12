package newscontroller

import (
	"muxapp/internal/bucket"
	data "muxapp/internal/database"
	"muxapp/internal/model"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

var ResponseJson = model.ResponseJson
var ResponseError = model.ResponseError
var Roleset string

func IndexNews(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	limit, err := strconv.Atoi(params["limit"])
	if err != nil {
		ResponseError(w, http.StatusBadRequest, map[string]interface{}{
			"code":    http.StatusBadRequest,
			"message": "Error",
		})
		return
	}
	offset, err := strconv.Atoi(params["offset"])
	if err != nil {
		ResponseError(w, http.StatusBadRequest, map[string]interface{}{
			"code":    http.StatusBadRequest,
			"message": "Error",
		})
		return
	}

	res := []model.News{}

	if err := data.DB.Limit(limit).Offset(offset).Order("created_at DESC").Find(&res).Error; err != nil {
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

var ShowNews = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		ResponseError(w, http.StatusBadRequest, map[string]interface{}{
			"code":    http.StatusBadRequest,
			"message": "Error",
		})
		return
	}

	res := model.News{}

	if err := data.DB.Order("created_at DESC").First(&res, id).Error; err != nil {
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
		"message": "Succses",
		"data":    res,
	})
})

var CreateNews = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	if Roleset != "operator" {
		ResponseError(w, http.StatusUnauthorized, map[string]interface{}{
			"code":    http.StatusUnauthorized,
			"message": "Error",
		})
		return
	}
	news := model.News{}

	// Get title and content from form
	title := r.FormValue("title")
	content := r.FormValue("content")

	// To rename file
	currentTime := time.Now().Format("2006#01#02#15#04#05")
	code := strings.ReplaceAll(currentTime, "#", "")

	// Get file
	file, handler, err := r.FormFile("thumbnail")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Procced upload file to minio
	url, err := bucket.Open(file, handler, "images", "thumbnail", code)
	if err != nil {
		ResponseError(w, http.StatusBadRequest, map[string]interface{}{
			"code":    http.StatusBadRequest,
			"message": "Error",
		})
		return
	}

	// Fill model news
	news.Thumbnail = url
	news.Content = content
	news.Title = title

	// Submit to database
	if err := data.DB.Create(&news).Error; err != nil {
		ResponseError(w, http.StatusInternalServerError, map[string]interface{}{
			"code":    http.StatusInternalServerError,
			"message": "Error",
		})
	}

	ResponseJson(w, http.StatusOK, map[string]interface{}{
		"code":    http.StatusOK,
		"message": "Succses",
	})
	return
})

var UpdateNews = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

	// Get title and content from form
	title := r.FormValue("title")
	content := r.FormValue("content")

	if err := data.DB.Where("news_id = ?", id).Updates(model.News{
		Title:   title,
		Content: content,
	}).Error; err != nil {
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

	return
})

var DeleteNews = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

	var res model.News
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
	return
})

func SearchNews(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	keyword := params["keyword"]

	var news []model.News
	if err := data.DB.Limit(10).Order("created_at DESC").Where("title like ?", "%"+keyword+"%").Find(&news).Error; err != nil {
		ResponseJson(w, http.StatusNotFound, map[string]interface{}{
			"code":    http.StatusNotFound,
			"message": err,
			"data":    nil,
		})
		return
	}
	ResponseJson(w, http.StatusOK, map[string]interface{}{
		"code":    http.StatusOK,
		"message": "Success",
		"data":    news,
	})
}

func GetNewsLimit(w http.ResponseWriter, r *http.Request) {
	res := []model.News{}

	if err := data.DB.Select("created_at, SUBSTRING(`title`, 1, 50) as title, SUBSTRING(`content`, 1, 50) as content, thumbnail ").Order("created_at DESC").Find(&res).Error; err != nil {
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
