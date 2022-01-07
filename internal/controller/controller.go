package controller

import (
	"net/http"
	"order-validation-v2/internal/usecase/orders"
	"order-validation-v2/internal/usecase/requirements"
	"order-validation-v2/internal/usecase/submissions"
	"order-validation-v2/internal/usecase/tasks"
	"order-validation-v2/internal/usecase/user"
	"order-validation-v2/pkg/logger"
	"os"

	"github.com/rs/cors"

	"github.com/gorilla/mux"
)

type Controller struct {
	router       *mux.Router
	order        orders.UseCase
	user         user.UseCase
	task         tasks.UseCase
	submissions  submissions.UseCase
	requirements requirements.UseCase
	logger       *logger.LoggerInstance
}

func NewController(o orders.UseCase, u user.UseCase, r requirements.UseCase, t tasks.UseCase, s submissions.UseCase, l *logger.LoggerInstance) *Controller {
	router := mux.NewRouter().StrictSlash(true)
	controller := &Controller{router: router, order: o, user: u, requirements: r, task: t, submissions: s, logger: l}
	return controller
}

func (c *Controller) RegisterHandler() {
	login := c.router.PathPrefix("/login").Subrouter()
	login.HandleFunc("/", c.Login).Methods("POST")

	userapp := c.router.PathPrefix("/orders").Subrouter()
	userapp.Use(c.validateUserJWT)
	userapp.HandleFunc("/", c.GetTasks).Methods("GET")
	userapp.HandleFunc("/profile", c.GetUserProfile).Methods("GET")
	userapp.HandleFunc("/profile/passwordchange", c.ChangePassword).Methods("POST")
	userapp.HandleFunc("/profile/usernamechange", c.ChangeUsername).Methods("POST")
	userapp.HandleFunc("/task={id}", c.GetSubmission).Methods("GET")
	userapp.HandleFunc("/submission", c.PostSubmission).Methods("POST")
	userapp.HandleFunc("/submission/id={id}", c.UpdateSubmission).Methods("POST")

	admin := c.router.PathPrefix("/admin").Subrouter()
	admin.Use(c.validateAdminJWT)
	admin.HandleFunc("/orders", c.GetAllUncompletedOrders).Methods("GET")
	admin.HandleFunc("/orders", c.AddNewOrder).Methods("POST")

	admin.HandleFunc("/orders/id={id}", c.GetStatusOfOrder).Methods("GET")
	admin.HandleFunc("/orders/id={id}", c.DeleteOrder).Methods("DELETE")
	admin.HandleFunc("/orders/id={id}", c.ModifyRequirements).Methods("PATCH")
	admin.HandleFunc("/orders/search:{query}", c.SearchOrders).Methods("GET")
	admin.HandleFunc("/user", c.NewUser).Methods("POST")
	admin.HandleFunc("/user", c.GetAllUsers).Methods("GET")
	admin.HandleFunc("/user/id={id}", c.DeleteUser).Methods("DELETE")
	admin.HandleFunc("/user/id={id}/tasks", c.GetTasksOfUser).Methods("GET")
	admin.HandleFunc("/tasks", c.GetAllAssignedTasks).Methods("GET")
	admin.HandleFunc("/tasks", c.AddNewTask).Methods("POST")
	admin.HandleFunc("/tasks/id={id}", c.DeleteTask).Methods("DELETE")
	admin.HandleFunc("/tasks/bulk", c.BulkAssignTasks).Methods("POST")
	admin.HandleFunc("/tasks/order={id}", c.GetTasksOnSpecificOrder).Methods("GET")
	admin.HandleFunc("/tasks/submitted", c.GetTaskstoReview).Methods("GET")
	admin.HandleFunc("/submission={id}/review", c.ReviewSubmission).Methods("POST")
}

func (c *Controller) Start() {
	cors := cors.AllowAll()
	port := os.Getenv("PORT")
	handler := cors.Handler(c.router)
	http.ListenAndServe(":"+port, handler)
}

func (c *Controller) StartLocally() {
	cors := cors.AllowAll()
	handler := cors.Handler(c.router)
	http.ListenAndServe(":"+"8080", handler)
}
