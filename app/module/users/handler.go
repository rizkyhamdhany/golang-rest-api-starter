package users

import (
	"net/http"
	"github.com/rizkyhamdhany/jobagency-back/app/handler"
	"github.com/jinzhu/gorm"
	"github.com/rizkyhamdhany/jobagency-back/app/cc"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/validator.v9"
)

var validate *validator.Validate

func init(){
	validate = validator.New()
}

var res = make(map[string] interface{})

// User Signup
func Signup(db *gorm.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		user := &User{
			Name : r.FormValue("name"),
			Email : r.FormValue("email"),
			Phone : r.FormValue("phone"),
			Wa : r.FormValue("wa"),
			Password : r.FormValue("password"),
		}
		err := validate.Struct(user)

		if err == nil {
			temp := User{}
			db.First(&temp, "email = ?", r.FormValue("email"))
			if temp.ID == 0 {
				password := []byte(r.FormValue("password"))
				hashedPassword, _ := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
				user.Password = string(hashedPassword)
				db.Create(&user)
				handler.RespondWithJson(w, http.StatusOK, "")
			} else {
				handler.RespondWithError(w, http.StatusInternalServerError, cc.MsgEmailAlreadyUsed)
			}
		} else {
			handler.RespondWithError(w, http.StatusInternalServerError, handler.GetValidatorMsg(err))
		}
	}
}

// User Login
func Login(db *gorm.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		email := r.FormValue("email")
		password := r.FormValue("password")
		if email != "" && password != "" {
			user := User{}
			db.First(&user, "email = ?", email)
			err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

			if err == nil{
				tokenString := getTokenFromUser(user)
				res["token"] = tokenString
				handler.RespondWithJson(w, http.StatusOK, res)
			} else {
				handler.RespondWithError(w, http.StatusNotFound, cc.MsgCredentialNotFound)
			}
		} else {
			handler.RespondWithError(w, http.StatusInternalServerError, cc.MsgCredentialRequired)
		}

	}
}

// User Information by Token
func GetMe(db *gorm.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		user := getUserFromToken(w, r)
		if (user.ID != 0){
			res["user"] = user
			handler.RespondWithJson(w, http.StatusOK, res)
		} else {
			handler.RespondWithError(w, http.StatusInternalServerError, cc.MsgSomethingWentWrong)
		}
	}
}