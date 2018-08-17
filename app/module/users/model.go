package users

import (
	"net/http"
	"strings"
	"github.com/dgrijalva/jwt-go"
	"fmt"
	"github.com/rizkyhamdhany/jobagency-back/app/cc"
	"time"
)

type User struct {
	ID        	uint `gorm:"primary_key" json:"id"`
	Name    	string	`json:"name" validate:"required"`
	Email 		string	`json:"email" validate:"required,email"`
	Phone 		string	`json:"phone" validate:"required"`
	Wa 			string	`json:"wa" validate:"required"`
	Password 	string	`json:"-" validate:"required,min=8"`
	Status 		string	`json:"status" validate:"-"`
	CreatedAt 	time.Time
	UpdatedAt 	time.Time
	DeletedAt 	*time.Time `sql:"index"`
}

func getTokenFromUser(user User) string{
	fmt.Println(user.ID)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"ID": user.ID,
		"Name": user.Name,
		"Email": user.Email,
		"Phone": user.Phone,
		"Wa": user.Wa,
		"Password": user.Password,
		"Status": user.Status,
	})
	tokenString, _ := token.SignedString([]byte(cc.JwtSecret))
	return tokenString
}

func getUserFromToken(w http.ResponseWriter, r *http.Request) User  {
	reqToken := r.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer ")
	reqToken = splitToken[1]

	token, _ := jwt.Parse(reqToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(cc.JwtSecret), nil
	})
	user := User{}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		id := (claims["ID"]).(float64)
		user.ID = uint(id)
		user.Name = claims["Name"].(string)
		user.Email = claims["Email"].(string)
		user.Phone = claims["Phone"].(string)
		user.Wa = claims["Wa"].(string)
		user.Password = claims["Password"].(string)
		user.Status = claims["Status"].(string)
	}

	return user
}