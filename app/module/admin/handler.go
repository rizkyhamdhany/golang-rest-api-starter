package admin

import (
	"github.com/jinzhu/gorm"
	"net/http"
	"golang.org/x/crypto/bcrypt"
	"github.com/rizkyhamdhany/kelase-micro/app/handler"
	"github.com/rizkyhamdhany/kelase-micro/app/cc"
)

func Login(db *gorm.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		email := r.FormValue("email")
		password := r.FormValue("password")
		if email != "" && password != "" {
			user := Admin{}
			db.First(&user, "email = ?", email)
			err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

			if err == nil{
				tokenString := getTokenFromUser(user)
				handler.RespondWithJson(w, http.StatusOK, tokenString)
			} else {
				handler.RespondWithError(w, http.StatusNotFound, cc.MsgCredentialNotFound)
			}
		} else {
			handler.RespondWithError(w, http.StatusInternalServerError, cc.MsgCredentialRequired)
		}

	}
}
