package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"example/api/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CreateTaskInput struct {
	Name        string              `json:"name" binding:"required"`
	Description string              `json:"description" binding:"required"`
	Start       int64               `json:"start"`
	Finish      int64               `json:"finish"`
	Subtask     []CreateTaskSubtask `json:"subtasks"`
}

type CreateTaskSubtask struct {
	Description string `json:"description"`
}

// POST /task
// create task
func CreateOneTask(c *gin.Context) {
	var input CreateTaskInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	start := sql.NullTime{}
	if err := start.Scan(time.Unix(input.Start, 0)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	finish := sql.NullTime{}
	if err := finish.Scan(time.Unix(input.Finish, 0)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	task := models.Task{
		Name:        input.Name,
		Description: input.Description,
		Start:       start,
		Finish:      finish,
	}

	if len(input.Subtask) > 0 {
		for _, st := range input.Subtask {
			task.SubTasks = append(task.SubTasks, models.SubTask{Description: st.Description})
		}
	}

	if err := models.DB.Create(&task).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"data": "Ok"})
	}
}

type UpdateTaskInput struct {
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Start       int64                `json:"start"`
	Finish      int64                `json:"finish"`
	SubTasks    []CUpdateTaskSubtask `json:"subtasks"`
}

type CUpdateTaskSubtask struct {
	ID          int64  `json:"id"`
	Description string `json:"description"`
}

// PATCH /task/:id
// update task
func UpdateTask(c *gin.Context) {
	var input UpdateTaskInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var task models.Task
	if err := models.DB.Where("id = ?", c.Param("id")).Find(&task).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Can't find task for update" + err.Error()})
	}

	if task.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Can't find task for update"})
		return
	}

	if input.Name != "" {
		task.Name = input.Name
	}

	if input.Description != "" {
		task.Description = input.Description
	}

	if input.Start != 0 {
		start := sql.NullTime{}
		if err := start.Scan(time.Unix(input.Start, 0)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		task.Start = start
	}

	if input.Finish != 0 {
		finish := sql.NullTime{}
		if err := finish.Scan(time.Unix(input.Finish, 0)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		task.Finish = finish
	}

	if len(input.SubTasks) > 0 {
		for _, st := range input.SubTasks {
			sbTask := models.SubTask{}
			if st.ID > 0 {
				models.DB.Model(&task).Where("id=?", st.ID).Association("SubTasks").Find(&sbTask)
				sbTask.Description = st.Description
			} else {
				sbTask.Description = st.Description
			}

			task.SubTasks = append(task.SubTasks, sbTask)
		}
	}

	if err := models.DB.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&task).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"data": "Ok"})
	}
}

// GET /task
// get all tasks
func Get(c *gin.Context) {
	var task2 []models.Task
	err := models.DB.Model(&models.Task{}).Preload("SubTasks").Find(&task2).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"data": task2})
	}
}

// GET /task/:id
// get task by it id
func GetOneTask(c *gin.Context) {
	var task models.Task
	err := models.DB.Model(&models.Task{}).Preload("SubTasks").Where("id = ?", c.Param("id")).First(&task).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"data": task})
	}
}

// DELETE /task/:id
// delete task with subtasks
func DeleteTask(c *gin.Context) {
	var task models.Task

	if err := models.DB.Where("id = ?", c.Param("id")).First(&task).Error; err == nil {
		models.DB.Select("SubTasks").Delete(&task)
		c.JSON(http.StatusOK, gin.H{"data": "Ok"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprint("Can't delet elment:", err)})
	}

}
