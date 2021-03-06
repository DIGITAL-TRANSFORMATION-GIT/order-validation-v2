package submissions

import (
	"order-validation-v2/internal/entity"
)

type Reader interface {
	GetByTaskID(taskID string) ([]*entity.Submission, error)
	Get(submissionID string) (*entity.Submission, error)
}

type Writer interface {
	Create(e *entity.Submission) (string, error)
	Update(e *entity.Submission) error
	Delete(id string) error
}

type Repository interface {
	Reader
	Writer
}

type UseCase interface {
	NewSubmission(message string, images []string, TaskID string) (string, error)
	EditSubmission(e *entity.Submission) error
	DeleteSubmission(id string) error
	GetSubmissionByTaskID(taskID string) ([]*entity.Submission, error)
	GetSubmission(submissionID string) (*entity.Submission, error)
}
