package handler

import (
	"net/http"
	"encoding/json"
	"github.com/rizkyhamdhany/jobagency-back/app/cc"
	"gopkg.in/go-playground/validator.v9"
)

type Res struct {
	Msg string `json:"msg" bson:"msg"`
	Status bool `json:"status" bson:"status"`
	Data *interface{} `json:"data" bson:"data"`
}

func RespondWithError(w http.ResponseWriter, code int, msg string) {
	var res = Res{msg, false, nil}
	response, _ := json.Marshal(res)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func RespondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	res := Res{cc.StatusOK, true, &payload}
	response, _ := json.Marshal(res)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func GetValidatorMsg(err error) string {
	msg := ""
	for _, err := range err.(validator.ValidationErrors) {
		if err.Tag() == "required" {
			msg += "The " + err.Field() + " field is required"
		}
		if err.Tag() == "email" {
			msg += "The " + err.Field() + " field is not correct format"
		}
		if err.Tag() == "min" {
			msg += "The " + err.Field() + " field min " + err.Param() + " character"
		}
		msg += "\n"
	}
	return msg
}