package helper

import (
	"encoding/json"
	data "muxapp/internal/database"
	nc "muxapp/internal/handler/controller/newscontroller"
	rc "muxapp/internal/handler/controller/rolecontroller"
	uc "muxapp/internal/handler/controller/usercontroller"
	"muxapp/internal/model"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var ResponseJson = model.ResponseJson
var ResponseError = model.ResponseError

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func Authenticate(w http.ResponseWriter, r *http.Request) {

	var user model.User

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		ResponseError(w, http.StatusInternalServerError, map[string]interface{}{
			"code":    http.StatusInternalServerError,
			"message": "Error",
			"data":    "Login failed",
		})
		return
	}

	var get_user model.User
	if err := data.DB.Where("email = ?", user.Email).Preload("Role").Find(&get_user).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			ResponseError(w, http.StatusNotFound, map[string]interface{}{
				"code":    http.StatusNotFound,
				"message": "Error",
				"data":    "Login failed",
			})
			return
		default:
			ResponseError(w, http.StatusInternalServerError, map[string]interface{}{
				"code":    http.StatusInternalServerError,
				"message": "Error",
				"data":    "Login failed",
			})
			return
		}
	}

	if len(get_user.Email) == 0 || len(get_user.Password) == 0 {
		ResponseError(w, http.StatusBadRequest, map[string]interface{}{
			"code":    http.StatusBadRequest,
			"message": "Error",
			"data":    "Login Failed",
		})
		return
	}

	match := CheckPasswordHash(user.Password, get_user.Password)

	if get_user.Email == user.Email && match == true {
		token, err := GetToken(int(get_user.UserID), get_user.Email, get_user.Role.RoleName)
		if err != nil {
			ResponseError(w, http.StatusInternalServerError, map[string]interface{}{
				"code":    http.StatusInternalServerError,
				"message": "Error",
				"data":    "Error generating JWT token ",
			})
		} else {
			w.Header().Set("Authorization", "Bearer "+token)
			ResponseJson(w, http.StatusOK, map[string]interface{}{
				"code":    http.StatusOK,
				"message": "Success",
				"data": map[string]interface{}{
					"user": map[string]interface{}{
						"email": get_user.Email,
						"role":  get_user.Role.RoleName,
					},
					"token": token,
				},
			})
		}
	} else {
		ResponseError(w, http.StatusBadRequest, map[string]interface{}{
			"code":    http.StatusBadRequest,
			"message": "Error",
			"data":    "Login failed",
		})

		return
	}
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tokenString := r.Header.Get("Authorization")
		if len(tokenString) == 0 {
			ResponseError(w, http.StatusUnauthorized, map[string]interface{}{
				"code":    http.StatusUnauthorized,
				"message": "Error",
				"data":    "Missing authorization header",
			})
			return
		}
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
		claims, err := VerifyToken(tokenString)
		if err != nil {
			ResponseError(w, http.StatusUnauthorized, map[string]interface{}{
				"code":    http.StatusUnauthorized,
				"message": "Error",
				"data":    "Missing or malformed jwt",
			})
			return
		}
		email := claims.(jwt.MapClaims)["email"].(string)
		role := claims.(jwt.MapClaims)["role"].(string)

		r.Header.Set("email", email)
		r.Header.Set("role", role)

		uc.Roleset = role
		rc.Roleset = role
		nc.Roleset = role

		next.ServeHTTP(w, r)
	})
}
