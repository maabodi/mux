package usercontroller

import (
	"encoding/json"
	data "muxapp/internal/database"
	"muxapp/internal/model"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var ResponseJson = model.ResponseJson
var ResponseError = model.ResponseError
var Roleset string

func IndexUser(w http.ResponseWriter, r *http.Request) {
	var res []model.User

	if err := data.DB.Preload("Role").Preload("Avatar").Preload("UserProfile").Find(&res).Error; err != nil {
		ResponseError(w, http.StatusInternalServerError, map[string]interface{}{
			"code":    http.StatusInternalServerError,
			"message": "Error",
		})
		return
	}

	var resUser []model.ResponseUser
	if len(res) > 0 {
		for i := 0; i < len(res); i++ {
			var e = res[i]
			resUser = append(resUser, model.ResponseUser{
				Name:        e.Name,
				Email:       e.Email,
				Role:        e.Role.RoleName,
				UserProfile: e.UserProfile,
				Avatar:      e.Avatar,
			})
		}
	}

	ResponseJson(w, http.StatusOK, map[string]interface{}{
		"code":    http.StatusOK,
		"message": "Success",
		"data":    resUser,
	})
}

func CreateUser(w http.ResponseWriter, r *http.Request) {

	var res model.User

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&res); err != nil {
		ResponseError(w, http.StatusBadRequest, map[string]interface{}{
			"code":    http.StatusBadRequest,
			"message": "Error",
		})
		return
	}

	defer r.Body.Close()

	bo := CheckEmail(res.Email)
	if bo == true {
		ResponseError(w, http.StatusBadRequest, map[string]interface{}{
			"code":    http.StatusBadRequest,
			"message": "Error",
		})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(res.Password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	res.Password = string(hash)

	if err := data.DB.Omit("Avatar").Create(&res).Association("UserProfile").Error; err != nil {
		ResponseError(w, http.StatusInternalServerError, map[string]interface{}{
			"code":    http.StatusInternalServerError,
			"message": "Error",
		})
		return
	}
	var avaurl string
	if res.RoleId == 3 {
		avaurl = "http://127.0.0.1:9000/images/default/admin.png"
	} else if res.RoleId == 2 {
		avaurl = "http://127.0.0.1:9000/images/default/operator.png"
	} else {
		avaurl = "http://127.0.0.1:9000/images/default/user.png"
	}

	var ava = &model.Avatar{
		AvatarUrl: avaurl,
		UserId:    res.UserID,
	}

	if err := data.DB.Create(&ava).Error; err != nil {
		ResponseError(w, http.StatusInternalServerError, map[string]interface{}{
			"code":    http.StatusInternalServerError,
			"message": "Error",
		})
		return
	}

	ResponseJson(w, http.StatusCreated, map[string]interface{}{
		"code":    http.StatusCreated,
		"message": "Success",
	})
}

var ShowUser = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		ResponseError(w, http.StatusBadRequest, map[string]interface{}{
			"code":    http.StatusBadRequest,
			"message": "Error",
		})
		return
	}

	var res model.User
	if err := data.DB.Preload("Role").Preload("UserProfile").Preload("Avatar").First(&res, id).Error; err != nil {
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
	var resp = model.ResponseUser{
		Name:        res.Name,
		Email:       res.Email,
		Role:        res.Role.RoleName,
		UserProfile: res.UserProfile,
		Avatar:      res.Avatar,
	}

	ResponseJson(w, http.StatusOK, map[string]interface{}{
		"code":    http.StatusOK,
		"message": "Success",
		"data":    resp,
	})
})

var UpdateUser = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

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

	var user model.User

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		ResponseError(w, http.StatusInternalServerError, map[string]interface{}{
			"code":    http.StatusInternalServerError,
			"message": "Error",
		})
		return
	}

	defer r.Body.Close()

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	user.Password = string(hash)

	if data.DB.Where("user_id = ?", id).Updates(&user).RowsAffected == 0 {
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
})

var DeleteUser = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

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

	var user model.User
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
})

func CheckEmail(email string) bool {

	var user model.User
	res := data.DB.Where("email = ?", email).First(&user)
	if res.RowsAffected > 0 {
		return true
	}

	return false
}

func IndexUserByRole(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	role := params["role"]
	var res []model.User

	if err := data.DB.Preload("Avatar").Preload("UserProfile").Where("role = ?", role).Find(&res).Error; err != nil {
		ResponseError(w, http.StatusInternalServerError, map[string]interface{}{
			"code":    http.StatusInternalServerError,
			"message": "Error",
		})
		return
	}

	var resUser []model.ResponseUser
	if len(res) > 0 {
		for i := 0; i < len(res); i++ {
			var e = res[i]
			resUser = append(resUser, model.ResponseUser{
				Name:        e.Name,
				Email:       e.Email,
				Role:        e.Role.RoleName,
				UserProfile: e.UserProfile,
				Avatar:      e.Avatar,
			})
		}
	}

	ResponseJson(w, http.StatusOK, map[string]interface{}{
		"code":    http.StatusOK,
		"message": "Success",
		"data":    resUser,
	})
}

func GetOneForDashboard(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		ResponseError(w, http.StatusBadRequest, map[string]interface{}{
			"code":    http.StatusBadRequest,
			"message": "Error",
		})
		return
	}

	res := model.User{}

	if err := data.DB.Preload("Avatar").First(&res, id).Error; err != nil {
		ResponseError(w, http.StatusInternalServerError, map[string]interface{}{
			"code":    http.StatusInternalServerError,
			"message": "Error",
		})
		return
	}

	var respon = &model.ResUser{
		Name:   res.Name,
		Email:  res.Email,
		Avatar: res.Avatar,
	}
	// resuser := model.ResUser{Name: res.Name, Email: res.Email, Avatar: res.UserProfile.Avatar.AvatarUrl}
	ResponseJson(w, http.StatusOK, map[string]interface{}{
		"code":    http.StatusOK,
		"message": "Success",
		"data":    respon,
	})
}
