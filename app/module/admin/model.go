package admin

import (
	"time"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"strings"
	"github.com/rizkyhamdhany/kelase-micro/app/cc"
)

type Admin struct {
	ID        	uint `gorm:"primary_key" json:"id"`
	Name    	string	`json:"name"`
	Email 		string	`json:"email"`
	Phone 		string	`json:"phone"`
	Role 		string	`json:"roles"`
	Password 	string	`json:"-"`
	CreatedAt 	time.Time
	UpdatedAt 	time.Time
	DeletedAt 	*time.Time `sql:"index"`
}

func getTokenFromUser(user Admin) string{
	fmt.Println(user.ID)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"ID": user.ID,
		"Name": user.Name,
		"Email": user.Email,
		"Phone": user.Phone,
		"Roles": user.Role,
		"Password": user.Password,
	})
	tokenString, _ := token.SignedString([]byte(cc.JwtSecret))
	return tokenString
}

func getUserFromToken(w http.ResponseWriter, r *http.Request) Admin  {
	reqToken := r.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer ")
	reqToken = splitToken[1]

	token, _ := jwt.Parse(reqToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(cc.JwtSecret), nil
	})
	user := Admin{}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		id := (claims["ID"]).(float64)
		user.ID = uint(id)
		user.Name = claims["Name"].(string)
		user.Email = claims["Email"].(string)
		user.Phone = claims["Phone"].(string)
		user.Role = claims["Roles"].(string)
		user.Password = claims["Password"].(string)
	}

	return user
}
