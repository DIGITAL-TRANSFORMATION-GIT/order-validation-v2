package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"order-validation-v2/internal/controller/models"
)

func (c *Controller) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	auth := r.Context().Value(ctxKey{})
	userID := fmt.Sprintf("%v", auth)
	user, err := c.user.GetUserbyID(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		c.logger.ErrorLogger.Println("Error while retrieving user info: ", err.Error())
		return
	}
	if user == nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("User Not Found"))
		return
	}
	response := models.BuildUserProfile(user)
	json.NewEncoder(w).Encode(response)

}

func (c *Controller) ChangePassword(w http.ResponseWriter, r *http.Request) {
	auth := r.Context().Value(ctxKey{})
	userID := fmt.Sprintf("%v", auth)
	req, err := ioutil.ReadAll(r.Body)
	var form models.ChangePasswordForm
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		c.logger.ErrorLogger.Println("Error while Logging in: ", err.Error())
		w.Write([]byte("Invalid Request"))
		return
	}
	err = json.Unmarshal(req, &form)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		c.logger.ErrorLogger.Println("Error while Logging in: ", err.Error())
		w.Write([]byte("Invalid Request"))
		return
	}
	authorize, user, err := c.user.ValidateAndRetrieveUser(userID, form.OldPassword)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		c.logger.ErrorLogger.Println("Error while validating user : ", err.Error())
		return
	}
	if authorize {
		user.Password = form.NewPassword
		err := c.user.UpdateUser(user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			c.logger.ErrorLogger.Println("Error while updating user : ", err.Error())
			return
		}
	}
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte("Invalid Old Password"))

}

func (c *Controller) ChangeUsername(w http.ResponseWriter, r *http.Request) {
	auth := r.Context().Value(ctxKey{})
	userID := fmt.Sprintf("%v", auth)
	req, err := ioutil.ReadAll(r.Body)
	var form models.ChangeUsernameForm
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		c.logger.ErrorLogger.Println("Error while Logging in: ", err.Error())
		w.Write([]byte("Invalid Request"))
		return
	}
	err = json.Unmarshal(req, &form)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		c.logger.ErrorLogger.Println("Error while Logging in: ", err.Error())
		w.Write([]byte("Invalid Request"))
		return
	}
	authorize, user, err := c.user.ValidateAndRetrieveUser(userID, form.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		c.logger.ErrorLogger.Println("Error while validating user : ", err.Error())
		return
	}
	if authorize {
		exists, err := c.user.ValidateUsername(form.NewUsername)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			c.logger.ErrorLogger.Println("Error while checking username availability : ", err.Error())
			return
		}
		if exists {
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte("Username exists"))
			return
		}
		user.Username = form.NewUsername
		err = c.user.UpdateUser(user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			c.logger.ErrorLogger.Println("Error while updating user : ", err.Error())
			return
		}
	}
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte("Invalid Old Password"))

}
