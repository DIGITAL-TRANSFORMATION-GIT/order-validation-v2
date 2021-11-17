package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"order-validation-v2/internal/controller/models"
	"order-validation-v2/internal/entity"
	"sync"
	"time"
)

func (c *Controller) BulkAssignTasks(w http.ResponseWriter, r *http.Request) {
	var newTasks models.BulkAddedTasks
	req, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid Request"))
		return
	}
	err = json.Unmarshal(req, &newTasks)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid Request"))
		return
	}
	assignedID := map[string]string{}
	var tasks []*entity.Task
	for _, task := range newTasks.Tasks {
		deadline, _ := time.Parse("2/Jan/2006 15:04:05", task.Deadline)
		t := entity.NewTask(task.RequirementID, task.UserID, task.Note, task.Prerequisite, deadline)
		assignedID[task.Num] = t.ID
		tasks = append(tasks, t)
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("%d Task has been added", len(tasks))))
	var wg sync.WaitGroup
	for _, task := range tasks {
		wg.Add(1)
		go c.assignPrerequiste(task, assignedID, wg)
	}
	wg.Wait()

}
func (c *Controller) AddNewTask(w http.ResponseWriter, r *http.Request) {
	var newTask models.NewTask
	req, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid Request"))
		return
	}
	err = json.Unmarshal(req, &newTask)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid Request"))
		return
	}
	deadline, _ := time.Parse("2/Jan/2006 15:04:05", newTask.Deadline)
	id, err := c.task.CreateTask(newTask.RequirementID, newTask.UserID, newTask.Note, newTask.Prerequisite, deadline)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		c.logger.ErrorLogger.Println("Error creating new task: ", err.Error())
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("Task %s has been created for user %s\n", id, newTask.UserID)))
}

func (c *Controller) GetAllAssignedTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := c.task.ListAllTasks()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		c.logger.ErrorLogger.Println("Error retrieving all tasks: ", err.Error())
	}
	response := models.BuildTasks(tasks)
	json.NewEncoder(w).Encode(response)
}
