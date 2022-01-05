package repository

import (
	"database/sql"
	"order-validation-v2/internal/entity"
)

type TaskPSQL struct {
	db *sql.DB
}

func NewTaskPSQL(db *sql.DB) *TaskPSQL {
	return &TaskPSQL{
		db: db,
	}
}

func (r *TaskPSQL) AddPrerequisite(taskID string, prerequisite string) error {
	stmt, err := r.db.Prepare(`INSERT INTO prerequisite (task_id, prerequisite) values ($1,$2)`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(taskID, prerequisite)
	if err != nil {
		return err
	}

	return nil

}
func (r *TaskPSQL) Create(t *entity.Task) (string, error) {
	stmt, err := r.db.Prepare(`
		INSERT INTO tasks (assigner_id, ID, user_id, requirement_id, note, fulfillment_status, allowed, deadline, num_of_prerequisite, total_reviewer) 
		values($1,$2,$3,$4,$5,$6,$7,$8, $9, $10)`)

	if err != nil {
		return t.ID, err
	}

	_, err = stmt.Exec(
		t.AssignerID,
		t.ID,
		t.UserID,
		t.RequirementID,
		t.Note,
		t.Status,
		t.Allowed,
		t.Deadline,
		t.NumOfPrerequisite,
		t.NumOfReviewer,
	)
	if err != nil {
		return t.ID, err
	}
	return t.ID, nil

}

func (r *TaskPSQL) GetByOrderID(orderID string) ([]*entity.TaskWithDetails, error) {
	stmt, err := r.db.Prepare(`SELECT tasks.id, tasks.note, users.username, tasks.deadline, requirements.request, 
								requirements.expected_outcome,orders.title 
								FROM tasks INNER JOIN requirements ON tasks.requirement_id=requirements.id 
								INNER JOIN users ON users.id = tasks.user_id
								INNER JOIN orders ON requirements.order_id = orders.id 
								WHERE orders.id = $1`)
	if err != nil {
		return nil, err
	}
	var tasks []*entity.TaskWithDetails
	rows, err := stmt.Query(orderID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var task entity.TaskWithDetails
		err = rows.Scan(&task.ID, &task.Note, &task.Username, &task.Deadline, &task.Request, &task.ExpectedOutcome, &task.OrderTitle)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, &task)
	}
	return tasks, nil

}

func (r *TaskPSQL) GetPrerequisites(taskID string) ([]string, error) {
	stmt, err := r.db.Prepare(`SELECT requirements.request
							  	FROM prerequisite INNER JOIN tasks on prerequisite.task_id = tasks.id
								INNER JOIN requirements ON tasks.requirement_id = requirements.id
								WHERE prerequisite = $1`)
	if err != nil {
		return nil, err
	}
	var prerequisites []string
	rows, err := stmt.Query(taskID)
	for rows.Next() {
		var prerequisite string
		err = rows.Scan(&prerequisite)
		if err != nil {
			return nil, err
		}
		prerequisites = append(prerequisites, prerequisite)
	}
	return prerequisites, nil
}

func (r *TaskPSQL) RemovePrerequisite(taskID string) ([]*entity.Task, error) {
	stmt, err := r.db.Prepare(`SELECT tasks.id, tasks.allowed, tasks.user_id, tasks.fulfillment_status, tasks.num_of_prerequisite, tasks.deadline
							  	FROM prerequisite INNER JOIN tasks on prerequisite.task_id = tasks.id
								WHERE prerequisite = $1`)

	if err != nil {
		return nil, err
	}
	var affectedTasks []*entity.Task
	rows, err := stmt.Query(taskID)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var t entity.Task
		err = rows.Scan(&t.ID, &t.Allowed, &t.UserID, &t.Status, &t.NumOfPrerequisite, &t.Deadline)
		if err != nil {
			return nil, err
		}
		affectedTasks = append(affectedTasks, &t)
	}
	_, err = r.db.Exec("DELETE FROM prerequisite WHERE prerequisite=$1", taskID)
	if err != nil {
		return nil, err
	}
	return affectedTasks, nil
}

func (r *TaskPSQL) Get(id string) (*entity.Task, error) {
	stmt, err := r.db.Prepare(`SELECT id, requirement_id, allowed, user_id, fulfillment_status, num_of_prerequisite, deadline 
								from tasks where id = $1`)
	var task entity.Task
	if err != nil {
		return nil, err
	}

	row := stmt.QueryRow(id)
	if row == nil {
		return nil, err
	}
	err = row.Scan(&task.ID, &task.RequirementID, &task.Allowed, &task.UserID,
		&task.Status, &task.NumOfPrerequisite, &task.Deadline)
	if err != nil {
		return nil, err
	}
	return &task, nil

}

func (r *TaskPSQL) GetbyUserID(userID string) ([]*entity.TaskWithDetails, error) {
	stmt, err := r.db.Prepare(`SELECT tasks.id, tasks.note, tasks.deadline, requirements.request, requirements.expected_outcome,  
								orders.title, orders.description, orders.deadline,tasks.fulfillment_status
								FROM tasks INNER JOIN requirements ON tasks.requirement_id=requirements.id 
								INNER JOIN orders ON requirements.order_id = orders.id 
								where user_id = $1 and tasks.allowed = true`)
	if err != nil {
		return nil, err
	}
	var tasks []*entity.TaskWithDetails
	rows, err := stmt.Query(userID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var t entity.TaskWithDetails
		err = rows.Scan(&t.ID, &t.Note, &t.Deadline, &t.Request, &t.ExpectedOutcome, &t.OrderTitle, &t.OrderDescription, &t.OrderDeadline,
			&t.Status)
		if err != nil {
			return nil, err
		}
		/*
			t.Messages, err = r.GetReviewMessages(t.ID)
			if err != nil {
				return nil, err
			}*/
		tasks = append(tasks, &t)
	}

	return tasks, nil

}
func (r *TaskPSQL) Update(e *entity.Task) error {
	_, err := r.db.Exec(`UPDATE tasks SET user_id = $1, fulfillment_status = $2, deadline = $3, num_of_prerequisite = $4,
						 allowed = $5, total_reviewer = $6 where id = $7`,
		e.UserID, e.Status, e.Deadline, e.NumOfPrerequisite, e.Allowed, e.NumOfReviewer, e.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *TaskPSQL) List() ([]*entity.TaskWithDetails, error) {
	stmt, err := r.db.Prepare(`SELECT tasks.id, tasks.note, tasks.deadline, users.username, requirements.request, requirements.expected_outcome,  
								orders.title, orders.description, orders.deadline, tasks.fulfillment_status, tasks.num_of_prerequisite 
								FROM tasks INNER JOIN requirements ON tasks.requirement_id=requirements.id 
								INNER JOIN users on users.id = tasks.user_id
								INNER JOIN orders ON requirements.order_id = orders.id `)
	if err != nil {
		return nil, err
	}
	var tasks []*entity.TaskWithDetails
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var t entity.TaskWithDetails
		err = rows.Scan(&t.ID, &t.Note, &t.Deadline, &t.Username, &t.Request, &t.ExpectedOutcome, &t.OrderTitle,
			&t.OrderDescription, &t.OrderDeadline, &t.Status, &t.NumOfPrerequisite)
		if err != nil {
			return nil, err
		}
		t.Messages, err = r.GetReviewMessages(t.ID)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, &t)
	}
	if len(tasks) == 0 {
		return nil, nil
	}

	return tasks, nil

}

func (r *TaskPSQL) GetTasksToReview(adminID string) ([]*entity.TaskWithDetails, error) {
	stmt, err := r.db.Prepare(`SELECT tasks.id, tasks.note, tasks.deadline, users.username, requirements.request, requirements.expected_outcome,  
								orders.title, orders.description, orders.deadline 
								FROM tasks INNER JOIN requirements ON tasks.requirement_id=requirements.id 
								INNER JOIN users ON users.id = tasks.user_id
								LEFT JOIN forwarded_review ON tasks.id = forwarded_review.task_id
								INNER JOIN orders ON requirements.order_id = orders.id 
								WHERE tasks.fulfillment_status = 1 and (tasks.assigner_id = $1 or forwarded_review.reviewer_id = $1)`)
	if err != nil {
		return nil, err
	}
	var tasks []*entity.TaskWithDetails
	rows, err := stmt.Query(adminID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var t entity.TaskWithDetails
		err = rows.Scan(&t.ID, &t.Note, &t.Deadline, &t.Username, &t.Request, &t.ExpectedOutcome, &t.OrderTitle, &t.OrderDescription, &t.OrderDeadline)
		if err != nil {
			return nil, err
		}
		/*
			t.Messages, err = r.GetReviewMessages(t.ID)
			if err != nil {
				return nil, err
			}*/
		tasks = append(tasks, &t)
	}
	if len(tasks) == 0 {
		return nil, nil
	}

	return tasks, nil

}

func (r *TaskPSQL) Delete(TaskID string) error {
	_, err := r.db.Exec("DELETE FROM task where requirement_id = $1", TaskID)
	if err != nil {
		return err
	}
	return nil
}

func (r *TaskPSQL) AddReviewer(TaskID string, NewReviewerID string) error {
	stmt, err := r.db.Prepare(`INSERT INTO forwarded_review (task_id, reviewer_id) VALUES ($1,$2)`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(TaskID, NewReviewerID)
	if err != nil {
		return err
	}
	return nil
}

func (r *TaskPSQL) AddReviewMessage(TaskID string, Message entity.Message) error {
	stmt, err := r.db.Prepare(`INSERT INTO submission_messages (task_id, user_id, msg) VALUES ($1,$2,$3)`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(TaskID, Message.UserID, Message.Message)
	if err != nil {
		return err
	}
	return nil
}

func (r *TaskPSQL) GetReviewMessages(TaskID string) ([]entity.Message, error) {
	stmt, err := r.db.Prepare(`SELECT username, message FROM submission_messages 
	INNER JOIN users ON users.id = submission_messages.user_id 
	WHERE task_id = $1`)
	if err != nil {
		return nil, err
	}
	var messages []entity.Message
	rows, err := stmt.Query(TaskID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var message entity.Message
		err = rows.Scan(&message.Username, &message.Message)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, nil
}

func (r *TaskPSQL) DeleteReviewer(UserID string) error {
	stmt, err := r.db.Prepare(`DELETE FROM forwarded_review WHERE reviewer_id = $1`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(UserID)
	if err != nil {
		return err
	}
	return nil
}
