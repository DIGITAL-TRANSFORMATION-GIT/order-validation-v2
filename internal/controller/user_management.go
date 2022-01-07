package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"order-validation-v2/internal/controller/models"

	"github.com/gorilla/mux"
)

func (c *Controller) NewUser(w http.ResponseWriter, r *http.Request) {
	var newUser models.NewUser
	req, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid Request"))
		c.logger.ErrorLogger.Println("Invalid Request: ", err.Error())
		return
	}
	err = json.Unmarshal(req, &newUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid Request"))
		c.logger.ErrorLogger.Println("Invalid Request, Can't unmarshal :", err.Error())
		return
	}
	exists, err := c.user.ValidateUsername(newUser.Username)
	if exists {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("Username Exists"))
		return
	}
	id, err := c.user.CreateUser(newUser.Username, newUser.Email, newUser.Password, newUser.Role)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		c.logger.ErrorLogger.Println("Error while creating new user: ", err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("User %s has been added with id %s\n", newUser.Username, id)))
}

func (c *Controller) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := c.user.ListUsers()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		c.logger.ErrorLogger.Println("Error retrieving all user: ", err.Error())
		return
	}
	var response []models.RetrievedUser
	for _, user := range users {
		response = append(response, models.BuildUserProfile(user))
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (c *Controller) DeleteUser(w http.ResponseWriter, r *http.Request) {
	err := c.user.DeleteUser(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		c.logger.ErrorLogger.Println("Error while deleting task : ", err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	return

}
