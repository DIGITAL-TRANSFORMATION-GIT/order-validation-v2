package tasks

import (
	"order-validation-v2/internal/entity"
	"time"
)

type Service struct {
	repo Repository
}

func NewService(r Repository) *Service {
	return &Service{
		repo: r,
	}
}

func (s *Service) ListAllTasks() ([]*entity.TaskWithDetails, error) {
	return s.repo.List()
}

func (s *Service) GetTasksofUser(userID string) ([]*entity.TaskWithDetails, error) {
	return s.repo.GetbyUserID(userID)

}

func (s *Service) Get(id string) (*entity.Task, error) {
	return s.repo.Get(id)
}

func (s *Service) UpdateTask(t *entity.Task) error {
	return s.repo.Update(t)
}

func (s *Service) DeleteTask(id string) error {
	return s.repo.Delete(id)
}

func (s *Service) CreateTask(assignerID string, requirementID int, userID string, Note string, prerequisiteTaskID []string, Deadline time.Time) (string, error) {
	task := entity.NewTask(assignerID, requirementID, userID, Note, prerequisiteTaskID, Deadline)
	taskID, err := s.repo.Create(task)
	for _, prerequisite := range prerequisiteTaskID {
		err = s.repo.AddPrerequisite(taskID, prerequisite)
		if err != nil {
			return "", err
		}
	}
	return taskID, err
}

func (s *Service) SaveTask(task *entity.Task) (string, error) {
	return s.repo.Create(task)
}

func (s *Service) RemovePrerequisite(prerequisiteID string) ([]*entity.Task, error) {
	return s.repo.RemovePrerequisite(prerequisiteID)
}

func (s *Service) GetTasksToReview(adminID string) ([]*entity.TaskWithDetails, error) {
	return s.repo.GetTasksToReview(adminID)
}

func (s *Service) AddReviewer(TaskID string, NewReviewerID string) error {
	return s.repo.AddReviewer(TaskID, NewReviewerID)
}

func (s *Service) GetTasksOnSpecificOrder(orderID string) ([]*entity.TaskWithDetails, error) {
	return s.repo.GetByOrderID(orderID)
}

func (s *Service) GetPrerequisites(TaskID string) ([]string, error) {
	return s.repo.GetPrerequisites(TaskID)
}

func (s *Service) DeleteReviewer(UserID string) error {
	return s.repo.DeleteReviewer(UserID)
}

func (s *Service) AddReviewMessage(TaskID string, Message entity.Message) error {
	return s.repo.AddReviewMessage(TaskID, Message)
}
