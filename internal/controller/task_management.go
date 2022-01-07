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

	"github.com/gorilla/mux"
)

func (c *Controller) BulkAssignTasks(w http.ResponseWriter, r *http.Request) {
	adminID := fmt.Sprintf("%v", r.Context().Value(ctxKey{}))
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
		t := entity.NewTask(adminID, task.RequirementID, task.UserID, task.Note, task.Prerequisite, deadline)
		assignedID[task.Num] = t.ID
		tasks = append(tasks, t)
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("%d Task has been added", len(tasks))))
	var wg sync.WaitGroup
	for _, task := range tasks {
		wg.Add(1)
		go c.assignPrerequiste(task, assignedID, &wg)
	}
	wg.Wait()

}
func (c *Controller) AddNewTask(w http.ResponseWriter, r *http.Request) {
	adminID := fmt.Sprintf("%v", r.Context().Value(ctxKey{}))
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
	id, err := c.task.CreateTask(adminID, newTask.RequirementID, newTask.UserID, newTask.Note, newTask.Prerequisite, deadline)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		c.logger.ErrorLogger.Println("Error creating new task: ", err.Error())
		return
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go c.updateRequirementStatus(newTask.RequirementID, &wg, 1)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("Task %s has been created for user %s\n", id, newTask.UserID)))
	wg.Wait()
}

func (c *Controller) GetAllAssignedTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := c.task.ListAllTasks()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		c.logger.ErrorLogger.Println("Error retrieving all tasks: ", err.Error())
	}
	for _, task := range tasks {
		if task.NumOfPrerequisite != 0 {
			prerequisites, err := c.task.GetPrerequisites(task.ID)
			if err != nil {
				c.logger.ErrorLogger.Println("Can't retrieve prerequisites: ", err.Error())
			}
			task.Prerequisites = prerequisites
		}
	}
	response := models.BuildTasks(tasks)
	json.NewEncoder(w).Encode(response)
}

func (c *Controller) GetTaskstoReview(w http.ResponseWriter, r *http.Request) {
	adminID := fmt.Sprintf("%v", r.Context().Value(ctxKey{}))
	tasks, err := c.task.GetTasksToReview(adminID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		c.logger.ErrorLogger.Println("Error retrieving all tasks: ", err.Error())
	}
	response := models.BuildTasks(tasks)
	json.NewEncoder(w).Encode(response)
}

func (c *Controller) GetTasksOfUser(w http.ResponseWriter, r *http.Request) {
	tasks, err := c.task.GetTasksofUser(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		c.logger.ErrorLogger.Println("Error while getting tasks: ", err.Error())
		return
	}
	if tasks == nil {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("No Tasks Present"))
		return
	}
	response := models.BuildTasks(tasks)
	json.NewEncoder(w).Encode(response)
}

func (c *Controller) GetTasksOnSpecificOrder(w http.ResponseWriter, r *http.Request) {
	tasks, err := c.task.GetTasksOnSpecificOrder(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		c.logger.ErrorLogger.Println("Error while getting tasks: ", err.Error())
		return
	}
	if len(tasks) == 0 {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("No Tasks Present"))
		return
	}
	response := models.BuildTasks(tasks)
	json.NewEncoder(w).Encode(response)
}

func (c *Controller) DeleteTask(w http.ResponseWriter, r *http.Request) {
	err := c.task.DeleteTask(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		c.logger.ErrorLogger.Println("Error while deleting task : ", err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	return
}
