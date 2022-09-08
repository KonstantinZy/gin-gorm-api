package models

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type Task struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	Deleted     gorm.DeletedAt `gorm:"index" json:"deleted"`
	Description string         `json:"description"`
	Name        string         `json:"name"`
	Start       sql.NullTime   `json:"start"`
	Finish      sql.NullTime   `json:"finish"`
	SubTasks    []SubTask      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type SubTask struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	Deleted     gorm.DeletedAt `gorm:"index" json:"deleted"`
	Description string         `json:"description"`
	TaskID      uint           `json:"-"`
}
