package controller

import (
	"fmt"
	"order-validation-v2/internal/controller/models"
	"order-validation-v2/internal/entity"
	"sync"
)

func (c *Controller) saveSubmission(submission models.Submission, wg *sync.WaitGroup) {
	image := models.DecodeSubmissionPayload(submission)
	_, err := c.submissions.NewSubmission(submission.Message, image, submission.TaskID)
	if err != nil {
		c.logger.ErrorLogger.Println("Error saving submission: ", err.Error())
		panic(err)
	}
	wg.Done()
}

func (c *Controller) updateTaskStatus(taskID string, wg *sync.WaitGroup, status uint8) {
	task, err := c.task.Get(taskID)
	if err != nil {
		c.logger.ErrorLogger.Println("Error retrieving task: ", err.Error())
		panic(err)
	}
	task.SetStatus(status)
	err = c.task.UpdateTask(task)
	if err != nil {
		c.logger.ErrorLogger.Println("Error updating task status: ", err.Error())
		panic(err)
	}
	wg.Done()

}

func (c *Controller) addComment(taskID string, userID string, message string, wg *sync.WaitGroup) {
	var reviewMessage entity.Message
	reviewMessage.UserID = userID
	reviewMessage.Message = message
	c.task.AddReviewMessage(taskID, reviewMessage)
	wg.Done()
}

func (c *Controller) processReviewForm(userID string, taskID string, wg *sync.WaitGroup, approved bool, forwardTo []string, message string) {
	task, err := c.task.Get(taskID)
	var wg2 sync.WaitGroup
	if err != nil {
		c.logger.ErrorLogger.Println("Error retrieving task: ", err.Error())
		panic(err)
	}
	wg2.Add(1)
	go c.addComment(taskID, userID, message, &wg2)
	task.NumOfReviewer += uint8(len(forwardTo))
	if approved {
		c.task.DeleteReviewer(userID)
		task.ReduceNumOfReviewer()
		c.task.UpdateTask(task)
		if len(forwardTo) > 0 {
			fmt.Println("Forwarding")
			wg2.Add(1)
			go c.forward(taskID, forwardTo, &wg2)
		} else if task.NumOfReviewer == 0 {
			fmt.Println("JOB DONE")
			wg2.Add(1)
			go c.updateTaskStatus(taskID, &wg2, 2)
		}
	} else {
		wg2.Add(1)
		go c.updateTaskStatus(taskID, &wg2, 0)
	}
	wg2.Wait()
	wg.Done()

}

func (c *Controller) forward(taskID string, adminIDs []string, wg *sync.WaitGroup) {
	for _, id := range adminIDs {
		fmt.Println("Forwarding to ", id)
		err := c.task.AddReviewer(taskID, id)
		if err != nil {
			c.logger.ErrorLogger.Println("Error forwarding: ", err.Error())
			panic(err)
		}
	}
	wg.Done()

}
func (c *Controller) updateRequirementStatus(requirementID int, wg *sync.WaitGroup, status int8) {
	req, err := c.requirements.GetRequirementbyID(requirementID)
	if err != nil {
		c.logger.ErrorLogger.Println("Error Getting requirement: ", err.Error())
		panic(err)
	}
	req.SetStatus(status)
	c.requirements.UpdateRequirement(req)
	wg.Done()
}

func (c *Controller) deletePrerequisite(prerequisiteTaskID string, wg *sync.WaitGroup) {
	c.logger.InfoLogger.Println("Removing Prerequisite: ", prerequisiteTaskID)
	affectedTasks, err := c.task.RemovePrerequisite(prerequisiteTaskID)
	var wg2 sync.WaitGroup
	if err != nil {
		c.logger.ErrorLogger.Println("Error deleting prerequisite: ", err.Error())
		panic(err)
	}
	for _, task := range affectedTasks {
		c.logger.InfoLogger.Println("Affected Task: ", task.ID)
		wg2.Add(1)
		go c.updateAffectedTasks(task, &wg2)
	}
	wg2.Wait()
	wg.Done()
}

func (c *Controller) updateAffectedTasks(affectedTask *entity.Task, wg *sync.WaitGroup) {
	affectedTask.ReducePrerequisite()
	c.logger.InfoLogger.Println("Remaining Prerequisite: ", affectedTask.NumOfPrerequisite)

	if affectedTask.NumOfPrerequisite == 0 {
		affectedTask.Allow()
	}
	err := c.task.UpdateTask(affectedTask)
	if err != nil {
		c.logger.ErrorLogger.Println("Error Updating task: ", err.Error())
		panic(err)
	}
	wg.Done()
}

func (c *Controller) assignPrerequiste(task *entity.Task, assignedID map[string]string, wg *sync.WaitGroup) {
	for count, prerequiste := range task.Prerequisites {
		task.Prerequisites[count] = assignedID[prerequiste]
	}
	c.task.SaveTask(task)
	wg.Done()
}
